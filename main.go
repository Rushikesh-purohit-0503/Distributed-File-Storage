package main

import (
	"fmt"
	"log"

	"github.com/Rushikesh-purohit-0503/Distributed-file-system/p2p"
)

func OnPeer(peer p2p.Peer) error {
	// peer.Close()
	fmt.Println("Hello there")
	return nil
}

func main() {

	tcpopts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.NOPDecoder{},
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpopts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("Received message: %+v\n", msg)
		}
	}()
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}

}
