package xormmodels

import (
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type ConnectionsList map[string]*xorm.Engine

type Connections struct {
	sync.Mutex
	ConnectionsList
}

var connections Connections

// Создаем подключение или возвращаем error в случае ошибки
// Если мы не передаем credentials - используем данные с env файла.
// 0 - логин / env - {connectionName}_SQL_LOGIN,
// 1 - пароль / env - {connectionName}_SQL_PASSWORD,
// 2 - host / env - {connectionName}_SQL_HOST,
// 3 - имя БД / env - {connectionName}_DATABASE
func GetConnection(connectionName string, credentials ...[]string) (*xorm.Engine, error) {
	if connections.ConnectionsList == nil {
		connections.ConnectionsList = make(ConnectionsList)
	}

	connections.Lock()
	defer connections.Unlock()

	if connection, has := connections.ConnectionsList[connectionName]; has {
		if err := connection.Ping(); err == nil {
			return connection, nil
		}
	}

	if connection, err := createConnection(connectionName, credentials...); err == nil {
		connections.ConnectionsList[connectionName] = connection
		return connection, nil
	} else {
		return nil, err
	}
}

func createConnection(connectionName string, credentials ...[]string) (*xorm.Engine, error) {
	var connectStr string

	if len(credentials) == 4 {
		connectStr = fmt.Sprintf(
			"%s:%s@(%s)/%s?charset=utf8&parseTime=True",
			credentials[0],
			credentials[1],
			credentials[2],
			credentials[3],
		)
	} else {
		connectStr = fmt.Sprintf(
			"%s:%s@(%s)/%s?charset=utf8&parseTime=True",
			os.Getenv(connectionName+"_SQL_LOGIN"),
			os.Getenv(connectionName+"_SQL_PASSWORD"),
			os.Getenv(connectionName+"_SQL_HOST"),
			os.Getenv(connectionName+"_DATABASE"),
		)
	}

	engine, err := xorm.NewEngine("mysql", connectStr)

	if err != nil {
		return nil, err
	}

	return engine, nil
}
