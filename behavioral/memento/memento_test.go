package memento

import "testing"

func TestMemento(t *testing.T) {
	editor := NewEditor()
	history := &History{}

	editor.Write("v1")
	history.Push(editor.Save())

	editor.Write("v2")
	history.Push(editor.Save())

	editor.Write("v3")

	snapshot, ok := history.Pop()
	if !ok {
		t.Fatal("expected snapshot")
	}
	editor.Restore(snapshot)
	if editor.Content() != "v2" {
		t.Fatalf("expected v2, got %s", editor.Content())
	}
}

func TestPopEmptyHistory(t *testing.T) {
	history := &History{}
	if _, ok := history.Pop(); ok {
		t.Fatal("expected no snapshot")
	}
}
