package fileFinder

import (
  "regexp"
  "defersify/internal/userSettings"
  "path/filepath"
  "fmt"
  "strings"
  "io/fs"
  "os"
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

func printFoundFiles(_files []string) {
  fmt.Println("Found the following files to defersify:")
  for _, val := range(_files) {
    fmt.Printf("\t%v\n", val)
  }
  fmt.Println("-----------------------------")
}

/*******************************************************************************
 * Function: FindDeferableFiles
 *
 * Description: This function searches the current directory and all 
 * subdirectories for files that match the naming format for defersification.
 ******************************************************************************/
func findDeferableFilesInDir(_startingPath string) []string {
  var returnVal []string
  pattern := regexp.MustCompile(regexForExtensions())

  err := filepath.WalkDir(_startingPath, func(path string, d fs.DirEntry, err error) error {
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

  return returnVal
}

/*******************************************************************************
 * Function: NuancedSearchForDeferableFiles
 *
 * Description: This function searches for all the submitted strings - file or
 * directory - for files that match the naming format for defersification.
 ******************************************************************************/

func SearchForDeferableFiles(_startingPaths []string) []string {
  var returnVal []string

  if len(_startingPaths) == 0 {
    returnVal = findDeferableFilesInDir(".")
  } else {
    for _, val := range(_startingPaths) {
      pathInfo, err := os.Stat(val)
      if (err != nil) {
        panic(err)
      }

      if (pathInfo.IsDir()) {
        nestedFiles := findDeferableFilesInDir(val)
        for _, nestedVal := range(nestedFiles) {
          returnVal = append(returnVal, nestedVal)
        }
      } else {
        if (regexp.MustCompile(regexForExtensions()).MatchString(pathInfo.Name())) {
          returnVal = append(returnVal, val)
        }
      }
    }
  }


  if (userSettings.Verbose) {
    printFoundFiles(returnVal)
  }
  return returnVal
}
