package main

import (
  "flag"
  "fmt"
  "strings"
)

var Extensions []string
var SupportedExtensions []string = []string{"c"}

func main() {
  /*****************************************************************************
   * Handling runtime flags
   ****************************************************************************/
  undo := flag.Bool("u", false, "Run \"un-defersify\" to revert files back to their original state.")
  verbose := flag.Bool("v", false, "run with verbose output")
  extensions := flag.String("ex", "c", "File extensions to subject to defersification")
  flag.Parse()

  /* Convert user submitted extensions into a list we can use */
  Extensions = strings.Split(*extensions, ",")

  if (true == *verbose) { 
    fmt.Println("Running defersify...") 
    fmt.Printf("\n")
    fmt.Println("You are running the program on files with extension:")
    for _, val := range(Extensions) {
      fmt.Printf("\t.%v\n", val)
      //if SupportedExtensions.Contains(
    }
  }

  fmt.Printf("\n")
  if (true == *undo) {
    undefersify()
  } else {
    defersify()
  }
}

func undefersify() {
  fmt.Println("Un-defersifying your code")
}

func defersify() {
  fmt.Println("Defersifying your code")
}
