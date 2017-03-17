package bettererror

import "testing"

func TestRegisterFacility(t *testing.T) {
	RegisterFacility(0x0001, "args")
	if facilities[0x0001] != "args" {
		t.Errorf("Wrong facility returned")
	}
	err := RegisterFacility(0x0001, "args")
	if err == nil || err.(*BetterError).Code() != 0x00000001 {
		t.Errorf("Expected Facility exist.")
	}
}

func TestCreateErrorUnknownFacility(t *testing.T) {
	err := NewBetterError(0x0002, 0x0001, "test")
	if err.Code() != 0x00000002 {
		t.Errorf("Expected Error: Unknown Facility, gotten %s", err.Error())
	}
}
