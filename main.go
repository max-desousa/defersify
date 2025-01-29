package main

import (
  "flag"
  "fmt"
)

func main() {
  undo := flag.Bool("u", false, "Run \"un-defersify\" to revert files back to their original state.")
  flag.Parse()
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
