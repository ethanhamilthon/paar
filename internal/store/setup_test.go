package store_test

import (
	"paar/internal/store"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	storage := store.NewStorage()
	infinite :=  time.Now().Add(time.Duration(1<<63 - 1))

	//Set key="key1" value="value1"
	storage.Store("key1", store.Values{Value: "value1", ExpireTo: infinite})

	//Get key="key1"
	value, ok := storage.Load("key1")
	if!ok {
		t.Errorf("key1 not found")
	}
	if value.Value!= "value1" {
		t.Errorf("value1 expected, but %s", value.Value)
	}
	if value.ExpireTo.Before(time.Now()) {
		t.Errorf("key1 expired")
	}

	//Delete key="key1"
	storage.Delete("key1")

	//Get key="key1"
	value, ok = storage.Load("key1")
	if ok {
		t.Errorf("key1 not deleted")
	}
	if value.Value != "" {
		t.Errorf("value1 expected, but %s", value.Value)
	}

	//Set 3 keys
	storage.Store("key2", store.Values{Value: "value2", ExpireTo: infinite})
	storage.Store("key3", store.Values{Value: "value3", ExpireTo: infinite})
	storage.Store("key4", store.Values{Value: "value4", ExpireTo: infinite})

	//Range
	keys := []string{}
	storage.Range(func(key string, value store.Values) bool {
		if value.ExpireTo.Before(time.Now()) {
			t.Errorf("key %s expired", key)
		}else{
			keys = append(keys, key)
		}
		return true
	})
	if len(keys)!= 3 {
		t.Errorf("keys expected 3, but %d", len(keys))
	}
	if keys[0]!= "key2" || keys[1]!= "key3" || keys[2]!= "key4" {
		t.Errorf("keys expected [key2, key3, key4], but %s", keys)
	}
}