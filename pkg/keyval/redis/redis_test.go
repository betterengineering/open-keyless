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

package redis_test

import (
	"bytes"
	"testing"

	"github.com/lodge93/open-keyless/pkg/keyval/redis"
)

const (
	TestKey = "foo"
)

type TestHarness struct {
	rkv *redis.RedisKeyVal
}

func SetupTest(t *testing.T) *TestHarness {
	rkv, err := redis.NewRedisKeyVal("127.0.0.1:6379", 10)
	if err != nil {
		t.Fatalf("could not instantiate test: '%s'", err)
	}

	return &TestHarness{
		rkv: rkv,
	}
}

func (testHarness *TestHarness) TearDownTest(t *testing.T) {
	err := testHarness.rkv.Delete(TestKey)
	if err != nil {
		t.Errorf("could not delete key: '%s'", err)
	}
}

func TestRedisKeyVal(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test since testing in short mode")
	}

	testHarness := SetupTest(t)
	defer testHarness.TearDownTest(t)

	expected := []byte("bar")
	err := testHarness.rkv.Set(TestKey, expected)
	if err != nil {
		t.Errorf("could not set key: '%s'", err)
	}

	actual, err := testHarness.rkv.Get(TestKey)
	if err != nil {
		t.Errorf("could not get key: '%s'", err)
	}

	if bytes.Compare(actual, expected) != 0 {
		t.Errorf("actual: '%s' does not match expected: '%s'", actual, expected)
	}
}
