package main

import (
	"fmt"
	"strings"
)

type lexState int

const (
	starting lexState = iota
	normal
	singleQuoting
	doubleQuoting
	escaping
	escapingInQuote
)

const (
	spaceRune   = ' '
	singleQuote = '\''
	doubleQuote = '"'
	escapeRune  = '\\'
)

func tokenize(str string) (tokens []string) {
	str = strings.TrimSpace(str)
	state := starting
	var builder strings.Builder

	for idx, rune := range str {
		switch state {
		case starting:
			switch rune {
			case spaceRune:
				continue
			case singleQuote:
				state = singleQuoting
			case doubleQuote:
				state = doubleQuoting
			case escapeRune:
				state = escaping
			default:
				state = normal
				builder.WriteRune(rune)
			}
		case singleQuoting:
			if rune == singleQuote {
				state = normal
			} else {
				builder.WriteRune(rune)
			}

		case doubleQuoting:
			switch rune {
			case doubleQuote:
				state = normal
			case escapeRune:
				state = escapingInQuote
			default:
				builder.WriteRune(rune)
			}

		case normal:
			switch rune {
			case singleQuote:
				state = singleQuoting
			case doubleQuote:
				state = doubleQuoting
			case spaceRune:
				if str[idx-1] == spaceRune {
					// skip
				} else {
					tokens = append(tokens, builder.String())
					builder.Reset()
				}
			case escapeRune:
				state = escaping
			default:
				builder.WriteRune(rune)
			}

		case escaping:
			builder.WriteRune(rune)
			state = normal
		case escapingInQuote:
			builder.WriteRune(rune)
			state = doubleQuoting
		default:
			fmt.Println("Unexpected state: ", state)
		}
	}
	if builder.Len() > 0 {
		tokens = append(tokens, builder.String())
	}
	return tokens
}
