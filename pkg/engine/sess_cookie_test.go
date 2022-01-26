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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCookie(t *testing.T) {
	config := `{"cookieName":"bsessionid","enableSetCookie":false,"gclifetime":3600,"ProviderConfig":"{\"cookieName\":\"bsessionid\",\"securityKey\":\"bhojpurcookiehashkey\"}"}`
	conf := new(ManagerConfig)
	if err := json.Unmarshal([]byte(config), conf); err != nil {
		t.Fatal("json decode error", err)
	}
	globalSessions, err := NewManager("cookie", conf)
	if err != nil {
		t.Fatal("init cookie session err", err)
	}
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		t.Fatal("set error,", err)
	}
	err = sess.Set(nil, "username", "bhojpur")
	if err != nil {
		t.Fatal("set error,", err)
	}
	if username := sess.Get(nil, "username"); username != "bhojpur" {
		t.Fatal("get username error")
	}
	sess.SessionRelease(nil, w)
	if cookiestr := w.Header().Get("Set-Cookie"); cookiestr == "" {
		t.Fatal("setcookie error")
	} else {
		parts := strings.Split(strings.TrimSpace(cookiestr), ";")
		for k, v := range parts {
			nameval := strings.Split(v, "=")
			if k == 0 && nameval[0] != "bsessionid" {
				t.Fatal("error")
			}
		}
	}
}

func TestDestorySessionCookie(t *testing.T) {
	config := `{"cookieName":"bsessionid","enableSetCookie":true,"gclifetime":3600,"ProviderConfig":"{\"cookieName\":\"bsessionid\",\"securityKey\":\"bhojpurcookiehashkey\"}"}`
	conf := new(ManagerConfig)
	if err := json.Unmarshal([]byte(config), conf); err != nil {
		t.Fatal("json decode error", err)
	}
	globalSessions, err := NewManager("cookie", conf)
	if err != nil {
		t.Fatal("init cookie session err", err)
	}

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	session, err := globalSessions.SessionStart(w, r)
	if err != nil {
		t.Fatal("session start err,", err)
	}

	// request again ,will get same sesssion id .
	r1, _ := http.NewRequest("GET", "/", nil)
	r1.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()
	newSession, err := globalSessions.SessionStart(w, r1)
	if err != nil {
		t.Fatal("session start err,", err)
	}
	if newSession.SessionID(nil) != session.SessionID(nil) {
		t.Fatal("get cookie session id is not the same again.")
	}

	// After destroy session , will get a new session id .
	globalSessions.SessionDestroy(w, r1)
	r2, _ := http.NewRequest("GET", "/", nil)
	r2.Header.Set("Cookie", w.Header().Get("Set-Cookie"))

	w = httptest.NewRecorder()
	newSession, err = globalSessions.SessionStart(w, r2)
	if err != nil {
		t.Fatal("session start error")
	}
	if newSession.SessionID(nil) == session.SessionID(nil) {
		t.Fatal("after destroy session and reqeust again ,get cookie session id is same.")
	}
}
