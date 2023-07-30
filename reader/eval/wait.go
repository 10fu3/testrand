package eval

import (
	"errors"
	"time"
)

type _wait struct{}

func (_ *_wait) Type() string {
	return "special_form.wait"
}

func (_ *_wait) String() string {
	return "#<syntax wait>"
}

func (_ *_wait) IsList() bool {
	return false
}

func (l *_wait) Equals(sexp SExpression) bool {
	return l.Type() == sexp.Type()
}

func (_ *_wait) Apply(env Environment, args SExpression) (SExpression, error) {
	if "cons_cell" != args.Type() {
		return nil, errors.New("need arguments")
	}
	arguments := args.(ConsCell)

	if !IsEmptyList(arguments.GetCdr()) {
		return nil, errors.New("need arguments length is 1")
	}

	waitTime, err := Eval(arguments.GetCar(), env)
	if err != nil {
		return nil, err
	}

	if waitTime.Type() != "number" {
		return nil, errors.New("need 1st argument must be number but got " + waitTime.Type())
	}
	durationTime := time.Millisecond * time.Duration(int(waitTime.(Number).GetValue()))
	time.Sleep(durationTime)

	return nil, nil
}

func NewWait() SExpression {
	return &_wait{}
}
