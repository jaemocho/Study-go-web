package app

import (
	"encoding/json"
	"fmt"
	"goweb/web13/model"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTodos(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(MakeHandler())
	defer ts.Close()

	// postform test 신규 데이터 생성
	resp, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"test todo"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	//todo 객체로 decode 안에 값을 채워준다.
	var todo model.Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "test todo")

	id := todo.ID

	// postform test 신규 데이터 생성
	resp, err = http.PostForm(ts.URL+"/todos", url.Values{"name": {"test todo2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	//todo 객체로 decode 안에 값을 채워준다.
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "test todo2")

	id2 := todo.ID

	// 전체 데이터 get test
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos := []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)

	for _, t := range todos {
		if t.ID == id {
			assert.Equal("test todo", t.Name)
		} else if t.ID == id2 {
			assert.Equal("test todo2", t.Name)
		} else {
			assert.Error(fmt.Errorf("testID should be id1 or id2"))
		}
	}

	// complete-todo test
	resp, err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id2) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	success := Success{}
	err = json.NewDecoder(resp.Body).Decode(&success)
	assert.NoError(err)
	assert.Equal(true, success.Success)

	// 전체 데이터 다시 get
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)

	for _, t := range todos {
		if t.ID == id2 {
			assert.True(t.Completed)
		}
	}

	// 삭제 테스트
	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(id), nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// 전체 데이터 다시 get
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 1)

	for _, t := range todos {
		assert.Equal(t.ID, id2)
	}

}
