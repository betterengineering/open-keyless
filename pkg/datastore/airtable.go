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

package datastore

import (
	"errors"

	airtable "github.com/fabioberger/airtable-go"
)

// AirtableDatastore is an implementation of the datastore interface for Airtable.
type AirtableDatastore struct {
	client *airtable.Client
}

// AirtableDatastoreConfig is a configuration struct for an Airtable datastore.
type AirtableDatastoreConfig struct {
	// Key is a valid Airtable API key.
	Key string

	// BaseID is a valid Airtable base id.
	BaseID string
}

// AirtableBadge provides a wrapper around the badge model for Airtable requests.
type AirtableBadge struct {
	AirtableID string `json:"id"`
	Fields     Badge  `json:"fields"`
}

// AirtableBadgeTmp is used to create a badge because the API does not like when AirtableID exists in the POST even
// though it is empty.
type AirtableBadgeTmp struct {
	AirtableID string `json:"-"`
	Fields     Badge  `json:"fields"`
}

// NewAirTableDatastore provides an initialized datastore for Airtable using the provided configuration.
func NewAirTableDataStore(config AirtableDatastoreConfig) (*AirtableDatastore, error) {
	client, err := airtable.New(config.Key, config.BaseID)
	if err != nil {
		return nil, err
	}

	return &AirtableDatastore{
		client: client,
	}, nil
}

// HasAccess returns true if the badge with the given ID should be given access.
func (ds *AirtableDatastore) HasAccess(id string) (bool, error) {
	airtableBadges, err := ds.listBadges()
	if err != nil {
		return false, err
	}

	for _, badge := range airtableBadges {
		if id == badge.Fields.ID && badge.Fields.Enabled == true {
			return true, nil
		}
	}

	return false, nil
}

// ListBadges returns a list of badges from the datastore.
func (ds *AirtableDatastore) ListBadges() ([]Badge, error) {
	airtableBadges, err := ds.listBadges()
	if err != nil {
		return nil, err
	}

	badges := []Badge{}
	for _, badge := range airtableBadges {
		badges = append(badges, badge.Fields)
	}

	return badges, nil
}

// CreateBadge creates a badge in the Airflow datastore with the provided values.
func (ds *AirtableDatastore) CreateBadge(id string, badgeType string, enabled bool) error {
	airtableBadge := AirtableBadgeTmp{
		Fields: Badge{
			ID:      id,
			Type:    badgeType,
			Enabled: enabled,
		},
	}

	return ds.client.CreateRecord("badges", &airtableBadge)
}

// EnableBadge enables a badge that exists in the datastore.
func (ds *AirtableDatastore) EnableBadge(id string) error {
	badge, err := ds.getBadgeByID(id)
	if err != nil {
		return err
	}

	updatedFields := map[string]interface{}{
		"enabled": true,
	}

	return ds.client.UpdateRecord("badges", badge.AirtableID, updatedFields, &AirtableBadge{})
}

// DisableBadge disables a badge that exists in the datastore.
func (ds *AirtableDatastore) DisableBadge(id string) error {
	badge, err := ds.getBadgeByID(id)
	if err != nil {
		return err
	}

	updatedFields := map[string]interface{}{
		"enabled": false,
	}

	return ds.client.UpdateRecord("badges", badge.AirtableID, updatedFields, &AirtableBadge{})
}

// DeleteBadge deletes a badge from the datastore.
func (ds *AirtableDatastore) DeleteBadge(id string) error {
	badge, err := ds.getBadgeByID(id)
	if err != nil {
		return err
	}

	return ds.client.DestroyRecord("badges", badge.AirtableID)
}

// GetBadge returns a badge with the given ID. If the badge does not exist, ErrBadgeDoesNotExist will be returned.
func (ds *AirtableDatastore) GetBadge(id string) (*Badge, error) {
	badge, err := ds.getBadgeByID(id)
	if err != nil {
		return nil, err
	}

	return &badge.Fields, nil
}

func (ds *AirtableDatastore) getBadgeByID(id string) (*AirtableBadge, error) {
	airtableBadges, err := ds.listBadges()
	if err != nil {
		return nil, err
	}

	for _, badge := range airtableBadges {
		if badge.Fields.ID == id {
			return &badge, nil
		}
	}

	return nil, errors.New(ErrBadgeDoesNotExist)
}

func (ds *AirtableDatastore) listBadges() ([]AirtableBadge, error) {
	airtableBadges := []AirtableBadge{}

	err := ds.client.ListRecords("badges", &airtableBadges)
	if err != nil {
		return nil, err
	}

	return airtableBadges, nil
}
