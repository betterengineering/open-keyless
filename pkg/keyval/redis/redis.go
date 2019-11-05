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

// Package redis provides an implementation of the keyval interface for redis.
package redis

import "github.com/mediocregopher/radix/v3"

// RedisKeyVal provides a redis implementation for the keyval interface.
type RedisKeyVal struct {
	client *radix.Pool
}

// NewRedisKeyVal provides an instantiated RedisKeyVal with the provided configuration.
func NewRedisKeyVal(address string, poolSize int) (*RedisKeyVal, error) {
	pool, err := radix.NewPool("tcp", address, poolSize)
	if err != nil {
		return nil, err
	}

	return &RedisKeyVal{
		client: pool,
	}, nil
}

// Get returns a value for a given key from the store.
func (kv *RedisKeyVal) Get(key string) ([]byte, error) {
	var value []byte
	err := kv.client.Do(radix.Cmd(&value, "GET", key))
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Set adds the given key to the store with the given value.
func (kv *RedisKeyVal) Set(key string, value []byte) error {
	return kv.client.Do(radix.Cmd(nil, "SET", key, string(value)))
}

// Delete removes the given key rom the store.
func (kv *RedisKeyVal) Delete(key string) error {
	return kv.client.Do(radix.Cmd(nil, "DEL", key))
}
