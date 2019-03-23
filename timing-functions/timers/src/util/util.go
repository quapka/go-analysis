package util

import (
	"encoding/csv"
	"errors"
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
