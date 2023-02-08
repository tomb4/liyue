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
		SetUserId(t int32)
		GetUserId() int32
		SetPingAt(t int64)
		GetPingAt() int64
		SetBizId(s string)
		GetBizId() string
		IsClosed() bool
	}

	conn struct {
		writeLock sync.Mutex
		closed    bool

		clientIp     string
		connectionId string
		connType     int32
		bizId        string
		userId       int32
		pingAt       int64
		createdAt    int64
	}
)

func (m *conn) SetConnType(t int32) {
	m.connType = t
}

func (m *conn) GetConnType() int32 {
	return m.connType
}

func (m *conn) SetUserId(t int32) {
	m.userId = t
}

func (m *conn) GetUserId() int32 {
	return m.userId
}

func (m *conn) SetPingAt(t int64) {
	m.pingAt = t
}

func (m *conn) GetPingAt() int64 {
	return m.pingAt
}

func (m *conn) SetBizId(t string) {
	m.bizId = t
}

func (m *conn) GetBizId() string {
	return m.bizId
}

func (m *conn) GetClientIp() string {
	return m.clientIp
}
func (m *conn) GetConnectionId() string {
	return m.connectionId
}

func (m *conn) IsDead() bool {
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
func (m *conn) IsClosed() bool {
	return m.closed
}
