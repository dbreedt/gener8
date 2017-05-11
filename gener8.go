package main

import (
  "flag"
  "io/ioutil"
  "os"
  "path/filepath"
  "fmt"
  "regexp"
)

type gener8 struct {
  in string
  out string
  pkg string
  kw1 string
  kw2 string
  kw3 string
  inData []byte
}

func (g *gener8) generate() {
  g.init()

  outData := string(g.inData)

  if g.pkg != "" {
    rxPkg := regexp.MustCompile(`\$pkg`)

    outData = rxPkg.ReplaceAllString(outData, g.pkg)
  }

  if g.kw1 != "" {
    rxKw1 := regexp.MustCompile(`\$kw1`)

    outData = rxKw1.ReplaceAllString(outData, g.kw1)
  }

  if g.kw2 != "" {
    rxKw2 := regexp.MustCompile(`\$kw2`)

    outData = rxKw2.ReplaceAllString(outData, g.kw2)
  }

  if g.kw3 != "" {
    rxKw3 := regexp.MustCompile(`\$kw3`)

    outData = rxKw3.ReplaceAllString(outData, g.kw3)
  }

  pwd, err := os.Getwd()

  check(err)

  path := filepath.Join(pwd, g.out)

  err = ioutil.WriteFile(path, []byte(outData), 0644) // r.w r.. r..

  check(err)
}

func (g *gener8) init() {
  flag.StringVar(&g.in, "in", "", "file to parse")
  flag.StringVar(&g.out, "out", "", "file to write the generated code to")
  flag.StringVar(&g.pkg, "pkg", "", "the value to replace $pkg with")
  flag.StringVar(&g.kw1, "kw1", "", "the value to replace $kw1 with")
  flag.StringVar(&g.kw2, "kw2", "", "the value to replace $kw2 with")
  flag.StringVar(&g.kw3, "kw3", "", "the value to replace $kw3 with")

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
  g:= &gener8{}
  g.generate()
}
