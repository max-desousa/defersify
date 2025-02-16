package fileFinder

import (
  "regexp"
  "defersify/internal/userSettings"
  "path/filepath"
  "fmt"
  "strings"
  "io/fs"
)

var Verbose bool = false

/*******************************************************************************
 * Function: regexForExtensions
 *
 * Description: This was just a helper function to not have to think about 
 * joining strings in one line with the slice. Makes the searching function a 
 * little cleaner. Ouput is a regular expression that will find all files 
 * starting with "deferable_" and ending with one of the supported extensions.
 ******************************************************************************/
func regexForExtensions() string {
  returnVal := "deferable_[a-zA-Z0-9-_]*\\.("
  returnVal += strings.Join(userSettings.Extensions, "|")
  returnVal += ")$"
  return returnVal
}


/*******************************************************************************
 * Function: FindDeferableFiles
 *
 * Description: This function searches the current directory and all 
 * subdirectories for files that match the naming format for defersification.
 ******************************************************************************/
func FindDeferableFiles() []string {
  var returnVal []string
  pattern := regexp.MustCompile(regexForExtensions())

  err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
    if err != nil {
      return err
    }

    if (!d.IsDir() && pattern.MatchString(d.Name())) {
      returnVal = append(returnVal, path)
    }
    return nil
  })

  if err != nil {
    fmt.Printf("Error walking the path: %v\n", err)
  }

  if (userSettings.Verbose) {
    fmt.Println("Found the following files to defersify:")
    for _, val := range(returnVal) {
      fmt.Printf("\t%v\n", val)
    }
    fmt.Println("-----------------------------")
  }

  return returnVal
}

