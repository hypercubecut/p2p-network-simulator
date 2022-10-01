package bitcoin

import (
	"p2psimulator/internal/bitcoin/msgtype"

	"github.com/bytedance/ns-x/v2/base"
)

type Packet struct {
	base.Packet

	MessageType msgtype.MessageType
	Payload     interface{}
	Source      *Node
	Destination *Node
}

func (m *Packet) Size() int {
	return 1
}

func (m *Packet) GetMessageType() msgtype.MessageType {
	return m.MessageType
}

func (m *Packet) GetPayload() interface{} {
	return m.Payload
}

func NewPacket(messageType msgtype.MessageType, payloadObj interface{}, src, des *Node) *Packet {
	return &Packet{
		MessageType: messageType,
		Payload:     payloadObj,
		Source:      src,
		Destination: des,
	}
}
