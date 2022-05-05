package xormmodels

import (
	"strings"

	"github.com/go-xorm/xorm"
)

type QueryJoin struct {
	Operator  string
	Table     string
	Condition string
}

type QueryModel struct {
	Select  string
	Selects []string
	Where   []string
	Join    QueryJoin
	Having  []string
	GroupBy []string
	Limit   int
	Offset  int
	Order   string
	Table   string
}

func NewQueryModel() *QueryModel {
	qm := new(QueryModel)
	qm.Where = make([]string, 0)
	qm.Selects = make([]string, 0)
	qm.GroupBy = make([]string, 0)
	qm.Having = make([]string, 0)
	return qm
}

func (qm *QueryModel) LatestId() {
	qm.Order = "id DESC"
}

func (qm *QueryModel) OrderByID() {
	qm.Order = "id ASC"
}

func (qm QueryModel) Fill(con *xorm.Session) {
	if qm.Select != "" {
		qm.Selects = append(qm.Selects, qm.Select)
	}

	if len(qm.Selects) > 0 {
		con.Select(strings.Join(qm.Selects, ","))
	}

	if qm.Where != nil {
		for _, val := range qm.Where {
			con.Where(val)
		}
	}

	if qm.Join.Operator != "" {
		con.Join(qm.Join.Operator, qm.Join.Table, qm.Join.Condition)
	}

	if len(qm.GroupBy) > 0 {
		con.GroupBy(strings.Join(qm.GroupBy, ","))
	}

	if len(qm.Having) > 0 {
		con.Having(strings.Join(qm.Having, " AND "))
	}

	if qm.Order != "" {
		con.OrderBy(qm.Order)
	}

	con.Limit(qm.Limit, qm.Offset)
}
