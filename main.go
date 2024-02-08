package main

import (
	"errors"
	"fmt"
)

func main() {
	source := "falsed and true is that"
	val, err := boolean([]rune(source))
	if err != nil {
		panic(err)
	}
	fmt.Printf("parsed boolean: %t", val)
}

func boolean(input []rune) (bool, error) {
	p := any(parseString("true"), parseString("false"))
	r, err := p(input)
	if err != nil {
		return false, err
	}
	return string(r.parsed) == "true", nil
}

// result is what a parser returns
type result struct {
	parsed    []rune
	remaining []rune
}

func (r result) String() string {
	return fmt.Sprintf("Parsed: [%s], Remaining: [%s]", string(r.parsed), string(r.remaining))
}

// parser is what a parser should look like
type parser func(input []rune) (result, error)

// many tries to parse input until failiure
func many(p1 parser) parser {
	return func(input []rune) (result, error) {
		aggregate := result{parsed: []rune{}, remaining: input}
		var err error
		for {
			result1, err := p1(aggregate.remaining)
			if err != nil {
				break
			}
			aggregate.remaining = result1.remaining
			aggregate.parsed = append(aggregate.parsed, result1.parsed...)
		}
		if len(aggregate.parsed) == 0 {
			return aggregate, err
		}
		return aggregate, nil
	}
}

// and parses input first with p1 and then with p2
func and(p1, p2 parser) parser {
	return func(input []rune) (result, error) {
		result1, err := p1(input)
		if err != nil {
			return result{parsed: []rune{}, remaining: input}, errors.New("and: failed parsing")
		}
		result2, err := p2(result1.remaining)
		if err != nil {
			return result{parsed: []rune{}, remaining: input}, errors.New("and: failed parsing")
		}
		result2.parsed = append(result1.parsed, result2.parsed...)
		return result2, nil
	}
}

// seq parses input following the sequence of parsers from parsers list
func seq(parsers ...parser) parser {
	return func(input []rune) (result, error) {
		aggregate := result{parsed: []rune{}, remaining: input}
		for _, p := range parsers {
			r, err := p(aggregate.remaining)
			if err != nil {
				return result{parsed: []rune{}, remaining: input}, err
			}
			aggregate.remaining = r.remaining
			aggregate.parsed = append(aggregate.parsed, r.parsed...)
		}
		return aggregate, nil
	}
}

// any parses input with the first parser that suceeds from list
func any(parsers ...parser) parser {
	return func(input []rune) (result, error) {
		for _, p := range parsers {
			result, err := p(input)
			if err == nil {
				return result, nil
			}
		}
		return result{parsed: []rune{}, remaining: input}, errors.New("any: failed parsing")
	}
}

// parseString parses input with matching string
func parseString(s string) parser {
	sRunes := []rune(s)
	return func(input []rune) (result, error) {
		if len(sRunes) > len(input) {
			return result{parsed: []rune{}, remaining: input}, errors.New("parseString: input shorter than parse string")
		}
		for i, r := range sRunes {
			if r != input[i] {
				return result{parsed: []rune{}, remaining: input}, errors.New("parseString: string failed parsing - bad input")
			}
		}
		return result{parsed: sRunes, remaining: input[len(sRunes):]}, nil
	}
}

// parseChar parses input with c
func parseChar(c rune) parser {
	return func(input []rune) (result, error) {
		if c == input[0] {
			return result{parsed: input[0:1], remaining: input[1:]}, nil
		}
		return result{parsed: []rune{c}, remaining: input}, errors.New("parseChar: character failed parsing - bad input")
	}
}
