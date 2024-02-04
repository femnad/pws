package main

import (
	"fmt"
	"log"

	"github.com/alexflint/go-arg"

	"github.com/femnad/pws/secret"
)

const (
	name    = "pws"
	version = "v0.1.0"
)

type args struct {
	Secret string `arg:"positional,required"`
}

func (args) Version() string {
	return fmt.Sprintf("%s %s", name, version)
}

func main() {
	var parsed args
	arg.MustParse(&parsed)

	err := secret.Copy(parsed.Secret)
	if err != nil {
		log.Fatal(err)
	}
}
