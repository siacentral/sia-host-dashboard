// +build ignore
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/shurcooL/vfsgen"
)

func main() {
	if err := os.MkdirAll("assets/assets.go", 0770); err != nil && !os.IsExist(err) {
		log.Fatalln(err)
	}

	err := vfsgen.Generate(http.Dir(os.Args[1]), vfsgen.Options{
		Filename:     "assets/assets.go",
		PackageName:  "assets",
		VariableName: "Web",
	})

	if err != nil {
		log.Fatalln(err)
	}
}
