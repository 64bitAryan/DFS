package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func CASPathTransformerFunction(key string) string {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	blockSize := 5

	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return strings.Join(paths, "/")
}

type PathTransformFunction func(string) string

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

func (s *Store) WriteStream(key string, r io.Reader) error {
	// transforming the key to the path
	path := s.PathTransformFunction(key)

	// making all the directories that path returns
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	// creating buffer, hashing the filename and
	// joining with path
	buf := new(bytes.Buffer)
	io.Copy(buf, r)

	fileNameBytes := md5.Sum(buf.Bytes())
	fileName := hex.EncodeToString(fileNameBytes[:])
	pathAndFileName := path + "/" + fileName

	// creating the file with path
	f, err := os.Create(pathAndFileName)

	if err != nil {
		return err
	}

	// copying from buf to file
	n, err := io.Copy(f, buf)
	if err != nil {
		return err
	}

	fmt.Printf("Written (%d) bytes to disk: %s", n, pathAndFileName)
	return nil
}
