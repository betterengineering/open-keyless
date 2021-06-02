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

package strike_test

import (
	"testing"
	"time"

	"github.com/betterengineering/open-keyless/internal/mocks"
	"github.com/betterengineering/open-keyless/pkg/strike"
	"github.com/golang/mock/gomock"
	"periph.io/x/periph/conn/gpio"
)

func TestDoorStrikeUnlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPin := mocks.NewMockPinIO(ctrl)
	mockPin.EXPECT().Out(gpio.Low).Return(nil).Times(1)
	mockPin.EXPECT().Out(gpio.High).Return(nil).Times(1)
	mockPin.EXPECT().Halt().Return(nil).Times(1)

	ds, err := strike.NewDoorStrike(mockPin)
	if err != nil {
		t.Fatalf("error setting up test - %s", err)
	}
	defer ds.Done()

	err = ds.Unlock(0 * time.Millisecond)
	if err != nil {
		t.Errorf("error unlocking the strike - %s", err)
	}

	time.Sleep(1 * time.Millisecond)
}

func TestDoorStrikeUnlockMultipleTimes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPin := mocks.NewMockPinIO(ctrl)
	mockPin.EXPECT().Out(gpio.Low).Return(nil).Times(1)
	mockPin.EXPECT().Out(gpio.High).Return(nil).Times(1)
	mockPin.EXPECT().Halt().Return(nil).Times(1)

	ds, err := strike.NewDoorStrike(mockPin)
	if err != nil {
		t.Fatalf("error setting up test - %s", err)
	}
	defer ds.Done()

	for x := 0; x < 20; x++ {
		err = ds.Unlock(1 * time.Millisecond)
		if err != nil {
			t.Errorf("error unlocking the strike - %s", err)
		}
	}

	time.Sleep(5 * time.Millisecond)
}

func TestDoorStrikeUnlockWithGap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPin := mocks.NewMockPinIO(ctrl)
	mockPin.EXPECT().Out(gpio.Low).Return(nil).Times(2)
	mockPin.EXPECT().Out(gpio.High).Return(nil).Times(2)
	mockPin.EXPECT().Halt().Return(nil).Times(1)

	ds, err := strike.NewDoorStrike(mockPin)
	if err != nil {
		t.Fatalf("error setting up test - %s", err)
	}
	defer ds.Done()

	err = ds.Unlock(1 * time.Millisecond)
	if err != nil {
		t.Errorf("error unlocking the strike - %s", err)
	}

	time.Sleep(10 * time.Millisecond)

	err = ds.Unlock(1 * time.Millisecond)
	if err != nil {
		t.Errorf("error unlocking the strike - %s", err)
	}

	time.Sleep(10 * time.Millisecond)
}

func TestDoorStrikeAfterDoneIsCalled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPin := mocks.NewMockPinIO(ctrl)
	mockPin.EXPECT().Out(gpio.Low).Return(nil).Times(0)
	mockPin.EXPECT().Out(gpio.High).Return(nil).Times(0)
	mockPin.EXPECT().Halt().Return(nil).Times(1)

	ds, err := strike.NewDoorStrike(mockPin)
	if err != nil {
		t.Fatalf("error setting up test - %s", err)
	}

	ds.Done()

	err = ds.Unlock(0)
	if err.Error() != strike.ErrStrikeNotInitialized {
		t.Errorf("expected error does not match - %s", err)
	}
}
