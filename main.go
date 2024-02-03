package main

import "fmt"

func main() {

	source := "11111111222223333"
	twoParser := lit('2')
	fourParser := lit('4')

	parser := or(twoParser, fourParser)
	fmt.Println(parser(source))
}

type parser func(string) string

func or(parser1, parser2 parser) parser {
	return func(s string) string {
		value1 := parser1(s)
		value2 := parser2(s)

		if value1 != s {
			return value1
		}
		return value2
	}
}

func lit(c rune) parser {
	return func(input string) string {
		found := 0
		for i, current := range input {
			if current == c {
				found = i
				break
			}
		}
		return input[found:]
	}
}
