package main

import (
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/emicklei/go-restful"

	"fmt"
	"github.com/go-chassis/go-restful-swagger20"
	"os"
	"path/filepath"
)

type Book struct {
	Id      string
	Title   string
	Author  string
	Student []Student
}

type ID string
type Age int64

type Student struct {
	Id   ID  `swag:"string"`
	Age  Age
	Name string
}

func modelTypeNameHandler(st reflect.Type) (string, bool)  {
	key := st.String()

	if len(st.Name()) == 0 {
		key = strings.Replace(key, "[]", "", -1)
	}

	if key=="main.Age" {
		return "number", true
	}


	return key, true
}

func main() {
	ws := new(restful.WebService)
	ws.Path("/book")
	ws.Consumes(restful.MIME_JSON, restful.MIME_XML)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	restful.Add(ws)
	ws.Route(ws.GET("/{medium}").To(getBookById).
		Doc("Search a books").
		Param(ws.PathParameter("medium", "digital or paperback").DataType("string")).
		Param(ws.QueryParameter("language", "en,nl,de").DataType("string")).Metadata("tags",[]string{"users","desc"}).
		Param(ws.HeaderParameter("If-Modified-Since", "last known timestamp").DataType("string").DataFormat("datetime")).
		Returns(200, "haha", Book{}))

	ws.Route(ws.PUT("/{medium}").To(modifyBookById).
		Operation("modifyBookById").
		Doc("modify a book").
		Param(ws.PathParameter("medium", "digital or paperback").DataType("string")).
		Reads(Book{Id: "2", Title: "go", Author: "lisi"}).
		Do(returns200, returns500))

	ws.Route(ws.POST("/add").To(addBook).
		Notes("add a book").
		Reads(Student{}).
		Do(returns200, returns500))

	ws.ApiVersion("1.0.1")

	val := os.Getenv("SWAGGERFILEPATH")
	fmt.Println(val)
	if val == "" {
		val, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}
	config := swagger.Config{
		WebServices:    restful.DefaultContainer.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: "http://localhost:8080",
		ApiPath:        "/apidocs.json",
		//FileStyle:	"json",
		OpenService: true,
		SwaggerPath: "/apidocs/",
		OutFilePath: filepath.Join(val, "api.yaml"),
		ModelTypeNameHandler: modelTypeNameHandler}

	config.Info.Description = "This is a sample server Book server"
	config.Info.Title = "swagger Book"
	swagger.RegisterSwaggerService(config, restful.DefaultContainer)

	log.Print("start listening on localhost:8080")
	server := &http.Server{Addr: ":8080", Handler: restful.DefaultContainer}
	log.Fatal(server.ListenAndServe())
}

func getBookById(req *restful.Request, resp *restful.Response) {
	book := Book{Id: "1", Title: "java", Author: "zhangsan"}
	id := req.PathParameter("medium")
	if id != book.Id {
		resp.WriteErrorString(http.StatusNotFound, "Book could not be found.")
	} else {
		resp.WriteEntity(book)
	}
}
func modifyBookById(req *restful.Request, resp *restful.Response) {}

func addBook(req *restful.Request, resp *restful.Response) {}

func returns200(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", map[string]Book{})
}

func returns500(b *restful.RouteBuilder) {
	b.Returns(http.StatusInternalServerError, "Bummer, something went wrong", nil)
}
