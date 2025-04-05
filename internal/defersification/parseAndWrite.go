package defersification

import (
  "fmt"
  "bufio"
  "os"
  "regexp"
  "strings"
)

/*******************************************************************************
 * Function: parseAndWriteFiles
 *
 * Description: This function will be the effective "main" for the 
 * defersification process. This function will be responsible for reading the
 * source file and writing to the new file.
 ******************************************************************************/
func parseAndWriteFiles(_readFile *os.File, _writeFile *os.File) {
  deferStack := make(map[int][]string)

  /* Variable tracking braces to determine what portions of the stack should be applied */
  var levelOfNest int = 0

  /* Varibles to track common instances of items that would span multiple lines,
   * or have a multiline context */
  withinBlockDefers := false
  withinBlockComment := false
  skipBraceDueToReturn := false

  /* Regular expressions for searching that's better than string.Contains */
  deferStatementRegex := regexp.MustCompile(`^\s*defer\s+`)
  deferStatementWithSpacingRegex := regexp.MustCompile(`defer\s+`)
  returnRegex := regexp.MustCompile(`^\s*return\s+`)

  /* Handle the file line by line */
  scanner := bufio.NewScanner(_readFile)
  for scanner.Scan() {
    line := scanner.Text()
    
    /* Check 1: Lines irrelevant to defer's or that affect level of nest */
    if (obviouslyFineLine(line)) {
      _writeFile.WriteString(line + "\n")
      continue
    }

    /* Check 2: Block comment detection */
    if (blockCommentDetector(line, &withinBlockComment)) {
      _writeFile.WriteString(line + "\n")
      continue
    }

    /* Check 3: Panic... no logice (yet?) for multiple opening braces */
    if (multipleOpeningBracesPanic(line)) {
      fmt.Println("Found a line with more than one opening brace... this is logic that is being put off for minimal viable product")
      return
    }

    /* Check 4: Defer block detection */
    if (detectBlockDefers(line, &withinBlockDefers, levelOfNest, &deferStack)) {
      continue
    }

    /* Meat and Potatoes: this is where the more substantive parts of the file
     * are handled. To make sure that any checks for correction at the end are
     * possible, the checks are done without "continue" statements. This means
     * an uglier chain of if/else's, but I like it as it is for now. Maybe if
     * the repo gets any love, I'll consider refactoring (or as I use it and 
     * find bugs... appropriate refactoring will of course occur. */

    /* The order of the checks OBVIOUSLY matters, a quick explanation would be:
     * 1. Check for opening brace to increment level of nest
     *   a. This takes precedence because functions need to be opened before
     *      anything else can happen.
     * 2. Check for defer statements to add to the stack
     *   a. Happens after opening brace because defer statements need to be nested,
     *      there are no global defers.
     * 3. Check for return statements to dump the stack
     *    a. Return statements are a higher priority for stack dumping and would
     *       be reached before closing braces
     * 4. Check for closing brace to decrement level of nest and dump the stack
     *    a. The last real check with the lowest priority
     * 5. Check for any other lines that don't fit the above criteria */  
    if ( strings.Contains(line, "{") ) {
      /* checking lines that contain an opening brace to monitor the stack level,
       * there is an assumption that code isn't included on the same line as an
       * opening brace. This is how I code, perhaps a future improvement. */
      levelOfNest++
      _writeFile.WriteString(line + "\n")

    } else if ( deferStatementRegex.MatchString(line) ) {
      /* Here the handling of defer statements is done, adding them to the 
       * relevant index of the deferStack */
      modifiedLine := deferStatementWithSpacingRegex.ReplaceAllString(line, "")
      deferStack[levelOfNest-1] = append(deferStack[levelOfNest-1], modifiedLine)

    } else if ( returnRegex.MatchString(line) ) {
      /* Here the handling of "return" statements is done. Dumping the entirety
       * of the deferStack. */

      /* additional check that the return statemtn doesn't also have a closing
       * brace for the section of code, and set a flag to skip the upcoming brace */
      if ( !strings.Contains(line, "}") ) {
        skipBraceDueToReturn = true
      }

      /* stack dump */
      //fmt.Printf("Return statement hit at level of nest: %v and a deferStack of %v\n", levelOfNest, deferStack)
      for i := levelOfNest; i > 0; i-- {
        for j := len(deferStack[i-1]); j > 0; j-- {
          _writeFile.WriteString(padDeferStatementForNest(levelOfNest, deferStack[i-1][j-1]) + "\n")
        }
      }

      /* write return statement to the file */
      _writeFile.WriteString(line + "\n")

      /* if the return is at the top level of the function, don't waste any time
       * and reset the stack so no errors have an opportunity to show up in
       * your output file */
      if ( levelOfNest == 1 ) {
        deferStack = make(map[int][]string)
      }

    } else if ( strings.Contains(line, "}") ) {
      
      /* Check if the brace should be skipped due to a previous return statement */
      if ( skipBraceDueToReturn ) {
        skipBraceDueToReturn = false
      } else {
        if ( levelOfNest > 0 ) {
          for i := len(deferStack[levelOfNest-1]); i > 0; i-- {
            _writeFile.WriteString(padDeferStatementForNest(levelOfNest, deferStack[levelOfNest-1][i-1]) + "\n")
          }
        }
      }
      /* decrement the lecel of nest, clear that level of the stack, and write
       * closing brace to output file. */
      deferStack[levelOfNest-1] = nil
      levelOfNest--
      _writeFile.WriteString(line + "\n")

    } else {
      /* catch all for lines that don't fit the above criteria - simply write
       * to output file */
      _writeFile.WriteString(line + "\n")
    }


    /* Final Check: Rectify stack in case of issues */
    if ( levelOfNest < 0 ) {
      /* Level of nest should never be less than 0... if it is, reset the stack */
      /* This was a thought to include in case a syntax error was found, it's a 
       * weak solution so the next function will hopefully be handled better and
       * compile */
      deferStack = make(map[int][]string)
    }
    
  }
}

