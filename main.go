package main

import (
	"log"

	"github.com/64bitAryan/distributedFileSystem/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {

	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAdder:  listenAddr,
		HandshakeFun: p2p.NOPHandshakeFunc,
		Decoder:      &p2p.DefaultDecoder{},
		//TODO: OnPeer func
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileOpts := FileServerOpts{
		StorageRoot:           listenAddr + "_network",
		PathTransformFunction: CASPathTransformerFunction,
		Transport:             *tcpTransport,
		BootStrapNodes:        nodes,
	}

	return NewFileServer(fileOpts)

}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	s2.Start()
}
