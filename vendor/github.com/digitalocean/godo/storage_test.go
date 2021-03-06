package godo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestStorageDrives_ListStorageDrives(t *testing.T) {
	setup()
	defer teardown()

	jBlob := `
	{
		"drives": [
			{
				"user_id": 42,
				"region": {"slug": "nyc3"},
				"id": "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
				"name": "my drive",
				"description": "my description",
				"size_gigabytes": 100,
				"droplet_ids": [10],
				"created_at": "2002-10-02T15:00:00.05Z"
			},
			{
				"user_id": 42,
				"region": {"slug": "nyc3"},
				"id": "96d414c6-295e-4e3a-ac59-eb9456c1e1d1",
				"name": "my other drive",
				"description": "my other description",
				"size_gigabytes": 100,
				"created_at": "2012-10-03T15:00:01.05Z"
			}
		],
		"links": {
	    "pages": {
	      "last": "https://api.digitalocean.com/v2/drives?page=2",
	      "next": "https://api.digitalocean.com/v2/drives?page=2"
	    }
	  },
	  "meta": {
	    "total": 28
	  }
	}`

	mux.HandleFunc("/v2/drives/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, jBlob)
	})

	drives, _, err := client.Storage.ListDrives(nil)
	if err != nil {
		t.Errorf("Storage.ListDrives returned error: %v", err)
	}

	expected := []Drive{
		{
			Region:        &Region{Slug: "nyc3"},
			ID:            "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
			Name:          "my drive",
			Description:   "my description",
			SizeGigaBytes: 100,
			DropletIDs:    []int{10},
			CreatedAt:     time.Date(2002, 10, 02, 15, 00, 00, 50000000, time.UTC),
		},
		{
			Region:        &Region{Slug: "nyc3"},
			ID:            "96d414c6-295e-4e3a-ac59-eb9456c1e1d1",
			Name:          "my other drive",
			Description:   "my other description",
			SizeGigaBytes: 100,
			CreatedAt:     time.Date(2012, 10, 03, 15, 00, 01, 50000000, time.UTC),
		},
	}
	if !reflect.DeepEqual(drives, expected) {
		t.Errorf("Storage.ListDrives returned %+v, expected %+v", drives, expected)
	}
}

func TestStorageDrives_Get(t *testing.T) {
	setup()
	defer teardown()
	want := &Drive{
		Region:        &Region{Slug: "nyc3"},
		ID:            "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
		Name:          "my drive",
		Description:   "my description",
		SizeGigaBytes: 100,
		CreatedAt:     time.Date(2002, 10, 02, 15, 00, 00, 50000000, time.UTC),
	}
	jBlob := `{
		"drive":{
			"region": {"slug":"nyc3"},
			"attached_to_droplet": null,
			"id": "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
			"name": "my drive",
			"description": "my description",
			"size_gigabytes": 100,
			"created_at": "2002-10-02T15:00:00.05Z"
		},
		"links": {
	    "pages": {
	      "last": "https://api.digitalocean.com/v2/drives?page=2",
	      "next": "https://api.digitalocean.com/v2/drives?page=2"
	    }
	  },
	  "meta": {
	    "total": 28
	  }
	}`

	mux.HandleFunc("/v2/drives/80d414c6-295e-4e3a-ac58-eb9456c1e1d1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, jBlob)
	})

	got, _, err := client.Storage.GetDrive("80d414c6-295e-4e3a-ac58-eb9456c1e1d1")
	if err != nil {
		t.Errorf("Storage.GetDrive returned error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Storage.GetDrive returned %+v, want %+v", got, want)
	}
}

