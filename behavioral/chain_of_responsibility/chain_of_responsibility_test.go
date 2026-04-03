package chain_of_responsibility

import (
	"testing"
)

// ============================================================================
// 支持工单处理测试
// ============================================================================

func TestFirstLineSupport(t *testing.T) {
	handler := NewFirstLineSupport()

	t.Run("HandleGeneralLevel1", func(t *testing.T) {
		ticket := NewSupportTicket("TKT-001", "general", "Password reset", 1)
		result := handler.Handle(ticket)
		if !result {
			t.Error("Expected FirstLineSupport to handle general level 1 ticket")
		}
	})

	t.Run("PassTechnicalTicket", func(t *testing.T) {
		ticket := NewSupportTicket("TKT-002", "technical", "Server error", 1)
		result := handler.Handle(ticket)
		// 应该无法处理，因为没有下一个处理器
		if result {
			t.Error("Expected FirstLineSupport to not handle technical ticket")
		}
	})

	t.Run("PassLevel2Ticket", func(t *testing.T) {
		ticket := NewSupportTicket("TKT-003", "general", "Complex issue", 2)
		result := handler.Handle(ticket)
		// 应该无法处理，因为没有下一个处理器
		if result {
			t.Error("Expected FirstLineSupport to not handle level 2 ticket")
		}
	})
}

func TestTechnicalSupport(t *testing.T) {
	handler := NewTechnicalSupport()

	t.Run("HandleTechnicalLevel1", func(t *testing.T) {
		ticket := NewSupportTicket("TKT-004", "technical", "Bug report", 1)
		result := handler.Handle(ticket)
		if !result {
			t.Error("Expected TechnicalSupport to handle technical level 1 ticket")
		}
	})

	t.Run("HandleTechnicalLevel2", func(t *testing.T) {
		ticket := NewSupportTicket("TKT-005", "technical", "System crash", 2)
		result := handler.Handle(ticket)
		if !result {
			t.Error("Expected TechnicalSupport to handle technical level 2 ticket")
		}
	})

	t.Run("PassBillingTicket", func(t *testing.T) {
		ticket := NewSupportTicket("TKT-006", "billing", "Invoice issue", 1)
		result := handler.Handle(ticket)
		if result {
			t.Error("Expected TechnicalSupport to not handle billing ticket")
		}
	})
}

func TestBillingSupport(t *testing.T) {
	handler := NewBillingSupport()

	t.Run("HandleBillingLevel1", func(t *testing.T) {
		ticket := NewSupportTicket("TKT-007", "billing", "Payment question", 1)
		result := handler.Handle(ticket)
		if !result {
			t.Error("Expected BillingSupport to handle billing level 1 ticket")
		}
	})

	t.Run("PassFraudTicket", func(t *testing.T) {
		ticket := NewSupportTicket("TKT-008", "billing", "Report fraud", 1)
		result := handler.Handle(ticket)
		if result {
			t.Error("Expected BillingSupport to not handle fraud ticket")
		}
	})
}

func TestSeniorSupport(t *testing.T) {
	handler := NewSeniorSupport()

	t.Run("HandleLevel2", func(t *testing.T) {
		ticket := NewSupportTicket("TKT-009", "general", "Complex issue", 2)
		result := handler.Handle(ticket)
		if !result {
			t.Error("Expected SeniorSupport to handle level 2 ticket")
		}
	})

	t.Run("PassLevel1", func(t *testing.T) {
		ticket := NewSupportTicket("TKT-010", "general", "Simple issue", 1)
		result := handler.Handle(ticket)
		if result {
			t.Error("Expected SeniorSupport to not handle level 1 ticket")
		}
	})
}

func TestManagerSupport(t *testing.T) {
	handler := NewManagerSupport()

	t.Run("HandleLevel3", func(t *testing.T) {
		ticket := NewSupportTicket("TKT-011", "general", "Escalated issue", 3)
		result := handler.Handle(ticket)
		if !result {
			t.Error("Expected ManagerSupport to handle level 3 ticket")
		}
	})

	t.Run("PassLevel2", func(t *testing.T) {
		ticket := NewSupportTicket("TKT-012", "general", "Normal issue", 2)
		result := handler.Handle(ticket)
		if result {
			t.Error("Expected ManagerSupport to not handle level 2 ticket")
		}
	})
}

