package main

import (
	"bytes"
	"log"
	"time"

	"github.com/64bitAryan/distributedFileSystem/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {

	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAdder:  listenAddr,
		HandshakeFun: p2p.NOPHandshakeFunc,
		Decoder:      &p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileOpts := FileServerOpts{
		StorageRoot:           listenAddr + "_network",
		PathTransformFunction: CASPathTransformerFunction,
		Transport:             tcpTransport,
		BootStrapNodes:        nodes,
	}

	s := NewFileServer(fileOpts)

	tcpTransport.OnPeer = s.OnPeer
	return s
}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()
	time.Sleep(3 * time.Second)

	go s2.Start()
	time.Sleep(3 * time.Second)

	data := bytes.NewReader([]byte("my big data file here!"))
	s2.StoreData("myprivatedata", data)

	select {}
}