func TestStorageDrives_Create(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &DriveCreateRequest{
		Region:        "nyc3",
		Name:          "my drive",
		Description:   "my description",
		SizeGibiBytes: 100,
	}

	want := &Drive{
		Region:        &Region{Slug: "nyc3"},
		ID:            "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
		Name:          "my drive",
		Description:   "my description",
		SizeGigaBytes: 100,
		CreatedAt:     time.Date(2002, 10, 02, 15, 00, 00, 50000000, time.UTC),
	}
	jBlob := `{
		"drive":{
			"region": {"slug":"nyc3"},
			"id": "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
			"name": "my drive",
			"description": "my description",
			"size_gigabytes": 100,
			"created_at": "2002-10-02T15:00:00.05Z"
		},
		"links": {}
	}`

	mux.HandleFunc("/v2/drives", func(w http.ResponseWriter, r *http.Request) {
		v := new(DriveCreateRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, createRequest) {
			t.Errorf("Request body = %+v, expected %+v", v, createRequest)
		}

		fmt.Fprint(w, jBlob)
	})

	got, _, err := client.Storage.CreateDrive(createRequest)
	if err != nil {
		t.Errorf("Storage.CreateDrive returned error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Storage.CreateDrive returned %+v, want %+v", got, want)
	}
}

func TestStorageDrives_Destroy(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/drives/80d414c6-295e-4e3a-ac58-eb9456c1e1d1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Storage.DeleteDrive("80d414c6-295e-4e3a-ac58-eb9456c1e1d1")
	if err != nil {
		t.Errorf("Storage.DeleteDrive returned error: %v", err)
	}
}

func TestStorageSnapshots_ListStorageSnapshots(t *testing.T) {
	setup()
	defer teardown()

	jBlob := `
	{
		"snapshots": [
			{
				"region": {"slug": "nyc3"},
				"id": "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
				"drive_id": "98d414c6-295e-4e3a-ac58-eb9456c1e1d1",
				"name": "my snapshot",
				"description": "my description",
				"size_gigabytes": 100,
				"created_at": "2002-10-02T15:00:00.05Z"
			},
			{
				"region": {"slug": "nyc3"},
				"id": "96d414c6-295e-4e3a-ac59-eb9456c1e1d1",
				"drive_id": "98d414c6-295e-4e3a-ac58-eb9456c1e1d1",
				"name": "my other snapshot",
				"description": "my other description",
				"size_gigabytes": 100,
				"created_at": "2012-10-03T15:00:01.05Z"
			}
		],
		"links": {
	    "pages": {
	      "last": "https://api.digitalocean.com/v2/drives?page=2",
	      "next": "https://api.digitalocean.com/v2/drives?page=2"
	    }
	  },
	  "meta": {
	    "total": 28
	  }
	}`

	mux.HandleFunc("/v2/drives/98d414c6-295e-4e3a-ac58-eb9456c1e1d1/snapshots", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, jBlob)
	})

	drives, _, err := client.Storage.ListSnapshots("98d414c6-295e-4e3a-ac58-eb9456c1e1d1", nil)
	if err != nil {
		t.Errorf("Storage.ListSnapshots returned error: %v", err)
	}

	expected := []Snapshot{
		{
			Region:        &Region{Slug: "nyc3"},
			ID:            "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
			DriveID:       "98d414c6-295e-4e3a-ac58-eb9456c1e1d1",
			Name:          "my snapshot",
			Description:   "my description",
			SizeGibiBytes: 100,
			CreatedAt:     time.Date(2002, 10, 02, 15, 00, 00, 50000000, time.UTC),
		},
		{
			Region:        &Region{Slug: "nyc3"},
			ID:            "96d414c6-295e-4e3a-ac59-eb9456c1e1d1",
			DriveID:       "98d414c6-295e-4e3a-ac58-eb9456c1e1d1",
			Name:          "my other snapshot",
			Description:   "my other description",
			SizeGibiBytes: 100,
			CreatedAt:     time.Date(2012, 10, 03, 15, 00, 01, 50000000, time.UTC),
		},
	}
	if !reflect.DeepEqual(drives, expected) {
		t.Errorf("Storage.ListSnapshots returned %+v, expected %+v", drives, expected)
	}
}

