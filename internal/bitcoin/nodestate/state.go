package nodestate

type State string

const (
	New        State = "new"
	Offline    State = "offline"
	Connecting State = "connecting"
	Connected  State = "connected"
)
