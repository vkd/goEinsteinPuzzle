package goeinstein

import (
	"fmt"
	"os"
)

type Type int

const (
	IntegerType Type = iota
	DoubleType
	StringType
	TableType
)

type Value interface {
	Close()
	AsTable() *Table
	GetType() Type
	AsInt() int
	AsDouble() float32
	AsString() string
	Clone() Value
}

type IntValue struct {
	value int
}

var _ Value = IntValue{}

func NewIntValue(val int) IntValue   { return IntValue{val} }
func (IntValue) Close()              {}
func (v IntValue) GetType() Type     { return IntegerType }
func (v IntValue) AsInt() int        { return v.value }
func (v IntValue) AsDouble() float32 { return float32(v.value) }
func (v IntValue) AsString() string  { return ToString(v.value) }
func (v IntValue) AsTable() *Table   { panic("Can't convert integer to table") }
func (v IntValue) Clone() Value      { return NewIntValue(v.value) }

type DoubleValue struct {
	value float32
}

var _ Value = DoubleValue{}

func NewDoubleValue(val float32) DoubleValue { return DoubleValue{val} }
func (DoubleValue) Close()                   {}
func (v DoubleValue) GetType() Type          { return DoubleType }
func (v DoubleValue) AsInt() int             { return int(v.value) }
func (v DoubleValue) AsDouble() float32      { return v.value }
func (v DoubleValue) AsString() string       { return ToString(v.value) }
func (v DoubleValue) AsTable() *Table        { panic("Can't convert double to table") }
func (v DoubleValue) Clone() Value           { return NewDoubleValue(v.value) }

type StringValue struct {
	value string
}

var _ Value = StringValue{}

func NewStringValue(val string) StringValue { return StringValue{val} }
func (StringValue) Close()                  {}
func (v StringValue) GetType() Type         { return StringType }
func (v StringValue) AsInt() int            { return StrToInt(v.value) }
func (v StringValue) AsDouble() float32     { return StrToDouble(v.value) }
func (v StringValue) AsString() string      { return v.value }
func (v StringValue) AsTable() *Table       { panic("Can't convert string to table") }
func (v StringValue) Clone() Value          { return NewStringValue(v.value) }

type TableValue struct {
	value *Table
}

var _ Value = TableValue{}

func NewTableValue(val *Table) TableValue { return TableValue{val} }
func (TableValue) Close()                 {}
func (v TableValue) GetType() Type        { return TableType }
func (v TableValue) AsInt() int           { panic("Can't convert table to int") }
func (v TableValue) AsDouble() float32    { panic("Can't convert table to double") }
func (v TableValue) AsString() string     { panic("Can't convert table to string") }
func (v TableValue) AsTable() *Table      { return v.value }
func (v TableValue) Clone() Value         { return NewTableValue(NewTableTable(v.value)) }

type Table struct {
	fields         map[string]Value
	lastArrayIndex int
}

func NewTableTable(table *Table) *Table {
	t := &Table{}
	t.Assign(table)
	return t
}

func NewTableFile(fileName string) *Table {
	t := &Table{
		fields: map[string]Value{},
	}
	t.lastArrayIndex = 0
	bs, err := os.ReadFile(fileName)
	if err != nil {
		panic(fmt.Errorf("Error opening file %q: %w", fileName, err)) //nolint:stylecheck
	}
	reader := NewUtfStreamReader(bs)
	lexal := NewLexal(reader)
	t.Parse(lexal, false, 0, 0)
	return t
}

func NewTableLexal(lexal *Lexal, line, pos int) *Table {
	t := &Table{}
	t.lastArrayIndex = 0
	t.Parse(lexal, true, line, pos)
	return t
}

func NewTable() *Table {
	t := &Table{}
	t.lastArrayIndex = 0
	return t
}

func (t *Table) Close() {}

func (t *Table) Assign(table *Table) {
	if t == table {
		return
	}

	t.fields = make(map[string]Value)
	t.lastArrayIndex = table.lastArrayIndex
	for k, v := range table.fields {
		t.fields[k] = v.Clone()
	}
}

func LexToValue(lexal *Lexal, lexeme *Lexeme) Value {
	panic("not implemented")
}

func (t *Table) GetInt(name string, dflt int) int {
	panic("not implemented")
}

func (t *Table) GetString(name string, dflt string) string {
	panic("not implemented")
}

func (t *Table) SetInt(name string, dflt int) int {
	panic("not implemented")
}

func (t *Table) SetString(name string, dflt string) string {
	panic("not implemented")
}

func (t *Table) Save(fileName string) {
	panic("not implemented")
}

func (t *Table) Parse(i interface{}, b bool, line, pos int) {
	panic("not implemented")
}

func NewUtfStreamReader(bs []byte) interface{} {
	panic("not implemented")
}

func NewLexal(i interface{}) interface{} {
	panic("not implemented")
}

type Lexal interface{}

type Lexeme interface{}
