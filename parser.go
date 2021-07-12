package vgui

import (
    "fmt"
    . "github.com/ahmetb/go-linq"
    "strings"
)

type Value struct {
    Name   *Token
    String *Token
    Body   []*Value
}

type Parser struct {
    tokens  []*Token
    index   int
    current *Token
}

func NewParserFromTokens(tokens []*Token) (*Parser, error) {
    if len(tokens) <= 0 {
        return nil, &ErrEndOfInput{}
    }
    return &Parser{tokens, 0, tokens[0]}, nil
}

type ErrIncorrectToken struct {
    tokenTypes []TokenType
    received   TokenType
}

func (e *ErrIncorrectToken) Error() string {
    var tokenTypes []string
    From(e.tokenTypes).SelectT(func(t TokenType) string {
        return t.String()
    }).ToSlice(&tokenTypes)
    return fmt.Sprintf("expected one of (%v) but got a %v instead.", strings.Join(tokenTypes, ", "), e.received.String())
}

func (l *Parser) eatAny(types ...TokenType) (*Token, error) {
    if l.current == nil {
        return nil, &ErrIncorrectToken{types, TokenNil}
    }
    if From(types).AnyWithT(func (t TokenType) bool { return l.current.TokenType == t }) {
        return nil, &ErrIncorrectToken{types, l.current.TokenType}
    }

    current := l.current
    l.index++
    if len(l.tokens) > l.index {
        l.current = l.tokens[l.index]
    } else {
        l.current = nil
    }
    return current, nil
}

func (l *Parser) nextValueList() ([]*Value, error) {
    var values []*Value
    for l.current != nil && (l.current.TokenType == TokenName || l.current.TokenType == TokenString) {
        value, err := l.nextValue()
        if err != nil {
            return values, err
        }
        values = append(values, value)
    }
    return values, nil
}

func (l *Parser) nextValue() (*Value, error) {
    name, err := l.eatAny(TokenString, TokenName)
    if err != nil {
        return nil, err
    }
    if l.current != nil && l.current.TokenType == TokenLeftBrace {
        valueBody, err := l.nextValueList()
        if err != nil {
            return nil, err
        }
        return &Value{name, nil, valueBody}, nil
    }
    str, err := l.eatAny(TokenString)
    if err != nil {
        return nil, err
    }
    return &Value{name, str, nil}, nil
}

func (l *Parser) ParseRoot() (*Value, error) {
    valueBody, err := l.nextValueList()
    if err != nil {
        return nil, err
    }
    return &Value{&Token{TokenName, "Root"}, nil, valueBody}, nil
}