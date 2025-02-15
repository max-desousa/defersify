package defersification

import (
  "fmt"
  "bufio"
  "os"
  "regexp"
  "strings"
)

func SeachForDefers(_filepath string) {
  deferStack := make(map[int][]string)
  var levelOfNest int = 0
  withinBlockDefers := false
  withinBlockComment := false
  lineNumber := 0

  deferStatementRegex := regexp.MustCompile(`^\s*defer\s+`)
  deferStatementWithSpacingRegex := regexp.MustCompile(`defer\s+`)
  startBlockDeferRegex := regexp.MustCompile(`^\s*defer\s*{?\s*$`)
  returnRegex := regexp.MustCompile(`^\s*return\s+`)
  


  scanner := bufio.NewScanner(readFile)

  for scanner.Scan() {
    fmt.Printf("The defer stack at line %v is: %v\n", lineNumber, deferStack)
    lineNumber++
    line := scanner.Text()

    /* write to the output file any line that won't affect the nature of a defer */
    if (obviouslyFineLine(line)) {
      fmt.Printf("Writing line %v to the file because it's obviolsy fine\n", lineNumber)
      writeFile.WriteString(line + "\n")
      continue
    }

    /* if the line has a multi line block comment handle that here using a flag
     * that passes all lines until block comment is closed */
    if ( strings.Contains(line, "/*") || withinBlockComment ) {
      withinBlockComment = true
      writeFile.WriteString(line + "\n")
      if ( strings.Contains(line, "*/") ) {
        withinBlockComment = false
      }
      continue
    }

    /* Not sure if there is a use case for multiple {'s in one line... therefore
     * throwing an error message if this is the case... perhaps a future feature */
    if strings.Count(line, "{") > 1 {
      fmt.Println("Found a line with more than one opening brace... this is logic that is being put off for minimal viable product")
      break
    }

    /* set flag if we're in a defer block */
    if ( startBlockDeferRegex.MatchString(line) ) {
      withinBlockDefers = true
      continue
    }
    if ( withinBlockDefers ) {
      if ( strings.Contains(line, "}") ) {
        withinBlockDefers = false
      } else {
        deferStack[levelOfNest-1] = append(deferStack[levelOfNest-1], line)
      }
      continue
    }

    if ( strings.Contains(line, "{") ) {

      levelOfNest++
      fmt.Printf("Incrementing level of nest to: %v due to line number %v\n", levelOfNest, lineNumber)
      writeFile.WriteString(line + "\n")

    } else if ( deferStatementRegex.MatchString(line) ) {

      modifiedLine := deferStatementWithSpacingRegex.ReplaceAllString(line, "")
      deferStack[levelOfNest-1] = append(deferStack[levelOfNest-1], modifiedLine)

    } else if ( returnRegex.MatchString(line) ) {

      fmt.Printf("Return statement hit at level of nest: %v and a deferStack of %v\n", levelOfNest, deferStack)
      for i := levelOfNest; i > 0; i-- {
        for j := len(deferStack[i-1]); j > 0; j-- {
          writeFile.WriteString(padDeferStatementForNest(levelOfNest, deferStack[i-1][j-1]) + "\n")
        }
      }
      writeFile.WriteString(line + "\n")
      if ( levelOfNest == 1 ) {
        deferStack = make(map[int][]string)
      }

    } else if ( strings.Contains(line, "}") ) {
      
      if ( levelOfNest > 0 ) {
        for i := len(deferStack[levelOfNest-1]); i > 0; i-- {
          writeFile.WriteString(deferStack[levelOfNest-1][i-1] + "\n")
          writeFile.WriteString(padDeferStatementForNest(levelOfNest, deferStack[levelOfNest-1][i-1]) + "\n")
        }
        deferStack[levelOfNest-1] = nil
        levelOfNest--
        fmt.Printf("Decrementing level of nest to: %v due to line number %v\n", levelOfNest, lineNumber)
      }
      writeFile.WriteString(line + "\n")

    } else {
      writeFile.WriteString(line + "\n")
    }

    if ( levelOfNest < 0 ) {
      deferStack = make(map[int][]string)
    }
  }
}

/*******************************************************************************
 * Function: detectCloseToBlockComment
 ******************************************************************************/
func detectCloseToBlockComment(_line string) bool {
  var returnVal bool = false

  blockCommentCloseRegex := regexp.MustCompile(`\*/`)

  if (blockCommentCloseRegex.MatchString(_line)) {
    returnVal = true
  }

  return returnVal
}

/*******************************************************************************
 * Function: detectOpenButNotCompleteBlockComment
 * 
 * Description: This function will detect if future lines will be comments.
 ******************************************************************************/
func detectOpenButNotCompleteBlockComment(_line string) bool {
  var returnVal bool = false

  blockCommentOpenRegex := regexp.MustCompile(`/\*`)
  blockCommentCompletedRegex := regexp.MustCompile(`/\*.*\*/`)

  /* Check that there is any sort of block comment in the line */
  if (blockCommentOpenRegex.MatchString(_line)) {
    if ( blockCommentCompletedRegex.MatchString(_line) ) {
      returnVal = false
    } else {
      returnVal = true
    }
  }

  return returnVal
}

/*******************************************************************************
* Function: obviouslyFineLine
 *
 * Description: This function is the first line of defense in determining if a
 * line from the source file can be written to the new file without any needs to
 * cache the information in the deferStack. This function will be used to bypass
 * any line that:
 * 1. Does not contain a valid defer statement
 * 2. Does not contain conflate start a block comment
 *   a. this is to prevent a "defer" statement from within a block comment from
 *      being considered a valid defer statement.
 * 3. There is no opening brace in the line
 ******************************************************************************/
func obviouslyFineLine(_line string) bool {
  /* assume line is safe first because it's easier to lock as false */
  var returnVal bool = true

  deferRegex := regexp.MustCompile(`^\s*defer\s+`)
  openBraceRegex := regexp.MustCompile(`.*[{}].*`)
  returnStatementRegex := regexp.MustCompile(`^\s*return\s+`)
  
  if ( !deferRegex.MatchString(_line) &&
       !detectOpenButNotCompleteBlockComment(_line) &&
       !openBraceRegex.MatchString(_line) &&
       !returnStatementRegex.MatchString(_line) ) {
    returnVal = true
  } else {
    returnVal = false
  }

  return returnVal
}


/*******************************************************************************
 * simpleFileRename
 *
 * Description: This function will take a file path and develops a similar new
 * name for a file that will effectively be a converted version of the original
 * 
 * outpus will be original.c -> original_defersified.c
 ******************************************************************************/
func simpleFileRename(_filepath string) string {
  var outputString string
  var indexOfFileExtensionPeriod int
  var fileExtension string
  var nonExtensionString string
 
  indexOfFileExtensionPeriod = strings.LastIndex(_filepath, "deferable_")

  nonExtensionString = _filepath[:indexOfFileExtensionPeriod]
  fileExtension = _filepath[(indexOfFileExtensionPeriod + len("deferable_")):]
  fmt.Printf("nonExtensionString: %v\n", nonExtensionString)
  fmt.Printf("fileExtension: %v\n", fileExtension)

  outputString = nonExtensionString + fileExtension
  return outputString
}

