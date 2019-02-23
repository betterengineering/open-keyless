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

// Package controller provides the primary application logic for the Open Keyless controller module.
package controller

import (
	"time"

	"github.com/lodge93/open-keyless/pkg/application"
	"github.com/lodge93/open-keyless/pkg/datastore"
	"github.com/lodge93/open-keyless/pkg/scanner"
	"github.com/lodge93/open-keyless/pkg/strike"
	log "github.com/sirupsen/logrus"
)

// Controller is the primary struct for Open Keyless controller.
type Controller struct {
	datastore   datastore.Datastore
	scanner     scanner.Scanner
	application *application.Application
	strike      strike.Strike
	ids         chan string
	errors      chan error
}

// NewController provides an initialized Controller with the provided configuration.
func NewController(config ControllerConfig) (*Controller, error) {
	app := application.NewApplication(config.ApplicationConfig, application.OpenKeylessController)

	ds, err := datastore.NewAirTableDataStore(config.AirtableConfig)
	if err != nil {
		log.WithFields(log.Fields{
			"application": app.AppType,
			"error":       err,
		}).Error("could not connect to the airtable datastore")
		return nil, err
	}

	str, err := strike.NewDefaultDoorStrike()
	if err != nil {
		log.WithFields(log.Fields{
			"application": app.AppType,
			"error":       err,
		}).Error("could not connect to door strike")
		return nil, err
	}

	ids := make(chan string, 100)
	errs := make(chan error, 100)

	scn, err := scanner.NewDefaultLibNFCScanner(ids, errs)
	if err != nil {
		log.WithFields(log.Fields{
			"application": app.AppType,
			"error":       err,
		}).Error("could not connect to the NFC scanner")
		return nil, err
	}

	return &Controller{
		datastore:   ds,
		application: app,
		scanner:     scn,
		strike:      str,
		ids:         ids,
		errors:      errs,
	}, nil
}

// Run will run the controller in a blocking fashion.
func (c *Controller) Run() {
	defer c.strike.Done()
	defer c.scanner.Done()

	c.scanner.Scan()
	c.application.PrintBanner()

	log.WithFields(log.Fields{
		"application": c.application.AppType,
	}).Info("scanning for badges")

	for {
		select {
		case id := <-c.ids:
			log.WithFields(log.Fields{
				"application": c.application.AppType,
				"id":          id,
			}).Debug("found badge id")
			c.processID(id)
		case err := <-c.errors:
			log.WithFields(log.Fields{
				"application": c.application.AppType,
				"error":       err,
			}).Error("encountered an error while scanning for badges")
		}
	}
}

func (c *Controller) processID(id string) {
	hasAccess, err := c.datastore.HasAccess(id)
	if err != nil {
		log.WithFields(log.Fields{
			"application": c.application.AppType,
			"error":       err,
		}).Error("error communicating with airtable")
		return
	}

	if hasAccess {
		c.grantAccess(id)
		return
	}

	log.WithFields(log.Fields{
		"application": c.application.AppType,
		"id":          id,
	}).Info("access denied for badge id")
}

func (c *Controller) grantAccess(id string) {
	log.WithFields(log.Fields{
		"application": c.application.AppType,
		"id":          id,
	}).Info("allowing access for badge id")

	err := c.strike.Unlock(time.Second * 3)
	if err != nil {
		log.WithFields(log.Fields{
			"application": c.application.AppType,
			"error":       err,
		}).Error("error unlocking strike for id")
	}
}
