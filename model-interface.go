package xormmodels

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-xorm/xorm"
)

type ModelInterface interface {
	GetConnection() (*xorm.Engine, error)
	TableName() string
}

func getPKValue(m ModelInterface) uint64 {
	ref := reflect.ValueOf(m).Elem()

	for i := 0; i < ref.NumField(); i++ {
		field := ref.Type().Field(i)
		val := ref.Field(i).Interface()
		if strings.ToLower(field.Name) == "id" {
			return val.(uint64)
		}
	}
	return 0
}

func SaveModel(m ModelInterface) error {
	con, err := m.GetConnection()
	if err != nil {
		return err
	}
	defer con.Close()

	id := getPKValue(m)

	var errDB error
	if id != 0 {
		_, errDB = con.Id(id).UseBool().AllCols().Update(m)
	} else {
		_, errDB = con.InsertOne(m)
	}

	if errDB != nil {
		return errDB
	}
	return nil
}

func All(m ModelInterface) error {
	con, err := m.GetConnection()
	if err != nil {
		return err
	}
	defer con.Close()

	if err := con.Find(m); err != nil {
		return err
	}
	return nil
}

func Find(m ModelInterface, id interface{}) error {
	con, err := m.GetConnection()
	if err != nil {
		return err
	}
	defer con.Close()

	if f, err := con.ID(id).Get(m); err != nil {
		return err
	} else {
		if !f {
			return errors.New(fmt.Sprintf("no model with id %d", id))
		}
	}
	return nil
}

func Delete(m ModelInterface) error {
	id := getPKValue(m)
	if id == 0 {
		return errors.New("no id field")
	}

	con, err := m.GetConnection()
	if err != nil {
		return err
	}
	defer con.Close()

	_, err = con.ID(id).Delete(m)
	return err
}

func Insert(m ModelInterface) (int64, error) {
	con, err := m.GetConnection()
	if err != nil {
		return 0, err
	}
	defer con.Close()

	return con.Insert(m)
}

func FindOne(m ModelInterface, mapa map[string]interface{}) (bool, error) {
	con, err := m.GetConnection()
	if err != nil {
		return false, err
	}
	defer con.Close()

	query := con.NewSession()

	for key, val := range mapa {
		query.Where(fmt.Sprintf("%s=?", key), val)
	}

	return query.Get(m)
}
