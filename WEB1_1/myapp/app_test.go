package myapp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexPathHandler(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	indexHandler(res, req)

	assert.Equal(http.StatusOK, res.Code)

	// if res.Code != http.StatusOK {
	// 	t.Fatal("Failed!! ", res.Code)
	// }

	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello World!", string(data))
}

func TestBarPathHandler_WithoutName(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar", nil)

	barHandler(res, req)

	assert.Equal(http.StatusOK, res.Code)

	// if res.Code != http.StatusOK {
	// 	t.Fatal("Failed!! ", res.Code)
	// }

	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("hello World!!", string(data))
}

func TestBarPathHandler_WithoutName2(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar", nil)

	barHandler(res, req)
	assert.Equal(http.StatusOK, res.Code)

	// if res.Code != http.StatusOK {
	// 	t.Fatal("Failed!! ", res.Code)
	// }

	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("hello World!!", string(data))
}

func TestBarPathHandler_WithName(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar?name=jjm", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)

	// if res.Code != http.StatusOK {
	// 	t.Fatal("Failed!! ", res.Code)
	// }

	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("hello jjm!", string(data))
}
