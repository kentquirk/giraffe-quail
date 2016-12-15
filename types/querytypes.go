package types

type OpType string

const (
	QUERY    OpType = "query"
	MUTATION OpType = "mutation"
)

type QueryField struct {
	Name         string
	Alias        string
	Arguments    *Scope
	SelectionSet SelectionSet
}

type SelectionSet struct {
	Fields []QueryField
}

type Operation struct {
	Type         OpType
	Name         string
	Variables    *Scope
	SelectionSet SelectionSet
}
