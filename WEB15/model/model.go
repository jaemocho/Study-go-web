package model

import "time"

type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"create_at"`
}

type DBHandler interface {
	GetTodos(sessionId string) []*Todo
	AddTodo(name string, sessionId string) *Todo
	RemoveTodo(id int) bool
	CompleteTodo(id int, complte bool) bool
	Close()
}

func NewDBHandler(filepath string) DBHandler {
	//sqlite 사용
	return newSqliteHandler(filepath)

	// postgres 사용
	//return newPQHandler(dbConn)
}
