package logic

import (
	"sync"
)

type (
	Connection interface {
		Read() (p []byte, err error)

		Write(p []byte) error

		Close() error

		IsDead() bool
		GetClientIp() string
		GetConnectionId() string
		SetConnType(t int32)
		GetConnType() int32
		SetUserId(t int64)
		GetUserId() int64
		SetPingAt(t int64)
		GetPingAt() int64
		SetBizId(s string)
		GetBizId() string
		IsClosed() bool
	}

	defaultConnection struct {
		writeLock sync.Mutex
		closed    bool

		clientIp     string
		connectionId string
		connType     int32
		bizId        string
		userId       int64
		pingAt       int64
		createdAt    int64
	}
)

func (m *defaultConnection) SetConnType(t int32) {
	m.connType = t
}

func (m *defaultConnection) GetConnType() int32 {
	return m.connType
}

func (m *defaultConnection) SetUserId(t int64) {
	m.userId = t
}

func (m *defaultConnection) GetUserId() int64 {
	return m.userId
}

func (m *defaultConnection) SetPingAt(t int64) {
	m.pingAt = t
}

func (m *defaultConnection) GetPingAt() int64 {
	return m.pingAt
}

func (m *defaultConnection) SetBizId(t string) {
	m.bizId = t
}

func (m *defaultConnection) GetBizId() string {
	return m.bizId
}

func (m *defaultConnection) GetClientIp() string {
	return m.clientIp
}
func (m *defaultConnection) GetConnectionId() string {
	return m.connectionId
}

func (m *defaultConnection) IsDead() bool {
	//now := utils.CurrentMillis()
	//如果已经登录过且有心跳
	//if m.pingAt != 0 {
	//	if now-m.pingAt > constant2.DeadConnectionInterval {
	//		return true
	//	}
	//	return false
	//}
	//if now-m.createdAt > constant2.DeadConnectionInterval {
	//	return true
	//}
	return false
}
func (m *defaultConnection) IsClosed() bool {
	return m.closed
}
