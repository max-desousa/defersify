package main

import (
  "flag"
  "runtime"
  "defersify/internal/userSettings"
  "defersify/internal/fileFinder"
  "defersify/internal/defersification"
  "strings"
)


/*******************************************************************************
 *
 ******************************************************************************/
func main() {
   /* Handling runtime flags */
  //undo := flag.Bool("u", false, "Run \"un-defersify\" to revert files back to their original state.")
  verbose := flag.Bool("v", false, "run with verbose output")
  extensions := flag.String("ex", "c", "File extensions to subject to defersification")
  output := flag.String("o", ".", "Output path if not in the same spot as the original files")
  flag.Parse()

  userSettings.Verbose = *verbose


  if (*output != ".") {
    var pathSeperator string

    switch (runtime.GOOS) {
    case "windows":
      pathSeperator = "\\" 
    default:
      pathSeperator = "/" 
    }

    if (!strings.HasSuffix(*output, pathSeperator))  {
      *output = *output + pathSeperator
    }
  }

  if (!userSettings.SetUserExtensionOptions(*extensions)) {
    panic("Error setting user extension options - run ./defersify --help or run in verbose mode for more information about error hit.")
  }

  files := fileFinder.SearchForDeferableFiles(flag.Args())

  for _, val := range(files) {
    defersification.DefersifyFile(val, *output)
  }
}

