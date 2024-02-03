package main

import (
	"errors"
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func main() {

	source := "11111111222223333"
	// orParser := or(lit('2'), lit('4'))
	parser := and(lit('1'), lit('2'))
	fmt.Println(parser(source))
}

type parser func(string) (string, error)

func and(parser1, parser2 parser) parser {
	return func(input string) (string, error) {
		value1, err := parser1(input)
		if err != nil {
			return input, err
		}
		spew.Dump(value1)
		value2, err := parser2(value1)
		if err != nil {
			return input, err
		}
		return value2, nil
	}
}

func or(parser1, parser2 parser) parser {
	return func(input string) (string, error) {
		value1, err := parser1(input)
		if err == nil {
			return value1, nil
		}
		value2, err := parser2(input)
		if err == nil {
			return value2, nil
		}
		return input, err
	}
}

func lit(c rune) parser {
	return func(input string) (string, error) {
		for i, current := range input {
			if current == c {
				return input[i:], nil
			}
		}
		return input, errors.New("lit: could not find character")
	}
}
