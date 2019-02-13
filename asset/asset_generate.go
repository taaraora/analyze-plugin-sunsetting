// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/httpfs/union"
	"github.com/shurcooL/vfsgen"
)

func main() {

	var fileSystem http.FileSystem = union.New(map[string]http.FileSystem{
		"/check":      http.Dir("./check"),
		"/settings": http.Dir("./settings"),
	})

	err := vfsgen.Generate(fileSystem, vfsgen.Options{
		PackageName:  "asset",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
