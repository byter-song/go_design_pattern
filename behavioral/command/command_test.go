package command

import (
	"strings"
	"testing"
)

func TestClassicCommand(t *testing.T) {
	light := NewLight()

	onButton := NewButton(NewTurnOnCommand(light))
	offButton := NewButton(NewTurnOffCommand(light))

	if got := onButton.Press(); got != "light on" {
		t.Fatalf("expected light on, got %s", got)
	}
	if !light.IsOn() {
		t.Fatal("expected light to be on")
	}

	if got := offButton.Press(); got != "light off" {
		t.Fatalf("expected light off, got %s", got)
	}
	if light.IsOn() {
		t.Fatal("expected light to be off")
	}
}

func TestFuncCommand(t *testing.T) {
	light := NewLight()
	cmd := FuncCommand(func() string {
		return light.TurnOn()
	})

	if got := cmd.Execute(); got != "light on" {
		t.Fatalf("expected light on, got %s", got)
	}
}

func TestScheduler(t *testing.T) {
	s := &Scheduler{}
	s.Add(LogCommand("a"))
	s.Add(LogCommand("b"))

	results := s.RunAll()
	if strings.Join(results, ",") != "log:a,log:b" {
		t.Fatalf("unexpected results: %v", results)
	}
}

func TestCommandInterfaceImplementation(t *testing.T) {
	var _ Command = (*TurnOnCommand)(nil)
	var _ Command = (*TurnOffCommand)(nil)
	var _ Command = FuncCommand(nil)
}
