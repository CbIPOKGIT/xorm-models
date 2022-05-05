package xormmodels

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
)

// Primary key of each model
type PK struct {
	Id uint64 `xorm:"pk autoincr" json:"id"`
}

func (p PK) GetPKValue() uint64 {
	return p.Id
}

// Timestamps
type Timestamps struct {
	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

// Soft deletes
type SoftDeletes struct {
	DeletedAt time.Time `xorm:"deleted" json:"deleted_at"`
}

type ModelInterface interface {
	GetConnection() (*xorm.Engine, error)
	TableName() string
}

type ModelInterfaceWithPK interface {
	ModelInterface
	GetPKValue() uint64
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

func SaveModel(m ModelInterfaceWithPK) error {
	con, err := m.GetConnection()
	if err != nil {
		return err
	}
	// defer con.Close()

	id := m.GetPKValue()

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
	// defer con.Close()

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
	// defer con.Close()

	if f, err := con.ID(id).Get(m); err != nil {
		return err
	} else {
		if !f {
			return errors.New(fmt.Sprintf("no model with id %d", id))
		}
	}
	return nil
}

func Delete(m ModelInterfaceWithPK) error {
	id := m.GetPKValue()

	con, err := m.GetConnection()
	if err != nil {
		return err
	}
	// defer con.Close()

	_, err = con.ID(id).Delete(m)
	return err
}

func Insert(m ModelInterface) (int64, error) {
	con, err := m.GetConnection()
	if err != nil {
		return 0, err
	}
	// defer con.Close()

	return con.Insert(m)
}

func FindOne(m ModelInterface, query *QueryModel) (bool, error) {
	con, err := m.GetConnection()
	if err != nil {
		return false, err
	}

	session := con.NewSession()
	if query != nil {
		query.Fill(session)
	}

	return session.Get(m)
}

func FindAll(m ModelInterface, query *QueryModel) error {
	con, err := m.GetConnection()
	if err != nil {
		return err
	}

	session := con.NewSession()
	if query != nil {
		query.Fill(session)
	}

	return session.Find(m)
}
