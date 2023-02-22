package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexPage(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("Hello World", string(data))
}

func TestDecoHandler(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// server의 log를 buffer에 찍도록 변경
	buf := &bytes.Buffer{}
	log.SetOutput(buf)

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	r := bufio.NewReader(buf)
	line, _, err := r.ReadLine()
	assert.NoError(err)
	assert.Contains(string(line), "[LOGGER1] Started")

	line, _, err = r.ReadLine()
	assert.NoError(err)
	assert.Contains(string(line), "[LOGGER1] Completed")
}

/*
test 코드만 먼저 작성하고 test 시

--- FAIL: TestDecoHandler (0.00s)
    c:\Users\조재모\go\src\goWeb\WEB7\main_test.go:44:
        	Error Trace:	c:\Users\조재모\go\src\goWeb\WEB7\main_test.go:44
        	Error:      	Received unexpected error:
        	            	EOF
        	Test:       	TestDecoHandler
    c:\Users\조재모\go\src\goWeb\WEB7\main_test.go:45:
        	Error Trace:	c:\Users\조재모\go\src\goWeb\WEB7\main_test.go:45
        	Error:      	"" does not contain "[LOGGER1] Started"
        	Test:       	TestDecoHandler
FAIL
FAIL	goweb/web7	0.244s
FAIL

decorator 2개 추가 했을 때

--- FAIL: TestDecoHandler (0.09s)
    c:\Users\조재모\go\src\goWeb\WEB7\main_test.go:45:
        	Error Trace:	c:\Users\조재모\go\src\goWeb\WEB7\main_test.go:45
        	Error:      	"2023/02/12 20:28:25 [LOGGER2] Started" does not contain "[LOGGER1] Started"
        	Test:       	TestDecoHandler
    c:\Users\조재모\go\src\goWeb\WEB7\main_test.go:49:
        	Error Trace:	c:\Users\조재모\go\src\goWeb\WEB7\main_test.go:49
        	Error:      	"2023/02/12 20:28:25 [LOGGER1] Started" does not contain "[LOGGER1] Completed"
        	Test:       	TestDecoHandler
FAIL
FAIL	goweb/web7	0.363s
FAIL

*/
