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

package datastore_test

import (
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/lodge93/open-keyless/pkg/datastore"
)

const (
	DisabledBadgeID  = "disabledBadgeID"
	EnabledBadgeID   = "enabledBadgeID"
	DefaultBadgeType = "card"
)

var DefaultBadges = []datastore.Badge{
	{
		ID:      DisabledBadgeID,
		Enabled: false,
		Type:    DefaultBadgeType,
	},
	{
		ID:      EnabledBadgeID,
		Enabled: true,
		Type:    DefaultBadgeType,
	},
}

func SetupTest(t *testing.T) *datastore.AirtableDatastore {
	if testing.Short() {
		t.Skip("skipping integration test due to short flag")
	}

	config := datastore.AirtableDatastoreConfig{
		Key:    os.Getenv("AIRTABLE_API_KEY"),
		BaseID: os.Getenv("AIRTABLE_BASE_ID"),
	}

	ds, err := datastore.NewAirTableDataStore(config)
	if err != nil {
		t.Fatalf("error setting up test - %s", err)
	}

	return ds
}

func SetupTestWithBadges(t *testing.T) *datastore.AirtableDatastore {
	ds := SetupTest(t)
	for _, badge := range DefaultBadges {
		err := ds.CreateBadge(badge.ID, badge.Type, badge.Enabled)
		if err != nil {
			t.Fatalf("could not setup test with badges - %s", err)
		}
	}

	return ds
}

func CleanupTest(t *testing.T, ds *datastore.AirtableDatastore) {
	for _, badge := range DefaultBadges {
		err := ds.DeleteBadge(badge.ID)
		if err != nil {
			t.Fatalf("could not cleanup test with badges - %s", err)
		}
	}
}

func TestAirTableDataStoreHasAccess(t *testing.T) {
	ds := SetupTestWithBadges(t)
	defer CleanupTest(t, ds)

	hasAccess, err := ds.HasAccess(EnabledBadgeID)
	if err != nil {
		t.Fatalf("could not complete request to airtable - %s", err)
	}

	if !hasAccess {
		t.Errorf("has access expected to be true")
	}
}

func TestAirTableDataStoreHasAccessDenied(t *testing.T) {
	ds := SetupTestWithBadges(t)
	defer CleanupTest(t, ds)

	hasAccess, err := ds.HasAccess(DisabledBadgeID)
	if err != nil {
		t.Fatalf("could not complete request to airtable - %s", err)
	}

	if hasAccess {
		t.Errorf("has access expected to be false")
	}
}

func TestAirTableDataStoreListBadges(t *testing.T) {
	ds := SetupTestWithBadges(t)
	defer CleanupTest(t, ds)

	badges, err := ds.ListBadges()
	if err != nil {
		t.Fatalf("could not complete request to airtable - %s", err)
	}

	sort.Slice(badges, func(i, j int) bool { return badges[i].ID < badges[j].ID })
	sort.Slice(DefaultBadges, func(i, j int) bool { return DefaultBadges[i].ID < DefaultBadges[j].ID })

	if !reflect.DeepEqual(badges, DefaultBadges) {
		t.Errorf("'%v' not equal to '%v'", badges, DefaultBadges)
	}
}

func TestAirTableDataStoreCreateAndDeleteBadge(t *testing.T) {
	ds := SetupTest(t)

	id := "foo"
	badgeType := "card"
	enabled := true

	err := ds.CreateBadge(id, badgeType, enabled)
	if err != nil {
		t.Fatalf("could not create badge - %s", err)
	}

	err = ds.DeleteBadge(id)
	if err != nil {
		t.Fatalf("could not delete badge - %s", err)
	}
}

func TestAirTableDataStoreEnableBadge(t *testing.T) {
	ds := SetupTestWithBadges(t)
	defer CleanupTest(t, ds)

	err := ds.EnableBadge(DisabledBadgeID)
	if err != nil {
		t.Fatalf("could not enable badge - %s", err)
	}

	badge, err := ds.GetBadge(DisabledBadgeID)
	if err != nil {
		t.Fatalf("could not get badge - %s", err)
	}

	if badge.Enabled != true {
		t.Error("badge was not enabled")
	}
}

func TestAirTableDataStoreGetBadge(t *testing.T) {
	ds := SetupTestWithBadges(t)
	defer CleanupTest(t, ds)

	badge, err := ds.GetBadge(DisabledBadgeID)
	if err != nil {
		t.Fatalf("could not get badge - %s", err)
	}

	if badge.ID != DisabledBadgeID {
		t.Errorf("badge ID '%s' does not match expected '%s'", badge.ID, DisabledBadgeID)
	}

	if badge.Enabled != false {
		t.Errorf("badge should not be enabled")
	}

	if badge.Type != DefaultBadgeType {
		t.Errorf("badge type is '%s' not '%s'", badge.Type, DefaultBadgeType)
	}
}

func TestAirTableDataStoreDisableBadge(t *testing.T) {
	ds := SetupTestWithBadges(t)
	defer CleanupTest(t, ds)

	err := ds.DisableBadge(EnabledBadgeID)
	if err != nil {
		t.Fatalf("could not disable badge - %s", err)
	}

	badge, err := ds.GetBadge(EnabledBadgeID)
	if err != nil {
		t.Fatalf("could not get badge - %s", err)
	}

	if badge.Enabled != false {
		t.Error("badge was not disabled")
	}
}
