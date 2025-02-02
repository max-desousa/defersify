package deferSearcher

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
  //withinBlockDefers := false
  withinBlockComment := false

  deferStatementRegex := regexp.MustCompile(`^\s*defer\s+`)
  deferStatementWithSpacingRegex := regexp.MustCompile(`defer\s+`)
  //deferRegex := regexp.MustCompile(`defer\s+`)
  returnRegex := regexp.MustCompile(`^\s*return\s+`)
  
  readFile, err := os.Open(_filepath)
  defer readFile.Close()

  if err != nil {
    fmt.Println("Error opening source file")
    return
  }

  writeFile, err := os.Create(simpleFileRename(_filepath))
  defer writeFile.Close()

  if err != nil {
    fmt.Println("Error creating new file")
    return
  }

  scanner := bufio.NewScanner(readFile)

  for scanner.Scan() {
    line := scanner.Text()

    if (obviouslyFineLine(line)) {
      writeFile.WriteString(line + "\n")
      continue
    }

    if ( strings.Contains(line, "/*") || withinBlockComment ) {
      withinBlockComment = true
      writeFile.WriteString(line + "\n")
      if ( strings.Contains(line, "*/") ) {
        withinBlockComment = false
      }
      continue
    }

    if strings.Count(line, "{") > 1 {
      fmt.Println("Found a line with more than one opening brace... this is logic that is being put off for minimal viable product")
      break
    }

    if ( strings.Contains(line, "{") ) {
      if (deferStatementRegex.MatchString(line)) {
      } else {
        levelOfNest++
        fmt.Println("Incrementing level of nest")
        writeFile.WriteString(line + "\n")
      }
    } else if ( deferStatementRegex.MatchString(line) ) {
      modifiedLine := deferStatementWithSpacingRegex.ReplaceAllString(line, "")
      deferStack[levelOfNest-1] = append(deferStack[levelOfNest-1], modifiedLine)
    } else if ( strings.Contains(line, "}") ) {
      if ( len(deferStack[levelOfNest-1]) > 0 ) {
        for _, val := range deferStack[levelOfNest-1] {
          writeFile.WriteString(val + "\n")
          fmt.Println("Writing from defer stack...")
        }
        //fmt.Println("Debugging start deferStack Analysis")
        //fmt.Println(deferStack)
        deferStack[levelOfNest-1] = nil
        //fmt.Println(deferStack)
        //fmt.Println("Debugging END deferStack Analysis")
      }
      levelOfNest--
      writeFile.WriteString(line + "\n")
    } else if ( returnRegex.MatchString(line) ) {
      for i := levelOfNest; i > 0; i-- {
        for j := len(deferStack[i-1]); j > 0; j-- {
        //for _, val := range deferStack[i-1] {
          //writeFile.WriteString(val + "\n")
          writeFile.WriteString(deferStack[i-1][j-1] + "\n")
          fmt.Println("Writing from defer stack...")
        }
        writeFile.WriteString(line + "\n")
        if ( levelOfNest == 1 ) {
          levelOfNest = 0
          deferStack = make(map[int][]string)
        }
      }
    } else {
      writeFile.WriteString(line + "\n")
    }

    //if ( returnRegex.MatchString(line) ) {
    //  fmt.Println("Test point - am I reaching here?")
    //  fmt.Printf("Test point - am I reaching here? - this is the perpetrating line\"%v\"\n", line)
    //  //for i := levelOfNest; i > 0; i-- {
    //  //  //for _, val := range deferStack[i-1] {
    //  //  //  lineOfCode := deferRegex.ReplaceAllString(val, "")
    //  //  //  writeFile.WriteString(lineOfCode + "\n")
    //  //  //  fmt.Println("Writing from defer stack...")
    //  //  //}
    //  //}
    //}
  }

  fmt.Println(deferStack)
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
 
  indexOfFileExtensionPeriod = strings.LastIndex(_filepath, ".")

  nonExtensionString = _filepath[:indexOfFileExtensionPeriod]
  fileExtension = _filepath[indexOfFileExtensionPeriod:]

  outputString = nonExtensionString + "_defersified" + fileExtension
  return outputString
}

