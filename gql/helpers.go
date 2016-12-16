package gql

import (
	"errors"
	"fmt"
	"io"

	"github.com/kentquirk/giraffe-quail/parser"
	"github.com/kentquirk/giraffe-quail/types"
	"github.com/kentquirk/giraffe-quail/typeschema"
)

type GQ struct {
	Types    *types.TypeRegistry
	Scope    *types.Scope
	Handlers map[string]Handler
}

// Handler is an interface that is used to execute queries
type Handler interface {
	Operate(field types.QueryField, global, local *types.Scope) error
}

func New(t *types.TypeRegistry, s *types.Scope) *GQ {
	return &GQ{Types: t, Scope: s, Handlers: make(map[string]Handler)}
}

func FromReader(name string, r io.Reader) (*GQ, error) {
	tr, gs, err := typeschema.LoadReader(name, r)
	return New(tr, gs), err
}

func FromFile(filename string) (*GQ, error) {
	tr, gs, err := typeschema.LoadFile(filename)
	return New(tr, gs), err
}

func (gq *GQ) ParseString(name string, qs string) (*types.Document, error) {
	qi, err := parser.Parse(name, []byte(qs))
	if err != nil {
		fmt.Println("Error parsing query." + err.Error())
		return types.NewDocument(), err
	}
	return qi.(*types.Document), nil
}

// Register stores the handler
func (gq *GQ) Register(name string, h Handler) error {
	gq.Handlers[name] = h
	return nil
}

// DoOp processes an individual operation
// Not to be confused with Doo Wop, which is a genre of pop music characterized
// by close harmony and nonsense syllables.
func (gq *GQ) DoOp(op *types.Operation) error {
	switch op.Type {
	case types.QUERY:
		for _, f := range op.SelectionSet.Fields {
			fmt.Println("name = ", f.Name)
			h, ok := gq.Handlers[f.Name]
			if !ok {
				return errors.New("No handler found for " + f.Name)
			}
			fmt.Printf("f:%#v\n", f)
			err := h.Operate(f, gq.Scope, op.Variables)
			if err != nil {
				return err
			}
		}
		return nil
	case types.MUTATION:
		return errors.New("Mutations not yet supported.")
	default:
		return errors.New("DoOp: Don't know about " + string(op.Type))
	}
}

// DoOps processes a collection of operations
func (gq *GQ) DoOps(doc *types.Document) error {
	for _, op := range doc.Operations {
		fmt.Printf("%v\n", op)
		if err := gq.DoOp(op); err != nil {
			return err
		}
	}
	return nil
}
