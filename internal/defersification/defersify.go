package defersification

import (
  "fmt"
  "os"
  "defersify/internal/userSettings"
)

func DefersifyFile(_filePath string) {
  if userSettings.Verbose {
    fmt.Println("Defersifying file: ", _filePath)
  }

  /* create object to use for reading source file */
  readFile, err := os.Open(_filePath)
  if err != nil {
    fmt.Println("Error opening source file")
    return
  }
  defer readFile.Close()

  /* create object to use for writing to new file */
  outputFileName := simpleFileRename(_filePath)

  if (userSettings.Verbose) {
    fmt.Println("Will write to file: ", outputFileName)
  }

  writeFile, err := os.Create(outputFileName)
  if err != nil {
    fmt.Println("Error creating new file")
    return
  }
  defer writeFile.Close()

  parseAndWriteFiles(readFile, writeFile)
}
