package defersification

import (
  "strings"
)

var IndentationString string

func padDeferStatementForNest(_levelOfNest int, _deferString string) string {
  var returnVal string
  var indentationString string
  var trimmedString string

  if (IndentationString == "") {
    indentationString = "  "
  } else {
    indentationString = IndentationString
  }  

  trimmedString = strings.TrimSpace(_deferString)
  returnVal = strings.Repeat(indentationString, _levelOfNest)
  returnVal += trimmedString
  return returnVal
}
