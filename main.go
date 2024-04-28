package main

import (
	"fmt"
	"log"

	"github.com/alexflint/go-arg"

	"github.com/femnad/pws/secret"
)

const (
	name    = "pws"
	version = "0.3.0"
)

type args struct {
	Overwrite bool   `arg:"-o,--overwrite" help:"Overwrite existing secret"`
	Secret    string `arg:"positional,required"`
	Vault     string `arg:"-v,--vault" help:"Vault name"`
}

func (args) Version() string {
	return fmt.Sprintf("%s v%s", name, version)
}

func main() {
	var parsed args
	arg.MustParse(&parsed)

	secretArgs := secret.Args{
		Name:      parsed.Secret,
		Overwrite: parsed.Overwrite,
		Vault:     parsed.Vault,
	}
	err := secret.Copy(secretArgs)
	if err != nil {
		log.Fatal(err)
	}
}