func TestDefaultHandler(t *testing.T) {
	handler := NewDefaultHandler()

	t.Run("CannotHandle", func(t *testing.T) {
		ticket := NewSupportTicket("TKT-013", "unknown", "Unknown issue", 5)
		result := handler.Handle(ticket)
		if result {
			t.Error("Expected DefaultHandler to not handle any ticket")
		}
	})
}

// ============================================================================
// 责任链构建器测试
// ============================================================================

func TestChainBuilder(t *testing.T) {
	t.Run("BuildChain", func(t *testing.T) {
		builder := NewChainBuilder()
		chain := builder.
			AddHandler(NewFirstLineSupport()).
			AddHandler(NewTechnicalSupport()).
			AddHandler(NewBillingSupport()).
			AddHandler(NewSeniorSupport()).
			AddHandler(NewManagerSupport()).
			AddHandler(NewDefaultHandler()).
			Build()

		if chain == nil {
			t.Fatal("Expected chain to be built")
		}

		// 测试一般问题被一线支持处理
		ticket := NewSupportTicket("TKT-014", "general", "Password reset", 1)
		result := chain.Handle(ticket)
		if !result {
			t.Error("Expected chain to handle general level 1 ticket")
		}
	})

	t.Run("TechnicalTicketFlow", func(t *testing.T) {
		builder := NewChainBuilder()
		chain := builder.
			AddHandler(NewFirstLineSupport()).
			AddHandler(NewTechnicalSupport()).
			AddHandler(NewDefaultHandler()).
			Build()

		// 技术问题应该被一线支持跳过，然后被技术支持处理
		ticket := NewSupportTicket("TKT-015", "technical", "Bug report", 1)
		result := chain.Handle(ticket)
		if !result {
			t.Error("Expected chain to handle technical ticket")
		}
	})

	t.Run("Level3TicketFlow", func(t *testing.T) {
		builder := NewChainBuilder()
		chain := builder.
			AddHandler(NewFirstLineSupport()).
			AddHandler(NewTechnicalSupport()).
			AddHandler(NewBillingSupport()).
			AddHandler(NewSeniorSupport()).
			AddHandler(NewManagerSupport()).
			AddHandler(NewDefaultHandler()).
			Build()

		// 级别 3 问题应该一路传递到经理
		ticket := NewSupportTicket("TKT-016", "general", "Critical issue", 3)
		result := chain.Handle(ticket)
		if !result {
			t.Error("Expected chain to handle level 3 ticket")
		}
	})

	t.Run("UnhandledTicket", func(t *testing.T) {
		builder := NewChainBuilder()
		chain := builder.
			AddHandler(NewFirstLineSupport()).
			AddHandler(NewDefaultHandler()).
			Build()

		// 级别 5 问题无法被处理
		ticket := NewSupportTicket("TKT-017", "unknown", "Unknown issue", 5)
		result := chain.Handle(ticket)
		if result {
			t.Error("Expected chain to not handle unknown level ticket")
		}
	})
}

// ============================================================================
// 函数类型责任链测试
// ============================================================================

func TestChain(t *testing.T) {
	t.Run("ExecuteChain", func(t *testing.T) {
		chain := NewChain()

		// 添加处理器
		chain.Add(func(req *Request) *Response {
			if req.Type == "A" {
				return &Response{Handled: true, Result: "Handled by A"}
			}
			return nil
		})

		chain.Add(func(req *Request) *Response {
			if req.Type == "B" {
				return &Response{Handled: true, Result: "Handled by B"}
			}
			return nil
		})

		// 测试类型 A
		resp := chain.Execute(&Request{Type: "A", Content: "test"})
		if !resp.Handled || resp.Result != "Handled by A" {
			t.Errorf("Expected 'Handled by A', got %v", resp)
		}

		// 测试类型 B
		resp = chain.Execute(&Request{Type: "B", Content: "test"})
		if !resp.Handled || resp.Result != "Handled by B" {
			t.Errorf("Expected 'Handled by B', got %v", resp)
		}

		// 测试未知类型
		resp = chain.Execute(&Request{Type: "C", Content: "test"})
		if resp.Handled {
			t.Error("Expected request to not be handled")
		}
	})

	t.Run("ChainStopsOnFirstHandler", func(t *testing.T) {
		chain := NewChain()
		handlerCount := 0

		chain.Add(func(req *Request) *Response {
			handlerCount++
			return &Response{Handled: true, Result: "First handler"}
		})

		chain.Add(func(req *Request) *Response {
			handlerCount++
			return &Response{Handled: true, Result: "Second handler"}
		})

		chain.Execute(&Request{Type: "any", Content: "test"})

		if handlerCount != 1 {
			t.Errorf("Expected 1 handler to be called, got %d", handlerCount)
		}
	})
}

