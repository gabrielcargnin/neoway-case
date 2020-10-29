package util_test

import (
	"bytes"
	"io/ioutil"
	"neoway-case/util"
	"testing"
)

var validatorTests = []struct {
	name string
	file string
}{
	{"cpf invalid", "test-files/validator1_test.txt"},
	{"cnpj incorrect", "test-files/validator2_test.txt"},
	{"cpf null", "test-files/validator3_test.txt"},
}

func TestInvalidValidation(t *testing.T) {
	for _, tt := range validatorTests {
		buf, _ := ioutil.ReadFile(tt.file)
		reader := bytes.NewReader(buf)
		t.Run(tt.name, func(t *testing.T) {
			consumptions, err := util.Parse(reader)
			if err = util.Validate(consumptions); err == nil {
				t.Error("Expected validation error")
			}
		})
	}
}
