package facade

import "testing"

func TestOrderFacadePlaceOrder(t *testing.T) {
	inventory := NewInventoryService(map[string]int{
		"BOOK":  10,
		"MOUSE": 5,
	})
	payment := NewPaymentGateway()
	shipping := NewShippingService()
	notification := NewNotificationService()
	facade := NewOrderFacade(inventory, payment, shipping, notification)

	items := []OrderItem{
		{SKU: "BOOK", Quantity: 2, UnitPrice: 49.9},
		{SKU: "MOUSE", Quantity: 1, UnitPrice: 99},
	}

	t.Run("Success", func(t *testing.T) {
		summary, err := facade.PlaceOrder("user-1", items, "Shanghai Pudong")
		if err != nil {
			t.Fatalf("PlaceOrder failed: %v", err)
		}

		if summary.OrderID == "" || summary.PaymentID == "" || summary.TrackingID == "" {
			t.Fatalf("expected non-empty identifiers, got %#v", summary)
		}
		if summary.Status != "confirmed" {
			t.Fatalf("expected confirmed status, got %s", summary.Status)
		}
		if summary.TotalAmount != 198.8 {
			t.Fatalf("expected total 198.8, got %.1f", summary.TotalAmount)
		}
		if inventory.Stock("BOOK") != 8 || inventory.Stock("MOUSE") != 4 {
			t.Fatalf("unexpected inventory state: BOOK=%d MOUSE=%d", inventory.Stock("BOOK"), inventory.Stock("MOUSE"))
		}
		if payment.ChargeCount() != 1 {
			t.Fatalf("expected 1 charge, got %d", payment.ChargeCount())
		}
		if shipping.ShipmentCount() != 1 {
			t.Fatalf("expected 1 shipment, got %d", shipping.ShipmentCount())
		}
		if len(notification.Messages()) != 1 {
			t.Fatalf("expected 1 notification, got %d", len(notification.Messages()))
		}
	})
}

func TestOrderFacadeRollback(t *testing.T) {
	t.Run("PaymentFailureReleasesInventory", func(t *testing.T) {
		inventory := NewInventoryService(map[string]int{"BOOK": 3})
		payment := NewPaymentGateway()
		payment.SetFailAbove(100)
		shipping := NewShippingService()
		notification := NewNotificationService()
		facade := NewOrderFacade(inventory, payment, shipping, notification)

		_, err := facade.PlaceOrder("user-2", []OrderItem{
			{SKU: "BOOK", Quantity: 2, UnitPrice: 80},
		}, "Hangzhou")
		if err == nil {
			t.Fatal("expected payment failure")
		}
		if inventory.Stock("BOOK") != 3 {
			t.Fatalf("expected inventory rollback, got %d", inventory.Stock("BOOK"))
		}
		if len(payment.Refunds()) != 0 {
			t.Fatalf("expected no refund on payment failure, got %d", len(payment.Refunds()))
		}
		if shipping.ShipmentCount() != 0 {
			t.Fatalf("expected no shipment, got %d", shipping.ShipmentCount())
		}
	})

	t.Run("ShippingFailureRefundsPayment", func(t *testing.T) {
		inventory := NewInventoryService(map[string]int{"BOOK": 4})
		payment := NewPaymentGateway()
		shipping := NewShippingService()
		shipping.SetFailAddressContains("BLOCKED")
		notification := NewNotificationService()
		facade := NewOrderFacade(inventory, payment, shipping, notification)

		_, err := facade.PlaceOrder("user-3", []OrderItem{
			{SKU: "BOOK", Quantity: 1, UnitPrice: 88},
		}, "BLOCKED-ZONE")
		if err == nil {
			t.Fatal("expected shipping failure")
		}
		if inventory.Stock("BOOK") != 4 {
			t.Fatalf("expected inventory rollback, got %d", inventory.Stock("BOOK"))
		}
		if len(payment.Refunds()) != 1 {
			t.Fatalf("expected 1 refund, got %d", len(payment.Refunds()))
		}
		if len(notification.Messages()) != 0 {
			t.Fatalf("expected no notification, got %d", len(notification.Messages()))
		}
	})
}

func TestOrderFacadeValidation(t *testing.T) {
	inventory := NewInventoryService(map[string]int{"BOOK": 1})
	payment := NewPaymentGateway()
	shipping := NewShippingService()
	notification := NewNotificationService()
	facade := NewOrderFacade(inventory, payment, shipping, notification)

	t.Run("EmptyItems", func(t *testing.T) {
		_, err := facade.PlaceOrder("user-4", nil, "Beijing")
		if err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("InsufficientStock", func(t *testing.T) {
		_, err := facade.PlaceOrder("user-5", []OrderItem{
			{SKU: "BOOK", Quantity: 2, UnitPrice: 20},
		}, "Beijing")
		if err == nil {
			t.Fatal("expected stock error")
		}
	})
}
