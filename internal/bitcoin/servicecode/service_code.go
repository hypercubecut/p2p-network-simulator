package servicecode

// Bitcoin service code
const (
	// Unnamed This node is not a full node. It may not be able to provide any data except for the transactions it originates.
	Unnamed = 0

	// FullNode This is a full node and can be asked for full blocks. It should implement all protocol
	// features available in its self-reported protocol version.
	FullNode = 1

	// MinerNode is a miner's node
	MinerNode = 2
)
