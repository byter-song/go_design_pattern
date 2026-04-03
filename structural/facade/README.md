# 外观模式 (Facade Pattern)

## 概述

外观模式通过一个统一入口封装多个复杂子系统，让客户端只面对一个更稳定、更简洁的 API。它不消灭子系统，而是把“如何编排这些子系统”的复杂度集中到一个外观对象中。

在 Go 里，外观模式通常不是通过庞大的类层次实现，而是通过一个聚合多个服务对象的结构体来完成流程编排。

## Go 语言实现特点

- 使用组合聚合多个子系统，而不是继承
- 每个子系统保持小而清晰的职责边界
- 外观只负责流程组织，不替代底层服务
- 依赖显式注入，测试时容易替换依赖

## 核心概念

```go
type InventoryService struct{}
type PaymentGateway struct{}
type ShippingService struct{}

type OrderFacade struct {
    inventory *InventoryService
    payment   *PaymentGateway
    shipping  *ShippingService
}

func (f *OrderFacade) PlaceOrder() error {
    // 1. 预占库存
    // 2. 扣款
    // 3. 创建物流单
    return nil
}
```

## 代码示例

### 电商下单外观

当前实现使用 `OrderFacade` 统一编排 4 个子系统：

1. `InventoryService`：检查并预占库存
2. `PaymentGateway`：执行扣款和退款
3. `ShippingService`：生成物流单
4. `NotificationService`：发送订单确认通知

客户端只需要调用：

```go
summary, err := facade.PlaceOrder("user-1", items, "Shanghai Pudong")
if err != nil {
    return err
}

fmt.Println(summary.OrderID, summary.TrackingID)
```

外观内部会处理失败回滚：

- 支付失败时释放库存
- 发货失败时退款并释放库存

## 使用场景

1. 多个服务需要按照固定顺序协作
2. 希望给上层业务提供更稳定的入口
3. 子系统复杂，但调用方只关心结果
4. 需要把补偿逻辑、回滚逻辑集中管理

## 优缺点

### 优点

- 降低调用复杂度
- 隐藏子系统编排细节
- 更容易复用标准业务流程
- 更容易为复杂流程补充监控和审计

### 缺点

- 外观可能逐渐膨胀成“大对象”
- 过度封装会掩盖底层能力
- 若职责边界不清，容易演变成万能入口

## 与其他模式的关系

- 与 **适配器模式** 不同：外观是“简化多个接口”，适配器是“转换不兼容接口”
- 与 **代理模式** 不同：外观主要做流程编排，代理主要控制访问
- 与 **中介者模式** 不同：外观面向客户端，中介者面向对象之间的协作

## 参考

- 实现代码：[`facade.go`](./facade.go)
- 测试代码：[`facade_test.go`](./facade_test.go)
