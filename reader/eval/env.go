package eval

import (
	"errors"
)

type Environment interface {
	GetValue(symbol Symbol) (SExpression, error)
	Define(symbol Symbol, sexp SExpression)
	Set(symbol Symbol, sexp SExpression) error
}

type environment struct {
	frame  map[string]SExpression
	parent Environment
}

func (e *environment) GetValue(symbol Symbol) (SExpression, error) {
	if value, ok := e.frame[symbol.GetValue()]; ok {
		return value, nil
	}
	if e.parent == nil {
		return nil, errors.New("UndefinedEvaluate")
	}
	return e.parent.GetValue(symbol)
}

func (e *environment) Define(symbol Symbol, sexp SExpression) {
	e.frame[symbol.GetValue()] = sexp
}

func (e *environment) Set(symbol Symbol, sexp SExpression) error {
	if _, ok := e.frame[symbol.GetValue()]; ok {
		e.frame[symbol.GetValue()] = sexp
		return nil
	}
	if e.parent == nil {
		return errors.New("UndefinedEvaluate")
	}
	return e.parent.Set(symbol, sexp)
}

func NewEnvironment(parent Environment) Environment {
	return &environment{
		frame:  map[string]SExpression{},
		parent: parent,
	}
}

func NewGlobalEnvironment() Environment {
	return &environment{
		frame: map[string]SExpression{
			"and":    NewAnd(),
			"or":     NewOr(),
			"if":     NewIf(),
			"eq?":    NewIsEq(),
			"define": NewDefine(),
		},
		parent: nil,
	}
}