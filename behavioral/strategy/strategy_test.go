package strategy

import (
	"reflect"
	"testing"
)

// ============================================================================
// 支付策略测试
// ============================================================================

func TestCreditCardStrategy(t *testing.T) {
	strategy := NewCreditCardStrategy("1234567890123456", "123", "12/25")

	t.Run("Pay", func(t *testing.T) {
		txID, err := strategy.Pay(100.50)
		if err != nil {
			t.Fatalf("Pay failed: %v", err)
		}

		if txID == "" {
			t.Error("Expected non-empty transaction ID")
		}
	})

	t.Run("GetName", func(t *testing.T) {
		name := strategy.GetName()
		if name != "Credit Card" {
			t.Errorf("Expected name 'Credit Card', got %q", name)
		}
	})
}

func TestPayPalStrategy(t *testing.T) {
	strategy := NewPayPalStrategy("user@example.com", "password")

	t.Run("Pay", func(t *testing.T) {
		txID, err := strategy.Pay(50.00)
		if err != nil {
			t.Fatalf("Pay failed: %v", err)
		}

		if txID == "" {
			t.Error("Expected non-empty transaction ID")
		}
	})

	t.Run("GetName", func(t *testing.T) {
		name := strategy.GetName()
		if name != "PayPal" {
			t.Errorf("Expected name 'PayPal', got %q", name)
		}
	})
}

func TestCryptoStrategy(t *testing.T) {
	strategy := NewCryptoStrategy("0x1234567890abcdef", "ETH")

	t.Run("Pay", func(t *testing.T) {
		txID, err := strategy.Pay(200.00)
		if err != nil {
			t.Fatalf("Pay failed: %v", err)
		}

		if txID == "" {
			t.Error("Expected non-empty transaction ID")
		}
	})

	t.Run("GetName", func(t *testing.T) {
		name := strategy.GetName()
		if name != "Crypto (ETH)" {
			t.Errorf("Expected name 'Crypto (ETH)', got %q", name)
		}
	})
}

// ============================================================================
// 购物车测试
// ============================================================================

func TestShoppingCart(t *testing.T) {
	t.Run("AddItemAndGetTotal", func(t *testing.T) {
		cart := NewShoppingCart()
		cart.AddItem(Item{Name: "Book", Price: 29.99})
		cart.AddItem(Item{Name: "Pen", Price: 5.99})

		total := cart.GetTotal()
		expected := 35.98

		if total != expected {
			t.Errorf("Expected total %.2f, got %.2f", expected, total)
		}
	})

	t.Run("CheckoutWithoutStrategy", func(t *testing.T) {
		cart := NewShoppingCart()
		cart.AddItem(Item{Name: "Book", Price: 29.99})

		_, err := cart.Checkout()
		if err == nil {
			t.Error("Expected error when no payment strategy set")
		}
	})

	t.Run("CheckoutEmptyCart", func(t *testing.T) {
		cart := NewShoppingCart()
		cart.SetPaymentStrategy(NewCreditCardStrategy("1234567890123456", "123", "12/25"))

		_, err := cart.Checkout()
		if err == nil {
			t.Error("Expected error when cart is empty")
		}
	})

	t.Run("CheckoutWithCreditCard", func(t *testing.T) {
		cart := NewShoppingCart()
		cart.AddItem(Item{Name: "Book", Price: 29.99})
		cart.SetPaymentStrategy(NewCreditCardStrategy("1234567890123456", "123", "12/25"))

		txID, err := cart.Checkout()
		if err != nil {
			t.Fatalf("Checkout failed: %v", err)
		}

		if txID == "" {
			t.Error("Expected non-empty transaction ID")
		}
	})

	t.Run("SwitchStrategy", func(t *testing.T) {
		cart := NewShoppingCart()
		cart.AddItem(Item{Name: "Book", Price: 29.99})

		// 先用信用卡
		cart.SetPaymentStrategy(NewCreditCardStrategy("1234567890123456", "123", "12/25"))
		txID1, err := cart.Checkout()
		if err != nil {
			t.Fatalf("First checkout failed: %v", err)
		}

		// 再添加商品，切换到 PayPal
		cart.AddItem(Item{Name: "Pen", Price: 5.99})
		cart.SetPaymentStrategy(NewPayPalStrategy("user@example.com", "password"))
		txID2, err := cart.Checkout()
		if err != nil {
			t.Fatalf("Second checkout failed: %v", err)
		}

		// 两次交易 ID 应该不同
		if txID1 == txID2 {
			t.Error("Expected different transaction IDs")
		}
	})
}

// ============================================================================
// 函数类型策略测试
// ============================================================================

