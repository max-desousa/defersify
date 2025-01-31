package deferSearcher

import (
  "fmt"
  "bufio"
  "os"
  "regexp"
)

func SeachForDefers(_filepath string) {
  //var withinComment bool = false
  //var deferStack [][]string
  //var levelOfDepth uint8 = 0
  //var newFileContents []string
  //var withinDefersBlock bool = false
  //var inFunctionBlock bool = false
 
  //deferRegex := regexp.MustCompile(`^\s*defer\s+`)
  //blockDefersRegex := regexp.MustCompile(`^\s*defer\s*{`)
  //blockCommentStartRegex := regexp.MustCompile(`^\s*/\*`)
  //blockCommentEndRegex := regexp.MustCompile(`\*/\s*$`)
  //completeBlockCommentRegex := regexp.MustCompile(`^\s*/\*.*\*/\s*$`)
  //openBraceRegex := regexp.MustCompile(`.*{.*`)

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

  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    line := scanner.Text()

    /***************************************************************************
     * Let's assume the MAJORITY of content provided in the software is not a
     * deferal of any sort. If the line doesn't contain any sort of characters
     * that would indicate a defer or a levelOfNest that needs to be accounted
     * for, then we can just write the line to the new file and move on to the
     * next line.
     **************************************************************************/
    if (obviouslyFineLine(line)) {
      writeFile.WriteString(line + "\n")
      continue
    }

    /***************************************************************************
     * else, let's do a more in depth analysis of the line to determine if it
     * is a defer or a block of defers.
     **************************************************************************/
    inFunctionBlock = openBraceRegex.MatchString(line)
    if (!inFunctionBlock) {
      continue
    }
  }


}

func detectOpenButNotCompleteBlockComment(_line string) bool {
  var returnVal bool = false

  blockCommentOpenRegex := regexp.MustCompile(`/\*`)
  blockCommentCloseRegex := regexp.MustCompile(`\*/`)
  blockCommentCompletedRegex := regexp.MustCompile(`/\*.*\*/`)

  /* Set return if we find an attempt to open a block comment but no attempt to close it */
  returnVal = (blockCommentOpenRegex.MatchString(_line) && !blockCommentCloseRegex.MatchString(_line))
  /* Or the results of a check to make sure that the attempt to close comse AFTER the attempt to open... i.e. valid */
  returnVal |= !blockCommentCompletedRegex.MatchString(_line)

  return returnVal
}


func obviouslyFineLine(_line string) bool {
  /* assume line is safe first because it's easier to lock as false */
  var returnVal bool = true

  deferRegex := regexp.MustCompile(`^\s*defer\s+`)
  
  /* fail search if possibly valid defer statement is found */
  returnVal &= !deferRegex.MatchString(_line)

  /* make sure that the line doesn't contain an open block comment */
  returnVal &= !detectOpenButNotCompleteBlockComment(_line)

  return returnVal
}



func generateNewFilePath(_filepath string) string {
  return _filepath + ".defers"
}

