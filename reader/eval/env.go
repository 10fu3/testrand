package eval

import (
	"errors"
	"github.com/google/uuid"
	"testrand/cmap"
	"testrand/reader/infra"
)

func NewEnvironmentForClosure(parent *Sexpression, frame cmap.ConcurrentMap[string, *Sexpression]) (*Sexpression, error) {

	if SexpressionTypeEnvironment != parent._sexp_type_id {
		return CreateNil(), errors.New("not environment")
	}

	env := Sexpression{
		_sexp_type_id:  SexpressionTypeEnvironment,
		_env_frame:     frame,
		_env_parent:    parent,
		_env_globalEnv: parent._env_globalEnv,
		_env_parentId:  parent._env_parentId,
	}
	return &env, nil
}

func GetDefaultFunction() cmap.ConcurrentMap[string, *Sexpression] {
	defaultFuncs := cmap.New[*Sexpression]()

	defaultFuncs.Set("car", NewCar())
	defaultFuncs.Set("cdr", NewCdr())
	defaultFuncs.Set("and", NewAnd())
	defaultFuncs.Set("or", NewOr())
	defaultFuncs.Set("if", NewIf())
	defaultFuncs.Set("eq?", NewIsEq())
	defaultFuncs.Set("not", NewIsNot())
	defaultFuncs.Set("define", NewDefine())
	defaultFuncs.Set("set", NewSet())
	defaultFuncs.Set("loop", NewLoop())
	defaultFuncs.Set("wait", NewWait())
	defaultFuncs.Set("+", NewAdd())
	defaultFuncs.Set("-", NewSub())
	defaultFuncs.Set("*", NewMul())
	defaultFuncs.Set("/", NewDiv())
	defaultFuncs.Set("%", NewMod())
	defaultFuncs.Set("begin", NewBegin())
	defaultFuncs.Set("lambda", NewLambda())
	defaultFuncs.Set("quote", NewQuote())
	defaultFuncs.Set("quasiquote", NewQuasiquote())
	defaultFuncs.Set("heavy", NewHeavy())
	defaultFuncs.Set("print", NewPrint())
	defaultFuncs.Set("println", NewPrintln())
	defaultFuncs.Set("transaction", NewTransaction())
	defaultFuncs.Set("global-set", NewGlobalSet())
	defaultFuncs.Set("global-get", NewGlobalGet())
	defaultFuncs.Set("global-get-all", NewGlobalGetAll())
	defaultFuncs.Set("global-clear-all", NewGlobalClearAll())
	defaultFuncs.Set("cd", NewCurrentDirectory())
	defaultFuncs.Set("read-file-line", NewFileReadLine())
	defaultFuncs.Set("read-file", NewReadFile())
	defaultFuncs.Set("new-hashmap", NewNativeHashmap())
	defaultFuncs.Set("put-hashmap", NewPutNativeHashmap())
	defaultFuncs.Set("get-hashmap", NewGetNativeHashmap())
	defaultFuncs.Set("pair-loop-hashmap", NewKeyValuePairForeachNativeHashmap())
	defaultFuncs.Set("get-now-time-micro", NewGetNowTimeNano())
	defaultFuncs.Set("string-append", NewStringAppend())
	defaultFuncs.Set("string-split", NewStringSplit())
	defaultFuncs.Set("string-len", NewStringLen())
	defaultFuncs.Set("foreach", NewForeach())
	defaultFuncs.Set("apply", NewApply())
	defaultFuncs.Set("eval", NewEval())
	defaultFuncs.Set("interaction-environment", NewInteractionEnvironment())
	defaultFuncs.Set("this-environment", NewThisEnvironment())
	defaultFuncs.Set("let", NewLet())
	defaultFuncs.Set("string->symbol", NewStringToSymbol())
	defaultFuncs.Set("symbol->string", NewSymbolName())
	defaultFuncs.Set("symbol-name", NewSymbolName())
	defaultFuncs.Set("to-string", NewToString())
	defaultFuncs.Set("gc", NewForceGC())
	defaultFuncs.Set("new-array", NewNativeArray())
	defaultFuncs.Set("get-array", NewGetIndexNativeArray())
	defaultFuncs.Set("set-array", NewSetIndexNativeArray())
	defaultFuncs.Set("array-len", NewLengthNativeArray())
	defaultFuncs.Set("array-append", NewAppendNativeArray())
	defaultFuncs.Set("array->list", NewNativeArrayToList())
	defaultFuncs.Set("list->array", NewListToNativeArray())
	defaultFuncs.Set("foreach-array", NewForeachNativeArray())
	defaultFuncs.Set("void", NewVoid())

	return defaultFuncs
}

func NewGlobalEnvironment() (*Sexpression, error) {

	id := uuid.NewString()
	superGlobalEnv, err := infra.SetupEtcd(id)
	env := Sexpression{
		_sexp_type_id:       SexpressionTypeEnvironment,
		_env_parentId:       id,
		_env_frame:          GetDefaultFunction(),
		_env_parent:         nil,
		_env_globalEnv:      nil,
		_env_superGlobalEnv: superGlobalEnv,
	}
	env._env_globalEnv = &env

	return &env, err
}

func NewGlobalEnvironmentById(id string) (*Sexpression, error) {

	superGlobalEnv, err := infra.SetupEtcd(id)
	env := Sexpression{
		_sexp_type_id:       SexpressionTypeEnvironment,
		_env_frame:          GetDefaultFunction(),
		_env_parent:         nil,
		_env_parentId:       id,
		_env_globalEnv:      nil,
		_env_superGlobalEnv: superGlobalEnv,
	}
	env._env_globalEnv = &env

	return &env, err
}
