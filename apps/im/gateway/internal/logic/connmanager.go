package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"liyue/apps/im/gateway/internal/svc"
	"sync"
	"time"

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
		handler         *GatewayLogic
		conversionMap   map[string][]int64 // TODO optimize this
	}
)

func OnceConnManager(ctxSvc *svc.ServiceContext) *ConnManager {
	_onceConnManager.Do(func() {
		_connManager = &ConnManager{
			Logger:          logx.WithContext(context.TODO()),
			unAuthConnStore: sync.Map{},
			connStore:       sync.Map{},
			connStoreLock:   sync.RWMutex{},
			handler:         NewGatewayLogic(context.TODO(), ctxSvc),
			conversionMap:   map[string][]int64{},
		}
	})
	return _connManager
}

func (c *ConnManager) SaveConn(conn Connection) {
	key := conn.GetUserId()
	old, exists := c.connStore.Load(key)
	//如果连接已存在且已经存在的连接与希望存入的连接不是同一个连接 则关闭老的连接
	if exists && old.(Connection).GetConnectionId() != conn.GetConnectionId() {
		err := old.(Connection).Close()
		if err != nil {
			c.Error("SaveConn err:", err)
		}
	}

	c.connStoreLock.Lock()
	defer c.connStoreLock.Unlock()

	c.connStore.Store(key, conn)
	c.unAuthConnStore.Delete(conn.GetConnectionId())
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

func (c *ConnManager) SendMessage(key int64, data []byte) error {
	cc, ok := c.connStore.Load(key)
	if !ok {
		c.Info("SendMessage not found", key)
		return errors.New("SendMessage not found")
	}

	conn := cc.(Connection)
	err := conn.Write(data)
	if err != nil {
		c.Error("SendMessage err", err)
		return err
	}

	return nil
}

func (c *ConnManager) Dispatch(conn Connection) {
	c.Info("Establish connection, id:", conn.GetConnectionId())
	c.unAuthConnStore.Store(conn.GetConnectionId(), conn)
	for {
		bytes, e := conn.Read()
		//如果是我们自己主动关闭的连接，则不需要打印错误信息了
		if e != nil && !conn.IsClosed() {
			if websocket.IsCloseError(e, websocket.CloseNoStatusReceived) {
				c.Info("normal read err:", e)
			} else {
				c.Info("read err:", e)
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
		c.Error("close err:", err)
	}
	c.DelConn(conn)
}

func (c *ConnManager) HandleMessage(data []byte, conn Connection) {
	packet, err := GetPacketByJson(data)
	if err != nil {
		c.Error("GetPacketByJson err:", err)
		return
	}

	if packet.CmdId != CmdLoginReq && conn.GetUserId() == 0 {
		c.Error("Please login first")
		return
	}

	var out MessageOut
	switch packet.CmdId {
	case CmdLoginReq:
		out = c.handleCmdLoginRep(conn, packet.Body)
	case CmdSendMessageReq:
		out = c.handleCmdSendMessageRep(conn, packet.Body)
	case CmdCreateConvReq:
		out = c.handleCmdCreateConvReq(packet.Body)
	default:
		return
	}

	outData, err := json.Marshal(out)
	if err != nil {
		c.Error("json err:", err)
		return
	}
	err = conn.Write(outData)
	if err != nil {
		c.Error("write err:", err)
	}
}

func (c *ConnManager) handleCmdCreateConvReq(data []byte) (out MessageOut) {
	var req MsgCreateConvReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		out.Msg = err.Error()
		return
	}

	if len(req.Uids) == 0 {
		out.Msg = "zero len"
		return
	}

	// TODO optimize: use redis
	convId := fmt.Sprintf("c%d", time.Now().Unix())
	c.conversionMap[convId] = req.Uids

	out.Code = CmdCreateConvResp
	out.Data = convId
	return
}

func (c *ConnManager) handleCmdSendMessageRep(conn Connection, data []byte) (out MessageOut) {
	var req MsgSendMessageRep
	err := json.Unmarshal(data, &req)
	if err != nil {
		out.Msg = err.Error()
		return
	}

	// TODO rpc send message

	uids, ok := c.conversionMap[req.ConvId]
	if !ok {
		out.Msg = "conversion not found"
		return
	}

	for _, uid := range uids {
		if uid == conn.GetUserId() {
			continue
		}
		err = c.SendMessage(uid, []byte(req.Msg))
		if err != nil {
			out.Msg = err.Error()
			continue
		}
	}

	out.Code = CmdSendMessageResp
	out.Msg = "send message ok"
	return
}

func (c *ConnManager) handleCmdLoginRep(conn Connection, data []byte) (out MessageOut) {
	var req MsgLoginRep
	err := json.Unmarshal(data, &req)
	if err != nil {
		out.Msg = err.Error()
		return
	}

	userInfo, err := c.handler.Login(req.Uid, req.Password)
	if err != nil {
		out.Msg = err.Error()
		return
	}

	conn.SetUserId(req.Uid)
	c.SaveConn(conn)

	out.Code = CmdLoginResp
	out.Msg = "login success"
	out.Data = userInfo.UserName
	return
}