/*******************************************************************************
 *Function: detectBlockDefers
 *
 * Description:
 ******************************************************************************/
 func detectBlockDefers(_line string, _withinBlockDefers *bool, _levelOfNest int, _deferStack *map[int][]string) bool {

   var returnVal bool = false
  startBlockDeferRegex := regexp.MustCompile(`^\s*defer\s*{?\s*$`)

  /* set flag if we're in a defer block */
  if ( startBlockDeferRegex.MatchString(_line) ) {
    *_withinBlockDefers = true
    returnVal = true
  }
  if ( *_withinBlockDefers ) {
    if ( strings.Contains(_line, "}") ) {
      *_withinBlockDefers = false
    } else {
      (*_deferStack)[_levelOfNest-1] = append((*_deferStack)[_levelOfNest-1], _line)
    }
    returnVal = true
  }
  return returnVal
}


/*******************************************************************************
 * Function: multipleOpeningBracesPanic
 *
 * Description:
 ******************************************************************************/
func multipleOpeningBracesPanic(_line string) bool {
  var returnVal bool = false
  /* Not sure if there is a use case for multiple {'s in one line... therefore
   * throwing an error message if this is the case... perhaps a future feature */
  if strings.Count(_line, "{") > 1 {
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
func simpleFileRename(_filepath string, _outputPath string) string {
  var outputString string
  var indexOfFileExtensionPeriod int
  var fileName string
  var pathString string
 
  indexOfFileExtensionPeriod = strings.LastIndex(_filepath, "deferable_")

  pathString = _filepath[:indexOfFileExtensionPeriod]
  fileName = _filepath[(indexOfFileExtensionPeriod + len("deferable_")):]
  //fmt.Printf("pathString: %v\n", pathString)
  //fmt.Printf("fileName: %v\n", fileName)

  if ("." == _outputPath) {
    /* when default output path is used */
    outputString = pathString + fileName
  } else {
    /* when non-default output path is used */
    outputString = _outputPath + fileName
  }
  return outputString
}


/*******************************************************************************
 * Function:
 *
 * Description:
 ******************************************************************************/
func blockCommentDetector(_line string, _withinBlockComment *bool) bool {
  var returnVal bool = false
  /* if the line has a multi line block comment handle that here using a flag
   * that passes all lines until block comment is closed */
  if ( strings.Contains(_line, "/*") || *_withinBlockComment ) {
    *_withinBlockComment = true
    returnVal = true
    if ( strings.Contains(_line, "*/") ) {
      *_withinBlockComment = false
    }
  }
  return returnVal
}
