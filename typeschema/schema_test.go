package typeschema

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(tm *testing.M) {
	TR = NewTypeRegistry()
	VR = NewValueRegistry()

	tm.Run()
}

func DumpErrors(err error) {
	list := err.(errList)
	for _, err := range list {
		pe := err.(*parserError)
		fmt.Printf("%+v\n", pe)
	}
}

func TestSingleEnum(t *testing.T) {
	s := `enum CatCommand { JUMP }`
	_, err := Parse("schematest", []byte(s))
	if err != nil {
		DumpErrors(err)
	}
	assert.Nil(t, err)
}

func TestMultiEnum(t *testing.T) {
	s := `enum DogCommands { SIT, DOWN, HEEL }`
	_, err := Parse("schematest", []byte(s))
	if err != nil {
		DumpErrors(err)
	}
	assert.Nil(t, err)
}

func TestRedefineEnum(t *testing.T) {
	s := `enum CatCommand { SLEEP }`
	_, err := Parse("schematest", []byte(s))
	assert.NotNil(t, err)
}

func TestRedefineEnumValue(t *testing.T) {
	s := `enum FishCommand { JUMP }`
	_, err := Parse("schematest", []byte(s))
	assert.NotNil(t, err)
}

func TestSimpleObj(t *testing.T) {
	s := `type Cat {
            name: String
            }
    `
	_, err := Parse("schematest", []byte(s))
	if err != nil {
		DumpErrors(err)
	}
	assert.Nil(t, err)
}

func TestObjImplements(t *testing.T) {
	s := `type Cat implements Pet {
          name: String
        }
    `
	_, err := Parse("schematest", []byte(s))
	if err != nil {
		DumpErrors(err)
	}
	assert.Nil(t, err)
}

func TestObjNonnullField(t *testing.T) {
	s := `type Cat {
          name: String!
        }
    `
	_, err := Parse("schematest", []byte(s))
	if err != nil {
		DumpErrors(err)
	}
	assert.Nil(t, err)
}

func TestFieldArgument(t *testing.T) {
	s := `type Cat {
          doesKnowCommand(catCommand: CatCommand!): Boolean!
        }
    `
	_, err := Parse("schematest", []byte(s))
	if err != nil {
		DumpErrors(err)
	}
	assert.Nil(t, err)
}

func TestFieldMultiArgument(t *testing.T) {
	s := `type Cat {
          doesKnowCommand(catCommand: CatCommand!, whispered: Boolean!): Boolean!
        }
    `
	_, err := Parse("schematest", []byte(s))
	if err != nil {
		DumpErrors(err)
	}
	assert.Nil(t, err)
}

func TestComplexObj(t *testing.T) {
	s := `type Cat implements Pet {
          name: String!
          nickname: String
          doesKnowCommand(catCommand: CatCommand!): Boolean!
          meowVolume: Int
        }
    `
	_, err := Parse("schematest", []byte(s))
	if err != nil {
		DumpErrors(err)
	}
	assert.Nil(t, err)
}

func TestInterface(t *testing.T) {
	s := `interface Pet {
              name: String!
            }
    `
	_, err := Parse("schematest", []byte(s))
	if err != nil {
		DumpErrors(err)
	}
	assert.Nil(t, err)
}

func TestUnion(t *testing.T) {
	s := `union CatOrDog = Cat | Dog `
	_, err := Parse("schematest", []byte(s))
	if err != nil {
		DumpErrors(err)
	}
	assert.Nil(t, err)
}

func checkOneFile(t *testing.T, filename string) {
	b, err := ioutil.ReadFile(filename)
	assert.Nil(t, err)
	_, err = Parse(filename, b)
	if err != nil {
		DumpErrors(err)
	}
	assert.Nil(t, err)
}

func TestPets(t *testing.T) {
	TR = NewTypeRegistry()
	VR = NewValueRegistry()
	checkOneFile(t, "tests/pets.schema")
}
