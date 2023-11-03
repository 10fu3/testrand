package eval

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"testrand/reader/globalEnv"
	"testrand/reader/infra"
)

type Environment interface {
	GetId() string
	GetValue(symbol Symbol) (SExpression, error)
	GetGlobalEnv() Environment
	GetSuperGlobalEnv() infra.ISuperGlobalEnv
	GetParentId() string
	Define(symbol Symbol, sexp SExpression)
	Set(symbol Symbol, sexp SExpression) error
	SExpression
}

type environment struct {
	id             string
	frame          map[string]SExpression
	parent         Environment
	globalEnv      Environment
	superGlobalEnv infra.ISuperGlobalEnv
	parentId       string
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

func (e *environment) GetGlobalEnv() Environment {
	return e.globalEnv
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

func (e *environment) GetSuperGlobalEnv() infra.ISuperGlobalEnv {
	return e.superGlobalEnv
}

func (e *environment) GetParentId() string {
	return e.parentId
}

func (e *environment) Type() string {
	return "environment"
}

func (e *environment) String() string {
	return fmt.Sprintf("#<environment %s>", e.id)
}

func (e *environment) IsList() bool {
	return false
}

func (e *environment) Equals(args SExpression) bool {
	if "environment" != args.Type() {
		return false
	}
	return e.id == args.(*environment).id
}

func NewEnvironment(parent Environment) (Environment, error) {
	env := &environment{
		id:             uuid.NewString(),
		frame:          map[string]SExpression{},
		parent:         parent,
		globalEnv:      parent.GetGlobalEnv(),
		superGlobalEnv: parent.GetSuperGlobalEnv(),
		parentId:       parent.GetParentId(),
	}
	globalEnv.Put(env.id, env)
	return env, nil
}

func GetDefaultFunction() map[string]SExpression {
	return map[string]SExpression{
		"car":                     NewCar(),
		"cdr":                     NewCdr(),
		"and":                     NewAnd(),
		"or":                      NewOr(),
		"if":                      NewIf(),
		"eq?":                     NewIsEq(),
		"not":                     NewIsNot(),
		"define":                  NewDefine(),
		"set":                     NewSet(),
		"loop":                    NewLoop(),
		"wait":                    NewWait(),
		"+":                       NewAdd(),
		"begin":                   NewBegin(),
		"lambda":                  NewLambda(),
		"quote":                   NewQuote(),
		"quasiquote":              NewQuasiquote(),
		"heavy":                   NewHeavy(),
		"print":                   NewPrint(),
		"println":                 NewPrintln(),
		"transaction":             NewTransaction(),
		"global-set":              NewGlobalSet(),
		"global-get":              NewGlobalGet(),
		"global-get-all":          NewGlobalGetAll(),
		"cd":                      NewCurrentDirectory(),
		"read-file-line":          NewFileReadLine(),
		"new-hashmap":             NewNativeHashmap(),
		"put-hashmap":             NewPutNativeHashmap(),
		"get-hashmap":             NewGetNativeHashmap(),
		"pair-loop-hashmap":       NewKeyValuePairNativeHashmap(),
		"get-now-time-micro":      NewGetNowTimeMicro(),
		"string-append":           NewStringAppend(),
		"string-split":            NewStringSplit(),
		"string-len":              NewStringLen(),
		"foreach":                 NewForeach(),
		"apply":                   NewApply(),
		"eval":                    NewEval(),
		"interaction-environment": NewInteractionEnvironment(),
		"this-environment":        NewThisEnvironment(),
		"let":                     NewLet(),
		"string->symbol":          NewStringToSymbol(),
		"symbol->string":          NewSymbolName(),
		"symbol-name":             NewSymbolName(),
		"to-string":               NewToString(),
	}
}

func NewGlobalEnvironment() (Environment, error) {

	id := uuid.NewString()
	superGlobalEnv, err := infra.SetupEtcd(id)
	env := &environment{
		id:             id,
		parentId:       id,
		frame:          GetDefaultFunction(),
		parent:         nil,
		globalEnv:      nil,
		superGlobalEnv: superGlobalEnv,
	}
	env.globalEnv = env
	globalEnv.Put(env.id, env)
	return env, err
}

func NewGlobalEnvironmentById(id string) (Environment, error) {

	superGlobalEnv, err := infra.SetupEtcd(id)

	env := &environment{
		id:             id,
		frame:          GetDefaultFunction(),
		parentId:       id,
		parent:         nil,
		globalEnv:      nil,
		superGlobalEnv: superGlobalEnv,
	}
	env.globalEnv = env
	globalEnv.Put(env.id, env)
	return env, err
}
