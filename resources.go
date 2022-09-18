package goeinstein

import (
	"embed"
	"fmt"
	"io"
	"path/filepath"
)

//go:embed res/*
var resFS embed.FS

var resources = NewResources("res/")

type Resources struct {
	base string
}

func NewResources(basePath string) *Resources {
	r := &Resources{}
	r.base = basePath
	return r
}

func (r Resources) GetRef(name string) []byte {
	name = filepath.Join(r.base, name)

	file, err := resFS.Open(name)
	if err != nil {
		panic(fmt.Errorf("open %q file: %w", name, err))
	}
	bs, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("read all %q file (file close: %v): %w", name, file.Close(), err))
	}
	err = file.Close()
	if err != nil {
		panic(fmt.Errorf("close %q file: %w", name, err))
	}
	return bs
}

func (r Resources) DelRef(_ []byte) {}
