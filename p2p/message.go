package p2p

import "net"

// RPC holds any arbitatry data that is being sent over the
// each transport. It is used to send data between nodes in the network.
type RPC struct {
	From net.Addr
	Payload []byte
}