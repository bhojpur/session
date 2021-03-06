package memcache

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// depend on github.com/bradfitz/gomemcache/memcache
//
// go install github.com/bradfitz/gomemcache/memcache
//
// Usage:
// import(
//   _ "github.com/bhojpur/session/pkg/provider/memcache"
//   session "github.com/bhojpur/session/pkg/engine"
// )
//
//	func init() {
//		globalSessions, _ = session.NewManager("memcache", ``{"cookieName":"bsessionid","gclifetime":3600,"ProviderConfig":"127.0.0.1:11211"}``)
//		go globalSessions.GC()
//	}

import (
	"context"
	"net/http"
	"strings"
	"sync"

	session "github.com/bhojpur/session/pkg/engine"

	"github.com/bradfitz/gomemcache/memcache"
)

var mempder = &MemProvider{}
var client *memcache.Client

// SessionStore memcache session store
type SessionStore struct {
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxlifetime int64
}

// Set value in memcache session
func (rs *SessionStore) Set(ctx context.Context, key, value interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values[key] = value
	return nil
}

// Get value in memcache session
func (rs *SessionStore) Get(ctx context.Context, key interface{}) interface{} {
	rs.lock.RLock()
	defer rs.lock.RUnlock()
	if v, ok := rs.values[key]; ok {
		return v
	}
	return nil
}

// Delete value in memcache session
func (rs *SessionStore) Delete(ctx context.Context, key interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.values, key)
	return nil
}

// Flush clear all values in memcache session
func (rs *SessionStore) Flush(context.Context) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values = make(map[interface{}]interface{})
	return nil
}

// SessionID get memcache session id
func (rs *SessionStore) SessionID(context.Context) string {
	return rs.sid
}

// SessionRelease save session values to memcache
func (rs *SessionStore) SessionRelease(ctx context.Context, w http.ResponseWriter) {
	b, err := session.EncodeGob(rs.values)
	if err != nil {
		return
	}
	item := memcache.Item{Key: rs.sid, Value: b, Expiration: int32(rs.maxlifetime)}
	client.Set(&item)
}

// MemProvider memcache session provider
type MemProvider struct {
	maxlifetime int64
	conninfo    []string
	poolsize    int
	password    string
}

// SessionInit init memcache session
// savepath like
// e.g. 127.0.0.1:9090
func (rp *MemProvider) SessionInit(ctx context.Context, maxlifetime int64, savePath string) error {
	rp.maxlifetime = maxlifetime
	rp.conninfo = strings.Split(savePath, ";")
	client = memcache.New(rp.conninfo...)
	return nil
}

// SessionRead read memcache session by sid
func (rp *MemProvider) SessionRead(ctx context.Context, sid string) (session.Store, error) {
	if client == nil {
		if err := rp.connectInit(); err != nil {
			return nil, err
		}
	}
	item, err := client.Get(sid)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			rs := &SessionStore{sid: sid, values: make(map[interface{}]interface{}), maxlifetime: rp.maxlifetime}
			return rs, nil
		}
		return nil, err
	}
	var kv map[interface{}]interface{}
	if len(item.Value) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = session.DecodeGob(item.Value)
		if err != nil {
			return nil, err
		}
	}
	rs := &SessionStore{sid: sid, values: kv, maxlifetime: rp.maxlifetime}
	return rs, nil
}

// SessionExist check memcache session exist by sid
func (rp *MemProvider) SessionExist(ctx context.Context, sid string) (bool, error) {
	if client == nil {
		if err := rp.connectInit(); err != nil {
			return false, err
		}
	}
	if item, err := client.Get(sid); err != nil || len(item.Value) == 0 {
		return false, err
	}
	return true, nil
}

// SessionRegenerate generate new sid for memcache session
func (rp *MemProvider) SessionRegenerate(ctx context.Context, oldsid, sid string) (session.Store, error) {
	if client == nil {
		if err := rp.connectInit(); err != nil {
			return nil, err
		}
	}
	var contain []byte
	if item, err := client.Get(sid); err != nil || len(item.Value) == 0 {
		// oldsid doesn't exists, set the new sid directly
		// ignore error here, since if it return error
		// the existed value will be 0
		item.Key = sid
		item.Value = []byte("")
		item.Expiration = int32(rp.maxlifetime)
		client.Set(item)
	} else {
		client.Delete(oldsid)
		item.Key = sid
		item.Expiration = int32(rp.maxlifetime)
		client.Set(item)
		contain = item.Value
	}

	var kv map[interface{}]interface{}
	if len(contain) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		var err error
		kv, err = session.DecodeGob(contain)
		if err != nil {
			return nil, err
		}
	}

	rs := &SessionStore{sid: sid, values: kv, maxlifetime: rp.maxlifetime}
	return rs, nil
}

// SessionDestroy delete memcache session by id
func (rp *MemProvider) SessionDestroy(ctx context.Context, sid string) error {
	if client == nil {
		if err := rp.connectInit(); err != nil {
			return err
		}
	}

	return client.Delete(sid)
}

func (rp *MemProvider) connectInit() error {
	client = memcache.New(rp.conninfo...)
	return nil
}

// SessionGC Impelment method, no used.
func (rp *MemProvider) SessionGC(context.Context) {
}

// SessionAll return all activeSession
func (rp *MemProvider) SessionAll(context.Context) int {
	return 0
}

func init() {
	session.Register("memcache", mempder)
}
