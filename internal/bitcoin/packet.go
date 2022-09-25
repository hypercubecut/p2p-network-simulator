package bitcoin

import (
	"encoding/json"
	"p2psimulator/internal/bitcoin/msgtype"

	"github.com/bytedance/ns-x/v2/base"
)

type Packet struct {
	base.Packet

	MessageType msgtype.MessageType
	Payload     []byte
	Source      *Node
	Destination *Node
}

func (m *Packet) Size() int {
	return len(m.Payload)
}

func (m *Packet) GetMessageType() msgtype.MessageType {
	return m.MessageType
}

func (m *Packet) GetPayload() []byte {
	return m.Payload
}

func NewPacket(messageType msgtype.MessageType, payloadObject interface{}, src, des *Node) *Packet {
	bytePayload, _ := json.Marshal(payloadObject)

	return &Packet{
		MessageType: messageType,
		Payload:     bytePayload,
		Source:      src,
		Destination: des,
	}
}
