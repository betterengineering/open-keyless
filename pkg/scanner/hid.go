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

package scanner

import (
	"bytes"
	"encoding/hex"
	"errors"
	"log"
	"sync"

	"github.com/karalabe/hid"
)

// HidScanner implements the scanner interface for HID based scanning devices.
type HidScanner struct {
	device  *hid.Device
	ids     chan string
	errors  chan error
	quit    chan bool
	wg      *sync.WaitGroup
	started bool
}

// NewDefaultHidScanner provides an instantiated hid scanner device.
func NewDefaultHidScanner(ids chan string, errs chan error) (*HidScanner, error) {
	devices := hid.Enumerate(0, 0)
	if len(devices) == 0 {
		return nil, errors.New("no HID device found")
	}

	device, err := devices[0].Open()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	return &HidScanner{
		device:  device,
		ids:     ids,
		errors:  errs,
		started: false,
		quit:    make(chan bool),
		wg:      &wg,
	}, nil
}

// Scan starts the scanner
func (hid *HidScanner) Scan() {
	if hid.started {
		return
	}

	hid.wg.Add(1)
	go hid.run()
}

// Done closes down the scanner.
func (hid *HidScanner) Done() error {
	hid.quit <- true
	hid.wg.Wait()
	hid.started = false
	return hid.device.Close()
}

func (hid *HidScanner) run() {
	defer hid.wg.Done()

	for {
		select {
		case <-hid.quit:
			return
		default:
			hid.scan()
		}
	}
}

func (hid *HidScanner) scan() {
	err := hid.write(0x8f)
	if err != nil {
		log.Println("error writing to device")
		return
	}

	cardData, err := hid.read()
	if err != nil {
		log.Println("error reading from device")
		return
	}
	if cardData == nil {
		return
	}

	id := hex.EncodeToString(cardData)
	hid.ids <- id
}

func (hid *HidScanner) write(cmd byte) error {
	msg := make([]byte, 8)
	msg[0] = cmd
	_, err := hid.device.SendFeatureReport(msg)
	return err
}

func (hid *HidScanner) read() ([]byte, error) {
	data := make([]byte, 8)
	emptyData := make([]byte, 8)

	_, err := hid.device.GetFeatureReport(data)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(data, emptyData) {
		return nil, nil
	}

	return data, nil
}
