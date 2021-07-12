package vgui

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestTokenizesTokenSequence(t *testing.T) {
    input := "{}[Flag]\"String\"Name"
    lexer, _ := NewLexerFromInput(input)
    tokens := lexer.GetTokens()
    assert.Len(t, tokens, 5)
    assert.Equal(t, tokens[0].TokenType, TokenLeftBrace)
    assert.Equal(t, tokens[1].TokenType, TokenRightBrace)
    assert.Equal(t, tokens[2].TokenType, TokenFlag)
    assert.Equal(t, tokens[3].TokenType, TokenString)
    assert.Equal(t, tokens[4].TokenType, TokenName)
    assert.Equal(t, tokens[0].Value, "{")
    assert.Equal(t, tokens[1].Value, "}")
    assert.Equal(t, tokens[2].Value, "Flag")
    assert.Equal(t, tokens[3].Value, "String")
    assert.Equal(t, tokens[4].Value, "Name")
}

func TestTokenizesNonTrivialSequenceWithWhitespace(t *testing.T) {
    input := "{\n" +
        "  }       [Flag]" +
        "\n\n\"String\"\n\n" +
        "  Name  "
    lexer, _ := NewLexerFromInput(input)
    tokens := lexer.GetTokens()
    assert.Len(t, tokens, 5)
    assert.Equal(t, tokens[0].TokenType, TokenLeftBrace)
    assert.Equal(t, tokens[1].TokenType, TokenRightBrace)
    assert.Equal(t, tokens[2].TokenType, TokenFlag)
    assert.Equal(t, tokens[3].TokenType, TokenString)
    assert.Equal(t, tokens[4].TokenType, TokenName)
    assert.Equal(t, tokens[0].Value, "{")
    assert.Equal(t, tokens[1].Value, "}")
    assert.Equal(t, tokens[2].Value, "Flag")
    assert.Equal(t, tokens[3].Value, "String")
    assert.Equal(t, tokens[4].Value, "Name")
}