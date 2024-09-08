package xormmodels

import (
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type ConnectionsList map[string]*xorm.Engine

type Connections struct {
	sync.RWMutex
	ConnectionsList
}

var connections Connections = Connections{}

func init() {
	connections.ConnectionsList = make(ConnectionsList)
}

// Создаем подключение или возвращаем error в случае ошибки
// Если мы не передаем credentials - используем данные с env файла.
// 0 - логин / env - {connectionName}_SQL_LOGIN,
// 1 - пароль / env - {connectionName}_SQL_PASSWORD,
// 2 - host / env - {connectionName}_SQL_HOST,
// 3 - имя БД / env - {connectionName}_DATABASE
//
// env variable DONT_PING_XORM_CONNECTION - 1 or 0. Ping existed connection or not
func GetConnection(connectionName string, credentials ...[]string) (*xorm.Engine, error) {
	connnection := getExistingConnection(connectionName)
	if connnection != nil {
		if err := connnection.Ping(); err == nil {
			return connnection, nil
		}
	}

	connections.Lock()
	defer connections.Unlock()

	if connection, err := createConnection(connectionName, credentials...); err == nil {
		connections.ConnectionsList[connectionName] = connection
		return connection, nil
	} else {
		return nil, err
	}
}

func getExistingConnection(connectionName string) *xorm.Engine {
	connections.RLock()
	defer connections.RUnlock()

	if connection, has := connections.ConnectionsList[connectionName]; has {
		return connection
	}

	return nil
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
