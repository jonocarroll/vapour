package main

import (
	"log"
	"os"

	"github.com/devOpifex/vapour/cli"
	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
	"github.com/devOpifex/vapour/transpiler"
	"github.com/devOpifex/vapour/walker"
)

func (v *vapour) transpile(conf cli.CLI) {
	v.root = conf.Indir
	err := v.readDir()

	if err != nil {
		log.Fatal("Failed to read vapour files")
	}

	// lex
	l := lexer.New(v.files)
	l.Run()

	if l.HasError() {
		l.Errors.Print()
		return
	}

	// parse
	p := parser.New(l)
	prog := p.Run()

	if p.HasError() {
		p.Errors().Print()
		return
	}

	// walk tree
	w := walker.New()
	w.Walk(prog)

	w.Errors().Print()

	if w.HasError() {
		return
	}

	if *conf.Check {
		return
	}

	// transpile
	trans := transpiler.New()
	trans.Transpile(prog)
	code := trans.GetCode()

	if *conf.Run {
		run(code)
		return
	}

	code = addHeader(code)

	// write
	path := *conf.Outdir + "/" + *conf.Outfile
	f, err := os.Create(path)

	if err != nil {
		log.Fatalf("Failed to create output file: %v", err.Error())
	}

	defer f.Close()

	_, err = f.WriteString(code)

	if err != nil {
		log.Fatalf("Failed to write output file: %v", err.Error())
	}

	// write types
	lines := trans.Env().GenerateTypes().String()
	f, err = os.Create(*conf.Types)

	if err != nil {
		log.Fatalf("Failed to create type file: %v", err.Error())
	}

	defer f.Close()

	_, err = f.WriteString(lines)

	if err != nil {
		log.Fatalf("Failed to write to types file: %v", err.Error())
	}
}

func (v *vapour) transpileFile(conf cli.CLI) {
	content, err := os.ReadFile(*conf.Infile)

	if err != nil {
		log.Fatal("Could not read vapour file")
	}

	// lex
	l := lexer.NewCode(*conf.Infile, string(content))
	l.Run()

	if l.HasError() {
		l.Errors.Print()
		return
	}

	// parse
	p := parser.New(l)
	prog := p.Run()

	if p.HasError() {
		p.Errors().Print()
		return
	}

	// walk tree
	w := walker.New()
	w.Walk(prog)
	w.Errors().Print()
	if w.HasError() {
		return
	}

	if *conf.Check {
		return
	}

	// transpile
	trans := transpiler.New()
	trans.Transpile(prog)
	code := trans.GetCode()

	if *conf.Run {
		run(code)
		return
	}

	code = addHeader(code)

	// write
	f, err := os.Create(*conf.Outfile)

	if err != nil {
		log.Fatal("Failed to create output file")
	}

	defer f.Close()

	_, err = f.WriteString(code)

	if err != nil {
		log.Fatal("Failed to write to output file")
	}
}
