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

package scanner_test

import (
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/fuzxxl/nfc/2.0/nfc"
	"github.com/golang/mock/gomock"
	"github.com/lodge93/open-keyless/internal/mocks"
	"github.com/lodge93/open-keyless/pkg/scanner"
)

func TestLibNFCScanner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ids := make(chan string, 100)
	errs := make(chan error, 100)
	device, err := givenInitializedDevice(ctrl)
	if err != nil {
		t.Fatalf("error setting up test - %s", err)
	}

	s, err := scanner.NewLibNFCScanner(device, ids, errs)
	if err != nil {
		t.Fatalf("error setting up test - %s", err)
	}

	var errorFound bool
	var idFound bool
	go func() {
		for {
			select {
			case <-ids:
				idFound = true
			case <-errs:
				errorFound = true
			}
		}
	}()

	s.Scan()
	time.Sleep(5 * time.Millisecond)
	s.Done()

	if errorFound {
		t.Errorf("error found while scanning for devices")
	}

	if !idFound {
		t.Errorf("there were no ids found while scanning")
	}
}

func TestLibNFCScannerWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ids := make(chan string, 100)
	errs := make(chan error, 100)
	device, err := givenInitializedErrorDevice(ctrl)
	if err != nil {
		t.Fatalf("error setting up test - %s", err)
	}

	s, err := scanner.NewLibNFCScanner(device, ids, errs)
	if err != nil {
		t.Fatalf("error setting up test - %s", err)
	}

	var errorFound bool
	var idFound bool
	go func() {
		for {
			select {
			case <-ids:
				idFound = true
			case <-errs:
				errorFound = true
			}
		}
	}()

	s.Scan()
	time.Sleep(1 * time.Millisecond)
	s.Done()

	if !errorFound {
		t.Errorf("there was not an error found when one was expected")
	}

	if idFound {
		t.Errorf("there were no ids found while scanning")
	}
}

func givenInitializedDevice(ctrl *gomock.Controller) (scanner.LibNFCDevice, error) {
	device := mocks.NewMockLibNFCDevice(ctrl)

	mod := nfc.Modulation{
		Type:     nfc.ISO14443a,
		BaudRate: 1,
	}

	targets, err := generateFakeTargets()
	if err != nil {
		return nil, err
	}

	device.EXPECT().InitiatorInit().Return(nil).Times(1)
	device.EXPECT().InitiatorListPassiveTargets(mod).Return(targets, nil).AnyTimes()
	device.EXPECT().Close().Return(nil).Times(1)

	return device, nil
}

func givenInitializedErrorDevice(ctrl *gomock.Controller) (scanner.LibNFCDevice, error) {
	device := mocks.NewMockLibNFCDevice(ctrl)

	mod := nfc.Modulation{
		Type:     nfc.ISO14443a,
		BaudRate: 1,
	}

	device.EXPECT().InitiatorInit().Return(nil).Times(1)
	device.EXPECT().InitiatorListPassiveTargets(mod).Return(nil, errors.New("there was an error")).AnyTimes()
	device.EXPECT().Close().Return(nil).Times(1)

	return device, nil
}

func generateFakeTargets() ([]nfc.Target, error) {
	var uid [10]byte

	b, err := hex.DecodeString("8604de7d")
	if err != nil {
		return nil, err
	}
	copy(uid[:], b)

	targets := []nfc.Target{
		&nfc.ISO14443aTarget{
			UID:    uid,
			UIDLen: len(b),
		},
	}

	return targets, nil
}
