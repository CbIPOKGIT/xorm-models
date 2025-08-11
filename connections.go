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

// ConnectionCredentials - структура для передачі даних для підключення.
//
// Якщо ми не передаємо credentials - використовуємо дані з env файла по наступному ключу:
//
//	логін / env - {connectionName}_SQL_LOGIN,
//	пароль / env - {connectionName}_SQL_PASSWORD,
//	host / env - {connectionName}_SQL_HOST,
//	імя БД / env - {connectionName}_DATABASE
type ConnectionCredentials struct {
	Login    string
	Password string
	Host     string
	Database string
}

var connections Connections = Connections{}

func init() {
	connections.ConnectionsList = make(ConnectionsList)
}

// GetConnection - get connection by name
//
// Створюємо підключення або повертаємо вже існуюче підключення
func GetConnection(connectionName string, credentials ...*ConnectionCredentials) (*xorm.Engine, error) {
	connnection := getExistingConnection(connectionName)
	if connnection != nil {
		if err := connnection.Ping(); err == nil {
			return connnection, nil
		}
	}

	return createConnection(connectionName, credentials...)
}

// getExistingConnection - якщо підключення вже існує - повертаємо його
func getExistingConnection(connectionName string) *xorm.Engine {
	connections.RLock()
	defer connections.RUnlock()

	if connection, has := connections.ConnectionsList[connectionName]; has {
		return connection
	}

	return nil
}

// createConnection - створюємо нове підключення та додаємо його в список
func createConnection(connectionName string, credentials ...*ConnectionCredentials) (*xorm.Engine, error) {
	connections.Lock()
	defer connections.Unlock()

	var connectStr string

	if len(credentials) > 0 {
		for _, credential := range credentials {
			connectStr = fmt.Sprintf(
				"%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True",
				credential.Login,
				credential.Password,
				credential.Host,
				credential.Database,
			)
		}

	} else {
		connectStr = fmt.Sprintf(
			"%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True",
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

	connections.ConnectionsList[connectionName] = engine
	return engine, nil
}
