package libetcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	etcd "go.etcd.io/etcd/client"
)

// StoreJSON saves a struct in JSON form under a key in Etcd
func StoreJSON(ctx context.Context, etcdKV etcd.KeysAPI, key string,
	value interface{}) error {

	// Marshal JSON
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshalling value to JSON: %s",
			err.Error())
	}

	// Store in Etcd
	_, err = etcdKV.Set(ctx, key, string(b[:]), nil)
	if err != nil {
		return fmt.Errorf("error saving value to Etcd", err.Error())
	}

	return nil
}

// LoadJSON retrieves a key from Etcd and decodes the value as JSON into
// a struct
func LoadJSON(ctx context.Context, etcdKV etcd.KeysAPI, key string,
	result interface{}) error {

	// Load value from Etcd
	resp, err := etcdKV.Get(ctx, key, &etcd.GetOptions{Quorum: true})
	if err != nil {
		return fmt.Errorf("error retrieving value from Etcd: %s",
			err.Error())
	}

	if resp.Node == nil {
		return errors.New("node at key was nil")
	}

	// Unmarshal
	err = json.Unmarshal(byte[Node.Value], result)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON value: %s",
			err.Error())
	}

	return nil
}
