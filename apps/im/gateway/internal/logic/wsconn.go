package logic

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/core/utils"
	"liyue/apps/im/gateway/internal/svc"
	"net/http"
)

type (
	websocketConn struct {
		*conn
		websocketConn *websocket.Conn
	}

	WebsocketServer struct {
		logx.Logger
		srv    *http.Server
		svcCtx *svc.ServiceContext
	}
)

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebsocketServer(svcCtx *svc.ServiceContext) WebsocketServer {
	s := WebsocketServer{
		Logger: logx.WithContext(context.TODO()),
		svcCtx: svcCtx,
	}

	router := mux.NewRouter()
	router.HandleFunc("/echo", s.echo)

	s.srv = &http.Server{
		Addr:    ":20001",
		Handler: router,
	}
	return s
}

func (s WebsocketServer) ListenAndServe() {
	threading.GoSafe(func() {
		err := s.srv.ListenAndServe()
		if err != nil {
			s.Error(err)
		}
	})
}

func (s WebsocketServer) echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.Close()
		return
	}

	webConn := &websocketConn{
		conn: &conn{
			connectionId: utils.NewUuid(),
			createdAt:    utils.CurrentMillis(),
			clientIp:     r.Header.Get("X-Forwarded-For"),
		},
		websocketConn: c,
	}
	if webConn.clientIp == "" {
		webConn.clientIp = r.RemoteAddr
	}

	OnceConnManager().Dispatch(webConn)
}

func (m *websocketConn) Read() (p []byte, err error) {
	_, bytes, err := m.websocketConn.ReadMessage()
	if err != nil {
		return nil, err
	}
	//if messageType != websocket.BinaryMessage {
	//	return nil, errors.New("bad type")
	//}
	return bytes, nil
}

func (m *websocketConn) Write(p []byte) error {
	m.writeLock.Lock()
	defer m.writeLock.Unlock()
	if m.websocketConn == nil {
		return errors.New("websocketConn conn is nil")
	}
	err := m.websocketConn.WriteMessage(websocket.BinaryMessage, p)
	return err
}

func (m *websocketConn) Close() error {
	if m.closed {
		return nil
	}
	err := m.websocketConn.Close()
	m.closed = true
	return err
}
