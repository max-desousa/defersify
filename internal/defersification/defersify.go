package defersification

import (
  "fmt"
  "os"
)

func DefersifyFile(_filePath string) {
  /* create object to use for reading source file */
  readFile, err := os.Open(_filePath)
  if err != nil {
    fmt.Println("Error opening source file")
    return
  }
  defer readFile.Close()

  /* create object to use for writing to new file */
  writeFile, err := os.Create(simpleFileRename(_filePath))
  if err != nil {
    fmt.Println("Error creating new file")
    return
  }
  defer writeFile.Close()

  parseAndWriteFiles(readFile, writeFile)
}
