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

package controller

import (
	"errors"

	"github.com/betterengineering/open-keyless/pkg/application"
	"github.com/betterengineering/open-keyless/pkg/datastore"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	// ErrAirtableAPIKeyNotFound is returned when the Airtable API key is not found in the config.
	ErrAirtableAPIKeyNotFound = "could not find the required airtable API key in the config"

	// ErrAirtableBaseIDNotFound is returned when the Airtable Base ID is not found in the config.
	ErrAirtableBaseIDNotFound = "could not find the required airtable base ID in the config"
)

// ControllerConfig provides configuration for the Controller application.
type ControllerConfig struct {
	// AirtableConfig is a configuration object for the Airtable Datastore.
	AirtableConfig datastore.AirtableDatastoreConfig

	// ApplicationConfig is used to configure metrics and logging for the controller.
	ApplicationConfig application.Config

	// TextFileConfig is used to configure the TextFile config.
	TextFileConfig datastore.TextFileConfig
}

// NewControllerConfig provides a populated controller config from a configuration file.
func NewControllerConfig() (ControllerConfig, error) {
	err := configureViper()
	if err != nil {
		return ControllerConfig{}, err
	}

	airtableConifg, err := populateAirtableConfig()
	if err != nil {
		return ControllerConfig{}, err
	}

	applicationConfig, err := populateApplicationConfig()
	if err != nil {
		return ControllerConfig{}, err
	}

	textFileConfig := populateTextFileConfig()

	return ControllerConfig{
		AirtableConfig:    airtableConifg,
		ApplicationConfig: applicationConfig,
		TextFileConfig:    textFileConfig,
	}, nil
}

func populateApplicationConfig() (application.Config, error) {
	logLevelString := viper.GetString("application.logging.level")
	if logLevelString == "" {
		logLevelString = "info"
	}

	logLevel, err := log.ParseLevel(logLevelString)
	if err != nil {
		return application.Config{}, err
	}

	metricsEnabled := false

	if viper.IsSet("application.metrics.enabled") {
		metricsEnabled = viper.GetBool("application.metrics.enabled")
	}

	adminInterface := viper.GetString("application.admin.interface")
	if adminInterface == "" {
		adminInterface = ":8081"
	}

	return application.Config{
		LogLevel:       logLevel,
		MetricsEnabled: metricsEnabled,
		AdminInterface: adminInterface,
	}, nil
}

func configureViper() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/open-keyless-controller/")
	viper.AddConfigPath("testdata")
	viper.AddConfigPath("$HOME/.open-keyless-controller")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}

func populateAirtableConfig() (datastore.AirtableDatastoreConfig, error) {
	apiKey := viper.GetString("datastore.airtable.key")
	if apiKey == "" {
		return datastore.AirtableDatastoreConfig{}, errors.New(ErrAirtableAPIKeyNotFound)
	}

	baseID := viper.GetString("datastore.airtable.base")
	if baseID == "" {
		return datastore.AirtableDatastoreConfig{}, errors.New(ErrAirtableBaseIDNotFound)
	}

	return datastore.AirtableDatastoreConfig{
		Key:    apiKey,
		BaseID: baseID,
	}, nil
}

func populateTextFileConfig() datastore.TextFileConfig {
	return datastore.TextFileConfig{
		Path: viper.GetString("datastore.textFile.path"),
	}
}
