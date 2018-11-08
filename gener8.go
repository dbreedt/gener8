package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type gener8 struct {
	skipFormat bool
	in         string
	out        string
	pkg        string
	kws        string
	inData     []byte
	trace      bool
}

func (g *gener8) generate() {
	outData := string(g.inData)

	if g.pkg != "" {
		g.traceOut("generate:Parsing pkg...")

		rxPkg := regexp.MustCompile(`\$pkg`)

		outData = rxPkg.ReplaceAllString(outData, g.pkg)

		g.traceOut("generate:Parsing pkg complete")
	}

	if g.kws != "" {
		g.traceOut("generate:Parsing kws...")

		keywords, err := parseKws(g.kws)

		check(err)

		g.traceOut("generate:Parsing pkg complete")

		for i := len(*keywords) - 1; i > -1; i-- {
			g.traceOut("generate:processing $kw%d", i+1)

			kw := (*keywords)[i]
			rxKw := regexp.MustCompile(fmt.Sprintf("\\$kw%d", i+1))
			outData = rxKw.ReplaceAllString(outData, kw)

			g.traceOut("generate:processing $kw%d complete", i+1)
		}
	}

	g.traceOut("generate:Getwd")

	pwd, err := os.Getwd()

	check(err)

	g.traceOut("generate:Construct out path")

	path := filepath.Join(pwd, g.out)

	outData = fmt.Sprintf("// %s Auto generated from %s by gener8\n\n", time.Now(), g.in) + outData

	g.traceOut("generate:WriteFile: path: %s", path)

	err = ioutil.WriteFile(path, []byte(outData), 0644) // r.w r.. r..

	check(err)

	g.traceOut("generate:WriteFile Complete")

	if !g.skipFormat {
		g.traceOut("generate:Formatting output file")

		cmd := exec.Command("gofmt", "-w", path)
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		if err != nil {
			check(fmt.Errorf("gofmt -w %s failed: %s", g.out, err))
		}

		g.traceOut("generate:Formatting complete")
	} else {
		g.traceOut("generate:SkipFormat")
	}
}

func parseKws(kws string) (*[]string, error) {
	keywords, err := csv.NewReader(strings.NewReader(kws)).Read()

	if err != nil {
		return nil, err
	}

	return &keywords, nil
}

func (g *gener8) setup() {
	flag.BoolVar(&g.skipFormat, "skip_format", false, "skip gofmt being run on the generated file")
	flag.BoolVar(&g.trace, "trace", false, "enables trace logging")
	flag.StringVar(&g.in, "in", "", "file to parse")
	flag.StringVar(&g.out, "out", "", "file to write the generated code to")
	flag.StringVar(&g.pkg, "pkg", "", "the value to replace $pkg with")
	flag.StringVar(&g.kws, "kws", "", "csv list of values to replace $kwn tokens with")

	flag.Parse()

	if g.in == "" || g.out == "" {
		fmt.Printf("in and out are required\n")
		os.Exit(1)
	}

	g.traceOut("setup:Complete flag.Parse")

	pwd, err := os.Getwd()

	check(err)

	path := filepath.Join(pwd, g.in)

	g.traceOut("setup:ReadFile: %s", path)

	g.inData, err = ioutil.ReadFile(path)

	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func (g *gener8) traceOut(format string, a ...interface{}) {
	if g.trace {
		fmt.Fprintf(os.Stderr, format+"\n", a...)
	}
}

func main() {
	g := &gener8{}
	g.setup()
	g.generate()
}
