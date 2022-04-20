package xormmodels

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

// Создаем подключение или возвращаем error в случае ошибки
// Если мы не передаем credentials - используем данные с env файла.
// 0 - логин / env - {connectionName}_SQL_LOGIN,
// 1 - пароль / env - {connectionName}_SQL_PASSWORD,
// 2 - host / env - {connectionName}_SQL_HOST,
// 3 - имя БД / env - {connectionName}_DATABASE
func GetConnection(connectionName string, credentials ...[]string) (*xorm.Engine, error) {

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

	engine.ShowSQL(false)

	return engine, nil
}
