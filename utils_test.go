package swagger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func compareJson(t *testing.T, actualJsonAsString string, expectedJsonAsString string) bool {
	success := false
	var actualMap map[string]interface{}
	json.Unmarshal([]byte(actualJsonAsString), &actualMap)
	var expectedMap map[string]interface{}
	err := json.Unmarshal([]byte(expectedJsonAsString), &expectedMap)
	if err != nil {
		var actualArray []interface{}
		json.Unmarshal([]byte(actualJsonAsString), &actualArray)
		var expectedArray []interface{}
		err := json.Unmarshal([]byte(expectedJsonAsString), &expectedArray)
		success = reflect.DeepEqual(actualArray, expectedArray)
		if err != nil {
			t.Fatalf("Unparsable expected JSON: %s, actual: %v, expected: %v", err, actualJsonAsString, expectedJsonAsString)
		}
	} else {
		success = reflect.DeepEqual(actualMap, expectedMap)
	}
	if !success {
		t.Log("---- expected -----")
		t.Log(withLineNumbers(expectedJsonAsString))
		t.Log("---- actual -----")
		t.Log(withLineNumbers(actualJsonAsString))
		t.Log("---- raw -----")
		t.Log(actualJsonAsString)
		t.Error("there are differences")
		return false
	}
	return true
}

func withLineNumbers(content string) string {
	var buffer bytes.Buffer
	lines := strings.Split(content, "\n")
	for i, each := range lines {
		buffer.WriteString(fmt.Sprintf("%d:%s\n", i, each))
	}
	return buffer.String()
}
func testFromStruct(t *testing.T, sample interface{}, expectedJson string) bool {
	return testJsonFromStructWithConfig(t, sample, expectedJson, &Config{})
}

func testJsonFromStructWithConfig(t *testing.T, sample interface{}, expectedJson string, config *Config) bool {
	m := modelsFromStructWithConfig(sample, config)
	data, _ := json.MarshalIndent(m, " ", " ")
	return compareJson(t, string(data), expectedJson)
}

func modelsFromStructWithConfig(sample interface{}, config *Config) map[string]*Items {
	models := map[string]*Items{}
	builder := modelBuilder{Definitions: &models, Config: config}
	builder.addModel(reflect.TypeOf(sample), "")
	return models
}
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix)
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}
