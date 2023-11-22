package eval

import (
	"context"
	"errors"
	"fmt"
	"github.com/dustinxie/lockfree"
	"github.com/dustinxie/lockfree/hashmap"
	"strconv"
	"strings"
	"testrand/cmap"
	"testrand/reader/infra"
)

var _symbolTable = lockfree.NewHashMap(hashmap.BucketSizeOption(18))

type _cell struct {
	_car *Sexpression
	_cdr *Sexpression
}

type Sexpression struct {
	_symbol       *Sexpression
	_boolean      bool
	_number       int64
	_string       string
	_sexp_type_id SexpressionType

	_apply_env          *Sexpression
	_apply_id           string
	_apply_body         *Sexpression
	_apply_formals      []*Sexpression
	_apply_formalsCount uint64
	_applyFunc          func(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error)

	env *Sexpression

	_env_frame          cmap.ConcurrentMap[string, *Sexpression]
	_env_parent         *Sexpression
	_env_parentId       string
	_env_globalEnv      *Sexpression
	_env_superGlobalEnv infra.ISuperGlobalEnv

	_cell *_cell

	_native_value interface{}

	_native_arr interface{}

	_native_map cmap.ConcurrentMap[string, interface{}]
}
type CallableConstructorArgs struct {
	Id           string
	FormalsCount int
	Fn           func(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error)
	Env          *Sexpression
}

func CreateNil() *Sexpression {
	return &Sexpression{_sexp_type_id: SexpressionTypeNil}
}

func CreateEmptyList() *Sexpression {
	return &Sexpression{_sexp_type_id: SexpressionTypeConsCell, _cell: &_cell{_car: &Sexpression{}, _cdr: &Sexpression{}}}
}

func CreateSymbol(symbol string) *Sexpression {
	if sexp, ok := _symbolTable.Get(symbol); ok {
		return sexp.(*Sexpression)
	}

	createdSymbol := &Sexpression{_sexp_type_id: SexpressionTypeSymbol, _symbol: &Sexpression{_sexp_type_id: SexpressionTypeString, _string: symbol}}

	_symbolTable.Set(symbol, createdSymbol)

	return createdSymbol
}

func Quote() *Sexpression {
	return CreateSymbol("quote")
}

func GetSymbol(symbol string) *Sexpression {
	return CreateSymbol(symbol)
}

func CreateInt(number int64) *Sexpression {
	return &Sexpression{_sexp_type_id: SexpressionTypeNumber, _number: number}
}

func CreateBool(boolean bool) *Sexpression {
	return &Sexpression{_sexp_type_id: SexpressionTypeBool, _boolean: boolean}
}

func CreateString(str string) *Sexpression {
	return &Sexpression{_sexp_type_id: SexpressionTypeString, _string: str}
}

func CreateConsCell(car *Sexpression, cdr *Sexpression) *Sexpression {
	return &Sexpression{_sexp_type_id: SexpressionTypeConsCell, _cell: &_cell{_car: car, _cdr: cdr}}
}

func CreateClosure(body *Sexpression,
	formals []*Sexpression,
	env *Sexpression,
	formalsCount uint64) (*Sexpression, error) {

	if env._sexp_type_id != SexpressionTypeEnvironment {
		return &Sexpression{}, errors.New("not environment")
	}

	if body._sexp_type_id != SexpressionTypeConsCell {
		return &Sexpression{}, errors.New("not cons cell")
	}

	return &Sexpression{
		_sexp_type_id:       SexpressionTypeClosure,
		_apply_body:         body,
		_apply_formals:      formals,
		_apply_env:          env,
		_apply_formalsCount: formalsCount,
		_applyFunc:          _closure_run,
	}, nil
}

func CreateSubroutine(name string, fn func(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error)) *Sexpression {
	return &Sexpression{_sexp_type_id: SexpressionTypeSubroutine, _apply_id: "subr " + name, _applyFunc: fn}
}

