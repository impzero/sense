package main

import (
	"errors"
	"fmt"
)

func main() {
	source := "hello world!"
	digits := "111122223333"

	helloParser := parseString("hello")
	helloAndSpaceParser := and(helloParser, parseChar(' '))
	parsed, err := helloAndSpaceParser([]rune(source))
	manyOneParser := many(parseChar('1'))
	parsedDigits, err := manyOneParser([]rune(digits))
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	fmt.Printf("parsed: %v\n", parsed)
	fmt.Printf("parsed 1s: %v", parsedDigits)
}

// result is what a parser returns
type result struct {
	parsed    []rune
	remaining []rune
}

func (r result) String() string {
	return fmt.Sprintf("Parsed: [%s], Remaining: [%s]", string(r.parsed), string(r.remaining))
}

type parser func(input []rune) (result, error)

func many(p1 parser) parser {
	return func(input []rune) (result, error) {
		aggregate := result{parsed: []rune{}, remaining: input}
		for {
			result1, err := p1(aggregate.remaining)
			if err != nil {
				break
			}
			aggregate.remaining = result1.remaining
			aggregate.parsed = append(aggregate.parsed, result1.parsed...)
		}
		return aggregate, nil
	}
}

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

func parseChar(c rune) parser {
	return func(input []rune) (result, error) {
		if c == input[0] {
			return result{parsed: input[0:1], remaining: input[1:]}, nil
		}
		return result{parsed: []rune{c}, remaining: input}, errors.New("parseChar: character failed parsing - bad input")
	}
}
