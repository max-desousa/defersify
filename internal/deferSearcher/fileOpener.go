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

  //scanner := bufio.NewScanner(file)

  //for scanner.Scan() {
  //  line := scanner.Text()
  //  inFunctionBlock = openBraceRegex.MatchString(line)
  //  if (!inFunctionBlock) {
  //    continue
  //  }

  //  
  //}
}

func generateNewFilePath(_filepath string) string {
  return _filepath + ".defers"
}
