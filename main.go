package main

import (
	"log"
	"time"

	"github.com/64bitAryan/distributedFileSystem/p2p"
)

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAdder:  ":3000",
		HandshakeFun: p2p.NOPHandshakeFunc,
		Decoder:      &p2p.DefaultDecoder{},
		//TODO: OnPeer func
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileOpts := FileServerOpts{
		StorageRoot:           "3000_network",
		PathTransformFunction: CASPathTransformerFunction,
		Transport:             *tcpTransport,
	}

	s := NewFileServer(fileOpts)

	go func() {
		time.Sleep(time.Second * 3)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
