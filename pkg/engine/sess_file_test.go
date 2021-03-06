package engine

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

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

const sid = "Session_id"
const sidNew = "Session_id_new"
const sessionPath = "./_session_runtime"

var (
	mutex sync.Mutex
)

func TestFileProvider_SessionInit(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(context.Background(), 180, sessionPath)
	if fp.maxlifetime != 180 {
		t.Error()
	}

	if fp.savePath != sessionPath {
		t.Error()
	}
}

func TestFileProvider_SessionExist(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(context.Background(), 180, sessionPath)

	exists, err := fp.SessionExist(context.Background(), sid)
	if err != nil {
		t.Error(err)
	}
	if exists {
		t.Error()
	}

	_, err = fp.SessionRead(context.Background(), sid)
	if err != nil {
		t.Error(err)
	}

	exists, err = fp.SessionExist(context.Background(), sid)
	if err != nil {
		t.Error(err)
	}
	if !exists {
		t.Error()
	}
}

func TestFileProvider_SessionExist2(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(context.Background(), 180, sessionPath)

	exists, err := fp.SessionExist(context.Background(), sid)
	if err != nil {
		t.Error(err)
	}
	if exists {
		t.Error()
	}

	exists, err = fp.SessionExist(context.Background(), "")
	if err == nil {
		t.Error()
	}
	if exists {
		t.Error()
	}

	exists, err = fp.SessionExist(context.Background(), "1")
	if err == nil {
		t.Error()
	}
	if exists {
		t.Error()
	}
}

func TestFileProvider_SessionRead(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(context.Background(), 180, sessionPath)

	s, err := fp.SessionRead(context.Background(), sid)
	if err != nil {
		t.Error(err)
	}

	_ = s.Set(nil, "sessionValue", 18975)
	v := s.Get(nil, "sessionValue")

	if v.(int) != 18975 {
		t.Error()
	}
}

func TestFileProvider_SessionRead1(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(context.Background(), 180, sessionPath)

	_, err := fp.SessionRead(context.Background(), "")
	if err == nil {
		t.Error(err)
	}

	_, err = fp.SessionRead(context.Background(), "1")
	if err == nil {
		t.Error(err)
	}
}

func TestFileProvider_SessionAll(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(context.Background(), 180, sessionPath)

	sessionCount := 546

	for i := 1; i <= sessionCount; i++ {
		_, err := fp.SessionRead(context.Background(), fmt.Sprintf("%s_%d", sid, i))
		if err != nil {
			t.Error(err)
		}
	}

	if fp.SessionAll(nil) != sessionCount {
		t.Error()
	}
}

func TestFileProvider_SessionRegenerate(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(context.Background(), 180, sessionPath)

	_, err := fp.SessionRead(context.Background(), sid)
	if err != nil {
		t.Error(err)
	}

	exists, err := fp.SessionExist(context.Background(), sid)
	if err != nil {
		t.Error(err)
	}
	if !exists {
		t.Error()
	}

	_, err = fp.SessionRegenerate(context.Background(), sid, sidNew)
	if err != nil {
		t.Error(err)
	}

	exists, err = fp.SessionExist(context.Background(), sid)
	if err != nil {
		t.Error(err)
	}
	if exists {
		t.Error()
	}

	exists, err = fp.SessionExist(context.Background(), sidNew)
	if err != nil {
		t.Error(err)
	}
	if !exists {
		t.Error()
	}
}

func TestFileProvider_SessionDestroy(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(context.Background(), 180, sessionPath)

	_, err := fp.SessionRead(context.Background(), sid)
	if err != nil {
		t.Error(err)
	}

	exists, err := fp.SessionExist(context.Background(), sid)
	if err != nil {
		t.Error(err)
	}
	if !exists {
		t.Error()
	}

	err = fp.SessionDestroy(context.Background(), sid)
	if err != nil {
		t.Error(err)
	}

	exists, err = fp.SessionExist(context.Background(), sid)
	if err != nil {
		t.Error(err)
	}
	if exists {
		t.Error()
	}
}

func TestFileProvider_SessionGC(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(context.Background(), 1, sessionPath)

	sessionCount := 412

	for i := 1; i <= sessionCount; i++ {
		_, err := fp.SessionRead(context.Background(), fmt.Sprintf("%s_%d", sid, i))
		if err != nil {
			t.Error(err)
		}
	}

	time.Sleep(2 * time.Second)

	fp.SessionGC(nil)
	if fp.SessionAll(nil) != 0 {
		t.Error()
	}
}

func TestFileSessionStore_Set(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(context.Background(), 180, sessionPath)

	sessionCount := 100
	s, _ := fp.SessionRead(context.Background(), sid)
	for i := 1; i <= sessionCount; i++ {
		err := s.Set(nil, i, i)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestFileSessionStore_Get(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(context.Background(), 180, sessionPath)

	sessionCount := 100
	s, _ := fp.SessionRead(context.Background(), sid)
	for i := 1; i <= sessionCount; i++ {
		_ = s.Set(nil, i, i)

		v := s.Get(nil, i)
		if v.(int) != i {
			t.Error()
		}
	}
}

func TestFileSessionStore_Delete(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(context.Background(), 180, sessionPath)

	s, _ := fp.SessionRead(context.Background(), sid)
	s.Set(nil, "1", 1)

	if s.Get(nil, "1") == nil {
		t.Error()
	}

	s.Delete(nil, "1")

	if s.Get(nil, "1") != nil {
		t.Error()
	}
}

func TestFileSessionStore_Flush(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(context.Background(), 180, sessionPath)

	sessionCount := 100
	s, _ := fp.SessionRead(context.Background(), sid)
	for i := 1; i <= sessionCount; i++ {
		_ = s.Set(nil, i, i)
	}

	_ = s.Flush(nil)

	for i := 1; i <= sessionCount; i++ {
		if s.Get(nil, i) != nil {
			t.Error()
		}
	}
}

func TestFileSessionStore_SessionID(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(context.Background(), 180, sessionPath)

	sessionCount := 85

	for i := 1; i <= sessionCount; i++ {
		s, err := fp.SessionRead(context.Background(), fmt.Sprintf("%s_%d", sid, i))
		if err != nil {
			t.Error(err)
		}
		if s.SessionID(nil) != fmt.Sprintf("%s_%d", sid, i) {
			t.Error(err)
		}
	}
}

func TestFileSessionStore_SessionRelease(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(context.Background(), 180, sessionPath)
	filepder.savePath = sessionPath
	sessionCount := 85

	for i := 1; i <= sessionCount; i++ {
		s, err := fp.SessionRead(context.Background(), fmt.Sprintf("%s_%d", sid, i))
		if err != nil {
			t.Error(err)
		}

		s.Set(nil, i, i)
		s.SessionRelease(nil, nil)
	}

	for i := 1; i <= sessionCount; i++ {
		s, err := fp.SessionRead(context.Background(), fmt.Sprintf("%s_%d", sid, i))
		if err != nil {
			t.Error(err)
		}

		if s.Get(nil, i).(int) != i {
			t.Error()
		}
	}
}