// ============================================================================
// 中间件风格责任链测试
// ============================================================================

func TestMiddlewareChain(t *testing.T) {
	t.Run("ExecuteMiddleware", func(t *testing.T) {
		ctx := NewContext(&Request{Type: "test", Content: "content"})

		executionOrder := []string{}

		ctx.Use(func(c *Context) bool {
			executionOrder = append(executionOrder, "middleware1")
			if c.Request.Type == "stop" {
				c.Response.Handled = true
				c.Response.Result = "Stopped at middleware1"
				return true
			}
			return c.Next()
		})

		ctx.Use(func(c *Context) bool {
			executionOrder = append(executionOrder, "middleware2")
			c.Response.Handled = true
			c.Response.Result = "Handled by middleware2"
			return true
		})

		ctx.Execute()

		if len(executionOrder) != 2 {
			t.Errorf("Expected 2 middleware calls, got %d", len(executionOrder))
		}
	})

	t.Run("StopInMiddleware", func(t *testing.T) {
		ctx := NewContext(&Request{Type: "stop", Content: "content"})

		executionOrder := []string{}

		ctx.Use(func(c *Context) bool {
			executionOrder = append(executionOrder, "middleware1")
			if c.Request.Type == "stop" {
				c.Response.Handled = true
				c.Response.Result = "Stopped"
				return true
			}
			return c.Next()
		})

		ctx.Use(func(c *Context) bool {
			executionOrder = append(executionOrder, "middleware2")
			return true
		})

		ctx.Execute()

		if len(executionOrder) != 1 {
			t.Errorf("Expected 1 middleware call, got %d", len(executionOrder))
		}
	})
}

// ============================================================================
// 审批流程测试
// ============================================================================

func TestApprovalChain(t *testing.T) {
	t.Run("TeamLeaderApproval", func(t *testing.T) {
		teamLeader := NewTeamLeader("Alice")

		request := &ExpenseRequest{ID: "EXP-001", Amount: 500, Reason: "Office supplies"}
		approved, message := teamLeader.Approve(request)

		if !approved {
			t.Error("Expected request to be approved")
		}
		if message != "Approved by Team Leader Alice" {
			t.Errorf("Unexpected message: %s", message)
		}
	})

	t.Run("DepartmentManagerApproval", func(t *testing.T) {
		teamLeader := NewTeamLeader("Alice")
		deptManager := NewDepartmentManager("Bob")
		teamLeader.SetNext(deptManager)

		request := &ExpenseRequest{ID: "EXP-002", Amount: 3000, Reason: "Equipment"}
		approved, message := teamLeader.Approve(request)

		if !approved {
			t.Error("Expected request to be approved")
		}
		if message != "Approved by Department Manager Bob" {
			t.Errorf("Unexpected message: %s", message)
		}
	})

	t.Run("CFOApproval", func(t *testing.T) {
		teamLeader := NewTeamLeader("Alice")
		deptManager := NewDepartmentManager("Bob")
		cfo := NewCFO("Charlie")
		teamLeader.SetNext(deptManager).SetNext(cfo)

		request := &ExpenseRequest{ID: "EXP-003", Amount: 10000, Reason: "Project"}
		approved, message := teamLeader.Approve(request)

		if !approved {
			t.Error("Expected request to be approved")
		}
		if message != "Approved by CFO Charlie" {
			t.Errorf("Unexpected message: %s", message)
		}
	})

	t.Run("NoApprover", func(t *testing.T) {
		teamLeader := NewTeamLeader("Alice")

		request := &ExpenseRequest{ID: "EXP-004", Amount: 2000, Reason: "Travel"}
		approved, message := teamLeader.Approve(request)

		if approved {
			t.Error("Expected request to not be approved")
		}
		if message != "No approver available" {
			t.Errorf("Unexpected message: %s", message)
		}
	})
}

// ============================================================================
// 接口兼容性测试
// ============================================================================

func TestInterfaceCompatibility(t *testing.T) {
	t.Run("SupportHandlerInterface", func(t *testing.T) {
		var _ SupportHandler = NewFirstLineSupport()
		var _ SupportHandler = NewTechnicalSupport()
		var _ SupportHandler = NewBillingSupport()
		var _ SupportHandler = NewSeniorSupport()
		var _ SupportHandler = NewManagerSupport()
		var _ SupportHandler = NewDefaultHandler()
	})

	t.Run("ApproverInterface", func(t *testing.T) {
		var _ Approver = NewTeamLeader("test")
		var _ Approver = NewDepartmentManager("test")
		var _ Approver = NewCFO("test")
	})
}

