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
	Directives   []*Directive
}

type SelectionSet struct {
	Fields []QueryField
}

type Operation struct {
	Type         OpType
	Name         string
	Variables    *Scope
	SelectionSet SelectionSet
	Directives   []*Directive
}

type Fragment struct {
}

type Directive struct {
	Name      string
	Arguments *Scope
}

// Definition just exists long enough to carry either an Operation or a Fragment
// to a Document in the parser
type Definition struct {
	Op   *Operation
	Frag *Fragment
}

type Document struct {
	Operations []*Operation
	Fragments  []*Fragment
}

func NewDocument() *Document {
	return &Document{
		Operations: make([]*Operation, 0),
		Fragments:  make([]*Fragment, 0),
	}
}

func (doc *Document) Add(def Definition) {
	if def.Op != nil {
		doc.Operations = append(doc.Operations, def.Op)
	}
	if def.Frag != nil {
		doc.Fragments = append(doc.Fragments, def.Frag)
	}
}
