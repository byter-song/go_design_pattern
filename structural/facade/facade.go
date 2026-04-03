// Package facade 展示了 Go 语言中实现外观模式的惯用方法。
//
// 外观模式为一组复杂的子系统提供统一且更易用的入口。
// 客户端只依赖外观对象，而不需要了解多个子系统之间的编排细节。
//
// Go 语言实现特点：
//   1. 通过组合多个服务对象构建外观
//   2. 使用小接口保持子系统职责清晰
//   3. 外观负责流程编排，而不是承载所有业务细节
//   4. 易于替换具体子系统，天然适合测试
package facade

import (
	"fmt"
	"strings"
	"time"
)

// ============================================================================
// 场景：电商下单外观
// ============================================================================

// OrderItem 表示订单中的商品。
type OrderItem struct {
	SKU       string
	Quantity  int
	UnitPrice float64
}

// OrderSummary 是下单后的汇总信息。
type OrderSummary struct {
	OrderID      string
	PaymentID    string
	TrackingID   string
	ReservedSKUs []string
	TotalAmount  float64
	Status       string
}

// InventoryService 负责库存管理。
type InventoryService struct {
	stock map[string]int
}

// NewInventoryService 创建库存服务。
func NewInventoryService(stock map[string]int) *InventoryService {
	copied := make(map[string]int, len(stock))
	for sku, qty := range stock {
		copied[sku] = qty
	}
	return &InventoryService{stock: copied}
}

// Reserve 预占库存。
func (s *InventoryService) Reserve(items []OrderItem) ([]string, error) {
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for %s", item.SKU)
		}
		if s.stock[item.SKU] < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for %s", item.SKU)
		}
	}

	reserved := make([]string, 0, len(items))
	for _, item := range items {
		s.stock[item.SKU] -= item.Quantity
		reserved = append(reserved, item.SKU)
	}

	return reserved, nil
}

// Release 释放库存。
func (s *InventoryService) Release(items []OrderItem) {
	for _, item := range items {
		s.stock[item.SKU] += item.Quantity
	}
}

// Stock 返回当前库存。
func (s *InventoryService) Stock(sku string) int {
	return s.stock[sku]
}

// PaymentReceipt 是支付结果。
type PaymentReceipt struct {
	PaymentID string
	Amount    float64
}

// PaymentGateway 负责支付处理。
type PaymentGateway struct {
	chargeCount int
	refunds     []string
	failAbove   float64
}

// NewPaymentGateway 创建支付网关。
func NewPaymentGateway() *PaymentGateway {
	return &PaymentGateway{}
}

// SetFailAbove 设置失败阈值。
func (g *PaymentGateway) SetFailAbove(limit float64) {
	g.failAbove = limit
}

// Charge 扣款。
func (g *PaymentGateway) Charge(userID string, amount float64) (PaymentReceipt, error) {
	if amount <= 0 {
		return PaymentReceipt{}, fmt.Errorf("invalid payment amount")
	}
	if g.failAbove > 0 && amount > g.failAbove {
		return PaymentReceipt{}, fmt.Errorf("payment rejected for %s", userID)
	}

	g.chargeCount++
	return PaymentReceipt{
		PaymentID: fmt.Sprintf("PAY-%d", time.Now().UnixNano()),
		Amount:    amount,
	}, nil
}

// Refund 退款。
func (g *PaymentGateway) Refund(paymentID string) {
	g.refunds = append(g.refunds, paymentID)
}

// ChargeCount 返回支付次数。
func (g *PaymentGateway) ChargeCount() int {
	return g.chargeCount
}

// Refunds 返回退款记录。
func (g *PaymentGateway) Refunds() []string {
	result := make([]string, len(g.refunds))
	copy(result, g.refunds)
	return result
}

// Shipment 是物流信息。
type Shipment struct {
	TrackingID string
	Address    string
}

// ShippingService 负责创建物流单。
type ShippingService struct {
	shipments      []Shipment
	failAddressSub string
}

// NewShippingService 创建物流服务。
func NewShippingService() *ShippingService {
	return &ShippingService{}
}

// SetFailAddressContains 设置地址失败关键词。
func (s *ShippingService) SetFailAddressContains(keyword string) {
	s.failAddressSub = keyword
}

// CreateShipment 创建物流单。
func (s *ShippingService) CreateShipment(userID string, items []OrderItem, address string) (Shipment, error) {
	if strings.TrimSpace(address) == "" {
		return Shipment{}, fmt.Errorf("shipping address is required")
	}
	if s.failAddressSub != "" && strings.Contains(address, s.failAddressSub) {
		return Shipment{}, fmt.Errorf("shipping service unavailable for %s", address)
	}

	shipment := Shipment{
		TrackingID: fmt.Sprintf("SHIP-%d", time.Now().UnixNano()),
		Address:    address,
	}
	s.shipments = append(s.shipments, shipment)
	return shipment, nil
}

// ShipmentCount 返回物流单数量。
func (s *ShippingService) ShipmentCount() int {
	return len(s.shipments)
}

// NotificationService 负责通知用户。
type NotificationService struct {
	messages []string
}

// NewNotificationService 创建通知服务。
func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// Send 发送通知。
func (s *NotificationService) Send(userID, message string) {
	s.messages = append(s.messages, fmt.Sprintf("%s:%s", userID, message))
}

// Messages 返回通知记录。
func (s *NotificationService) Messages() []string {
	result := make([]string, len(s.messages))
	copy(result, s.messages)
	return result
}

// OrderFacade 封装下单流程。
type OrderFacade struct {
	inventory    *InventoryService
	payment      *PaymentGateway
	shipping     *ShippingService
	notification *NotificationService
}

// NewOrderFacade 创建下单外观。
func NewOrderFacade(
	inventory *InventoryService,
	payment *PaymentGateway,
	shipping *ShippingService,
	notification *NotificationService,
) *OrderFacade {
	return &OrderFacade{
		inventory:    inventory,
		payment:      payment,
		shipping:     shipping,
		notification: notification,
	}
}

// PlaceOrder 执行完整下单流程。
func (f *OrderFacade) PlaceOrder(userID string, items []OrderItem, address string) (OrderSummary, error) {
	if len(items) == 0 {
		return OrderSummary{}, fmt.Errorf("order must contain at least one item")
	}

	reserved, err := f.inventory.Reserve(items)
	if err != nil {
		return OrderSummary{}, err
	}

	total := calculateTotal(items)
	receipt, err := f.payment.Charge(userID, total)
	if err != nil {
		f.inventory.Release(items)
		return OrderSummary{}, err
	}

	shipment, err := f.shipping.CreateShipment(userID, items, address)
	if err != nil {
		f.payment.Refund(receipt.PaymentID)
		f.inventory.Release(items)
		return OrderSummary{}, err
	}

	summary := OrderSummary{
		OrderID:      fmt.Sprintf("ORDER-%d", time.Now().UnixNano()),
		PaymentID:    receipt.PaymentID,
		TrackingID:   shipment.TrackingID,
		ReservedSKUs: reserved,
		TotalAmount:  total,
		Status:       "confirmed",
	}

	f.notification.Send(userID, fmt.Sprintf("order %s confirmed", summary.OrderID))
	return summary, nil
}

func calculateTotal(items []OrderItem) float64 {
	total := 0.0
	for _, item := range items {
		total += float64(item.Quantity) * item.UnitPrice
	}
	return total
}
