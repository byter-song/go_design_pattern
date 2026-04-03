// Package mediator 展示了 Go 语言中实现中介者模式的惯用方法。
//
// 中介者模式通过一个中心对象协调多个对象之间的交互，减少对象间直接依赖。
package mediator

import "fmt"

// Mediator 定义消息转发能力。
type Mediator interface {
	Send(sender, target, message string) error
}

// ChatRoom 是一个简单中介者。
type ChatRoom struct {
	users map[string]*User
}

// NewChatRoom 创建聊天室。
func NewChatRoom() *ChatRoom {
	return &ChatRoom{users: make(map[string]*User)}
}

// Register 注册用户。
func (r *ChatRoom) Register(user *User) {
	r.users[user.name] = user
}

// Send 转发消息。
func (r *ChatRoom) Send(sender, target, message string) error {
	user, ok := r.users[target]
	if !ok {
		return fmt.Errorf("target %s not found", target)
	}
	user.Receive(fmt.Sprintf("%s -> %s", sender, message))
	return nil
}

// User 是同事对象。
type User struct {
	name     string
	room     Mediator
	inbox    []string
}

// NewUser 创建用户。
func NewUser(name string, room Mediator) *User {
	return &User{name: name, room: room}
}

// Send 发送消息。
func (u *User) Send(target, message string) error {
	return u.room.Send(u.name, target, message)
}

// Receive 接收消息。
func (u *User) Receive(message string) {
	u.inbox = append(u.inbox, message)
}

// Inbox 返回收件箱快照。
func (u *User) Inbox() []string {
	result := make([]string, len(u.inbox))
	copy(result, u.inbox)
	return result
}
