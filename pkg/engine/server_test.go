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
	"crypto/aes"
	"encoding/json"
	"testing"
)

func Test_gob(t *testing.T) {
	a := make(map[interface{}]interface{})
	a["username"] = "bhojpur"
	a[12] = 234
	a["user"] = User{"bhoj", "pur"}
	b, err := EncodeGob(a)
	if err != nil {
		t.Error(err)
	}
	c, err := DecodeGob(b)
	if err != nil {
		t.Error(err)
	}
	if len(c) == 0 {
		t.Error("decodeGob empty")
	}
	if c["username"] != "bhojpur" {
		t.Error("decode string error")
	}
	if c[12] != 234 {
		t.Error("decode int error")
	}
	if c["user"].(User).Username != "bhoj" {
		t.Error("decode struct error")
	}
}

type User struct {
	Username string
	NickName string
}

func TestGenerate(t *testing.T) {
	str := generateRandomKey(20)
	if len(str) != 20 {
		t.Fatal("generate length is not equal to 20")
	}
}

func TestCookieEncodeDecode(t *testing.T) {
	hashKey := "testhashKey"
	blockkey := generateRandomKey(16)
	block, err := aes.NewCipher(blockkey)
	if err != nil {
		t.Fatal("NewCipher:", err)
	}
	securityName := string(generateRandomKey(20))
	val := make(map[interface{}]interface{})
	val["name"] = "bhojpur"
	val["gender"] = "male"
	str, err := encodeCookie(block, hashKey, securityName, val)
	if err != nil {
		t.Fatal("encodeCookie:", err)
	}
	dst, err := decodeCookie(block, hashKey, securityName, str, 3600)
	if err != nil {
		t.Fatal("decodeCookie", err)
	}
	if dst["name"] != "bhojpur" {
		t.Fatal("dst get map error")
	}
	if dst["gender"] != "male" {
		t.Fatal("dst get map error")
	}
}

func TestParseConfig(t *testing.T) {
	s := `{"cookieName":"bsessionid","gclifetime":3600}`
	cf := new(ManagerConfig)
	cf.EnableSetCookie = true
	err := json.Unmarshal([]byte(s), cf)
	if err != nil {
		t.Fatal("parse json error,", err)
	}
	if cf.CookieName != "bsessionid" {
		t.Fatal("parseconfig get cookiename error")
	}
	if cf.Gclifetime != 3600 {
		t.Fatal("parseconfig get gclifetime error")
	}

	cc := `{"cookieName":"bsessionid","enableSetCookie":false,"gclifetime":3600,"ProviderConfig":"{\"cookieName\":\"bsessionid\",\"securityKey\":\"bhojpurcookiehashkey\"}"}`
	cf2 := new(ManagerConfig)
	cf2.EnableSetCookie = true
	err = json.Unmarshal([]byte(cc), cf2)
	if err != nil {
		t.Fatal("parse json error,", err)
	}
	if cf2.CookieName != "bsessionid" {
		t.Fatal("parseconfig get cookiename error")
	}
	if cf2.Gclifetime != 3600 {
		t.Fatal("parseconfig get gclifetime error")
	}
	if cf2.EnableSetCookie {
		t.Fatal("parseconfig get enableSetCookie error")
	}
	cconfig := new(cookieConfig)
	err = json.Unmarshal([]byte(cf2.ProviderConfig), cconfig)
	if err != nil {
		t.Fatal("parse ProviderConfig err,", err)
	}
	if cconfig.CookieName != "bsessionid" {
		t.Fatal("ProviderConfig get cookieName error")
	}
	if cconfig.SecurityKey != "bhojpurcookiehashkey" {
		t.Fatal("ProviderConfig get securityKey error")
	}
}