func CreateSpecialForm(name string, fn func(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error)) *Sexpression {
	return &Sexpression{_sexp_type_id: SexpressionTypeSpecialForm, _apply_id: "syntax " + name, _applyFunc: fn}
}

func CreateNativeArray(arr interface{}) *Sexpression {
	return &Sexpression{_native_arr: arr, _sexp_type_id: SexpressionTypeNativeArray}
}

func CreateNativeHashmap(m cmap.ConcurrentMap[string, interface{}]) *Sexpression {
	return &Sexpression{_native_map: m, _sexp_type_id: SexpressionTypeNativeHashmap}
}

func CreateNativeValue(v interface{}) *Sexpression {
	return &Sexpression{_native_value: v, _sexp_type_id: SexpressionTypeNativeValue}
}

func (s *Sexpression) SexpressionTypeId() SexpressionType {
	return s._sexp_type_id
}

func (s *Sexpression) IsList() bool {
	if s._sexp_type_id != SexpressionTypeConsCell {
		return false
	}
	if s._cell == nil {
		return false
	}
	look := s
	for {
		if look._cell._cdr._sexp_type_id == SexpressionTypeNil {
			return true
		}
		if look._cell._cdr._sexp_type_id != SexpressionTypeConsCell {
			return false
		}
		look = look._cell._cdr
	}
}

func (s *Sexpression) String() string {
	switch s._sexp_type_id {
	case SexpressionTypeNumber:
		return strconv.FormatInt(s._number, 10)
	case SexpressionTypeBool:
		if s._boolean {
			return "#t"
		}
		return "#f"
	case SexpressionTypeString:
		return fmt.Sprintf("\"%s\"", s._string)
	case SexpressionTypeNil:
		return "#nil"
	case SexpressionTypeConsCell:
		if SexpressionTypeSymbol == s._cell._car._sexp_type_id {
			if s._cell._car._symbol == Quote() &&
				s._cell._cdr._sexp_type_id == SexpressionTypeConsCell &&
				s._cell._cdr._cell._cdr._sexp_type_id == SexpressionTypeNil {
				return fmt.Sprintf("'%s", s._cell._cdr._cell._car.String())
			}
		}
		var joinedString strings.Builder
		joinedString.WriteString("(")
		var lookCell *_cell = s._cell

		for {
			if lookCell._car._sexp_type_id != SexpressionTypeNil {
				joinedString.WriteString(lookCell._car.String())
				if lookCell._cdr._sexp_type_id == SexpressionTypeConsCell {
					if lookCell._cdr._cell._car._sexp_type_id != SexpressionTypeNil && lookCell._cdr._cell._cdr._sexp_type_id != SexpressionTypeNil {
						joinedString.WriteString(" ")
					}
				}
			}
			if lookCell._cdr._sexp_type_id != SexpressionTypeConsCell {
				if lookCell._cdr._sexp_type_id != SexpressionTypeNil {
					joinedString.WriteString(" . " + lookCell._cdr.String())
				}
				joinedString.WriteString(")")
				break
			}
			lookCell = lookCell._cdr._cell
		}
		return joinedString.String()
	case SexpressionTypeSubroutine:
		return fmt.Sprintf("#<subr %s>", s._apply_id)
	case SexpressionTypeSpecialForm:
		return fmt.Sprintf("#<syntax %s>", s._apply_id)
	case SexpressionTypeClosure:
		return "#<closure>"
	case SexpressionTypeNativeHashmap:
		return fmt.Sprintf("#<native-hashmap>")
	case SexpressionTypeNativeArray:
		return fmt.Sprintf("#<native-array>")
	case SexpressionTypeEnvironment:
		return "#<environment>"

	default:
		return s._symbol._string
	}
}

