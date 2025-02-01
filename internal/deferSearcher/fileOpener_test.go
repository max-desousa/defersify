package deferSearcher

import (
  "testing"
)



/*******************************************************************************
 * Test Group: Block Comment Detection
 ******************************************************************************/
func Test_BlockCommentDetection_None0(t *testing.T) {
  var returnVal bool

  returnVal =  detectOpenButNotCompleteBlockComment("func main() {")

  if (returnVal) {
    t.Errorf("Expected false, got true")
  } 
}

func Test_BlockCommentDetection_StandardComment(t *testing.T) {
  var returnVal bool

  returnVal =  detectOpenButNotCompleteBlockComment("int var = 0xFF; //This is a standar comment")

  if (returnVal) {
    t.Errorf("Expected false, got true")
  } 
}

func Test_BlockCommentDetection_CompleteInlineBlockComment(t *testing.T) {
  var returnVal bool

  returnVal =  detectOpenButNotCompleteBlockComment("int var = 0xFF; /* This is an in-line block comment */")

  if (returnVal) {
    t.Errorf("Expected false, got true")
  } 
}

func Test_BlockCommentDetection_BackwardsInlineBlock(t *testing.T) {
  var returnVal bool

  returnVal =  detectOpenButNotCompleteBlockComment("*/ int var = 0xFF; /* This is an in-line block comment")

  if (!returnVal) {
    t.Errorf("Expected true, got false")
  } 
}

func Test_BlockCommentDetection_StandardHangingBlockComment(t *testing.T) {
  var returnVal bool

  returnVal =  detectOpenButNotCompleteBlockComment("int var = 0xFF; /* This is a hanging in-line block comment")

  if (!returnVal) {
    t.Errorf("Expected true, got false")
  } 
}



/*******************************************************************************
 * Test Group: Obviously Fine Line
 ******************************************************************************/
func Test_FineLineDetection_BasicLine(t *testing.T) {
  var returnVal bool

  returnVal =  obviouslyFineLine("int var = 0xFF;")

  if (!returnVal) {
    t.Errorf("Expected true, got false")
  } 
}

func Test_FineLineDetection_BasicDeferStatement(t *testing.T) {
  var returnVal bool

  returnVal =  obviouslyFineLine("defer file.Close();")

  if (returnVal) {
    t.Errorf("Expected false, got true")
  } 
}

func Test_FineLineDetection_BlockDeferStatement(t *testing.T) {
  var returnVal bool

  returnVal =  obviouslyFineLine("defer {")

  if (returnVal) {
    t.Errorf("Expected false, got true")
  } 
}

func Test_FineLineDetection_FunctionDefinition(t *testing.T) {
  var returnVal bool

  returnVal =  obviouslyFineLine("int main() {")

  if (returnVal) {
    t.Errorf("Expected false, got true")
  } 
}

func Test_FineLineDetection_ClosingBrace(t *testing.T) {
  var returnVal bool

  returnVal =  obviouslyFineLine("return ReturnVal}")

  if (returnVal) {
    t.Errorf("Expected false, got true")
  } 
}

func Test_FineLineDetection_IncompleteBlockComment(t *testing.T) {
  var returnVal bool

  returnVal =  obviouslyFineLine("uint_8t foo = 0xFF; /* defer This is an in-line block comment")

  if (returnVal) {
    t.Errorf("Expected false, got true")
  } 
}



/*******************************************************************************
 * Test Group: Obviously Fine Line
 ******************************************************************************/
