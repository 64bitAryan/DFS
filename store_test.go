package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mybestpictures"
	pathname := CASPathTransformerFunction(key)
	fmt.Println(pathname)
}

func TestDelete(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunction: CASPathTransformerFunction,
	}
	s := NewStore(opts)
	key := "mybestpictures"
	data := []byte("some jpg bytes")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunction: CASPathTransformerFunction,
	}
	s := NewStore(opts)
	key := "mybestpictures"
	data := []byte("some jpg bytes")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	r, err := s.Read(key)

	if err != nil {
		t.Error(err)
	}

	b, err := io.ReadAll(r)

	if string(b) != string(data) {
		t.Errorf("want %s has %s ", data, b)
		t.Error(err)
	}
}