func (s *Sexpression) Equals(sexp *Sexpression) bool {
	switch s._sexp_type_id {
	case SexpressionTypeNumber:
		return s._number == sexp._number
	case SexpressionTypeBool:
		return s._boolean == sexp._boolean
	case SexpressionTypeString:
		return s._string == sexp._string
	case SexpressionTypeNil:
		return s._sexp_type_id == sexp._sexp_type_id
	case SexpressionTypeConsCell:
		if sexp._sexp_type_id != SexpressionTypeConsCell {
			return false
		}
		return s._cell._car.Equals(sexp._cell._car) && s._cell._cdr.Equals(sexp._cell._cdr)
	case SexpressionTypeSubroutine:
		return s._apply_id == sexp._apply_id
	case SexpressionTypeSpecialForm:
		return s._apply_id == sexp._apply_id
	case SexpressionTypeClosure:
		panic("not implemented")
	case SexpressionTypeNativeHashmap:
		panic("not implemented")
	case SexpressionTypeNativeArray:
		panic("not implemented")
	case SexpressionTypeEnvironment:
		panic("not implemented")
	default:
		return s._symbol == sexp._symbol
	}
}

func (s *Sexpression) IsCallable() bool {
	return s._sexp_type_id == SexpressionTypeSubroutine || s._sexp_type_id == SexpressionTypeSpecialForm || s._sexp_type_id == SexpressionTypeClosure
}

func (s *Sexpression) GetFormalsCount() uint64 {
	return s._apply_formalsCount
}

func (s *Sexpression) GetEnvParentId() string {
	return s._env_parentId
}

func (s *Sexpression) IsNumber() bool {
	return s._sexp_type_id == SexpressionTypeNumber
}

func (s *Sexpression) IsString() bool {
	return s._sexp_type_id == SexpressionTypeString
}

func (s *Sexpression) IsBool() bool {
	return s._sexp_type_id == SexpressionTypeBool
}

func (s *Sexpression) IsSymbol() bool {
	return s._sexp_type_id == SexpressionTypeSymbol
}

func (s *Sexpression) IsNil() bool {
	return s._sexp_type_id == SexpressionTypeNil
}

func (s *Sexpression) IsConsCell() bool {
	return s._sexp_type_id == SexpressionTypeConsCell
}

func (s *Sexpression) IsSubroutine() bool {
	return s._sexp_type_id == SexpressionTypeSubroutine
}

func (s *Sexpression) IsSpecialForm() bool {
	return s._sexp_type_id == SexpressionTypeSpecialForm
}

func (s *Sexpression) IsClosure() bool {
	return s._sexp_type_id == SexpressionTypeClosure
}

func (s *Sexpression) IsNativeHashmap() bool {
	return s._sexp_type_id == SexpressionTypeNativeHashmap
}

func (s *Sexpression) IsNativeArray() bool {
	return s._sexp_type_id == SexpressionTypeNativeArray
}

func (s *Sexpression) IsEnvironment() bool {
	return s._sexp_type_id == SexpressionTypeEnvironment
}

func (s *Sexpression) IsNativeValue() bool {
	return s._sexp_type_id == SexpressionTypeNativeValue
}

func (s *Sexpression) GetValueFromFrame(key *Sexpression) (*Sexpression, bool) {
	if !s.IsEnvironment() {
		return CreateNil(), false
	}

	if !key.IsSymbol() {
		return CreateNil(), false
	}

	v, ok := s._env_frame.Get(key._symbol._string)

	if !ok {
		parentFrame := s._env_parent
		for {
			if parentFrame == nil {
				return CreateNil(), false
			}
			parentGet, parentOk := parentFrame._env_frame.Get(key._symbol._string)
			if !parentOk {
				parentFrame = parentFrame._env_parent
				continue
			}
			return parentGet, true
		}
	}
	return v, true
}

