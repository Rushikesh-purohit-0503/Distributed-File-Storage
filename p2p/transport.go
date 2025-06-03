package p2p

// peer is represantaion of the remote node.
type Peer interface{}

// Transport is anything that handles communication
// between nodes. This can be of the form of (TCP
// UDP, websockets, ...).
type Transport interface {
	ListenAndAccept() error
}
