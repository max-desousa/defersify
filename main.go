package main

import (
  "flag"
  "fmt"
  "strings"
  "slices"
  "path/filepath"
  "regexp"
  "io/fs"
  "defersify/internal/deferSearcher"
)

var Extensions []string
var SupportedExtensions []string = []string{"c", "go"}
var Verbose bool

/*******************************************************************************
 * Function: handleUserSubmittedExtensions
 * 
 * Description: This function takes the extensions submitted for defersification
 * by the user. Given differnt programs structured in a way that may be 
 * incompatible with a default method. This function checks agianst the 
 * predefined list of supported languages. (Ex: think python - no braces default
 * algo developed for c would be wrong for applying to python)
 ******************************************************************************/
func handleUserSubmittedExtensions(_userFlagInput string) bool {
  var returnVal bool = true
  var intermediaryExtensionList []string

   /* handling the fact that splitting a string with commas at the extreme leads
    * to empty slots in the slice... which I don't want*/
  if (',' == _userFlagInput[(len(_userFlagInput)-1)]) {
    /* remove trailing comma */
    _userFlagInput = _userFlagInput[0:(len(_userFlagInput)-1)]
  }

  /* remove comma if it's first char */
  if (',' == _userFlagInput[0]) {
    _userFlagInput = _userFlagInput[1:]
  }

  /* save string of extensions into a csv slice - but this is temporary because
   * we want to make sure there are no duplicates */
  intermediaryExtensionList = strings.Split(_userFlagInput, ",")

  for _, val := range(intermediaryExtensionList) {

    /* check that extension in slice is supported */
    if (!slices.Contains(SupportedExtensions, val)) {

      fmt.Printf("User has submitted file extension %v for defersification - this language is not yet supported\n", val)
      returnVal = false
      break

    } else {

      /* add temporary copy of extenssion to final slice if it's not already
       * there (duplicate handling) */
      if (!slices.Contains(Extensions, val)) {
        Extensions = append(Extensions, val)

      }
    }
  }

  /* verbose display of selected extensions */
  if (Verbose && returnVal) {
    fmt.Printf("The file extensions defersify will look for are:\n")
    for _, val := range(Extensions) {
      fmt.Printf("\t.%v\n", val)
    }
  }

  return returnVal
}


/*******************************************************************************
 * Function: 
 *
 * Description: 
 ******************************************************************************/
func regexForExtensions() string {
  returnVal := "\\.("
  returnVal += strings.Join(Extensions, "|")
  returnVal += ")$"
  return returnVal
}


/*******************************************************************************
 * Function: 
 *
 * Description: 
 ******************************************************************************/
func findFilesToDefersify() []string {
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

  return returnVal
}


/*******************************************************************************
 *
 ******************************************************************************/
func main() {
   /* Handling runtime flags */
  //undo := flag.Bool("u", false, "Run \"un-defersify\" to revert files back to their original state.")
  verbose := flag.Bool("v", false, "run with verbose output")
  extensions := flag.String("ex", "c", "File extensions to subject to defersification")
  flag.Parse()

  Verbose = *verbose
  handleUserSubmittedExtensions(*extensions)
  filesToDefersify := findFilesToDefersify()
  fmt.Println(filesToDefersify)
  deferSearcher.SeachForDefers(filesToDefersify[0])
  deferSearcher.SeachForDefers(filesToDefersify[1])
}

