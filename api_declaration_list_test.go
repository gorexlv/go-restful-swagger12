package swagger

import "testing"

func TestApiDeclarationList(t *testing.T) {
	api := APIDefinition{}
	api.BasePath = "/books"
	apiDecl := ApiDeclarationList{}
	apiDecl.Put("/books", api)
	k, ok := apiDecl.At("/books")
	if !ok {
		t.Error("want model back")
	}
	if got, want := k.BasePath, "/books"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
