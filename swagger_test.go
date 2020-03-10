package swagger

import (
	"encoding/json"
	"github.com/emicklei/go-restful"
	"os"
	"strings"
	"testing"
)

//测试Info
func TestInfoStruct(t *testing.T) {
	config := Config{
		Info: Info{
			Title:       "Title",
			Description: "Description",
			Version:     "Version",
		},
	}
	sws := newSwaggerService(config)
	listing := APIDefinition{
		Swagger:  swaggerVersion,
		Info:     sws.config.Info,
		BasePath: "",
		Paths:    nil,
	}
	str, err := json.MarshalIndent(listing, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	compareJson(t, string(str), `
	{
		"swagger": "2.0",
		"paths": null,
		"basePath": "",
		"info": {
			"title": 	"Title",
			"description": 	"Description",
			"version":	"Version"
		}
	}
	`)
}

//测试Api
func TestServiceToApi(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/tests")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_XML)
	ws.Route(ws.GET("/a").To(dummy).Writes(sample{}))
	ws.Route(ws.PUT("/b").To(dummy).Writes(sample{}))
	ws.Route(ws.POST("/c").To(dummy).Writes(sample{}))
	ws.Route(ws.DELETE("/d").To(dummy).Writes(sample{}))

	ws.Route(ws.GET("/d").To(dummy).Writes(sample{}))
	ws.Route(ws.PUT("/c").To(dummy).Writes(sample{}))
	ws.Route(ws.POST("/b").To(dummy).Writes(sample{}))
	ws.Route(ws.DELETE("/a").To(dummy).Writes(sample{}))
	ws.ApiVersion("1.2.3")
	cfg := Config{
		WebServicesUrl:   "http://here.com",
		ApiPath:          "/apipath",
		WebServices:      []*restful.WebService{ws},
		PostBuildHandler: func(in *ApiDeclarationList) {},
	}
	sws := newSwaggerService(cfg)
	decl := sws.composeDeclaration(ws, "/tests")

	if decl.Info.Version != "1.2.3" {
		t.Errorf("got %v want %v", decl.Swagger, "1.2.3")
	}
	if decl.BasePath != "/tests" {
		t.Errorf("got %v want %v", decl.BasePath, "/tests")
	}
	if len(decl.Paths) != 4 {
		t.Errorf("got %v want %v", len(decl.Paths), 4)
	}
	pathOrder := ""
	for path, _ := range decl.Paths {
		pathOrder += path
	}

	if len(pathOrder) != 8 {
		t.Errorf("got %v want %v", len(pathOrder), 8)
	}
}

func dummy(req *restful.Request, res *restful.Response) {}

type sample struct {
	id       string `swagger:"required"` // TODO
	items    []item
	rootItem item `json:"root" description:"root desc"`
}
type item struct {
	itemName string `json:"name"`
}
type TestItem struct {
	Id, Name string
}
type User struct {
	Id, Name string
}
type Responses struct {
	Code  int
	Users *[]User
	Items *[]TestItem
}

//测试responses
func TestComposeResponses(t *testing.T) {
	responseErrors := map[int]restful.ResponseError{}
	responseErrors[400] = restful.ResponseError{Code: 400, Message: "Bad Request", Model: TestItem{}}
	responseErrors[200] = restful.ResponseError{
		Headers: map[string]restful.Header{
			"X-Test-Integer": {
				Items: &restful.Items{
					Type:   "integer",
					Format: "int32",
				},
				Description: "test integer header",
			},
			"X-Test-Array": {
				Items: &restful.Items{
					Type: "array",
					Items: &restful.Items{
						Type: "array",
						Items: &restful.Items{
							Format: "int64",
							Type:   "integer",
						},
					},
				},
				Description: "test array header",
			},
		},
	}
	route := restful.Route{ResponseErrors: responseErrors}
	decl := new(APIDefinition)
	decl.Definitions = map[string]*Items{}
	msgs := composeResponses(route, decl, &Config{})
	if msgs["400"].Description != "Bad Request" {
		t.Errorf("got %s want Bad Request", msgs["400"].Description)
	}
	if msgs["400"].Schema.Ref != "#/definitions/TestItem" {
		t.Errorf("got %s want #/definitions/TestItem", msgs["400"].Schema.Ref)
	}
	if msgs["200"].Headers["X-Test-Integer"].Type != "integer" {
		t.Errorf("got %s want integer", msgs["200"].Headers["X-Test-Integer"].Type)
	}
	if msgs["200"].Headers["X-Test-Array"].Items.Items.Format != "int64" {
		t.Errorf("got %s want int64", msgs["200"].Headers["X-Test-Array"].Items.Items.Format)
	}
}

//测试Definitions
func TestAddModel(t *testing.T) {
	sws := newSwaggerService(Config{})
	api := APIDefinition{
		Definitions: map[string]*Items{}}

	sws.addModelFromSampleTo(Responses{Items: &[]TestItem{}}, &api.Definitions)
	model, ok := atMap("Responses", &api.Definitions)
	if !ok {
		t.Fatal("missing Responses model")
	}
	if model.Type != "object" {
		t.Fatal("wrong model type  " + model.Type.(string))
	}
	str := ""
	for key, _ := range model.Properties {
		str = str + key
	}
	if !strings.Contains(str, "Code") {
		t.Fatal("missing code")
	}
	if !strings.Contains(str, "Users") {
		t.Fatal("missing User")
	}
	if !strings.Contains(str, "Items") {
		t.Fatal("missing Items")
	}
	if model.Properties["Code"].Type != "integer" {
		t.Fatal("wrong code type:" + model.Properties["Code"].Type.(string))
	}
	if model.Properties["Users"].Type != "array" {
		t.Fatal("wrong Users type:" + model.Properties["Users"].Type.(string))
	}
	if model.Properties["Items"].Type != "array" {
		t.Fatal("wrong Items type:" + model.Properties["Items"].Type.(string))
	}
	if model.Properties["Users"].Items == nil {
		t.Fatal("missing Users items")
	}
	if model.Properties["Items"].Items == nil {
		t.Fatal("missing Items items")
	}
	if model.Properties["Users"].Items.Ref != "#/definitions/User" {
		t.Fatal("wrong Users Ref:" + model.Properties["Users"].Items.Ref)
	}
	if model.Properties["Items"].Items.Ref != "#/definitions/TestItem" {
		t.Fatal("wrong Items Ref:" + model.Properties["Items"].Items.Ref)
	}

	model1, ok1 := atMap("User", &api.Definitions)
	if !ok1 {
		t.Fatal("missing User model")
	}
	if model1.Type != "object" {
		t.Fatal("wrong model User type  " + model1.Type.(string))
	}
	str1 := ""
	for key, _ := range model1.Properties {
		str1 = str1 + key
	}
	if !strings.Contains(str1, "Id") {
		t.Fatal("missing User Id")
	}
	if !strings.Contains(str1, "Name") {
		t.Fatal("missing User Name")
	}
	if model1.Properties["Id"].Type != "string" {
		t.Fatal("wrong User Id type:" + model1.Properties["Id"].Type.(string))
	}
	if model1.Properties["Name"].Type != "string" {
		t.Fatal("wrong User Name type:" + model1.Properties["Name"].Type.(string))
	}

	model2, ok2 := atMap("TestItem", &api.Definitions)
	if !ok2 {
		t.Fatal("missing TestItem model")
	}
	if model2.Type != "object" {
		t.Fatal("wrong model TestItem type  " + model2.Type.(string))
	}
	str2 := ""
	for key, _ := range model2.Properties {
		str2 = str2 + key
	}
	if !strings.Contains(str2, "Id") {
		t.Fatal("missing TestItem Id")
	}
	if !strings.Contains(str2, "Name") {
		t.Fatal("missing TestItem Name")
	}
	if model2.Properties["Id"].Type != "string" {
		t.Fatal("wrong TestItem Id type:" + model2.Properties["Id"].Type.(string))
	}
	if model2.Properties["Name"].Type != "string" {
		t.Fatal("wrong TestItem Name type:" + model2.Properties["Name"].Type.(string))
	}
}

//测试将openapi协议以json形式保存在本地
func TestWriteJsonToFile(t *testing.T) {
	//测试前请先设定环境变量
	val := os.Getenv("SWAGGERFILEPATH")
	os.Remove(val)
	os.Mkdir(val, 0777)
	ws := new(restful.WebService)
	ws.Path("/file")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/write").To(dummy).Writes(sample{}))
	cfg := Config{
		WebServices: []*restful.WebService{ws},
		FileStyle:   "json",
		OutFilePath: val,
	}
	sws := newSwaggerService(cfg)
	sws.WriteToFile()
	files, err := ListDir(val, "json")
	if err != nil || len(files) != 1 {
		t.Fatal("No local json file was generated")
	}
}

//测试将openapi协议以yaml形式保存在本地
func TestWriteYamlToFile(t *testing.T) {
	val := os.Getenv("SWAGGERFILEPATH")
	os.RemoveAll(val)
	os.Mkdir(val, 0777)
	ws := new(restful.WebService)
	ws.Path("/file")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/write").To(dummy).Writes(sample{}))
	cfg := Config{
		WebServices: []*restful.WebService{ws},
		OutFilePath: val,
	}
	sws := newSwaggerService(cfg)
	sws.WriteToFile()
	files, err := ListDir(val, "yaml")
	if err != nil || len(files) != 1 {
		t.Fatal("No local yaml file was generated")
	}
}
