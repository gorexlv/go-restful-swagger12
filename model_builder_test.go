package swagger

import (
	"reflect"
	"testing"
)

func TestSampleToModelAsJson(t *testing.T) {
	testFromStruct(t, sample{items: []item{}}, `{
  "sample": {
    "type": "object",
    "properties": {
       "id": {
         "type": "string"
     },
       "items": {
         "type": "array",
         "items": {
           "$ref": "#/definitions/item"
       }
     },
      "root": {
        "$ref": "#/definitions/item"
       }
     }
  },
  "item": {
    "type": "object",
    "properties": {
      "name": {
        "type": "string"
       }
     }
   }
 }`)
}

//test anonymous struct
func TestAnonymousPtrStruct(t *testing.T) {
	type X struct {
		A *struct {
			B int
		}
	}

	expected := `{
  "X": {
   "type": "object",
   "properties": {
     "A": {
       "$ref": "#/definitions/X.A"
    }
   }
  },
  "X.A": {
    "type": "object",
    "properties": {
      "B": {
        "type": "integer",
        "format": "int32"
      }
    }
  }
 }`
	testFromStruct(t, X{}, expected)
}

type File struct {
	History     []File
	HistoryPtrs []*File
}

//test recursion struct
func TestRecursiveStructure(t *testing.T) {
	testFromStruct(t, File{}, `{
  "File": {
   "type": "object",
   "properties": {
    "History": {
     "type": "array",
     "items": {
      "$ref": "#/definitions/File"
     }
    },
    "HistoryPtrs": {
     "type": "array",
     "items": {
      "$ref": "#/definitions/File"
     }
    }
   }
  }
 }`)
}

func TestAtMap(t *testing.T) {
	name := []string{"Book", "Student", "Class", "abc", "ABC"}
	mapItem := map[string]*Items{}
	mapItem["Book"] = &Items{}
	mapItem["Student"] = &Items{}
	mapItem["Class"] = &Items{}
	mapItem["abc"] = &Items{}
	mapItem["ABC"] = &Items{}
	for _, n := range name {
		if _, ok := atMap(n, &mapItem); !ok {
			t.Errorf("get %v but not %v", false, true)
		}
	}
}

func TestIsPrimitiveType(t *testing.T) {
	modelName := []string{"uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16",
		"int32", "int64", "float32", "float64", "bool", "string", "byte", "rune", "time.Time"}
	b := modelBuilder{}
	for _, name := range modelName {
		if ok := b.isPrimitiveType(name); !ok {
			t.Errorf("get %v but not %v", false, true)
		}
	}
}

type model struct {
	strType   string
	intType   int
	boolType  bool
	sliceType []string
	ptrType   *string
	mapType   map[string]Example
	exaType   Example
}

type Example struct {
	Id, Name string
}

func TestKeyFrom(t *testing.T) {
	typed := []string{"string", "int", "bool", "string", "*string", "map[string]swagger.Example", "swagger.Example"}
	b := modelBuilder{}
	model := model{}
	st := reflect.TypeOf(model)
	for i := 0; i < st.NumField(); i++ {
		tp := b.keyFrom(st.Field(i).Type)
		if tp != typed[i] {
			t.Errorf("get %v but not %v", tp, typed[i])
		}
	}
}

func Test(t *testing.T) {
	type X struct {
		A int
		B int `json:"C"`
		D int `json:"-"`
	}

	expected := `{
	  "X": {
		"type": "object",
	   "properties": {
		"A": {
		 "type": "integer",
		 "format": "int32"
		},
		"C": {
		 "type": "integer",
		 "format": "int32"
		}
	   }
	  }
	 }`

	testFromStruct(t, X{}, expected)
}
