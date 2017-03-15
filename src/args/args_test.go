package args

import (
	"bettererror"
	"testing"
)

var errorCount int

func help(vs ...string) {
	errorCount = errorCount + 1
}

func TestEvalArgsTooManyParams(t *testing.T) {
	a := NewArg()
	errorCount = 0
	a.RegisterArg("help", ArgFunc(help), 0, "--")
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
	if err.Code() != 0x00010002 {
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
