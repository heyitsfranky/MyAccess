package MyAccess

import (
	"encoding/json"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/heyitsfranky/MyConfig"
)

var data *InitData
var memDB *memcache.Client

type InitData struct {
	AccessMemcacheAddress string `yaml:"access_memcache_address"`
}

func Init(configPath string) error {
	if data == nil {
		err := MyConfig.Init(configPath, &data)
		if err != nil {
			return err
		}
		memDB = memcache.New(data.AccessMemcacheAddress)
		testKey := "test_key_for_connection_check"
		_, err = memDB.Get(testKey)
		if err != nil && err != memcache.ErrCacheMiss {
			return fmt.Errorf("failed to connect to memcache: %v", err)
		}
	}
	return nil
}

func ReadJSON(key string) (map[string]interface{}, error) {
	if data == nil {
		return nil, fmt.Errorf("must first call Init() to initialize the memcache settings")
	}
	item, err := memDB.Get(key)
	if err != nil {
		if err != memcache.ErrCacheMiss {
			return nil, err
		}
		return nil, nil
	}
	value := make(map[string]interface{})
	err = json.Unmarshal(item.Value, &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func Read(key string) (interface{}, error) {
	if data == nil {
		return nil, fmt.Errorf("must first call Init() to initialize the memcache settings")
	}
	item, err := memDB.Get(key)
	if err != nil {
		if err != memcache.ErrCacheMiss {
			return nil, err
		}
		return nil, nil
	}
	return item.Value, nil
}

func ReadString(key string) (string, error) {
	if data == nil {
		return "", fmt.Errorf("must first call Init() to initialize the memcache settings")
	}
	item, err := memDB.Get(key)
	if err != nil {
		if err != memcache.ErrCacheMiss {
			return "", err
		}
		return "", nil
	}
	return string(item.Value), nil
}