func JoinList(left, right *Sexpression) (*Sexpression, error) {

	if !left.IsList() {
		return &Sexpression{}, errors.New("left is not a list")
	}

	if !right.IsList() {
		return &Sexpression{}, errors.New("right is not a list")
	}

	baseRoot := left
	baseLook := baseRoot

	copyRoot := &Sexpression{
		_cell: &_cell{
			_car: left,
			_cdr: right,
		},
	}

	copyLook := copyRoot

	for {
		var cdr = baseLook._cell._cdr
		if IsEmptyList(cdr) {
			copyLook._cell._car = baseLook._cell._car
			copyLook._cell._cdr = right
			return copyRoot, nil
		} else {
			copyLook._cell._car = baseLook._cell._car
			//copyLook.Car = baseLook.GetCar()

			copyLook._cell._cdr = &Sexpression{
				_cell: &_cell{
					_car: &Sexpression{},
					_cdr: &Sexpression{},
				},
			}
			copyLook = copyLook._cell._cdr
			baseLook = baseLook._cell._cdr
		}
	}
}

func ToArray(sexp *Sexpression) ([]*Sexpression, uint64, error) {

	if sexp._sexp_type_id != SexpressionTypeConsCell {
		return []*Sexpression{}, 0, errors.New("need list")
	}

	list := make([]*Sexpression, 0, 0)
	look := sexp
	var tail *Sexpression
	var tailCell *_cell
	size := uint64(0)
	for {
		if SexpressionTypeConsCell != look._sexp_type_id {
			return nil, 0, errors.New("need list")
		}
		tail = look
		if SexpressionTypeConsCell != tail._sexp_type_id {
			return nil, 0, errors.New("need list")
		}
		tailCell = tail._cell
		if SexpressionTypeNil == tailCell._car._sexp_type_id && SexpressionTypeNil == tailCell._cdr._sexp_type_id {
			break
		}
		list = append(list, look._cell._car)
		look = look._cell._cdr
		size++
	}
	return list, size, nil
}

func ToConsCell(list []*Sexpression) *Sexpression {
	var head = &Sexpression{_sexp_type_id: SexpressionTypeConsCell, _cell: &_cell{
		_car: &Sexpression{},
		_cdr: &Sexpression{},
	}}
	var look = head

	for _, sexp := range list {
		look._cell._car = sexp
		look._cell._cdr = &Sexpression{_sexp_type_id: SexpressionTypeConsCell, _cell: &_cell{
			_car: &Sexpression{},
			_cdr: &Sexpression{},
		}}
		look = look._cell._cdr
	}
	if look._cell._cdr._sexp_type_id != SexpressionTypeNil {
		look._cell._cdr = &Sexpression{_sexp_type_id: SexpressionTypeNil}
	}

	return head
}

func IsEmptyList(list *Sexpression) bool {

	if SexpressionTypeConsCell != list.SexpressionTypeId() {
		return false
	}

	return SexpressionTypeNil == list._cell._car._sexp_type_id && SexpressionTypeNil == list._cell._cdr._sexp_type_id
}

type SexpressionType int

const (
	SexpressionTypeNil SexpressionType = iota
	SexpressionTypeSymbol
	SexpressionTypeNumber
	SexpressionTypeBool
	SexpressionTypeString
	SexpressionTypeConsCell
	SexpressionTypeSubroutine
	SexpressionTypeSpecialForm
	SexpressionTypeClosure
	SexpressionTypeNativeHashmap
	SexpressionTypeNativeArray
	SexpressionTypeEnvironment
	SexpressionTypeNativeValue
)

var SexpressionTypeNames = map[SexpressionType]string{
	SexpressionTypeNil:           "nil",
	SexpressionTypeSymbol:        "symbol",
	SexpressionTypeNumber:        "number",
	SexpressionTypeBool:          "bool",
	SexpressionTypeString:        "string",
	SexpressionTypeConsCell:      "cons_cell",
	SexpressionTypeSubroutine:    "subroutine",
	SexpressionTypeSpecialForm:   "special_form",
	SexpressionTypeClosure:       "closure",
	SexpressionTypeNativeHashmap: "native_hashmap",
	SexpressionTypeNativeArray:   "native_array",
	SexpressionTypeEnvironment:   "environment",
	SexpressionTypeNativeValue:   "native_value",
}

func (s SexpressionType) String() string {
	return SexpressionTypeNames[s]
}
