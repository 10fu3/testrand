package eval

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"testrand/reader/infra"
)

type Environment interface {
	GetValue(symbol Symbol) (SExpression, error)
	GetGlobalEnv() Environment
	GetSuperGlobalEnv() infra.ISuperGlobalEnv
	GetParentId() string
	Define(symbol Symbol, sexp SExpression)
	Set(symbol Symbol, sexp SExpression) error
	SExpression
}

type environment struct {
	frame          map[string]SExpression
	parent         Environment
	globalEnv      Environment
	superGlobalEnv infra.ISuperGlobalEnv
	parentId       string
	mutex          sync.RWMutex
}

func (e *environment) GetValue(symbol Symbol) (SExpression, error) {
	e.mutex.RLock()
	if value, ok := e.frame[symbol.GetValue()]; ok {
		e.mutex.RUnlock()
		return value, nil
	}
	e.mutex.RUnlock()
	if e.parent == nil {
		return nil, errors.New("UndefinedEvaluate")
	}
	return e.parent.GetValue(symbol)
}

func (e *environment) GetGlobalEnv() Environment {
	return e.globalEnv
}

func (e *environment) Define(symbol Symbol, sexp SExpression) {
	e.mutex.Lock()
	e.frame[symbol.GetValue()] = sexp
	e.mutex.Unlock()
}

func (e *environment) Set(symbol Symbol, sexp SExpression) error {
	e.mutex.RLock()
	if _, ok := e.frame[symbol.GetValue()]; ok {
		e.mutex.RUnlock()
		e.mutex.Lock()
		e.frame[symbol.GetValue()] = sexp
		e.mutex.Unlock()
		return nil
	}
	e.mutex.RUnlock()
	if e.parent == nil {
		return errors.New("UndefinedEvaluate")
	}
	return e.parent.Set(symbol, sexp)
}

func (e *environment) GetSuperGlobalEnv() infra.ISuperGlobalEnv {
	return e.superGlobalEnv
}

func (e *environment) GetParentId() string {
	return e.parentId
}

func (e *environment) TypeId() string {
	return "environment"
}

func (e *environment) AtomId() AtomType {
	return AtomTypeEnvironment
}

func (e *environment) String() string {
	return fmt.Sprintf("#<environment>")
}

func (e *environment) IsList() bool {
	return false
}

func (e *environment) Equals(args SExpression) bool {
	if "environment" != args.TypeId() {
		return false
	}
	return args.(*environment) == e
}

func NewEnvironment(parent Environment) (Environment, error) {
	env := &environment{
		frame:          map[string]SExpression{},
		parent:         parent,
		globalEnv:      parent.GetGlobalEnv(),
		superGlobalEnv: parent.GetSuperGlobalEnv(),
		parentId:       parent.GetParentId(),
	}
	return env, nil
}

func NewEnvironmentForClosure(parent Environment, frame map[string]SExpression) (Environment, error) {
	env := &environment{
		frame:          frame,
		parent:         parent,
		globalEnv:      parent.GetGlobalEnv(),
		superGlobalEnv: parent.GetSuperGlobalEnv(),
		parentId:       parent.GetParentId(),
	}
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
		"-":                       NewMinus(),
		"*":                       NewMultiply(),
		"/":                       NewDivide(),
		"mod":                     NewMod(),
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
		"global-clear-all":        NewGlobalClearAll(),
		"cd":                      NewCurrentDirectory(),
		"read-file-line":          NewFileReadLine(),
		"read-file":               NewFileRead(),
		"new-hashmap":             NewNativeHashmap(),
		"put-hashmap":             NewPutNativeHashmap(),
		"get-hashmap":             NewGetNativeHashmap(),
		"hashmap->list":           NewKeyValuePairNativeHashmapToConsCell(),
		"pair-loop-hashmap":       NewKeyValuePairNativeHashmap(),
		"get-now-time-micro":      NewGetNowTimeNano(),
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
		"gc":                      NewForceGC(),
		"new-array":               NewNativeArray(),
		"get-array":               NewGetNativeArray(),
		"set-array":               NewSetNativeArray(),
		"array-len":               NewLengthNativeArray(),
		"array-append":            NewAppendNativeArray(),
		"array->list":             NewNativeArrayToList(),
		"list->array":             NewListToNativeArray(),
		"foreach-array":           NewForeachNativeArray(),
		"void":                    NewVoid(),
	}
}

func NewGlobalEnvironment() (Environment, error) {

	id := uuid.NewString()
	superGlobalEnv, err := infra.SetupEtcd(id)
	env := &environment{
		parentId:       id,
		frame:          GetDefaultFunction(),
		parent:         nil,
		globalEnv:      nil,
		superGlobalEnv: superGlobalEnv,
	}
	env.globalEnv = env

	return env, err
}

func NewGlobalEnvironmentById(id string) (Environment, error) {

	superGlobalEnv, err := infra.SetupEtcd(id)

	env := &environment{
		frame:          GetDefaultFunction(),
		parentId:       id,
		parent:         nil,
		globalEnv:      nil,
		superGlobalEnv: superGlobalEnv,
	}
	env.globalEnv = env
	return env, err
}
