package main

import (
	"fmt"
	"log"

	"github.com/64bitAryan/distributedFileSystem/p2p"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	// fmt.Println("doing some login with the peer outside of TCP transport")
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAdder:  ":3000",
		HandshakeFun: p2p.NOPHandshakeFunc,
		Decoder:      &p2p.DefaultDecoder{},
		OnPeer:       OnPeer,
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
