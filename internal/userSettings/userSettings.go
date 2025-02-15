package userSettings

import (
  "fmt"
  "slices"
  "strings"
)

var supportedExtensions []string = []string{"c", "cpp", "cc"}
var Extensions []string
var Verbose bool = false

func GetSupportedExtensions() []string {
  return supportedExtensions
}

/*******************************************************************************
 * Function: SetUserExtensionOptions
 * 
 * Description: This function takes the extensions submitted for defersification
 * by the user. Given differnt programs structured in a way that may be 
 * incompatible with a default method. This function checks agianst the 
 * predefined list of supported languages. (Ex: think python - no braces default
 * algo developed for c would be wrong for applying to python)
 ******************************************************************************/
func SetUserExtensionOptions(_userFlagInput string) bool {
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
    if (!slices.Contains(supportedExtensions, val)) {

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
