package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

const (
	chunkSize = 64000
)

func (g *gener8) generate() {
	subs := []string{}

	if g.pkg != "" {
		subs = append(subs, "$pkg", g.pkg)
	}

	if g.kws != "" {
		keywords, err := parseKws(g.kws)
		check(err)

		for i := len(keywords) - 1; i > -1; i-- {
			key := fmt.Sprintf("$kw%d", i+1)
			subs = append(subs, key, keywords[i])

		}
	}
	replacer := strings.NewReplacer(subs...)

	g.traceOut("generate:Getwd")

	pwd, err := os.Getwd()

	check(err)

	g.traceOut("generate:Construct out path")

	path := filepath.Join(pwd, g.out)

	tmpFile, err := ioutil.TempFile("", "gener8")

	check(err)

	defer os.Remove(tmpFile.Name())

	g.traceOut("generate:WriteTempFile: path: %s", tmpFile.Name())

	_, err = fmt.Fprintf(tmpFile, "// Auto generated from %s by gener8\n\n", g.in)
	check(err)

	_, err = replacer.WriteString(tmpFile, string(g.inData))
	check(err)

	err = tmpFile.Close()
	check(err)

	if !g.skipFormat {
		g.traceOut("generate:Formatting tmpFile")

		cmd := exec.Command("gofmt", "-w", tmpFile.Name())
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		if err != nil {
			check(fmt.Errorf("gofmt -w %s failed: %s", g.out, err))
		}

		g.traceOut("generate:Formatting complete")
	} else {
		g.traceOut("generate:SkipFormat - tmpFile")
	}

	if !compareFiles(tmpFile.Name(), path) {
		g.traceOut("generate:WriteFile: path: %s", path)
		rawOutData, err := ioutil.ReadFile(tmpFile.Name())
		check(err)

		err = ioutil.WriteFile(path, rawOutData, 0644) // r.w r.. r..
		check(err)

		g.traceOut("generate:WriteFile Complete")
	} else {
		g.traceOut("generate:WriteFile Skipped...")
	}
}

func parseKws(kws string) ([]string, error) {
	keywords, err := csv.NewReader(strings.NewReader(kws)).Read()

	if err != nil {
		return nil, err
	}

	return keywords, nil
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

func compareFiles(file1, file2 string) bool {
	f1, err := os.Open(file1)

	if err != nil && os.IsNotExist(err) {
		return false
	}

	check(err)

	f2, err := os.Open(file2)

	if err != nil && os.IsNotExist(err) {
		return false
	}

	check(err)

	b1 := make([]byte, chunkSize)
	b2 := make([]byte, chunkSize)

	for {
		br1, err1 := f1.Read(b1)
		br2, err2 := f2.Read(b2)

		b1 = b1[:br1]
		b2 = b2[:br2]

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				check(err1)
				check(err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

func main() {
	g := &gener8{}
	g.setup()
	g.generate()
}
