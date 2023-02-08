package logic

import (
	"context"
	"sync"

	"liyue/apps/im/gateway/internal/ecode"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	_connManager     *ConnManager
	_onceConnManager sync.Once
)

type (
	ConnManager struct {
		logx.Logger
		unAuthConnStore sync.Map
		connStore       sync.Map
		connStoreLock   sync.RWMutex
	}
)

func OnceConnManager() *ConnManager {
	_onceConnManager.Do(func() {
		_connManager = &ConnManager{
			Logger:          logx.WithContext(context.TODO()),
			unAuthConnStore: sync.Map{},
			connStore:       sync.Map{},
			connStoreLock:   sync.RWMutex{},
		}
	})
	return _connManager
}

func (c *ConnManager) DelConn(conn Connection) {
	if conn.GetUserId() == 0 {
		c.unAuthConnStore.Delete(conn.GetConnectionId())
		return
	}

	key := conn.GetUserId()
	//保证查询出来的与最终删除的原子性
	c.connStoreLock.Lock()
	defer c.connStoreLock.Unlock()
	v, ok := c.connStore.Load(key)
	if !ok {
		return
	}
	if v.(Connection).GetConnectionId() != conn.GetConnectionId() {
		return
	}
	c.connStore.Delete(key)
}

func (c *ConnManager) Dispatch(conn Connection) {
	c.Info("Establish connection, id:", conn.GetConnectionId())
	c.unAuthConnStore.Store(conn.GetConnectionId(), conn)
	for {
		bytes, e := conn.Read()
		//如果是我们自己主动关闭的连接，则不需要打印错误信息了
		if e != nil && !conn.IsClosed() {
			if websocket.IsCloseError(e, websocket.CloseNoStatusReceived) {
				c.Info("read err 正常关闭报错", e)
			} else {
				c.Info("read err", e)
			}
			//如果是一些常规错误，则忽略此次报文
			if e == ecode.ErrMessageType {
				continue
			} else {
				break
			}
		}
		if conn.IsDead() || conn.IsClosed() {
			break
		}
		c.HandleMessage(bytes, conn)
	}
	err := conn.Close()
	if err != nil {
		c.Error("close err", err)
	}
	c.DelConn(conn)
}

func (c *ConnManager) HandleMessage(bts []byte, conn Connection) {
	// TODO
	err := conn.Write(bts)
	if err != nil {
		c.Error("write err", err)
	}
}
