package reader

import (
	"bufio"
	"errors"
	"strconv"
	"testrand/reader/eval"
	"testrand/reader/lexer"
	"testrand/reader/token"
)

type reader struct {
	lexer.Lexer
	token.Token
	nestingLevel int
}

type Reader interface {
	Read() (eval.SExpression, error)
}

func (r *reader) getCdr() (eval.SExpression, error) {
	if r.Token.GetKind() == token.TokenKindRPAREN {
		return eval.NewConsCell(eval.NewNil(), eval.NewNil()), nil
	}
	if r.Token.GetKind() == token.TokenKindDot {
		nextToken, err := r.Lexer.GetNextToken()
		if err != nil {
			return nil, err
		}
		r.Token = nextToken
		sexp, err := r.sExpression()
		if err != nil {
			return nil, err
		}
		return sexp, nil
	}
	car, err := r.sExpression()
	if err != nil {
		return nil, err
	}
	cdr, err := r.getCdr()
	return eval.NewConsCell(car, cdr), nil
}

func (r *reader) sExpression() (eval.SExpression, error) {
	if r.Token.GetKind() == token.TokenKindNumber {
		value := r.GetInt()
		if r.nestingLevel != 0 {
			nextToken, err := r.GetNextToken()
			if err != nil {
				return nil, err
			}
			r.Token = nextToken
		}
		return eval.NewInt(value), nil
	}
	if r.Token.GetKind() == token.TokenKindSymbol {
		value := r.GetSymbol()
		if r.nestingLevel != 0 {
			nextToken, err := r.GetNextToken()
			if err != nil {
				return nil, err
			}
			r.Token = nextToken
		}
		return eval.NewSymbol(value), nil
	}
	if r.Token.GetKind() == token.TokenKindBoolean {
		value := r.GetBool()
		if r.nestingLevel != 0 {
			nextToken, err := r.GetNextToken()
			if err != nil {
				return nil, err
			}
			r.Token = nextToken
		}
		return eval.NewBool(value), nil
	}
	if r.Token.GetKind() == token.TokenKindQuote {
		nextToken, err := r.GetNextToken()
		if err != nil {
			return nil, err
		}
		r.Token = nextToken
		sexp, err := r.sExpression()
		if err != nil {
			return nil, err
		}
		return eval.NewConsCell(eval.NewSymbol("quote"), eval.NewConsCell(sexp, eval.NewNil())), nil
	}
	if r.Token.GetKind() == token.TokenKindLparen {
		r.nestingLevel += 1
		nextToken, err := r.Lexer.GetNextToken()
		if err != nil {
			return nil, err
		}
		r.Token = nextToken
		if r.Token.GetKind() == token.TokenKindRPAREN {
			r.nestingLevel -= 1
			if r.nestingLevel != 0 {
				nextToken, err = r.Lexer.GetNextToken()
				if err != nil {
					return nil, err
				}
				r.Token = nextToken
			}
			return eval.NewConsCell(eval.NewNil(), eval.NewNil()), nil
		}
		car, err := r.sExpression()
		if err != nil {
			return nil, err
		}
		cdr, err := r.getCdr()
		if err != nil {
			return nil, err
		}
		r.nestingLevel -= 1
		if r.nestingLevel != 0 {
			nextToken, err = r.GetNextToken()
			if err != nil {
				return nil, err
			}
			r.Token = nextToken
		}
		return eval.NewConsCell(car, cdr), nil
	}
	return nil, errors.New("Invalid expression: " + strconv.Itoa(int(r.GetKind())))
}

func (r *reader) Read() (eval.SExpression, error) {
	r.nestingLevel = 0
	t, err := r.Lexer.GetNextToken()
	if err != nil {
		return nil, err
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
