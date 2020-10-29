package util_test

import (
	"bytes"
	"io/ioutil"
	"neoway-case/util"
	"testing"
)

var parserTests = []struct {
	name string
	file string
}{
	{"wrong columns number", "test-files/parser1_test.txt"},
	{"private null", "test-files/parser2_test.txt"},
	{"wrong date format", "test-files/parser3_test.txt"},
}

func TestInvalidParse(t *testing.T) {
	for _, tt := range parserTests {
		buf, _ := ioutil.ReadFile(tt.file)
		reader := bytes.NewReader(buf)
		t.Run(tt.name, func(t *testing.T) {
			_, err := util.Parse(reader)
			if err == nil {
				t.Error("Expected err != nil")
			}
		})
	}
}
