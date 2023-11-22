package eval

import (
	"bufio"
	"errors"
	"testrand/reader/lexer"
	"testrand/reader/token"
)

type reader struct {
	lexer.Lexer
	token.Token
	nestingLevel int
}

type Reader interface {
	Read() (*Sexpression, error)
}

func (r *reader) getCdr() (*Sexpression, error) {
	if r.Token.GetKind() == token.TokenKindRPAREN {
		return CreateEmptyList(), nil
	}
	if r.Token.GetKind() == token.TokenKindDot {
		nextToken, err := r.Lexer.GetNextToken()
		if err != nil {
			return CreateNil(), err
		}
		r.Token = nextToken
		sexp, err := r.sExpression()
		if err != nil {
			return CreateNil(), err
		}
		return sexp, nil
	}
	car, err := r.sExpression()
	if err != nil {
		return CreateNil(), err
	}
	cdr, err := r.getCdr()
	return CreateConsCell(car, cdr), nil
}

func (r *reader) sExpression() (*Sexpression, error) {
	if r.Token.GetKind() == token.TokenKindNumber {
		value := r.GetInt()
		if r.nestingLevel != 0 {
			nextToken, err := r.GetNextToken()
			if err != nil {
				return CreateNil(), err
			}
			r.Token = nextToken
		}
		return CreateInt(value), nil
	}

	if r.Token.GetKind() == token.TokenKindString {
		value := r.GetString()
		if r.nestingLevel != 0 {
			nextToken, err := r.GetNextToken()
			if err != nil {
				return CreateNil(), err
			}
			r.Token = nextToken
		}
		return CreateString(value), nil
	}

	if r.Token.GetKind() == token.TokenKindSymbol {
		value := r.GetSymbol()
		if r.nestingLevel != 0 {
			nextToken, err := r.GetNextToken()
			if err != nil {
				return CreateNil(), err
			}
			r.Token = nextToken
		}
		return CreateSymbol(value), nil
	}
	if r.Token.GetKind() == token.TokenKindBoolean {
		value := r.GetBool()
		if r.nestingLevel != 0 {
			nextToken, err := r.GetNextToken()
			if err != nil {
				return CreateNil(), err
			}
			r.Token = nextToken
		}
		return CreateBool(value), nil
	}
	if r.Token.GetKind() == token.TokenKindNil {
		if r.nestingLevel != 0 {
			nextToken, err := r.GetNextToken()
			if err != nil {
				return CreateNil(), err
			}
			r.Token = nextToken
		}
		return CreateNil(), nil
	}
	if r.Token.GetKind() == token.TokenKindQuote {
		nextToken, err := r.GetNextToken()
		if err != nil {
			return CreateNil(), err
		}
		r.Token = nextToken
		sexp, err := r.sExpression()
		if err != nil {
			return CreateNil(), err
		}
		return CreateConsCell(CreateSymbol("quote"), CreateConsCell(sexp, CreateConsCell(CreateNil(), CreateNil()))), nil
	}

	if r.Token.GetKind() == token.TokenKindUnquote {
		nextToken, err := r.GetNextToken()
		if err != nil {
			return CreateNil(), err
		}

		r.Token = nextToken
		sexp, err := r.sExpression()
		if err != nil {
			return CreateNil(), err
		}
		return CreateConsCell(CreateSymbol("unquote"), CreateConsCell(sexp, CreateConsCell(CreateNil(), CreateNil()))), nil
	}

	if r.Token.GetKind() == token.TokenKindUnquoteSplicing {
		nextToken, err := r.GetNextToken()
		if err != nil {
			return CreateNil(), err
		}
		r.Token = nextToken
		sexp, err := r.sExpression()
		if err != nil {
			return CreateNil(), err
		}
		return CreateConsCell(CreateSymbol("unquote-splicing"), CreateConsCell(sexp, CreateConsCell(CreateNil(), CreateNil()))), nil
	}

	if r.Token.GetKind() == token.TokenKindQuasiquote {
		nextToken, err := r.GetNextToken()
		if err != nil {
			return CreateNil(), err
		}
		r.Token = nextToken
		sexp, err := r.sExpression()
		if err != nil {
			return CreateNil(), err
		}
		return CreateConsCell(CreateSymbol("quasiquote"), CreateConsCell(sexp, CreateConsCell(CreateNil(), CreateNil()))), nil
	}
	if r.Token.GetKind() == token.TokenKindLparen {
		r.nestingLevel += 1
		nextToken, err := r.Lexer.GetNextToken()
		if err != nil {
			return CreateNil(), err
		}
		r.Token = nextToken
		if r.Token.GetKind() == token.TokenKindRPAREN {
			r.nestingLevel -= 1
			if r.nestingLevel != 0 {
				nextToken, err = r.Lexer.GetNextToken()
				if err != nil {
					return CreateNil(), err
				}
				r.Token = nextToken
			}
			return CreateConsCell(CreateNil(), CreateNil()), nil
		}
		car, err := r.sExpression()
		if err != nil {
			return CreateNil(), err
		}
		cdr, err := r.getCdr()
		if err != nil {
			return CreateNil(), err
		}
		r.nestingLevel -= 1
		if r.nestingLevel != 0 {
			nextToken, err = r.GetNextToken()
			if err != nil {
				return CreateNil(), err
			}
			r.Token = nextToken
		}
		return CreateConsCell(car, cdr), nil
	}
	return CreateNil(), errors.New("Invalid expression: " + r.Token.String())
}

func (r *reader) Read() (*Sexpression, error) {
	r.nestingLevel = 0
	t, err := r.Lexer.GetNextToken()
	if err != nil {
		return CreateNil(), err
	}
	r.Token = t
	return r.sExpression()
}

func New(in *bufio.Reader) Reader {
	return &reader{
		Lexer:        lexer.New(in),
		Token:        nil,
		nestingLevel: 0,
	}
}
