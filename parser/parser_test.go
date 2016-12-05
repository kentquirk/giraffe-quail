package parser

import (
	"testing"

	"github.com/kentquirk/giraffe-quail/types"
	"github.com/kentquirk/giraffe-quail/typeschema"
	"github.com/stretchr/testify/assert"
)

func TestMain(tm *testing.M) {
	TR = types.NewTypeRegistry()
	GlobalScope = types.NewScope()
	var err error

	TR, GlobalScope, err = typeschema.LoadSchemaFromFile("tests/starwars.schema")
	if err != nil {
		panic("Couldn't load schema: " + err.Error())
	}
	tm.Run()
}

func TestSimpleHero(t *testing.T) {
	s := `query HeroNameQuery {
              hero {
                name
              }
            }
         `
	_, err := Parse("querytest", []byte(s))
	if err != nil {
		DumpErrors(err)
	}
	assert.Nil(t, err)
}

func TestNotQuery(t *testing.T) {
	s := `foo test {}`
	_, err := Parse("querytest", []byte(s))
	assert.NotNil(t, err)
}

func TestUnclosedQuote(t *testing.T) {
	s := `query TestQuery {
              Luke: human(id:"123) {
                name
              }
            }`
	_, err := Parse("querytest", []byte(s))
	assert.NotNil(t, err)
}

func TestPets(t *testing.T) {
	err := LoadQueryFromFile("tests/testqueries.gql")
	if err != nil {
		DumpErrors(err)
	}
	assert.Nil(t, err)
}
