package main

import "testing"

func TestAlarmCommand(t *testing.T) {
	if cmd := AlarmCommand("darwin"); cmd == nil || cmd.Args[0] != "afplay" {
		t.Errorf("darwin = %v", cmd)
	}
	if cmd := AlarmCommand("plan9"); cmd != nil {
		t.Errorf("unknown OS should give nil, got %v", cmd.Args)
	}
}
