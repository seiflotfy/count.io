package storage

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/utils"

	"github.com/boltdb/bolt"
)

func setupTests() {
	os.Setenv("SKZ_DATA_DIR", "/tmp/skizze_storage_data")
	os.Setenv("SKZ_INFO_DIR", "/tmp/skizze_storage_info")
	path, err := os.Getwd()
	utils.PanicOnError(err)
	path = filepath.Dir(path)
	configPath := filepath.Join(path, "config/default.toml")
	os.Setenv("SKZ_CONFIG", configPath)
	tearDownTests()
}

func tearDownTests() {
	os.RemoveAll(config.GetConfig().DataDir)
	os.RemoveAll(config.GetConfig().InfoDir)
	os.Mkdir(config.GetConfig().DataDir, 0777)
	os.Mkdir(config.GetConfig().InfoDir, 0777)
}

func TestNoCounters(t *testing.T) {
	setupTests()
	defer tearDownTests()
	//FIXME: size of cache should be read from config
	m1 := newManager()
	m2 := newManager()
	m1.Create("marvel")
	data1 := []byte("wolverine")
	m1.SaveData("marvel", data1, 0)
	data2, err := m2.LoadData("marvel", 0, 0)
	if err != nil {
		t.Error("Expected no error loading data, got", err)
	}
	if bytes.Compare(data1, data2) != 0 {
		t.Error("Expected data2 == "+string(data1)+" got", data2)
	}
}

func TestGetAllInfo(t *testing.T) {
	setupTests()
	defer tearDownTests()
	func() {
		db, err := getInfoDB()
		if err != nil {
			t.Fatal(err)
		}

		err = db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("info"))
			err := bucket.Put([]byte("thing"), []byte(`{
				"id": "thing",
				"type": "default",
				"capacity": 12345
			}`))
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			t.Fatal(err)
		}
	}()
	m := newManager()
	infoData, err := m.LoadAllInfo()
	if err != nil {
		t.Fatal(err)
	}

	if len(infoData) != 1 {
		t.Error("Expected exactly one infoData, got", len(infoData))
	}
}

func TestDeleteInfo(t *testing.T) {
	setupTests()
	defer tearDownTests()
	func() {
		db, err := getInfoDB()
		if err != nil {
			t.Fatal(err)
		}

		err = db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("info"))
			err := bucket.Put([]byte("thing"), []byte(`{
				"id": "thing",
				"type": "default",
				"capacity": 12345
			}`))
			if err != nil {
				return err
			}
			err = bucket.Put([]byte("venom"), []byte(`{
				"id": "venom",
				"type": "default",
				"capacity": 67890
			}`))
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			t.Fatal(err)
		}
	}()
	m := newManager()

	m.DeleteInfo("venom")
	infoData, err := m.LoadAllInfo()
	if err != nil {
		t.Fatal(err)
	}

	if len(infoData) != 1 {
		t.Error("Expected exactly one infoData, got", len(infoData))
	}
}

func TestSaveAndDeleteData(t *testing.T) {
	setupTests()
	defer tearDownTests()
	m := newManager()
	m.Create("phoenix")
	m.SaveData("phoenix", []byte("phoenix"), 0)
	path := filepath.Join(config.GetConfig().DataDir, "phoenix")
	if _, err := os.Stat(path); err != nil {
		t.Error("Expected data in,", path, "got", err)
	}

	err := m.DeleteData("phoenix")
	if err != nil {
		t.Error("Expected no error deleting data, got", err)
	}
	if _, err := os.Stat(path); err == nil {
		t.Error("Expected no data in,", path, "got", err)
	}
}
