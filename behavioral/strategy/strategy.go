
// Package strategy 展示了 Go 语言中实现策略模式的惯用方法。
//
// 策略模式定义了一系列算法，并将每个算法封装起来，使它们可以互相替换。
// 策略模式让算法的变化独立于使用算法的客户。
//
// Go 语言实现特点：
//   1. 使用 interface 定义策略契约，而非抽象类
//   2. 函数类型（func type）可以作为轻量级策略实现
//   3. 利用闭包实现带状态的策略
//   4. 策略可以在运行时动态替换
//
// 与 Java/C++ 的区别：
//   - 不需要复杂的类层次结构
//   - 不需要继承，只需要实现接口方法
//   - 函数类型让简单策略更加轻量
//
// 适用场景：
//   - 需要动态选择算法的场景
//   - 避免大量条件语句（if-else/switch）
//   - 算法需要独立变化，与使用方解耦
package strategy

import (
	"fmt"
	"sort"
	"time"
)

// ============================================================================
// 场景：支付策略
// ============================================================================

// PaymentStrategy 是支付策略接口。
// 在 Go 中，我们使用接口而非抽象类来定义策略契约。
type PaymentStrategy interface {
	// Pay 执行支付
	Pay(amount float64) (string, error)
	// GetName 返回策略名称
	GetName() string
}

// ============================================================================
// 具体策略实现
// ============================================================================

// CreditCardStrategy 信用卡支付策略
type CreditCardStrategy struct {
	cardNumber string
	cvv        string
	expiryDate string
}

// NewCreditCardStrategy 创建信用卡支付策略
func NewCreditCardStrategy(cardNumber, cvv, expiryDate string) *CreditCardStrategy {
	return &CreditCardStrategy{
		cardNumber: cardNumber,
		cvv:        cvv,
		expiryDate: expiryDate,
	}
}

// Pay 实现 PaymentStrategy 接口
func (c *CreditCardStrategy) Pay(amount float64) (string, error) {
	// 模拟支付处理
	fmt.Printf("[CreditCard] Paying $%.2f using card ending in %s\n", 
		amount, c.cardNumber[len(c.cardNumber)-4:])
	return fmt.Sprintf("CC-%d", time.Now().UnixNano()), nil
}

// GetName 返回策略名称
func (c *CreditCardStrategy) GetName() string {
	return "Credit Card"
}

// PayPalStrategy PayPal 支付策略
type PayPalStrategy struct {
	email    string
	password string
}

// NewPayPalStrategy 创建 PayPal 支付策略
func NewPayPalStrategy(email, password string) *PayPalStrategy {
	return &PayPalStrategy{
		email:    email,
		password: password,
	}
}

// Pay 实现 PaymentStrategy 接口
func (p *PayPalStrategy) Pay(amount float64) (string, error) {
	fmt.Printf("[PayPal] Paying $%.2f using account %s\n", amount, p.email)
	return fmt.Sprintf("PP-%d", time.Now().UnixNano()), nil
}

// GetName 返回策略名称
func (p *PayPalStrategy) GetName() string {
	return "PayPal"
}

// CryptoStrategy 加密货币支付策略
type CryptoStrategy struct {
	walletAddress string
	cryptoType    string
}

// NewCryptoStrategy 创建加密货币支付策略
func NewCryptoStrategy(walletAddress, cryptoType string) *CryptoStrategy {
	return &CryptoStrategy{
		walletAddress: walletAddress,
		cryptoType:    cryptoType,
	}
}

// Pay 实现 PaymentStrategy 接口
func (c *CryptoStrategy) Pay(amount float64) (string, error) {
	fmt.Printf("[Crypto] Paying $%.2f in %s from wallet %s...%s\n",
		amount, c.cryptoType, c.walletAddress[:6], c.walletAddress[len(c.walletAddress)-4:])
	return fmt.Sprintf("CRYPTO-%d", time.Now().UnixNano()), nil
}

// GetName 返回策略名称
func (c *CryptoStrategy) GetName() string {
	return fmt.Sprintf("Crypto (%s)", c.cryptoType)
}

// ============================================================================
// 上下文：使用策略的对象
// ============================================================================

// ShoppingCart 购物车，使用支付策略
type ShoppingCart struct {
	items           []Item
	paymentStrategy PaymentStrategy
}

// Item 购物车商品
type Item struct {
	Name  string
	Price float64
}

// NewShoppingCart 创建购物车
func NewShoppingCart() *ShoppingCart {
	return &ShoppingCart{
		items: make([]Item, 0),
	}
}

// AddItem 添加商品
func (s *ShoppingCart) AddItem(item Item) {
	s.items = append(s.items, item)
}

// GetTotal 计算总价
func (s *ShoppingCart) GetTotal() float64 {
	total := 0.0
	for _, item := range s.items {
		total += item.Price
	}
	return total
}

// SetPaymentStrategy 设置支付策略
// 这是策略模式的核心：可以在运行时动态切换策略
func (s *ShoppingCart) SetPaymentStrategy(strategy PaymentStrategy) {
	s.paymentStrategy = strategy
}

// Checkout 结账
func (s *ShoppingCart) Checkout() (string, error) {
	if s.paymentStrategy == nil {
		return "", fmt.Errorf("no payment strategy set")
	}
	
	total := s.GetTotal()
	if total == 0 {
		return "", fmt.Errorf("cart is empty")
	}
	
	fmt.Printf("\n[ShoppingCart] Checking out with %s\n", s.paymentStrategy.GetName())
	fmt.Printf("[ShoppingCart] Total amount: $%.2f\n", total)
	
	return s.paymentStrategy.Pay(total)
}

