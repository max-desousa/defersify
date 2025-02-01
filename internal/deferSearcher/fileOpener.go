package deferSearcher

import (
  "fmt"
  "bufio"
  "os"
  "regexp"
)

func SeachForDefers(_filepath string) {

  readFile, err := os.Open(_filepath)
  defer readFile.Close()

  if err != nil {
    fmt.Println("Error opening source file")
    return
  }

  writeFile, err := os.Create("defers.txt")
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
  
  if ( !deferRegex.MatchString(_line) &&
       !detectOpenButNotCompleteBlockComment(_line) &&
       !openBraceRegex.MatchString(_line) ) {
    returnVal = true
  } else {
    returnVal = false
  }

  return returnVal
}



func generateNewFilePath(_filepath string) string {
  return _filepath + ".defers"
}

