package myapp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestFooHandler_WithoutJson(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/foo", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code)
}

func TestFooHandler_WithJson(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(`{"first_name":"jjm", "last_name":"kim", "email":"jameo.cho@hanmail.net"}`))

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)

	user := new(User)
	err := json.NewDecoder(res.Body).Decode(user)

	assert.Nil(err)
	assert.Equal("jjm", user.FirstName)
	assert.Equal("kim", user.LastName)

}
