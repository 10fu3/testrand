package eval

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type SExpression interface {
	TypeId() string
	AtomId() SExpressionType
	String() string
	IsList() bool
	Equals(sexp SExpression) bool
}

type Subroutine interface {
	SExpression
	Apply(arg SExpression)
}

type Symbol interface {
	SExpression
	GetValue() string
}

type symbol struct {
	name string
}

func (s *symbol) AtomId() SExpressionType {
	return SExpressionTypeSymbol
}

func (s *symbol) TypeId() string {
	return "symbol"
}

func (s *symbol) IsList() bool {
	return false
}

func (s *symbol) String() string {
	return s.name
}

func (s *symbol) Equals(sexp SExpression) bool {
	if sexp.TypeId() != "symbol" {
		return false
	}
	return s.name == (sexp).(Symbol).GetValue()
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

func (i *_int) AtomId() SExpressionType {
	return SExpressionTypeNumber
}

func (i *_int) TypeId() string {
	return "number"
}

func (i *_int) IsList() bool {
	return false
}

func (i *_int) Equals(sexp SExpression) bool {
	if "number" != sexp.TypeId() {
		return false
	}
	return i.GetValue() == sexp.(Number).GetValue()
}

type Number interface {
	GetValue() int64
	String() string
	SExpression
}

func NewInt(val int64) Number {
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

func (b *_bool) Equals(sexp SExpression) bool {
	if "bool" != sexp.TypeId() {
		return false
	}

	return b.Value == sexp.(Bool).GetValue()
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

func (b *_bool) TypeId() string {
	return "bool"
}

func (b *_bool) AtomId() SExpressionType {
	return SExpressionTypeBool
}

func (b *_bool) IsList() bool {
	return false
}

var trueSexp = &_bool{Value: true}
var falseSexp = &_bool{Value: false}

func NewBool(b bool) Bool {
	if b {
		return trueSexp
	}
	return falseSexp
}

type Str interface {
	GetValue() string
	String() string
	SExpression
}

type _string struct {
	Value string
}

func NewString(s string) Str {
	return &_string{Value: s}
}

func (s *_string) Equals(sexp SExpression) bool {
	if "string" != sexp.TypeId() {
		return false
	}
	return s.Value == sexp.(Str).GetValue()
}

func (s *_string) GetValue() string {
	return s.Value
}

func (s *_string) String() string {
	return fmt.Sprintf("\"%s\"", s.Value)
}

func (s *_string) TypeId() string {
	return "string"
}

func (s *_string) AtomId() SExpressionType {
	return SExpressionTypeString
}

func (s *_string) IsList() bool {
	return false
}

type Nil interface {
	SExpression
}

type _nil struct {
}

func (n *_nil) Equals(sexp SExpression) bool {
	return "nil" == sexp.TypeId()
}

func (n *_nil) TypeId() string {
	return "nil"
}

func (n *_nil) AtomId() SExpressionType {
	return SExpressionTypeNil
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
}

type _cons_cell struct {
	Car SExpression
	Cdr SExpression
}

func (cell *_cons_cell) Equals(sexp SExpression) bool {
	if "cons_cell" != sexp.TypeId() {
		return false
	}
	c := sexp.(ConsCell)
	if !cell.Car.Equals(c.GetCar()) {
		return false
	}
	return !cell.Cdr.Equals(c.GetCdr())
}

func NewConsCell(car SExpression, cdr SExpression) ConsCell {
	return &_cons_cell{
		Car: car,
		Cdr: cdr,
	}
}

func JoinList(left, right SExpression) (ConsCell, error) {

	if !left.IsList() {
		return nil, errors.New("left is not a list")
	}

	if !right.IsList() {
		return nil, errors.New("right is not a list")
	}

	baseRoot := left.(*_cons_cell)
	baseLook := baseRoot

	copyRoot := &_cons_cell{
		Car: NewNil(),
		Cdr: NewNil(),
	}
	copyLook := copyRoot

	for {
		if IsEmptyList(baseLook.GetCdr()) {
			copyLook.Car = baseLook.GetCar()
			copyLook.Cdr = right
			return copyRoot, nil
		} else {
			copyLook.Car = baseLook.GetCar()
			copyLook.Cdr = NewConsCell(NewNil(), NewNil())
			copyLook = copyLook.Cdr.(*_cons_cell)
			baseLook = baseLook.GetCdr().(*_cons_cell)
		}
	}
}

func (cell *_cons_cell) TypeId() string {
	return "cons_cell"
}

func (cell *_cons_cell) AtomId() SExpressionType {
	return SExpressionTypeConsCell
}

func (cell *_cons_cell) String() string {
	if "symbol" == cell.Car.TypeId() && "quote" == ((cell.Car).(Symbol)).GetValue() && cell.Cdr.TypeId() == "cons_cell" && "nil" == ((cell.Cdr).(ConsCell)).GetCdr().TypeId() {
		return fmt.Sprintf("'%s", ((cell.Cdr).(ConsCell)).GetCar().String())
	}
	var joinedString strings.Builder
	joinedString.WriteString("(")
	var lookCell ConsCell = cell

	for {
		if lookCell.GetCar().TypeId() != "nil" {
			joinedString.WriteString(lookCell.GetCar().String())
			if lookCell.GetCdr().TypeId() == "cons_cell" {
				if lookCell.GetCdr().(ConsCell).GetCar().TypeId() != "nil" && lookCell.GetCdr().(ConsCell).GetCdr().TypeId() != "nil" {
					joinedString.WriteString(" ")
				}
			}
		}

		if lookCell.GetCdr().TypeId() != "cons_cell" {
			if lookCell.GetCdr().TypeId() != "nil" {
				joinedString.WriteString(" . " + lookCell.GetCdr().String())
			}
			joinedString.WriteString(")")
			break
		}
		lookCell = (lookCell.GetCdr()).(ConsCell)
	}
	return joinedString.String()
}

func _toArray(sexp ConsCell) ([]SExpression, uint64, error) {
	var list []SExpression
	look := sexp

	var count = uint64(0)
	for {
		if SExpressionTypeConsCell != look.AtomId() {
			return nil, 0, errors.New("need list")
		}
		if SExpressionTypeNil == look.GetCar().AtomId() && SExpressionTypeNil == look.GetCdr().AtomId() {
			break
		}
		if count < uint64(cap(list)) {
			list = list[:count+1] // slice の延長
			list[count] = look.(ConsCell).GetCar()
		} else if cap(list) < len(list)+1 {
			list = append(list, look.(ConsCell).GetCar(), nil)
		}

		look = look.GetCdr().(ConsCell)
		count++
	}
	return list, count, nil
}

func ToArray(sexp SExpression) ([]SExpression, uint64, error) {
	if sexp.AtomId() != SExpressionTypeConsCell {
		return nil, 0, errors.New("need list")
	}
	return _toArray(sexp.(ConsCell))
}

func (cell *_cons_cell) IsList() bool {
	if "cons_cell" == cell.Cdr.TypeId() {
		if IsEmptyList(cell.Cdr) {
			return true
		}
		return cell.Cdr.IsList()
	}
	return false
}

func (cell *_cons_cell) GetCar() SExpression {
	return cell.Car
}

func (cell *_cons_cell) GetCdr() SExpression {
	return cell.Cdr
}

func ToConsCell(list []SExpression, argsLength uint64) ConsCell {
	var cells = make([]_cons_cell, argsLength+1)
	var look = &cells[0]

	for i, sexp := range list {
		look.Car = sexp
		look.Cdr = &(cells[i+1])
		if uint64(i) == argsLength-1 {
			look.Cdr = NewConsCell(NewNil(), NewNil())
		}
		look = (look.Cdr).(*_cons_cell)
	}

	return &cells[0]
}

func IsEmptyList(list SExpression) bool {
	if "cons_cell" != list.TypeId() {
		return false
	}
	tmp := (list).(ConsCell)

	return "nil" == tmp.GetCar().TypeId() && "nil" == tmp.GetCdr().TypeId()
}

type SExpressionType int

const (
	SExpressionTypeSymbol SExpressionType = iota
	SExpressionTypeNumber
	SExpressionTypeBool
	SExpressionTypeString
	SExpressionTypeNil
	SExpressionTypeConsCell
	SExpressionTypeSubroutine
	SExpressionTypeSpecialForm
	SExpressionTypeClosure
	SExpressionTypeNativeHashmap
	SExpressionTypeNativeArray
	SExpressionTypeEnvironment
)
