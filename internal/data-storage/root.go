/*
 * Created on Mon Dec 16 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package datastorage

import "errors"

type dataStorageRoot map[string]interface{}

func (d dataStorageRoot) Get(key string) (interface{}, error) {
	if _, ok := d[key]; !ok {
		return nil, errors.New("key not found")
	}
	return d[key], nil
}

func (d dataStorageRoot) Set(key string, value interface{}) {
	d[key] = value
}

var storage *dataStorageRoot

func init() {
	storage = new(dataStorageRoot)
}

func GetStorage() *dataStorageRoot {
	return storage
}
