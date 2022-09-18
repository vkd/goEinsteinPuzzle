package goeinstein

type TokenType int8

//nolint:golint
const (
	Word TokenType = iota
	Para
	Eof
)

type Token struct {
	tokenType TokenType
	content   string
}

func NewToken(tp TokenType) *Token {
	t := &Token{}
	t.tokenType = tp
	return t
}

func NewTokenContent(tp TokenType, content string) *Token {
	t := &Token{
		content: content,
	}
	t.tokenType = tp
	return t
}

func (t *Token) GetType() TokenType { return t.tokenType }
func (t *Token) GetContent() string { return t.content }

func (t *Token) ToString() string {
	//nolint:gocritic
	if Word == t.tokenType {
		return "Word: '" + t.content + "'"
	} else if Para == t.tokenType {
		return "Para"
	} else if Eof == t.tokenType {
		return "Eof"
	} else {
		return "Unknown"
	}
}

func WIsSpace(ch byte) bool {
	return ' ' == ch || '\n' == ch || '\r' == ch || '\t' == ch || '\f' == ch || '\v' == ch
}

type Tokenizer struct {
	text       string
	currentPos int
	stack      []*Token
}

func NewTokenizer(s string) *Tokenizer {
	t := &Tokenizer{
		text: s,
	}
	t.currentPos = 0
	return t
}

func (t *Tokenizer) SkipSpaces(notSearch bool) bool {
	ln := len(t.text)
	var foundDoubleReturn bool
	for ln > t.currentPos && WIsSpace(t.text[t.currentPos]) {
		t.currentPos++
		if !notSearch && '\n' == t.text[t.currentPos-1] && t.currentPos < ln && '\n' == t.text[t.currentPos] {
			notSearch = true
			foundDoubleReturn = true
		}
	}
	return foundDoubleReturn
}

func (t *Tokenizer) GetNextToken() *Token {
	if 0 < len(t.stack) {
		token := t.stack[0]
		t.stack = t.stack[1:]
		return token
	}
	ln := len(t.text)
	if t.SkipSpaces(t.currentPos == 0) && t.currentPos < ln {
		return NewToken(Para)
	}
	if t.currentPos >= ln {
		return NewToken(Eof)
	}
	wordStart := t.currentPos
	for ln > t.currentPos && !WIsSpace(t.text[t.currentPos]) {
		t.currentPos++
	}
	return NewTokenContent(Word, t.text[wordStart:t.currentPos])
}

func (t *Tokenizer) Unget(token *Token) {
	t.stack = append(t.stack, token)
}

func (t *Tokenizer) IsFinished() bool {
	if 0 < len(t.stack) {
		return false
	}
	return t.currentPos >= len(t.text)
}
