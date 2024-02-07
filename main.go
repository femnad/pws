package main

import (
	"fmt"
	"log"

	"github.com/alexflint/go-arg"

	"github.com/femnad/pws/secret"
)

const (
	name    = "pws"
	version = "v0.2.0"
)

type args struct {
	Overwrite bool   `arg:"-o,--overwrite" help:"Overwrite existing secret"`
	Secret    string `arg:"positional,required"`
}

func (args) Version() string {
	return fmt.Sprintf("%s %s", name, version)
}

func main() {
	var parsed args
	arg.MustParse(&parsed)

	err := secret.Copy(parsed.Secret, parsed.Overwrite)
	if err != nil {
		log.Fatal(err)
	}
}
