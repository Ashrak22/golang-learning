package args

import (
	"bettererror"
	"testing"
)

var errorCount int

func help(vs ...string) {
	errorCount = errorCount + 1
}

func TestEvalArgsTooFewParams(t *testing.T) {
	a := NewArg()
	errorCount = 0
	a.RegisterArg("help", ArgFunc(help), 2, "--")
	var args []string
	args = append(args, "args", "--help", "abc")
	err := a.EvalArgs(args).(*bettererror.BetterError)
	if err.Code() != 0x00010003 {
		t.Fatalf("Expected error: Too Many params, received: %s", err.Error())
	}
}

func TestEvalArgsWrongArg(t *testing.T) {
	a := NewArg()
	errorCount = 0
	a.RegisterArg("help", ArgFunc(help), 0, "--")
	var args []string
	args = append(args, "args", "--hepp")
	err := a.EvalArgs(args).(*bettererror.BetterError)
	if err.Code() != 0x00010004 {
		t.Fatalf("Expected error: Wrong arg, received: %s", err.Error())
	}
}

func TestEvalArgsNoParams(t *testing.T) {
	a := NewArg()
	errorCount = 0
	a.RegisterArg("help", ArgFunc(help), 0, "--")
	var args []string
	args = append(args, "args")
	err := a.EvalArgs(args)
	if err != nil {
		t.Fatalf("Expected no error, received: %s", err.Error())
	}
}

func TestEvalArgsPossibleParam(t *testing.T) {
	a := NewArg()
	errorCount = 0
	a.RegisterArg("print", ArgFunc(help), 1, "--")
	err := a.RegisterArg("print", ArgFunc(help), 1, "--").(*bettererror.BetterError)
	if err.Code() != 0x00010001 {
		t.Fatalf("Expected Arg already registered, received: %s", err.Error())
	}
}

func TestEvalArgsCorrect(t *testing.T) {
	a := NewArg()
	errorCount = 0
	a.RegisterArg("help", ArgFunc(help), 0, "--")
	var args []string
	args = append(args, "args", "--help")
	err := a.EvalArgs(args)
	if err != nil {
		t.Fatalf("Expected no error, received: %s", err.Error())
	}
}

var abs int

func prints(vs ...string) {
	abs = len(vs)
}

func TestMoreParams(t *testing.T) {
	a := NewArg()
	a.RegisterArg("print", ArgFunc(prints), 1, "--")
	a.RegisterArg("printb", ArgFunc(prints), 1, "--")
	var args []string
	args = append(args, "args", "--print", "20", "--printb", "123")
	err := a.EvalArgs(args)
	if abs != 1 {
		t.Fatalf("Expected 1 param passed to callback, received %d", abs)
	}
	if err != nil {
		t.Fatalf("Expected no error, received: %s", err.Error())
	}
}
