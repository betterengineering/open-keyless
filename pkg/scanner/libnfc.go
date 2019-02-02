// Copyright 2019 Mark Spicer
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

package scanner

import (
	"encoding/hex"
	"errors"
	"sync"

	"github.com/fuzxxl/nfc/2.0/nfc"
)

const (
	// ErrUnsupportedTagType is returned when a tag with an unsupported tag type was read from the scanner.
	ErrUnsupportedTagType = "read a device, but could not cast it to the proper tag type"
)

// LibNFCDevice is an interface used to generate a mock for nfc.Device.
type LibNFCDevice interface {
	InitiatorInit() error
	Close() error
	InitiatorListPassiveTargets(mod nfc.Modulation) ([]nfc.Target, error)
}

// LibNFCScanner implements the scanner interface for a libnfc compatible device.
type LibNFCScanner struct {
	device  LibNFCDevice
	mod     nfc.Modulation
	ids     chan string
	errors  chan error
	quit    chan bool
	wg      *sync.WaitGroup
	started bool
}

// NewDefaultLibNFCScanner provides an initialized libnfc scanner. The provided channels can be read off of in order to
// get a stream of id strings or errors from the device. The id and error channel should be buffered, otherwise the
// scanner will block until ids/errors are read off of the respective channel. Be sure to call Close when you are done
// with the scanner to clean up.
func NewDefaultLibNFCScanner(ids chan string, errs chan error) (*LibNFCScanner, error) {
	device, err := nfc.Open("")
	if err != nil {
		return nil, err
	}

	return NewLibNFCScanner(device, ids, errs)
}

// NewLibNFCScanner provides an initialized libnfc scanner with the provided LibNFCDevice. The provided channels can be
// read off of in order to get a stream of id strings or errors from the device. The id and error channel should be
// buffered, otherwise the scanner will block until ids/errors are read off of the respective channel. Be sure to call
// Close when you are done with the scanner to clean up.
func NewLibNFCScanner(device LibNFCDevice, ids chan string, errs chan error) (*LibNFCScanner, error) {
	err := device.InitiatorInit()
	if err != nil {
		return nil, err
	}

	mod := nfc.Modulation{
		Type:     nfc.ISO14443a,
		BaudRate: 1,
	}

	var wg sync.WaitGroup

	return &LibNFCScanner{
		device:  device,
		mod:     mod,
		ids:     ids,
		errors:  errs,
		started: false,
		quit:    make(chan bool),
		wg:      &wg,
	}, nil
}

// Scan starts the scanner if it has not already been started.
func (s *LibNFCScanner) Scan() {
	if s.started {
		return
	}

	s.wg.Add(1)
	go s.run()
}

// Done will stop all goroutines and close the LibNFCDevice.
func (s *LibNFCScanner) Done() error {
	s.quit <- true
	s.wg.Wait()
	s.started = false
	return s.device.Close()
}

func (s *LibNFCScanner) run() {
	defer s.wg.Done()

	for {
		select {
		case <-s.quit:
			return
		default:
			s.scan()
		}
	}
}

func (s *LibNFCScanner) scan() {
	targets, err := s.device.InitiatorListPassiveTargets(s.mod)
	if err != nil {
		s.errors <- err
		return
	}

	for _, target := range targets {
		t, ok := target.(*nfc.ISO14443aTarget)
		if !ok {
			s.errors <- errors.New(ErrUnsupportedTagType)
			continue
		}

		id := s.convert(t.UID[:t.UIDLen])
		s.ids <- id
	}
}

func (s *LibNFCScanner) convert(uid []byte) string {
	return hex.EncodeToString(uid)
}
