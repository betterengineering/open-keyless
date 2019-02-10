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

// Package datastore provides an interface and implementations for interacting with a badge datastore.
package datastore

const (
	// ErrBadgeDoesNotExist is returned when the badge requested does not exist in the datastore.
	ErrBadgeDoesNotExist = "the badge requested does not exist"
)

// Datastore is an interface for accessing a badge datastore.
type Datastore interface {
	HasAccess(id string) (bool, error)
	ListBadges() ([]Badge, error)
	CreateBadge(id string, badgeType string, enabled bool) error
	EnableBadge(id string) error
	DisableBadge(id string) error
	DeleteBadge(id string) error
	GetBadge(id string) (*Badge, error)
}

// Badge is a model for a badge in the datastore.
type Badge struct {
	// ID is the id burned into the RFID badge.
	ID string `json:"id"`

	// Enabled determines if the badge should be considered active.
	Enabled bool `json:"enabled"`

	// Type is the type of badge. E.x. card, sticker, keychain.
	Type string `json:"type"`
}
