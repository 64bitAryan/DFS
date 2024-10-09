package main

import "github.com/64bitAryan/distributedFileSystem/p2p"

type FileServerOpts struct {
	StorageRoot           string
	PathTransformFunction PathTransformFunction
	Transport             p2p.TCPTransport
}

type FileServer struct {
	FileServerOpts

	Store *Store
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:                  opts.StorageRoot,
		PathTransformFunction: opts.PathTransformFunction,
	}
	return &FileServer{
		FileServerOpts: opts,
		Store:          NewStore(storeOpts),
	}
}

func (s *FileServerOpts) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	return nil
}
