package vgui

import (
    "strings"
    "unicode"
)

type TokenType int

const (
    TokenNil TokenType = iota
    TokenLeftBrace
    TokenRightBrace
    TokenFlag
    TokenString
    TokenName
)

type runePredicate func(value rune) bool

type Token struct {
    TokenType TokenType
    Value     string
}

type Lexer struct {
    input   []rune
    index   int
    current rune
}

func (t TokenType) String() string {
    return [...]string{"{", "}", "Flag", "String", "Name"}[t]
}

func NewToken(tokenType TokenType, value string) *Token {
    return &Token{tokenType, value}
}

func NewLexerFromInput(input string) (*Lexer, error) {
    if len(input) <= 0 {
        return nil, &ErrEmptyInput{}
    }
    runes := []rune(input)
    return &Lexer{runes, 0, runes[0]}, nil
}

func (l *Lexer) peek(n int) (rune, error) {
    if len(l.input) > l.index+n {
        return l.input[l.index+n], nil
    }
    return rune(0), &ErrEndOfInput{}
}

func (l *Lexer) advance() {
    l.index += 1
    current, _ := l.peek(0)
    l.current = current
}

func (l *Lexer) skipWhitespace() {
    for current, err := l.peek(0); err == nil && unicode.IsSpace(current); current, err = l.peek(0) {
        l.advance()
    }
}

func (l *Lexer) skipLine() {
    for current, err := l.peek(0); err == nil && current == '\n'; current, err = l.peek(0) {
        l.advance()
    }
}

func (l *Lexer) getStringToken() *Token {
    token := l.getGenericToken(TokenString, func(value rune) bool {
        return value != '"'
    })
    l.advance()
    return token
}

func (l *Lexer) getNameToken() *Token {
    return l.getGenericToken(TokenName, func(value rune) bool {
        return unicode.IsDigit(value) || unicode.IsLetter(value)
    })
}

func (l *Lexer) getFlagToken() *Token {
    token := l.getGenericToken(TokenFlag, func(value rune) bool {
        return value != ']'
    })
    l.advance()
    return token
}

func (l *Lexer) getGenericToken(tokenType TokenType, predicate runePredicate) *Token {
    value := strings.Builder{}
    for current, err := l.peek(0); err == nil && predicate(current); current, err = l.peek(0) {
        value.WriteRune(current)
        l.advance()
    }

    return NewToken(tokenType, value.String())
}

func (l *Lexer) GetTokens() []*Token {
    var tokens []*Token
    for len(l.input) > l.index {
        l.skipWhitespace()

        // skip over comments
        if next, err := l.peek(1); l.current == '/' && err != nil && next == '/' {
            l.skipLine()
            continue
        }

        // skip over single line directives, these are handled by the preprocessor
        if l.current == '#' {
            l.skipLine()
            continue
        }

        if l.current == '"' {
            l.advance()
            tokens = append(tokens, l.getStringToken())
            continue
        }

        if unicode.IsDigit(l.current) || unicode.IsLetter(l.current) {
            tokens = append(tokens, l.getNameToken())
            continue
        }

        if l.current == '[' {
            l.advance()
            tokens = append(tokens, l.getFlagToken())
            continue
        }

        if l.current == '{' {
            l.advance()
            tokens = append(tokens, NewToken(TokenLeftBrace, "{"))
            continue
        }

        if l.current == '}' {
            l.advance()
            tokens = append(tokens, NewToken(TokenRightBrace, "}"))
            continue
        }
    }
    return tokens
}