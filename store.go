package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

func CASPathTransformerFunction(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	blockSize := 5

	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}
}

type PathTransformFunction func(string) PathKey

type PathKey struct {
	PathName string
	FileName string
}

func (p *PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

var DefaultPathTransformFun = func(key string) string {
	return key
}

type StoreOpts struct {
	PathTransformFunction PathTransformFunction
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunction(key)

	_, err := os.Stat(pathKey.FullPath())
	if err == fs.ErrNotExist {
		return false
	}

	return true
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunction(key)
	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()
	return os.RemoveAll(pathKey.FullPath())
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunction(key)
	return os.Open(pathKey.FullPath())
}

func (s *Store) writeStream(key string, r io.Reader) error {
	// transforming the key to the path
	pathKey := s.PathTransformFunction(key)

	// making all the directories that path returns
	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	pathAndFileName := pathKey.FullPath()

	// creating the file with path
	f, err := os.Create(pathAndFileName)

	if err != nil {
		return err
	}

	// copying from buf to file
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	fmt.Printf("Written (%d) bytes to disk: %s", n, pathAndFileName)
	return nil
}
