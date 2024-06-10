package store

import (
	"encoding/json"
	"os"
	"sync"
)

type Disk struct{
	sm *sync.Map
}

func NewDisk(sm *sync.Map) *Disk {
	return &Disk{
		sm: sm,
	}
}

type DataDTO struct {
	Data map[string]Values
}


func (d *Disk) syncMapToMap(syncMap *sync.Map) map[string]Values {
	result := make(map[string]Values)
	syncMap.Range(func(key, value interface{}) bool {
		k, ok1 := key.(string)
		if !ok1 {
			return false
		}
		result[k] = value.(Values)
		return true
	})
	return result
}

func (d *Disk) Load(path string) (map[string]Values, error) {
	if !fileExists(path) {
		return make(map[string]Values), nil
	}
	file, err := os.Open(path)
	if err != nil {
		return make(map[string]Values) ,err
	}
	defer file.Close()

	var m DataDTO
	err = json.NewDecoder(file).Decode(&m)
	if err != nil {
		return make(map[string]Values), err
	}
	return m.Data, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (d *Disk) Save(path string) error {
	m := DataDTO{
		Data: d.syncMapToMap(d.sm),
	}
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}