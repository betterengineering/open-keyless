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

package buzzer

import (
	"errors"
	"log"
	"sync"
	"time"

	"periph.io/x/periph/conn/gpio"

	"periph.io/x/periph/conn/gpio/gpioreg"
)

const (
	// ErrCouldNotInitializeGPIOPin is returned when the GPIO pin for the buzzer can not be initialized.
	ErrCouldNotInitializeGPIOPin = "could not initialize the GPIO pin for the buzzer"

	BuzzerNotInitialized = "the buzzer cannot be buzzed before being initialized"
)

type Buzzer struct {
	initialized bool
	pin         gpio.PinIO
	ctrl        chan time.Duration
	quit        chan bool
	wg          *sync.WaitGroup
}

func NewBuzzer() (*Buzzer, error) {
	var wg sync.WaitGroup

	pin := gpioreg.ByName("23")
	if pin == nil {
		return nil, errors.New(ErrCouldNotInitializeGPIOPin)
	}

	buzz := &Buzzer{
		initialized: true,
		pin:         pin,
		ctrl:        make(chan time.Duration, 100),
		quit:        make(chan bool),
		wg:          &wg,
	}

	buzz.wg.Add(1)
	go buzz.buzzerController()

	return buzz, nil
}

func (buzz *Buzzer) Buzz() error {
	if !buzz.initialized {
		return errors.New(BuzzerNotInitialized)
	}

	buzz.ctrl <- 1 * time.Second

	return nil
}

func (buzz *Buzzer) Done() {
	buzz.pin.Halt()
	buzz.initialized = false
	buzz.quit <- true
	buzz.wg.Wait()
}

func (buzz *Buzzer) buzzerController() {
	defer buzz.wg.Done()

	for {
		select {
		case durr := <-buzz.ctrl:
			err := buzz.pin.Out(gpio.High)
			if err != nil {
				log.Printf("error settting buzzer - %s", err)
				continue
			}

			time.Sleep(durr)
			err = buzz.pin.Out(gpio.Low)
			if err != nil {
				log.Printf("error halting buzzer - %s", err)
				continue
			}
		case <-buzz.quit:
			return
		}
	}
}
