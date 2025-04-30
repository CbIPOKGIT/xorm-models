package xormmodels

import (
	"fmt"
	"strings"
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

type Model interface {
	GetConnection() (*xorm.Engine, error)
	TableName() string
}

type ModelPK interface {
	Model
	GetPKValue() uint64
}

type ModelIsMapped interface {
	ToMap() map[string]any
}

type ModelWithUpsertKeys interface {
	UniqueKeys() []string
	UpdateKeys() []string
}

type UpsertModel interface {
	Model
	ModelWithUpsertKeys
	ModelIsMapped
}

func Save(m Model) error {
	con, err := m.GetConnection()
	if err != nil {
		return err
	}

	var id uint64
	if modelPK, is := m.(ModelPK); is {
		id = modelPK.GetPKValue()
	}

	var errDB error
	if id != 0 {
		_, errDB = con.ID(id).UseBool().AllCols().Update(m)
	} else {
		_, errDB = con.InsertOne(m)
	}

	return errDB
}

func All(m Model) error {
	con, err := m.GetConnection()
	if err != nil {
		return err
	}

	if err := con.Find(m); err != nil {
		return err
	}
	return nil
}

func Find(m ModelPK, id any) error {
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

func Delete(m ModelPK) error {
	id := m.GetPKValue()

	con, err := m.GetConnection()
	if err != nil {
		return err
	}

	_, err = con.ID(id).Delete(m)
	return err
}

func Insert(m Model) (int64, error) {
	con, err := m.GetConnection()
	if err != nil {
		return 0, err
	}
	// defer con.Close()

	return con.Insert(m)
}

func Upsert(m UpsertModel, withTimestamps bool) error {
	return UpsertMultiple([]UpsertModel{m}, withTimestamps)
}

func UpsertMultiple(models []UpsertModel, withTimestamps bool) error {
	if len(models) == 0 {
		return nil
	}

	model := models[0]

	keys := append(model.UniqueKeys(), model.UpdateKeys()...)
	operands := make([]string, len(keys))
	for i := range keys {
		operands[i] = "?"
	}

	var sql string
	{
		insertKeys := keys[:]
		if withTimestamps {
			insertKeys = append(insertKeys, "created_at", "updated_at")
		}
		sql = fmt.Sprintf(
			`INSERT INTO %s (%s) VALUES `,
			model.TableName(),
			strings.Join(insertKeys, ","),
		)
	}
	args := make([]any, 0, len(models)*len(keys))

	insertSlice := make([]string, len(models))
	for i, m := range models {
		insertOperands := operands[:]
		if withTimestamps {
			insertOperands = append(insertOperands, "NOW()", "NOW()")
		}
		insertSlice[i] = fmt.Sprintf("(%s)", strings.Join(insertOperands, ","))

		modelMap := m.ToMap()
		for _, k := range keys {
			if v, ok := modelMap[k]; ok {
				args = append(args, v)
			} else {
				args = append(args, nil)
			}
		}
	}

	sql += strings.Join(insertSlice, ", ")

	updateValues := make([]string, len(model.UpdateKeys()))
	for i, k := range model.UpdateKeys() {
		updateValues[i] = fmt.Sprintf("%[1]s=VALUES(%[1]s)", k)
	}
	if withTimestamps {
		updateValues = append(updateValues, "updated_at=NOW()")
	}

	sql += fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", strings.Join(updateValues, ", "))

	con, err := model.GetConnection()
	if err != nil {
		return err
	}

	_, err = con.Exec(append([]any{sql}, args...)...)
	return err
}
