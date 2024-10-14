package main

import (
	"fmt"

	"github.com/64bitAryan/distributedFileSystem/p2p"
)

type FileServerOpts struct {
	StorageRoot           string
	PathTransformFunction PathTransformFunction
	Transport             p2p.TCPTransport
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
	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)
		case <-s.quitCh:
			return

		}
	}
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.loop()

	return nil
}
