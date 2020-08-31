package swagger

import (
	"testing"
)

type mathTest struct {
	beforName, afterName string
}

func TestGetModelName(t *testing.T) {
	modelName := []mathTest{
		mathTest{"main.Book", "#/definitions/Book"},
		mathTest{"swagger.Student", "#/definitions/Student"},
		mathTest{"Teacher", "#/definitions/Teacher"},
	}
	for _, n := range modelName {
		afterName := getModelName(n.beforName)
		if afterName != n.afterName {
			t.Errorf("get %v but not %v", afterName, n.afterName)
		}
	}
}

func TestGetOtherName(t *testing.T) {
	otherName := []mathTest{
		mathTest{"int", "integer"},
		mathTest{"int8", "integer"},
		mathTest{"int16", "integer"},
		mathTest{"int32", "integer"},
		mathTest{"int64", "integer"},
		mathTest{"float32", "number"},
		mathTest{"float64", "number"},
		mathTest{"string", "string"},
		mathTest{"bool", "boolean"},
		mathTest{"uint", "integer"},
		mathTest{"uint8", "integer"},
		mathTest{"uint16", "integer"},
		mathTest{"uint32", "integer"},
		mathTest{"uint64", "integer"},
		mathTest{"byte", "integer"},
		mathTest{"time.Time", "string"},
		mathTest{"file", "file"},
		mathTest{"time", "time"},
		mathTest{"date", "date"},
	}
	for _, o := range otherName {
		afterName := getOtherName(o.beforName)
		if afterName != o.afterName {
			t.Errorf("get %v but not %v", afterName, o.afterName)
		}
	}
}

func TestGetFormat(t *testing.T) {
	formatName := []mathTest{
		mathTest{"int", "int32"},
		mathTest{"int32", "int32"},
		mathTest{"int64", "int64"},
		mathTest{"float32", "float"},
		mathTest{"float64", "double"},
		mathTest{"byte", "byte"},
		mathTest{"uint8", "byte"},
		mathTest{"time.Time", "date-time"},
		mathTest{"*time.Time", "date-time"},
		mathTest{"datetime", "date-time"},
		mathTest{"dateTime", "date-time"},
		mathTest{"string", ""},
		mathTest{"file", ""},
	}
	for _, f := range formatName {
		afterName := getFormat(f.beforName)
		if afterName != f.afterName {
			t.Errorf("get %v but not %v", afterName, f.afterName)
		}
	}
}
