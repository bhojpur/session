package redis_sentinel

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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	session "github.com/bhojpur/session/pkg/engine"
)

func TestRedisSentinel(t *testing.T) {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "bsessionid",
		EnableSetCookie: true,
		Gclifetime:      3600,
		Maxlifetime:     3600,
		Secure:          false,
		CookieLifeTime:  3600,
		ProviderConfig:  "127.0.0.1:6379,100,,0,master",
	}
	globalSessions, e := session.NewManager("redis_sentinel", sessionConfig)
	if e != nil {
		t.Log(e)
		return
	}
	// todo test if e==nil
	go globalSessions.GC()

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		t.Fatal("session start failed:", err)
	}
	defer sess.SessionRelease(nil, w)

	// SET AND GET
	err = sess.Set(nil, "username", "bhojpur")
	if err != nil {
		t.Fatal("set username failed:", err)
	}
	username := sess.Get(nil, "username")
	if username != "bhojpur" {
		t.Fatal("get username failed")
	}

	// DELETE
	err = sess.Delete(nil, "username")
	if err != nil {
		t.Fatal("delete username failed:", err)
	}
	username = sess.Get(nil, "username")
	if username != nil {
		t.Fatal("delete username failed")
	}

	// FLUSH
	err = sess.Set(nil, "username", "bhojpur")
	if err != nil {
		t.Fatal("set failed:", err)
	}
	err = sess.Set(nil, "password", "1qaz2wsx")
	if err != nil {
		t.Fatal("set failed:", err)
	}
	username = sess.Get(nil, "username")
	if username != "bhojpur" {
		t.Fatal("get username failed")
	}
	password := sess.Get(nil, "password")
	if password != "1qaz2wsx" {
		t.Fatal("get password failed")
	}
	err = sess.Flush(nil)
	if err != nil {
		t.Fatal("flush failed:", err)
	}
	username = sess.Get(nil, "username")
	if username != nil {
		t.Fatal("flush failed")
	}
	password = sess.Get(nil, "password")
	if password != nil {
		t.Fatal("flush failed")
	}

	sess.SessionRelease(nil, w)

}

func TestProvider_SessionInit(t *testing.T) {

	savePath := `
{ "save_path": "my save path", "idle_timeout": "3s"}
`
	cp := &Provider{}
	cp.SessionInit(context.Background(), 12, savePath)
	assert.Equal(t, "my save path", cp.SavePath)
	assert.Equal(t, 3*time.Second, cp.idleTimeout)
	assert.Equal(t, int64(12), cp.maxlifetime)
}
