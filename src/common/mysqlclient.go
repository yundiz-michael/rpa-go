package common

import (
	"bytes"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type MysqlClient struct {
	Env *YamlFile
	Db  *sqlx.DB
}

func (m *MysqlClient) Init() {
	_db := Consul.ReadDbConfig(m.Env.Database.Mysql)
	var dsn bytes.Buffer
	dsn.WriteString(_db.Account.User)
	dsn.WriteString(":")
	dsn.WriteString(_db.Account.Password)
	dsn.WriteString("@tcp(")
	dsn.WriteString(_db.Hosts[0])
	dsn.WriteString(")/")
	dsn.WriteString(_db.ID)
	dsn.WriteString("?charset=utf8mb4&parseTime=True")
	connectUri := dsn.String()

	var err error
	m.Db, err = sqlx.Connect("mysql", connectUri)
	if err != nil {
		LoggerStd.Panic(connectUri, zap.NamedError("error", err))
	}
	m.Db.SetMaxOpenConns(10)
	m.Db.SetMaxIdleConns(1)
	m.Db.SetConnMaxLifetime(time.Minute * 3)
}

func (m *MysqlClient) Select(dest interface{}, query string, args ...interface{}) error {
	return m.Db.Select(dest, query, args...)
}

func (m *MysqlClient) Get(dest interface{}, query string, args ...interface{}) error {
	return m.Db.Get(dest, query, args...)
}

func (m *MysqlClient) Update(query string, args ...interface{}) (sql.Result, error) {
	return m.Db.Exec(query, args...)
}