func TestStorageSnapshots_Get(t *testing.T) {
	setup()
	defer teardown()
	want := &Snapshot{
		Region:        &Region{Slug: "nyc3"},
		ID:            "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
		DriveID:       "98d414c6-295e-4e3a-ac58-eb9456c1e1d1",
		Name:          "my snapshot",
		Description:   "my description",
		SizeGibiBytes: 100,
		CreatedAt:     time.Date(2002, 10, 02, 15, 00, 00, 50000000, time.UTC),
	}
	jBlob := `{
		"snapshot":{
			"region": {"slug": "nyc3"},
			"id": "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
			"drive_id": "98d414c6-295e-4e3a-ac58-eb9456c1e1d1",
			"name": "my snapshot",
			"description": "my description",
			"size_gigabytes": 100,
			"created_at": "2002-10-02T15:00:00.05Z"
		},
		"links": {
	    "pages": {
				"last": "https://api.digitalocean.com/v2/drives/98d414c6-295e-4e3a-ac58-eb9456c1e1d1/snapshots?page=2",
				"next": "https://api.digitalocean.com/v2/drives/98d414c6-295e-4e3a-ac58-eb9456c1e1d1/snapshots?page=2"
	    }
	  },
	  "meta": {
	    "total": 28
	  }
	}`

	mux.HandleFunc("/v2/snapshots/80d414c6-295e-4e3a-ac58-eb9456c1e1d1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, jBlob)
	})

	got, _, err := client.Storage.GetSnapshot("80d414c6-295e-4e3a-ac58-eb9456c1e1d1")
	if err != nil {
		t.Errorf("Storage.GetSnapshot returned error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Storage.GetSnapshot returned %+v, want %+v", got, want)
	}
}

func TestStorageSnapshots_Create(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &SnapshotCreateRequest{
		DriveID:     "98d414c6-295e-4e3a-ac58-eb9456c1e1d1",
		Name:        "my snapshot",
		Description: "my description",
	}

	want := &Snapshot{
		Region:        &Region{Slug: "nyc3"},
		ID:            "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
		DriveID:       "98d414c6-295e-4e3a-ac58-eb9456c1e1d1",
		Name:          "my snapshot",
		Description:   "my description",
		SizeGibiBytes: 100,
		CreatedAt:     time.Date(2002, 10, 02, 15, 00, 00, 50000000, time.UTC),
	}
	jBlob := `{
		"snapshot":{
			"region": {"slug": "nyc3"},
			"id": "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
			"drive_id": "98d414c6-295e-4e3a-ac58-eb9456c1e1d1",
			"name": "my snapshot",
			"description": "my description",
			"size_gigabytes": 100,
			"created_at": "2002-10-02T15:00:00.05Z"
		},
		"links": {
	    "pages": {
	      "last": "https://api.digitalocean.com/v2/drives/98d414c6-295e-4e3a-ac58-eb9456c1e1d1/snapshots?page=2",
	      "next": "https://api.digitalocean.com/v2/drives/98d414c6-295e-4e3a-ac58-eb9456c1e1d1/snapshots?page=2"
	    }
	  },
	  "meta": {
	    "total": 28
	  }
	}`

	mux.HandleFunc("/v2/drives/98d414c6-295e-4e3a-ac58-eb9456c1e1d1/snapshots", func(w http.ResponseWriter, r *http.Request) {
		v := new(SnapshotCreateRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, createRequest) {
			t.Errorf("Request body = %+v, expected %+v", v, createRequest)
		}

		fmt.Fprint(w, jBlob)
	})

	got, _, err := client.Storage.CreateSnapshot(createRequest)
	if err != nil {
		t.Errorf("Storage.CreateSnapshot returned error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Storage.CreateSnapshot returned %+v, want %+v", got, want)
	}
}

func TestStorageSnapshots_Destroy(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/snapshots/80d414c6-295e-4e3a-ac58-eb9456c1e1d1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Storage.DeleteSnapshot("80d414c6-295e-4e3a-ac58-eb9456c1e1d1")
	if err != nil {
		t.Errorf("Storage.DeleteSnapshot returned error: %v", err)
	}
}
