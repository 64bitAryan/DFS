package main

import (
	"log"

	"github.com/64bitAryan/distributedFileSystem/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAdder:  ":3000",
		HandshakeFun: p2p.NOPHandshakeFunc,
		Decoder:      &p2p.GOBDecoder{},
	}
	tr := p2p.NewTCPTransport(tcpOpts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
