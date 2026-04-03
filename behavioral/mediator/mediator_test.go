package mediator

import "testing"

func TestMediator(t *testing.T) {
	room := NewChatRoom()
	alice := NewUser("alice", room)
	bob := NewUser("bob", room)

	room.Register(alice)
	room.Register(bob)

	if err := alice.Send("bob", "hello"); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	inbox := bob.Inbox()
	if len(inbox) != 1 || inbox[0] != "alice -> hello" {
		t.Fatalf("unexpected inbox: %v", inbox)
	}
}

func TestMediatorTargetNotFound(t *testing.T) {
	room := NewChatRoom()
	alice := NewUser("alice", room)
	room.Register(alice)

	if err := alice.Send("missing", "hello"); err == nil {
		t.Fatal("expected error")
	}
}