func TestSortStrategies(t *testing.T) {
	data := []int{64, 34, 25, 12, 22, 11, 90}
	expected := []int{11, 12, 22, 25, 34, 64, 90}

	t.Run("BubbleSort", func(t *testing.T) {
		result := BubbleSort(data)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("BubbleSort failed. Expected %v, got %v", expected, result)
		}
	})

	t.Run("QuickSort", func(t *testing.T) {
		result := QuickSort(data)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("QuickSort failed. Expected %v, got %v", expected, result)
		}
	})

	t.Run("StandardSort", func(t *testing.T) {
		result := StandardSort(data)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("StandardSort failed. Expected %v, got %v", expected, result)
		}
	})
}

func TestSorter(t *testing.T) {
	data := []int{64, 34, 25, 12, 22, 11, 90}
	expected := []int{11, 12, 22, 25, 34, 64, 90}

	t.Run("WithStrategy", func(t *testing.T) {
		sorter := NewSorter(BubbleSort)
		result := sorter.Sort(data)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Sort failed. Expected %v, got %v", expected, result)
		}
	})

	t.Run("SwitchStrategy", func(t *testing.T) {
		sorter := NewSorter(BubbleSort)

		// 先用冒泡排序
		result1 := sorter.Sort(data)
		if !reflect.DeepEqual(result1, expected) {
			t.Errorf("First sort failed. Expected %v, got %v", expected, result1)
		}

		// 切换到快速排序
		sorter.SetStrategy(QuickSort)
		result2 := sorter.Sort(data)
		if !reflect.DeepEqual(result2, expected) {
			t.Errorf("Second sort failed. Expected %v, got %v", expected, result2)
		}
	})

	t.Run("DefaultStrategy", func(t *testing.T) {
		sorter := NewSorter(nil)
		result := sorter.Sort(data)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Default sort failed. Expected %v, got %v", expected, result)
		}
	})
}

// ============================================================================
// 闭包策略测试
// ============================================================================

func TestDiscountStrategies(t *testing.T) {
	t.Run("PercentageDiscount", func(t *testing.T) {
		discount := NewPercentageDiscount(20)
		result := discount(100.0)
		expected := 80.0

		if result != expected {
			t.Errorf("Expected %.2f, got %.2f", expected, result)
		}
	})

	t.Run("FixedDiscount", func(t *testing.T) {
		discount := NewFixedDiscount(30)
		result := discount(100.0)
		expected := 70.0

		if result != expected {
			t.Errorf("Expected %.2f, got %.2f", expected, result)
		}
	})

	t.Run("FixedDiscountExceedsAmount", func(t *testing.T) {
		discount := NewFixedDiscount(150)
		result := discount(100.0)
		expected := 0.0

		if result != expected {
			t.Errorf("Expected %.2f, got %.2f", expected, result)
		}
	})

	t.Run("ThresholdDiscountApplied", func(t *testing.T) {
		discount := NewThresholdDiscount(100, 20)
		result := discount(150.0)
		expected := 130.0

		if result != expected {
			t.Errorf("Expected %.2f, got %.2f", expected, result)
		}
	})

	t.Run("ThresholdDiscountNotApplied", func(t *testing.T) {
		discount := NewThresholdDiscount(100, 20)
		result := discount(50.0)
		expected := 50.0

		if result != expected {
			t.Errorf("Expected %.2f, got %.2f", expected, result)
		}
	})
}

func TestPriceCalculator(t *testing.T) {
	t.Run("WithPercentageDiscount", func(t *testing.T) {
		calc := NewPriceCalculator()
		calc.SetDiscountStrategy(NewPercentageDiscount(10))

		result := calc.CalculateFinalPrice(100.0)
		expected := 90.0

		if result != expected {
			t.Errorf("Expected %.2f, got %.2f", expected, result)
		}
	})

	t.Run("WithFixedDiscount", func(t *testing.T) {
		calc := NewPriceCalculator()
		calc.SetDiscountStrategy(NewFixedDiscount(25))

		result := calc.CalculateFinalPrice(100.0)
		expected := 75.0

		if result != expected {
			t.Errorf("Expected %.2f, got %.2f", expected, result)
		}
	})

	t.Run("SwitchDiscountStrategy", func(t *testing.T) {
		calc := NewPriceCalculator()

		// 先用百分比折扣
		calc.SetDiscountStrategy(NewPercentageDiscount(10))
		result1 := calc.CalculateFinalPrice(100.0)
		if result1 != 90.0 {
			t.Errorf("Expected 90.0, got %.2f", result1)
		}

		// 切换到固定折扣
		calc.SetDiscountStrategy(NewFixedDiscount(30))
		result2 := calc.CalculateFinalPrice(100.0)
		if result2 != 70.0 {
			t.Errorf("Expected 70.0, got %.2f", result2)
		}
	})

	t.Run("DefaultNoDiscount", func(t *testing.T) {
		calc := NewPriceCalculator()

		result := calc.CalculateFinalPrice(100.0)
		expected := 100.0

		if result != expected {
			t.Errorf("Expected %.2f, got %.2f", expected, result)
		}
	})
}

