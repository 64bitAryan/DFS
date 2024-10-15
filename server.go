package main

import (
	"fmt"
	"log"

	"github.com/64bitAryan/distributedFileSystem/p2p"
)

type FileServerOpts struct {
	StorageRoot           string
	PathTransformFunction PathTransformFunction
	Transport             p2p.TCPTransport
	BootStrapNodes        []string
}

type FileServer struct {
	FileServerOpts

	store  *Store
	quitCh chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:                  opts.StorageRoot,
		PathTransformFunction: opts.PathTransformFunction,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitCh:         make(chan struct{}),
	}
}

func (s *FileServer) Stop() {
	close(s.quitCh)
}

func (s *FileServer) loop() {
	defer func() {
		s.Transport.Close()
		fmt.Println("File server stopped due to user quit channel")
	}()

	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)
		case <-s.quitCh:
			return

		}
	}
}

func (s *FileServer) bootStrapNetwork() error {
	for _, addr := range s.BootStrapNodes {
		if len(addr) == 0 {
			continue
		}

		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)
			}
		}(addr)
	}
	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.bootStrapNetwork()
	s.loop()

	return nil
}
