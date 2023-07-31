package eval

import (
	"errors"
	"github.com/google/uuid"
	"testrand/reader/globalEnv"
)

type Environment interface {
	GetId() string
	GetValue(symbol Symbol) (SExpression, error)
	Define(symbol Symbol, sexp SExpression)
	Set(symbol Symbol, sexp SExpression) error
}

type environment struct {
	id     string
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

func (e *environment) GetId() string {
	return e.id
}

func NewEnvironment(parent Environment) Environment {
	env := &environment{
		id:     uuid.NewString(),
		frame:  map[string]SExpression{},
		parent: parent,
	}
	globalEnv.Put(env.id, env)
	return env
}

func NewGlobalEnvironment() Environment {
	env := &environment{
		frame: map[string]SExpression{
			"and":        NewAnd(),
			"or":         NewOr(),
			"if":         NewIf(),
			"eq?":        NewIsEq(),
			"not":        NewIsNot(),
			"define":     NewDefine(),
			"set":        NewSet(),
			"loop":       NewLoop(),
			"wait":       NewWait(),
			"+":          NewAdd(),
			"begin":      NewBegin(),
			"lambda":     NewLambda(),
			"quote":      NewQuote(),
			"quasiquote": NewQuasiquote(),
			"heavy":      NewHeavy(),
			"print":      NewPrintln(),
		},
		parent: nil,
	}
	globalEnv.Put(env.id, env)
	return env
}
