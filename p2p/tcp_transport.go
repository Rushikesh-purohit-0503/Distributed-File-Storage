package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represent the remote node over TCPestablished connection.
type TCPPeer struct {
	// conn is underlying connection of the peer.
	conn net.Conn
	// if we dial & retrieve a conn => outbound = true.
	// if we accept and retreive a conn => outbound == false.
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}
type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan *RPC
	mu       sync.RWMutex
}

// Consume implements the Transport interface.
// It returns (read-only) channel for reading the incoming messages from another peer in network.
func (t *TCPTransport) Consume() <-chan *RPC {
	return t.rpcch
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan *RPC, 100), // Buffered channel to hold RPC messages
	}
}

// close implements the Peer interface.
// It closes the underlying connection of the peer.
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {

	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}
	//Read loop
	rpc := RPC{}
	// buff := make([]byte, 2000)
	for {
		// Read data from the connection into the buffer
		// n, err := conn.Read(buff)
		// if err != nil {
		// 	fmt.Printf("TCP read error: %s\n", err)
		// 	conn.Close()
		// 	return
		// }
		for {
			if err := t.Decoder.Decode(conn, &rpc); err != nil {
				fmt.Printf("TCP error: %s\n", err)
				continue
			}
			rpc.From = conn.RemoteAddr()
			t.rpcch <- &rpc
		}
	}
}
