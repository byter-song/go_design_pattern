package state

import "testing"

func TestOrderStateFlow(t *testing.T) {
	order := NewOrder()

	if order.StateName() != "pending" {
		t.Fatalf("expected pending, got %s", order.StateName())
	}

	if err := order.Pay(); err != nil {
		t.Fatalf("Pay failed: %v", err)
	}
	if order.StateName() != "paid" {
		t.Fatalf("expected paid, got %s", order.StateName())
	}

	if err := order.Ship(); err != nil {
		t.Fatalf("Ship failed: %v", err)
	}
	if order.StateName() != "shipped" {
		t.Fatalf("expected shipped, got %s", order.StateName())
	}

	if err := order.Complete(); err != nil {
		t.Fatalf("Complete failed: %v", err)
	}
	if order.StateName() != "completed" {
		t.Fatalf("expected completed, got %s", order.StateName())
	}
}

func TestInvalidTransitions(t *testing.T) {
	t.Run("ShipBeforePay", func(t *testing.T) {
		order := NewOrder()
		if err := order.Ship(); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("PayTwice", func(t *testing.T) {
		order := NewOrder()
		_ = order.Pay()
		if err := order.Pay(); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestStateInterfaceImplementation(t *testing.T) {
	var _ OrderState = pendingState{}
	var _ OrderState = paidState{}
	var _ OrderState = shippedState{}
	var _ OrderState = completedState{}
}
