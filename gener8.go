package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type gener8 struct {
	in     string
	out    string
	pkg    string
	kws    string
	inData []byte
}

func (g *gener8) generate() {
	g.init()

	outData := string(g.inData)

	if g.pkg != "" {
		rxPkg := regexp.MustCompile(`\$pkg`)

		outData = rxPkg.ReplaceAllString(outData, g.pkg)
	}

	if g.kws != "" {

		keywords, err := parseKws(g.kws)
		if err != nil {
			panic(err)
		}

		for i, kw := range *keywords {
			rxKw := regexp.MustCompile(fmt.Sprintf("\\$kw%d", i+1))
			outData = rxKw.ReplaceAllString(outData, kw)
		}
	}

	pwd, err := os.Getwd()

	check(err)

	path := filepath.Join(pwd, g.out)

	outData = fmt.Sprintf("// %s Auto generated from %s by gener8\n\n", time.Now(), g.in) + outData

	err = ioutil.WriteFile(path, []byte(outData), 0644) // r.w r.. r..

	check(err)
}

func parseKws(kws string) (*[]string, error) {
	keywords, err := csv.NewReader(strings.NewReader(kws)).Read()
	if err != nil {
		return nil, err
	}
	return &keywords, nil
}

func (g *gener8) init() {
	flag.StringVar(&g.in, "in", "", "file to parse")
	flag.StringVar(&g.out, "out", "", "file to write the generated code to")
	flag.StringVar(&g.pkg, "pkg", "", "the value to replace $pkg with")
	flag.StringVar(&g.kws, "kws", "", "csv list of values to replace $kwn tokens with")

	flag.Parse()

	if g.in == "" || g.out == "" {
		fmt.Printf("in and out are required\n")
		os.Exit(1)
	}

	pwd, err := os.Getwd()

	check(err)

	path := filepath.Join(pwd, g.in)

	g.inData, err = ioutil.ReadFile(path)

	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func main() {
	g := &gener8{}
	g.generate()
}
