package types

type QueryField struct {
	Name         string
	Alias        string
	Arguments    *Scope
	SelectionSet SelectionSet
}

type SelectionSet struct {
	Fields []QueryField
}
