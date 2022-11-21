package bitcoin

import (
	"fmt"
	"p2psimulator/internal/bitcoin/msgtype"
	"time"

	"github.com/bytedance/ns-x/v2/base"
)

func (n *Node) handleP2PWithBroadcast(request *Packet, nodes map[string]base.Node, now time.Time) []base.Event {
	reqDTO := request.Payload.(*P2PMessage)

	if n.name == reqDTO.Dest {
		n.logger.Info(fmt.Sprintf("%s peer received message, %s", n.name, now))
		return nil
	}

	_, ok := n.cache[reqDTO.MsgID]
	if ok {
		return nil
	}

	n.cache[reqDTO.MsgID] = struct{}{}

	return n.sendPacketToPeerNodes(&Packet{
		MessageType: msgtype.P2PWithBroadcastMessageType,
		Payload:     reqDTO,
		Source:      n,
	}, nodes, now)
}
