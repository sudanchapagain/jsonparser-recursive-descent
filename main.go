package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type TokenType int

const (
	TokenString TokenType = iota
	TokenNumber
	TokenBoolean
	TokenNull
	TokenLeftBrace
	TokenRightBrace
	TokenLeftBracket
	TokenRightBracket
	TokenColon
	TokenComma
	TokenEOF
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input    string
	pos      int
	current  rune
}

func NewLexer(input string) *Lexer {
	lexer := &Lexer{input: input, pos: 0}
	lexer.advance()
	return lexer
}

func (l *Lexer) advance() {
	if l.pos < len(l.input) {
		l.current = rune(l.input[l.pos])
		l.pos++
	} else {
		l.current = 0
	}
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.current) {
		l.advance()
	}
}

func (l *Lexer) nextToken() Token {
	l.skipWhitespace()

	switch l.current {
	case '{':
		l.advance()
		return Token{Type: TokenLeftBrace, Value: "{"}
	case '}':
		l.advance()
		return Token{Type: TokenRightBrace, Value: "}"}
	case '[':
		l.advance()
		return Token{Type: TokenLeftBracket, Value: "["}
	case ']':
		l.advance()
		return Token{Type: TokenRightBracket, Value: "]"}
	case ':':
		l.advance()
		return Token{Type: TokenColon, Value: ":"}
	case ',':
		l.advance()
		return Token{Type: TokenComma, Value: ","}
	case '"':
		return l.readString()
	default:
		if unicode.IsDigit(l.current) || l.current == '-' {
			return l.readNumber()
		} else if unicode.IsLetter(l.current) {
			return l.readKeyword()
		}
	}

	return Token{Type: TokenEOF, Value: ""}
}

func (l *Lexer) readString() Token {
	var sb strings.Builder
	l.advance()

	for l.current != '"' && l.current != 0 {
		sb.WriteRune(l.current)
		l.advance()
	}
	l.advance()

	return Token{Type: TokenString, Value: sb.String()}
}

func (l *Lexer) readNumber() Token {
	start := l.pos - 1
	for unicode.IsDigit(l.current) || l.current == '.' || l.current == '-' {
		l.advance()
	}
	value := l.input[start : l.pos-1]
	return Token{Type: TokenNumber, Value: value}
}

func (l *Lexer) readKeyword() Token {
	start := l.pos - 1
	for unicode.IsLetter(l.current) {
		l.advance()
	}
	value := l.input[start : l.pos-1]

	switch value {
	case "true", "false":
		return Token{Type: TokenBoolean, Value: value}
	case "null":
		return Token{Type: TokenNull, Value: value}
	}
	panic("Unexpected keyword: " + value)
}

type Parser struct {
	lexer *Lexer
	token Token
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.token = p.lexer.nextToken()
}

func (p *Parser) parseJSON() interface{} {
	switch p.token.Type {
	case TokenLeftBrace:
		return p.parseObject()
	case TokenLeftBracket:
		return p.parseArray()
	default:
		panic("Invalid JSON start")
	}
}

func (p *Parser) parseObject() map[string]interface{} {
	obj := make(map[string]interface{})
	p.nextToken()

	for p.token.Type != TokenRightBrace {
		if p.token.Type != TokenString {
			panic("Expected string key in object")
		}
		key := p.token.Value
		p.nextToken()

		if p.token.Type != TokenColon {
			panic("Expected ':' after key")
		}
		p.nextToken()

		value := p.parseValue()
		obj[key] = value

		if p.token.Type == TokenComma {
			p.nextToken()
		} else if p.token.Type != TokenRightBrace {
			panic("Expected ',' or '}' in object")
		}
	}

	p.nextToken()
	return obj
}

func (p *Parser) parseArray() []interface{} {
	arr := []interface{}{}
	p.nextToken()

	for p.token.Type != TokenRightBracket {
		arr = append(arr, p.parseValue())

		if p.token.Type == TokenComma {
			p.nextToken()
		} else if p.token.Type != TokenRightBracket {
			panic("Expected ',' or ']' in array")
		}
	}

	p.nextToken()
	return arr
}

func (p *Parser) parseValue() interface{} {
	switch p.token.Type {
	case TokenString:
		val := p.token.Value
		p.nextToken()
		return val
	case TokenNumber:
		val, _ := strconv.ParseFloat(p.token.Value, 64)
		p.nextToken()
		return val
	case TokenBoolean:
		val := p.token.Value == "true"
		p.nextToken()
		return val
	case TokenNull:
		p.nextToken()
		return nil
	case TokenLeftBrace:
		return p.parseObject()
	case TokenLeftBracket:
		return p.parseArray()
	default:
		panic("Unexpected token: " + p.token.Value)
	}
}

func main() {
	jsonInput := `{
		"name": "nepal",
		"age": 0,
		"country": true,
		"districts": ["Kathmandu", "Lalitpur"],
		"address": { "continent": "Asia", "Location": "South Asia" }
	}`

	lexer := NewLexer(jsonInput)
	parser := NewParser(lexer)

	parsedData := parser.parseJSON()
	fmt.Println(parsedData)
}
