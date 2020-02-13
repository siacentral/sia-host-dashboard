// +build ignore
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(http.Dir(os.Args[1]), vfsgen.Options{
		Filename:     "assets/assets.go",
		PackageName:  "assets",
		VariableName: "Web",
	})

	if err != nil {
		log.Fatalln(err)
	}
}
