package myapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL)

	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	data, _ := ioutil.ReadAll(resp.Body)

	assert.Equal("Hello World", string(data))
}

func TestUsers(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users")

	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	data, _ := ioutil.ReadAll(resp.Body)

	assert.Contains(string(data), "No Users")

}

func TestUsers_WithUserData(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "applications/json", strings.NewReader((`{"first_name":"jaemo", "last_name":"cho", "email":"ddd@hanmail.net"}`)))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	resp, err = http.Post(ts.URL+"/users", "applications/json", strings.NewReader((`{"first_name":"jaemo1", "last_name":"cho1", "email":"ddd@hanmail.net"}`)))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	resp, err = http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	// data, err := ioutil.ReadAll(resp.Body)
	// assert.NoError(err)
	// assert.NotZero(len(data))

	users := []*User{}
	err = json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(err)
	assert.Equal(2, len(users))

}

func TestGetUserInfo(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users/54")

	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User Id:54")

}

func TestCreateUserInfo(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "applications/json", strings.NewReader((`{"first_name":"jaemo", "last_name":"cho", "email":"ddd@hanmail.net"}`)))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID) // User가 처음 생성 될 때 default 값이 0이라 생성된 값이 반환되서 오면id가 0이면 안된다.

	id := user.ID
	resp, err = http.Get(ts.URL + "/users/" + strconv.Itoa(id))
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	user2 := new(User)
	err = json.NewDecoder(resp.Body).Decode(user2)
	assert.NoError(err)
	assert.Equal(user.ID, user2.ID)
	assert.Equal(user.FirstName, user2.FirstName)
}

func TestDeleteUserInfo(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// DELETE PUT 은 제공하지 않아 NewRequest 로 선언해서 사용
	// http.Get, http.Post는
	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err := http.DefaultClient.Do(req)

	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	data, _ := ioutil.ReadAll(resp.Body)
	// log.Print(string(data))
	assert.Contains(string(data), "No User ID:1")

	// 등록
	resp, err = http.Post(ts.URL+"/users", "applications/json", strings.NewReader((`{"first_name":"jaemo", "last_name":"cho", "email":"ddd@hanmail.net"}`)))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	// 확인
	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	// 삭제
	req, _ = http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err = http.DefaultClient.Do(req)

	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	data, _ = ioutil.ReadAll(resp.Body)
	// log.Print(string(data))
	assert.Contains(string(data), "Deleted User ID:1")

}

func TestUpdateUserInfo(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// DELETE PUT 은 제공하지 않아 NewRequest 로 선언해서 사용
	// http.Get, http.Post는
	req, _ := http.NewRequest("PUT", ts.URL+"/users", strings.NewReader((`{"id":1, "first_name":"jaemo1", "last_name":"cho1", "email":"ddd1@hanmail.net"}`)))
	resp, err := http.DefaultClient.Do(req)

	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User ID:1")

	// 생성
	resp, err = http.Post(ts.URL+"/users", "applications/json", strings.NewReader((`{"first_name":"jaemo", "last_name":"cho", "email":"ddd@hanmail.net"}`)))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	updateStr := fmt.Sprintf(`{"id":%d, "first_name":"jaemo1", "last_name":"cho1", "email":"ddd1@hanmail.net"}`, user.ID)

	req, _ = http.NewRequest("PUT", ts.URL+"/users", strings.NewReader(updateStr))
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	updateUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(updateUser)
	assert.NoError(err)
	assert.Equal(updateUser.ID, user.ID)
	assert.Equal("jaemo1", updateUser.FirstName)
	assert.Equal("cho1", updateUser.LastName)
	assert.Equal("ddd1@hanmail.net", updateUser.Email)

}
