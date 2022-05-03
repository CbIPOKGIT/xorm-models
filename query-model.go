package xormmodels

type QueryModel struct {
	Select string
	Where  map[string]interface{}
	Limit  int
	Offset int
	Order  string
}

func NewQueryModel() *QueryModel {
	qm := new(QueryModel)
	qm.Where = make(map[string]interface{})
	return qm
}

func (qm *QueryModel) LatestId() {
	qm.Order = "id ASC"
}
