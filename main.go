package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func main() {
	gopath, _ := filepath.Abs("gopath")
	config := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Tests: true,
		Env:   append(os.Environ(), "GOPATH="+gopath, "GO111MODULE=off"),
	}

	pkgs, err := packages.Load(config, "github.com/docker/docker/daemon/cluster/executor/container")
	if err != nil {
		log.Fatal(err)
	} else if packages.PrintErrors(pkgs) > 0 {
		log.Fatal(errors.New("errors encountered while loading packages"))
	}

	prog, _ := ssautil.AllPackages(pkgs, ssa.InstantiateGenerics)
	prog.Build()

	log.Println("Built SSA program")

	// Extract all main package candidates from the SSA program.
	mains := ssautil.MainPackages(prog.AllPackages())
	if len(mains) == 0 {
		log.Println("No main packages detected")
		return
	}

	a_config := &pointer.Config{
		Mains:          mains,
		BuildCallGraph: true,
	}

	for fun := range ssautil.AllFunctions(prog) {
		for _, block := range fun.Blocks {
			for _, insn := range block.Instrs {
				if val, ok := insn.(ssa.Value); ok && pointer.CanPoint(val.Type()) {
					a_config.AddQuery(val)
				}
			}
		}
	}

	_, err = pointer.Analyze(a_config)
	if err != nil {
		fmt.Println("Failed pointer analysis")
		fmt.Println(err)
		os.Exit(1)
	}

	log.Println("OK")
}
