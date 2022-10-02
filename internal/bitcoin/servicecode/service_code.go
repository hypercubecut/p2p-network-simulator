package servicecode

// Bitcoin service code
const (
	// Unnamed This node is not a full node. It may not be able to provide any data except for the transactions it originates.
	Unnamed = 0

	// NODE_NETWORK This is a full node and can be asked for full blocks. It should implement all protocol
	// features available in its self-reported protocol version.
	NODE_NETWORK = 1

	// NODE_GETUTXO This is a full node capable of responding to the getutxo protocol request.
	// This is not supported by any currently-maintained Bitcoin node. See BIP64 for details on how this is implemented.
	NODE_GETUTXO = 2

	// NODE_BLOOM This is a full node capable and willing to handle bloom-filtered connections. See BIP111 for details.
	NODE_BLOOM = 4

	// NODE_WITNESS This is a full node that can be asked for blocks and transactions including witness data.
	// See BIP144 for details.
	NODE_WITNESS = 8

	// NODE_XTHIN This is a full node that supports Xtreme Thinblocks.
	// This is not supported by any currently-maintained Bitcoin node.
	NODE_XTHIN = 16

	// NODE_NETWORK_LIMITEDT this is the same as NODE_NETWORK but the node has at least the last 288 blocks (last 2 days).
	// See BIP159 for details on how this is implemented.
	NODE_NETWORK_LIMITED = 1024
)
