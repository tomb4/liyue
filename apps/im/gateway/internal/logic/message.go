package logic

import "encoding/json"

const (
	CmdLoginRep int16 = iota + 1
	CmdLoginResp
	CmdSendMessageRep
	CmdSendMessageResp
)

type (
	MsgLoginRep struct {
		Uid      int64  `json:"uid"`
		Password string `json:"password"`
	}

	MsgSendMessageRep struct {
		To  int64  `json:"to"`
		Msg string `json:"msg"`
	}
)

type (
	Packet struct {
		CmdId   int16
		BodyLen int
		Body    []byte
	}

	JsonMessage struct {
		CmdId int16  `json:"cmdId"`
		Body  string `json:"body"`
	}

	MessageOut struct {
		Code int16    `json:"code"`
		Msg  string `json:"msg"`
		Data string `json:"data"`
	}
)

func GetPacketByJson(bts []byte) (*Packet, error) {
	var jm JsonMessage
	err := json.Unmarshal(bts, &jm)
	if err != nil {
		return nil, err
	}

	return &Packet{
		CmdId:   jm.CmdId,
		BodyLen: len(jm.Body),
		Body:    []byte(jm.Body),
	}, nil
}
