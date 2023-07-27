package eval

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type SExpression interface {
	Type() string
	String() string
	IsList() bool
}

type Atom interface {
	SExpression
	Equals(atom Atom) bool
}

type Subroutine interface {
	SExpression
	Apply(arg SExpression)
}

type Symbol interface {
	Atom
	GetValue() string
}

type symbol struct {
	name string
}

func (s *symbol) Type() string {
	return "symbol"
}

func (s *symbol) IsList() bool {
	return false
}

func (s *symbol) String() string {
	return s.name
}

func (s *symbol) Equals(atom Atom) bool {
	if atom.Type() != "atom" {
		return false
	}
	return s.name == (atom).(Symbol).GetValue()
}

func (s *symbol) GetValue() string {
	return s.name
}

func NewSymbol(sym string) Symbol {
	return &symbol{name: sym}
}

type _int struct {
	Value int64
}

func (i *_int) GetValue() int64 {
	return i.Value
}

func (i *_int) String() string {
	return strconv.FormatInt(i.Value, 10)
}

func (i *_int) Type() string {
	return "number"
}

func (i *_int) IsList() bool {
	return false
}

type Int interface {
	GetValue() int64
	String() string
	SExpression
}

func NewInt(val int64) Int {
	return &_int{
		Value: val,
	}
}

type Bool interface {
	GetValue() bool
	String() string
	SExpression
}

type _bool struct {
	Value bool
}

func (b *_bool) GetValue() bool {
	return b.Value
}

func (b *_bool) String() string {
	if b.Value {
		return "#t"
	}
	return "#f"
}

func (b *_bool) Type() string {
	return "bool"
}

func (b *_bool) IsList() bool {
	return false
}

func NewBool(b bool) Bool {
	return &_bool{Value: b}
}

type Nil interface {
	SExpression
}

type _nil struct {
}

func (n *_nil) Type() string {
	return "nil"
}

func (n *_nil) String() string {
	return "#nil"
}

func (n *_nil) IsList() bool {
	return false
}

func NewNil() Nil {
	return &_nil{}
}

type ConsCell interface {
	SExpression
	GetCar() SExpression
	GetCdr() SExpression
	ToArray(sexp SExpression) ([]SExpression, error)
}

type _cons_cell struct {
	Car SExpression
	Cdr SExpression
}

func NewConsCell(car SExpression, cdr SExpression) ConsCell {
	return &_cons_cell{
		Car: car,
		Cdr: cdr,
	}
}

func (cell *_cons_cell) Type() string {
	return "cons_cell"
}

func (cell *_cons_cell) String() string {
	if "symbol" == cell.Car.Type() && "quote" == ((cell.Car).(Symbol)).GetValue() && cell.Cdr.Type() == "cons_cell" && "nil" == ((cell.Cdr).(ConsCell)).GetCdr().Type() {
		return fmt.Sprintf("'%s", ((cell.Cdr).(ConsCell)).GetCar().String())
	}
	var joinedString strings.Builder
	joinedString.WriteString("(")
	var lookCell ConsCell = cell

	for {
		if lookCell.GetCar().Type() != "nil" {
			joinedString.WriteString(lookCell.GetCar().String())
			if lookCell.GetCdr().Type() == "cons_cell" {
				if lookCell.GetCdr().(ConsCell).GetCar().Type() != "nil" && lookCell.GetCdr().(ConsCell).GetCdr().Type() != "nil" {
					joinedString.WriteString(" ")
				}
			}
		}

		if lookCell.GetCdr().Type() != "cons_cell" {
			if lookCell.GetCdr().Type() != "nil" {
				joinedString.WriteString(" . " + lookCell.GetCdr().String())
			}
			joinedString.WriteString(")")
			break
		}
		lookCell = (lookCell.GetCdr()).(ConsCell)
	}
	return joinedString.String()
}

func (cell *_cons_cell) ToArray(sexp SExpression) ([]SExpression, error) {
	var cons = sexp
	var list []SExpression
	var temp ConsCell = nil
	for cons.Type() != "nil" {
		if cons.Type() != "cons_cell" {
			return nil, errors.New("type error: " + cons.Type())
		}
		temp = (cons).(ConsCell)
		list = append(list, temp.GetCar())
		cons = temp.GetCdr()
	}
	return list, nil
}

func (cell *_cons_cell) IsList() bool {
	return cell.GetCdr().IsList()
}

func (cell *_cons_cell) GetCar() SExpression {
	return cell.Car
}

func (cell *_cons_cell) GetCdr() SExpression {
	return cell.Cdr
}

func ToConsCell(list []SExpression) ConsCell {
	var head = (NewConsCell(NewNil(), NewNil())).(*_cons_cell)
	var look = head
	var beforeLook *_cons_cell = nil

	for _, sexp := range list {
		look.Car = sexp
		look.Cdr = NewConsCell(NewNil(), NewNil())
		beforeLook = look
		look = (look.Cdr).(*_cons_cell)
	}
	if beforeLook != nil {
		beforeLook.Cdr = NewNil()
	}
	return head
}
