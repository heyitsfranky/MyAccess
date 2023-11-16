package MyAccess

import (
	"encoding/json"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
)

const config_path = "config.yaml"

func TestInit(t *testing.T) {
	err := Init(config_path)
	if err != nil {
		t.Fatalf("Unexpected error during initialization: %v", err)
	}
}

func TestReadJSON(t *testing.T) {
	err := Init(config_path)
	if err != nil {
		t.Fatalf("Error during initialization: %v", err)
	}
	key := "nonexistent_key"
	value, err := ReadJSON(key)
	if err != nil {
		t.Fatalf("Unexpected error while reading JSON: %v", err)
	}
	if value != nil {
		t.Fatalf("Expected nil value for non-existing key, but got: %v", value)
	}
	existingKey := "existing_key"
	existingValue := map[string]interface{}{"name": "John", "age": 30}
	jsonValue, _ := json.Marshal(existingValue)
	memDB.Set(&memcache.Item{Key: existingKey, Value: jsonValue})
	value, err = ReadJSON(existingKey)
	if err != nil {
		t.Fatalf("Unexpected error while reading JSON: %v", err)
	}
	if !compareJSON(value, existingValue) {
		t.Fatalf("Expected value: %v, but got: %v", existingValue, value)
	}
}

func TestRead(t *testing.T) {
	err := Init(config_path)
	if err != nil {
		t.Fatalf("Error during initialization: %v", err)
	}
	key := "nonexistent_key"
	value, err := Read(key)
	if err != nil {
		t.Errorf("Unexpected error while reading value: %v", err)
	}
	if value != nil {
		t.Errorf("Expected nil value for non-existing key, but got: %v", value)
	}
	existingKey := "existing_key"
	existingValue := "Hello, World!"
	memDB.Set(&memcache.Item{Key: existingKey, Value: []byte(existingValue)})
	value, err = Read(existingKey)
	if err != nil {
		t.Errorf("Unexpected error while reading value: %v", err)
	}
	if string(value.([]byte)) != existingValue {
		t.Errorf("Expected value: %v, but got: %v", existingValue, string(value.([]byte)))
	}
}

func TestReadString(t *testing.T) {
	err := Init(config_path)
	if err != nil {
		t.Fatalf("Error during initialization: %v", err)
	}
	key := "nonexistent_key"
	value, err := ReadString(key)
	if err != nil {
		t.Errorf("Unexpected error while reading value: %v", err)
	}
	if value != "" {
		t.Errorf("Expected empty string for non-existing key, but got: %v", value)
	}
	existingKey := "existing_key"
	existingValue := "Hello, World!"
	memDB.Set(&memcache.Item{Key: existingKey, Value: []byte(existingValue)})
	value, err = ReadString(existingKey)
	if err != nil {
		t.Errorf("Unexpected error while reading value: %v", err)
	}
	if value != existingValue {
		t.Errorf("Expected value: %v, but got: %v", existingValue, value)
	}
}

func compareJSON(a, b map[string]interface{}) bool {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}