// ============================================================================
// 函数类型策略（Go 特有轻量级实现）
// ============================================================================

// SortStrategy 是一个函数类型，用于排序策略。
// 这是 Go 特有的轻量级策略实现方式，不需要定义结构体。
type SortStrategy func([]int) []int

// BubbleSort 冒泡排序策略
func BubbleSort(data []int) []int {
	result := make([]int, len(data))
	copy(result, data)
	
	n := len(result)
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	return result
}

// QuickSort 快速排序策略
func QuickSort(data []int) []int {
	result := make([]int, len(data))
	copy(result, data)
	
	var quickSort func([]int, int, int)
	quickSort = func(arr []int, low, high int) {
		if low < high {
			pi := partition(arr, low, high)
			quickSort(arr, low, pi-1)
			quickSort(arr, pi+1, high)
		}
	}
	
	quickSort(result, 0, len(result)-1)
	return result
}

func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1
	
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// StandardSort 使用标准库排序
func StandardSort(data []int) []int {
	result := make([]int, len(data))
	copy(result, data)
	sort.Ints(result)
	return result
}

// Sorter 排序器，使用函数类型策略
type Sorter struct {
	strategy SortStrategy
}

// NewSorter 创建排序器
func NewSorter(strategy SortStrategy) *Sorter {
	return &Sorter{strategy: strategy}
}

// SetStrategy 设置排序策略
func (s *Sorter) SetStrategy(strategy SortStrategy) {
	s.strategy = strategy
}

// Sort 执行排序
func (s *Sorter) Sort(data []int) []int {
	if s.strategy == nil {
		// 默认使用标准库排序
		return StandardSort(data)
	}
	return s.strategy(data)
}

// ============================================================================
// 闭包策略（带状态的策略）
// ============================================================================

// DiscountStrategy 折扣策略函数类型
type DiscountStrategy func(float64) float64

// NewPercentageDiscount 创建百分比折扣策略
// 返回一个闭包，捕获折扣百分比
func NewPercentageDiscount(percentage float64) DiscountStrategy {
	return func(amount float64) float64 {
		discount := amount * (percentage / 100)
		finalAmount := amount - discount
		fmt.Printf("[Discount] %.0f%% off: $%.2f -> $%.2f\n", 
			percentage, amount, finalAmount)
		return finalAmount
	}
}

// NewFixedDiscount 创建固定金额折扣策略
func NewFixedDiscount(discountAmount float64) DiscountStrategy {
	return func(amount float64) float64 {
		if discountAmount > amount {
			fmt.Printf("[Discount] Fixed $%.2f off: $%.2f -> $0.00\n", 
				discountAmount, amount)
			return 0
		}
		finalAmount := amount - discountAmount
		fmt.Printf("[Discount] Fixed $%.2f off: $%.2f -> $%.2f\n", 
			discountAmount, amount, finalAmount)
		return finalAmount
	}
}

// NewThresholdDiscount 创建阈值折扣策略（满减）
// 这是一种更复杂的策略，结合了多个条件
func NewThresholdDiscount(threshold, discount float64) DiscountStrategy {
	return func(amount float64) float64 {
		if amount >= threshold {
			finalAmount := amount - discount
			fmt.Printf("[Discount] Spend $%.2f, save $%.2f: $%.2f -> $%.2f\n", 
				threshold, discount, amount, finalAmount)
			return finalAmount
		}
		fmt.Printf("[Discount] No discount (need $%.2f more): $%.2f\n", 
			threshold-amount, amount)
		return amount
	}
}

// PriceCalculator 价格计算器
type PriceCalculator struct {
	discountStrategy DiscountStrategy
}

// NewPriceCalculator 创建价格计算器
func NewPriceCalculator() *PriceCalculator {
	return &PriceCalculator{
		// 默认无折扣
		discountStrategy: func(amount float64) float64 { return amount },
	}
}

// SetDiscountStrategy 设置折扣策略
func (p *PriceCalculator) SetDiscountStrategy(strategy DiscountStrategy) {
	p.discountStrategy = strategy
}

// CalculateFinalPrice 计算最终价格
func (p *PriceCalculator) CalculateFinalPrice(originalPrice float64) float64 {
	return p.discountStrategy(originalPrice)
}

// ============================================================================
// 策略注册表（运行时策略选择）
// ============================================================================

// StrategyRegistry 策略注册表
type StrategyRegistry struct {
	strategies map[string]PaymentStrategy
}

// NewStrategyRegistry 创建策略注册表
func NewStrategyRegistry() *StrategyRegistry {
	return &StrategyRegistry{
		strategies: make(map[string]PaymentStrategy),
	}
}

// Register 注册策略
func (r *StrategyRegistry) Register(name string, strategy PaymentStrategy) {
	r.strategies[name] = strategy
}

// Get 获取策略
func (r *StrategyRegistry) Get(name string) (PaymentStrategy, bool) {
	strategy, ok := r.strategies[name]
	return strategy, ok
}

// List 列出所有策略名称
func (r *StrategyRegistry) List() []string {
	names := make([]string, 0, len(r.strategies))
	for name := range r.strategies {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
