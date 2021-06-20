// Copyright 2021 Mark Spicer
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
// documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
// Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package datastore

import (
	"errors"
	"io/ioutil"
	"log"
	"strings"
)

// TextFile implements the datastore interface with a file.
type TextFile struct {
	ids []string
}

// TextFileConfig is a configuration struct for a a TextFile datastore.
type TextFileConfig struct {
	Path string
}

// NewTextFile provides an instantiated datastore.
func NewTextFile(cfg TextFileConfig) (*TextFile, error) {
	content, err := ioutil.ReadFile(cfg.Path)
	if err != nil {
		return nil, err
	}

	ids := strings.Split(string(content), "\n")
	return &TextFile{
		ids: ids,
	}, nil
}

func (txt *TextFile) HasAccess(id string) (bool, error) {
	log.Printf("Verifying access for %s", id)
	for _, cur := range txt.ids {
		if id == cur {
			log.Printf("access granted for %s", id)
			return true, nil
		}
	}

	log.Printf("access denied for %s", id)
	return false, nil
}

func (txt *TextFile) ListBadges() ([]Badge, error) {
	return nil, errors.New("not implemented")
}

func (txt *TextFile) CreateBadge(id string, badgeType string, enabled bool) error {
	return errors.New("not implemented")
}

func (txt *TextFile) EnableBadge(id string) error {
	return errors.New("not implemented")
}

func (txt *TextFile) DisableBadge(id string) error {
	return errors.New("not implemented")
}

func (txt *TextFile) DeleteBadge(id string) error {
	return errors.New("not implemented")
}

func (txt *TextFile) GetBadge(id string) (*Badge, error) {
	return nil, errors.New("not implemented")
}
