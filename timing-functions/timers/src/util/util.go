package util

import (
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"math/big"
	"strings"
)

// Reads file as csv separated with semicolon.
func LoadInputData(fileName string) ([][]string, error) {

	// read whole file into data
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	// new csv reader, separator = ;
	r := csv.NewReader(strings.NewReader(string(data)))
	r.Comma = ';'

	// read all
	result, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return result, err
}

// Converts a number represented as string into bytes.
// Error if the string does not represent a number in the given base.
func StringToIntBytes(input string, base int) ([]byte, error) {

	value := new(big.Int)

	value, ok := value.SetString(input, base)
	if !ok {
		return nil, errors.New("Cannot convert \"" + input + "\" into a number.")
	}

	return value.Bytes(), nil
}

// Constant reader always returns the same byte.
type constantReader struct {
	value byte
}

func NewConstantReader(value byte) io.Reader {
	reader := new(constantReader)
	reader.value = value
	return reader
}

func (r constantReader) Read(p []byte) (n int, err error) {

	for i := range p {
		p[i] = r.value
	}

	return len(p), nil
}
