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

// Package strike provides mechanisms to interface with an electric door strike via a Raspberry Pi GPIO interface. This
// package was designed with two requirements in mind. The first being that the strike should only need to be unlocked
// for a given duration of time and it should not be up to the consumer to ensure the strike is locked again. The second
// being that if the strike is requested to be unlocked while it is already unlocked, the duration elapsed should be the
// time already elapsed plus the last unlock duration.
package strike

import (
	"errors"
	"sync"
	"time"

	"periph.io/x/periph/conn/gpio"
)

const (
	// ErrStrikeNotInitialized is returned when unlock is called after Done() has been called on the strike.
	ErrStrikeNotInitialized = "strike cannot be unlocked after Done() has been called"
)

// Strike is an interface for an electric door strike. This interface was created for the sole purpose of generating a
// mock for this package which can be found in the internal/mocks package.
type Strike interface {
	Unlock(dur time.Duration) error
	Done()
}

// DoorStrike is an implementation of the strike interface to control an electric door strike with a Raspberry Pi via
// the GPIO interface.
type DoorStrike struct {
	initialized    bool
	pin            gpio.PinIO
	timeCtrlChan   chan time.Duration
	strikeCtrlChan chan bool
	quitStrike     chan bool
	quitTime       chan bool
	wg             *sync.WaitGroup
}

// NewDoorStrike returns an initialized DoorStrike ready for use.
func NewDoorStrike(pin gpio.PinIO) (*DoorStrike, error) {
	var wg sync.WaitGroup

	ds := &DoorStrike{
		initialized:    true,
		pin:            pin,
		timeCtrlChan:   make(chan time.Duration, 100),
		strikeCtrlChan: make(chan bool, 100),
		quitStrike:     make(chan bool),
		quitTime:       make(chan bool),
		wg:             &wg,
	}

	ds.wg.Add(2)
	go ds.timeControlLoop()
	go ds.strikeControlLoop()

	return ds, nil
}

// Done will clean up the electric strike cleanly and will allow for the GPIO pin to be reused.
func (ds *DoorStrike) Done() {
	ds.pin.Halt()
	ds.initialized = false
	ds.quitStrike <- true
	ds.quitTime <- true
	ds.wg.Wait()
}

// Unlock unlocks the electric door strike for the provided duration. Unlock is thread safe and can be called
// simultaneously from multiple threads. If the duration of a previous call to Unlock has not elapsed, the total
// duration will be the elapsed duration of the previous call plus the new duration.
func (ds *DoorStrike) Unlock(dur time.Duration) error {
	if !ds.initialized {
		return errors.New(ErrStrikeNotInitialized)
	}

	ds.timeCtrlChan <- dur
	return nil
}

func (ds *DoorStrike) timeControlLoop() {
	defer ds.wg.Done()

	timerStarted := false
	timer := time.NewTimer(0 * time.Millisecond)

	for {
		select {
		case dur := <-ds.timeCtrlChan:
			if timerStarted {
				timer.Reset(dur)
				continue
			}

			timer = time.NewTimer(dur)
			timerStarted = true
			ds.strikeCtrlChan <- true
		case <-timer.C:
			ds.strikeCtrlChan <- false
			timerStarted = false
		case <-ds.quitTime:
			return
		}
	}
}

func (ds *DoorStrike) strikeControlLoop() {
	defer ds.wg.Done()

	for {
		select {
		case open := <-ds.strikeCtrlChan:
			if open {
				ds.unlock()
				continue
			}

			ds.lock()
		case <-ds.quitStrike:
			return
		}
	}
}

func (ds *DoorStrike) lock() error {
	return ds.pin.Out(gpio.Low)
}

func (ds *DoorStrike) unlock() error {
	return ds.pin.Out(gpio.High)
}
