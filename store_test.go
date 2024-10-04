package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mybestpictures"
	pathname := CASPathTransformerFunction(key)
	fmt.Println(pathname)
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunction: CASPathTransformerFunction,
	}
	s := NewStore(opts)
	data := bytes.NewReader([]byte("some jpg bytes"))

	if err := s.WriteStream("myspecialpicture", data); err != nil {
		t.Error(err)
	}
}
