package core

import (
	"encoding/json"

	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

type Repository interface {
	Replace()
}

type repository struct {
	session    *r.Session
	tableName  string
	actionFlag uint32
	actionKey  string
	message    interface{}
}

func (m *repository) Replace() {
	if m.actionFlag != 2 {
		var item map[string]interface{}
		j, _ := json.Marshal(m.message)
		json.Unmarshal(j, &item)
		r.Table(m.tableName).Get(m.actionKey).Replace(item).RunWrite(m.session)
	} else {
		r.Table(m.tableName).Get(m.actionKey).Replace(nil).RunWrite(m.session)
	}

}

func NewRepository(session *r.Session, tableName string, actionFlag uint32, actionKey string, message interface{}) Repository {
	return &repository{
		session,
		tableName,
		actionFlag,
		actionKey,
		message,
	}
}
