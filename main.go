package main

import (
  "fmt"
  "flag"
  "defersify/internal/userSettings"
  "defersify/internal/fileFinder"
)


/*******************************************************************************
 *
 ******************************************************************************/
func main() {
   /* Handling runtime flags */
  //undo := flag.Bool("u", false, "Run \"un-defersify\" to revert files back to their original state.")
  verbose := flag.Bool("v", false, "run with verbose output")
  extensions := flag.String("ex", "c", "File extensions to subject to defersification")
  flag.Parse()

  userSettings.Verbose = *verbose
  if (userSettings.SetUserExtensionOptions(*extensions)) {
    files := fileFinder.FindDeferableFiles()
    for _, val := range(files) {
      fmt.Println(val)
    }
  } else {
    fmt.Println("Error setting user extension options")
    fmt.Println("The current supported extensions are: ", userSettings.GetSupportedExtensions())
  }
}

