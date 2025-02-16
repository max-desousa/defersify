package main

import (
  "flag"
  "defersify/internal/userSettings"
  "defersify/internal/fileFinder"
  "defersify/internal/defersification"
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

  if (!userSettings.SetUserExtensionOptions(*extensions)) {
    panic("Error setting user extension options - run ./defersify --help or run in verbose mode for more information about error hit.")
  }

  files := fileFinder.SearchForDeferableFiles(flag.Args())

  for _, val := range(files) {
    defersification.DefersifyFile(val)
  }
}

