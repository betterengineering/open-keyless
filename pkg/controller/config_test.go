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

package controller_test

import (
	"reflect"
	"testing"

	"github.com/betterengineering/open-keyless/pkg/application"
	"github.com/sirupsen/logrus"

	"github.com/betterengineering/open-keyless/pkg/controller"
	"github.com/betterengineering/open-keyless/pkg/datastore"
)

func TestNewControllerConfig(t *testing.T) {
	actual, err := controller.NewControllerConfig()
	if err != nil {
		t.Fatalf("could not create controller config - %s", err)
	}

	expected := controller.ControllerConfig{
		AirtableConfig: datastore.AirtableDatastoreConfig{
			Key:    "foo",
			BaseID: "bar",
		},
		ApplicationConfig: application.Config{
			LogLevel:       logrus.WarnLevel,
			MetricsEnabled: true,
			AdminInterface: ":9091",
		},
		TextFileConfig: datastore.TextFileConfig{
			Path: "/foo/ids.txt",
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected '%+v' does not equal actual '%+v'", expected, actual)
	}
}
