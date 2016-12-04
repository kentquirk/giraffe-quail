package typeschema

import (
	"testing"

	"github.com/kentquirk/giraffe-quail/types"
	"github.com/stretchr/testify/assert"
)

func TestMain(tm *testing.M) {
	TR = types.NewTypeRegistry()
	VR = types.NewValueRegistry()

	tm.Run()
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
	s := `type Bear {
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
	s := `type Ferret implements Pet {
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
	s := `type Pig {
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
	s := `type Lion {
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
	s := `type HouseCat implements Pet {
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

func TestPets(t *testing.T) {
	_, _, err := LoadSchemaFromFile("tests/pets.schema")
	if err != nil {
		DumpErrors(err)
	}
	assert.Nil(t, err)
}