// ============================================================================
// 使用场景测试
// ============================================================================

func TestSupportSystemScenario(t *testing.T) {
	// 场景：完整的支持工单处理系统

	builder := NewChainBuilder()
	supportChain := builder.
		AddHandler(NewFirstLineSupport()).
		AddHandler(NewTechnicalSupport()).
		AddHandler(NewBillingSupport()).
		AddHandler(NewSeniorSupport()).
		AddHandler(NewManagerSupport()).
		AddHandler(NewDefaultHandler()).
		Build()

	testCases := []struct {
		name     string
		ticket   *SupportTicket
		expected bool
	}{
		{
			name:     "General Level 1 - First Line",
			ticket:   NewSupportTicket("TKT-101", "general", "Password reset", 1),
			expected: true,
		},
		{
			name:     "Technical Level 1 - Technical Support",
			ticket:   NewSupportTicket("TKT-102", "technical", "Bug report", 1),
			expected: true,
		},
		{
			name:     "Billing Level 1 - Billing Support",
			ticket:   NewSupportTicket("TKT-103", "billing", "Invoice question", 1),
			expected: true,
		},
		{
			name:     "General Level 2 - Senior Support",
			ticket:   NewSupportTicket("TKT-104", "general", "Complex issue", 2),
			expected: true,
		},
		{
			name:     "General Level 3 - Manager",
			ticket:   NewSupportTicket("TKT-105", "general", "Escalated issue", 3),
			expected: true,
		},
		{
			name:     "Unknown Level 5 - No Handler",
			ticket:   NewSupportTicket("TKT-106", "unknown", "Unknown issue", 5),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := supportChain.Handle(tc.ticket)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestExpenseApprovalScenario(t *testing.T) {
	// 场景：费用审批流程

	// 构建审批链
	teamLeader := NewTeamLeader("Alice")
	deptManager := NewDepartmentManager("Bob")
	cfo := NewCFO("Charlie")

	teamLeader.SetNext(deptManager).SetNext(cfo)

	testCases := []struct {
		name           string
		amount         float64
		expectedApprover string
	}{
		{
			name:           "Small expense - Team Leader",
			amount:         500,
			expectedApprover: "Team Leader",
		},
		{
			name:           "Medium expense - Department Manager",
			amount:         3000,
			expectedApprover: "Department Manager",
		},
		{
			name:           "Large expense - CFO",
			amount:         10000,
			expectedApprover: "CFO",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := &ExpenseRequest{
				ID:     "EXP-TEST",
				Amount: tc.amount,
				Reason: "Test expense",
			}

			approved, message := teamLeader.Approve(request)
			if !approved {
				t.Error("Expected request to be approved")
			}

			if !contains(message, tc.expectedApprover) {
				t.Errorf("Expected approval by %s, got: %s", tc.expectedApprover, message)
			}
		})
	}
}

func TestValidationChainScenario(t *testing.T) {
	// 场景：请求验证链

	chain := NewChain()

	// 验证请求不为空
	chain.Add(func(req *Request) *Response {
		if req.Content == "" {
			return &Response{Handled: true, Result: "Validation failed: content is empty"}
		}
		return nil
	})

	// 验证请求类型
	chain.Add(func(req *Request) *Response {
		validTypes := map[string]bool{"A": true, "B": true, "C": true}
		if !validTypes[req.Type] {
			return &Response{Handled: true, Result: "Validation failed: invalid type"}
		}
		return nil
	})

	// 处理请求
	chain.Add(func(req *Request) *Response {
		return &Response{Handled: true, Result: "Request processed successfully"}
	})

	testCases := []struct {
		name     string
		request  *Request
		expected string
	}{
		{
			name:     "Valid request",
			request:  &Request{Type: "A", Content: "valid content"},
			expected: "Request processed successfully",
		},
		{
			name:     "Empty content",
			request:  &Request{Type: "A", Content: ""},
			expected: "Validation failed: content is empty",
		},
		{
			name:     "Invalid type",
			request:  &Request{Type: "D", Content: "content"},
			expected: "Validation failed: invalid type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := chain.Execute(tc.request)
			if resp.Result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, resp.Result)
			}
		})
	}
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
