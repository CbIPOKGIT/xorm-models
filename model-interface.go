package xormmodels

import (
	"fmt"
	"time"

	"xorm.io/xorm"
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
	DeletedAt *time.Time `xorm:"deleted" json:"deleted_at"`
}

type ModelInterface interface {
	GetConnection() (*xorm.Engine, error)
	TableName() string
}

type ModelInterfaceWithPK interface {
	ModelInterface
	GetPKValue() uint64
}

func SaveModel(m ModelInterfaceWithPK, withDeleted ...bool) error {
	con, err := m.GetConnection()
	if err != nil {
		return err
	}

	id := m.GetPKValue()

	var errDB error
	if id != 0 {
		query := con.ID(id).UseBool().AllCols()
		if len(withDeleted) > 0 && withDeleted[0] {
			query.Unscoped()
		}
		_, errDB = query.Update(m)
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

	if f, err := con.ID(id).Get(m); err != nil {
		return err
	} else {
		if !f {
			return fmt.Errorf("no model with id %d", id)
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