// ============================================================================
// 策略注册表测试
// ============================================================================

func TestStrategyRegistry(t *testing.T) {
	registry := NewStrategyRegistry()

	t.Run("RegisterAndGet", func(t *testing.T) {
		ccStrategy := NewCreditCardStrategy("1234567890123456", "123", "12/25")
		registry.Register("creditcard", ccStrategy)

		strategy, ok := registry.Get("creditcard")
		if !ok {
			t.Error("Expected to find creditcard strategy")
		}

		if strategy.GetName() != "Credit Card" {
			t.Errorf("Expected 'Credit Card', got %q", strategy.GetName())
		}
	})

	t.Run("GetNonExistent", func(t *testing.T) {
		_, ok := registry.Get("nonexistent")
		if ok {
			t.Error("Expected not to find nonexistent strategy")
		}
	})

	t.Run("List", func(t *testing.T) {
		registry.Register("paypal", NewPayPalStrategy("user@example.com", "password"))
		registry.Register("crypto", NewCryptoStrategy("0x1234567890abcdef", "BTC"))

		names := registry.List()
		if len(names) != 3 { // creditcard, paypal, crypto
			t.Errorf("Expected 3 strategies, got %d", len(names))
		}
	})
}

// ============================================================================
// 接口兼容性测试
// ============================================================================

func TestInterfaceCompatibility(t *testing.T) {
	t.Run("PaymentStrategyInterface", func(t *testing.T) {
		var _ PaymentStrategy = NewCreditCardStrategy("1234567890123456", "123", "12/25")
		var _ PaymentStrategy = NewPayPalStrategy("user@example.com", "password")
		var _ PaymentStrategy = NewCryptoStrategy("0x1234567890abcdef", "ETH")
	})
}

// ============================================================================
// 使用场景测试
// ============================================================================

func TestRealWorldScenario(t *testing.T) {
	// 场景：电商网站，用户可以选择不同的支付方式

	cart := NewShoppingCart()
	cart.AddItem(Item{Name: "Laptop", Price: 999.99})
	cart.AddItem(Item{Name: "Mouse", Price: 29.99})

	// 使用策略注册表管理支付方式
	registry := NewStrategyRegistry()
	registry.Register("creditcard", NewCreditCardStrategy("4111111111111111", "123", "12/25"))
	registry.Register("paypal", NewPayPalStrategy("customer@example.com", "password"))
	registry.Register("crypto", NewCryptoStrategy("0xabcdef1234567890", "ETH"))

	// 用户选择信用卡支付
	strategy, _ := registry.Get("creditcard")
	cart.SetPaymentStrategy(strategy)

	txID, err := cart.Checkout()
	if err != nil {
		t.Fatalf("Checkout failed: %v", err)
	}

	if txID == "" {
		t.Error("Expected non-empty transaction ID")
	}
}

func TestSortingScenario(t *testing.T) {
	// 场景：数据处理系统，根据数据规模选择不同排序算法

	smallData := []int{5, 2, 8, 1, 9}
	largeData := make([]int, 1000)
	for i := range largeData {
		largeData[i] = 1000 - i
	}

	// 小数据量使用冒泡排序（简单，低开销）
	smallSorter := NewSorter(BubbleSort)
	smallResult := smallSorter.Sort(smallData)

	expectedSmall := []int{1, 2, 5, 8, 9}
	if !reflect.DeepEqual(smallResult, expectedSmall) {
		t.Errorf("Small data sort failed. Expected %v, got %v", expectedSmall, smallResult)
	}

	// 大数据量使用快速排序（高效）
	largeSorter := NewSorter(QuickSort)
	largeResult := largeSorter.Sort(largeData)

	// 验证排序正确性
	for i := 0; i < len(largeResult)-1; i++ {
		if largeResult[i] > largeResult[i+1] {
			t.Error("Large data sort failed: not sorted correctly")
			break
		}
	}
}

func TestECommerceDiscountScenario(t *testing.T) {
	// 场景：电商平台，根据促销活动应用不同折扣策略

	calc := NewPriceCalculator()

	// 普通用户：无折扣
	originalPrice := 100.0
	result := calc.CalculateFinalPrice(originalPrice)
	if result != 100.0 {
		t.Errorf("Expected 100.0, got %.2f", result)
	}

	// 会员日：20% 折扣
	calc.SetDiscountStrategy(NewPercentageDiscount(20))
	result = calc.CalculateFinalPrice(originalPrice)
	if result != 80.0 {
		t.Errorf("Expected 80.0, got %.2f", result)
	}

	// 满减活动：满 100 减 30
	calc.SetDiscountStrategy(NewThresholdDiscount(100, 30))
	result = calc.CalculateFinalPrice(originalPrice)
	if result != 70.0 {
		t.Errorf("Expected 70.0, got %.2f", result)
	}
}
