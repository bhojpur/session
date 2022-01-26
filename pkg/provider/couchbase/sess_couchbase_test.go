package couchbase

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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider_SessionInit(t *testing.T) {
	// using old style
	savePath := `http://host:port/,Pool,Bucket`
	cp := &Provider{}
	cp.SessionInit(context.Background(), 12, savePath)
	assert.Equal(t, "http://host:port/", cp.SavePath)
	assert.Equal(t, "Pool", cp.Pool)
	assert.Equal(t, "Bucket", cp.Bucket)
	assert.Equal(t, int64(12), cp.maxlifetime)

	savePath = `
{ "save_path": "my save path", "pool": "mypool", "bucket": "mybucket"}
`
	cp = &Provider{}
	cp.SessionInit(context.Background(), 12, savePath)
	assert.Equal(t, "my save path", cp.SavePath)
	assert.Equal(t, "mypool", cp.Pool)
	assert.Equal(t, "mybucket", cp.Bucket)
	assert.Equal(t, int64(12), cp.maxlifetime)
}
