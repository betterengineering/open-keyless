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

// Package application provides common utilities to create open-keyless applications with.
package application

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// Application is an object used to provide common application functionality for open-keyless based applications.
type Application struct {
	AppType string
	Config  Config
}

// NewApplication provides an instantiated Application object with the provided configuration and type.
func NewApplication(config Config, appType string) *Application {
	app := &Application{
		AppType: appType,
		Config:  config,
	}

	if config.MetricsEnabled {
		app.enablePrometheusMetrics()
	}

	app.configureLogging()

	return app
}

// PrintBanner prints a banner message. This should be called once the application has been fully started.
func (app *Application) PrintBanner() {
	banner := app.getBannerText()
	version := "0.1.0"
	fmt.Fprint(os.Stderr, fmt.Sprintf(banner, version))
}

func (app *Application) getBannerText() string {
	switch app.AppType {
	case OpenKeylessController:
		return `   ____                      __ __           __              
  / __ \____  ___  ____     / //_/__  __  __/ /__  __________
 / / / / __ \/ _ \/ __ \   / ,< / _ \/ / / / / _ \/ ___/ ___/
/ /_/ / /_/ /  __/ / / /  / /| /  __/ /_/ / /  __(__  |__  ) 
\____/ .___/\___/_/ /_/  /_/ |_\___/\__, /_/\___/____/____/  
    /_/                            /____/                    

Open Keyless Controller
Version %s
`
	default:
		return ""
	}
}

func (app *Application) enablePrometheusMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(app.Config.AdminInterface, nil)
}

func (app *Application) configureLogging() {
	log.SetLevel(app.Config.LogLevel)
}
