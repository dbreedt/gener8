package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParsing(t *testing.T) {
	kws := "a,b,c"

	keywords, err := parseKws(kws)

	if err != nil {
		t.Error(err)
	}

	if len(*keywords) != 3 {
		t.Error("Found ", len(*keywords), " should be 3")
	}

	kws = "a,b,\"d,e\""

	keywords, err = parseKws(kws)

	if err != nil {
		t.Error(err)
	}

	if len(*keywords) != 3 {
		t.Error("Found ", len(*keywords), " should be 3")
	}

	kws = "AAA,AaA,aaaAAAaaa,"

	keywords, err = parseKws(kws)

	if err != nil {
		t.Error(err)
	}

	if len(*keywords) != 4 {
		t.Error("Found ", len(*keywords), " should be 3")
	}
}

func TestGenerate(t *testing.T) {
	pwd, err := os.Getwd()

	if err != nil {
		t.Error(err)
	}

	out := fmt.Sprintf("test_%d.gr8.tmp", time.Now().Unix())
	testFileOutput := filepath.Join(pwd, out)

	defer func() {
		os.Remove(testFileOutput)
	}()

	g := &gener8{
		skipFormat: true,
		kws:        "1,2,3,4,5,6,7,8,9,10,11",
		out:        out,
		inData: []byte(`
$kw10
$kw11
$kw9
$kw8
$kw7
$kw5
$kw6
$kw1
$kw3
$kw2
$kw4
`),
	}

	g.generate()

	inData, err := ioutil.ReadFile(testFileOutput)

	if err != nil {
		t.Errorf("can't read output file: %s. Reason: %+v", testFileOutput, err)
	}

	lines := strings.Split(string(inData), "\n")

	if len(lines) < 14 {
		t.Errorf("Expecting lines to have len > 13. inData:\n%s", string(inData))
	}

	lines = lines[3:]

	if lines[0] != "10" || lines[1] != "11" || lines[2] != "9" || lines[3] != "8" ||
		lines[4] != "7" || lines[5] != "5" || lines[6] != "6" || lines[7] != "1" ||
		lines[8] != "3" || lines[9] != "2" || lines[10] != "4" {
		t.Errorf("generation failed: lines: %v", lines)
	}
}

func TestGenerateNotEnoughInputParams(t *testing.T) {
	pwd, err := os.Getwd()

	if err != nil {
		t.Error(err)
	}

	out := fmt.Sprintf("test_%d.gr8.tmp", time.Now().Unix())
	testFileOutput := filepath.Join(pwd, out)

	defer func() {
		os.Remove(testFileOutput)
	}()

	g := &gener8{
		skipFormat: true,
		kws:        "1,2",
		out:        out,
		inData: []byte(`
$kw1
$kw2
$kw3
$kw4
`),
	}

	g.generate()

	inData, err := ioutil.ReadFile(testFileOutput)

	if err != nil {
		t.Errorf("can't read output file: %s. Reason: %+v", testFileOutput, err)
	}

	lines := strings.Split(string(inData), "\n")

	if len(lines) < 8 {
		t.Errorf("Expecting lines to have len > 8. inData:\n%s", string(inData))
	}

	lines = lines[3:]

	if lines[0] != "1" || lines[1] != "2" || lines[2] != "$kw3" || lines[3] != "$kw4" {
		t.Errorf("generation failed: lines: %v", lines)
	}
}

func TestGenerateWayTooManyInputParams(t *testing.T) {
	pwd, err := os.Getwd()

	if err != nil {
		t.Error(err)
	}

	out := fmt.Sprintf("test_%d.gr8.tmp", time.Now().Unix())
	testFileOutput := filepath.Join(pwd, out)

	defer func() {
		os.Remove(testFileOutput)
	}()

	g := &gener8{
		skipFormat: true,
		kws:        "1,2,3,4,5,6,7,8,9,10",
		out:        out,
		inData: []byte(`
$kw1
$kw2
$kw3
$kw4
`),
	}

	g.generate()

	inData, err := ioutil.ReadFile(testFileOutput)

	if err != nil {
		t.Errorf("can't read output file: %s. Reason: %+v", testFileOutput, err)
	}

	lines := strings.Split(string(inData), "\n")

	if len(lines) < 8 {
		t.Errorf("Expecting lines to have len > 8. inData:\n%s", string(inData))
	}

	lines = lines[3:]

	if lines[0] != "1" || lines[1] != "2" || lines[2] != "3" || lines[3] != "4" {
		t.Errorf("generation failed: lines: %v", lines)
	}
}
