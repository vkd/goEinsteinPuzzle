package goeinstein

import (
	"encoding/json"
	"fmt"
	"os"
)

type TableStorage struct {
	// table *Table
	mem map[string]Value
}

var _ Storage = (*TableStorage)(nil)

func NewTableStorage() *TableStorage {
	t := &TableStorage{
		mem: make(map[string]Value),
	}
	t.parse()
	// t.table = NewTableFile(t.GetFileName())
	return t
}

func (t *TableStorage) Close() {
	t.Flush()
}

func (t *TableStorage) GetFileName() string {
	return "./einstein/conf.cfg"
}

func (t *TableStorage) GetInt(name string, dflt int) int {
	if v, ok := t.mem[name]; ok {
		return v.AsInt()
	}
	return dflt
	// return t.table.GetInt(name, dflt)
}

func (t *TableStorage) GetString(name string, dflt string) string {
	if v, ok := t.mem[name]; ok {
		return v.AsString()
	}
	return dflt
	// return t.table.GetString(name, dflt)
}

func (t *TableStorage) SetInt(name string, value int) {
	t.mem[name] = NewIntValue(value)
	// t.table.SetInt(name, value)
}

func (t *TableStorage) SetString(name string, value string) {
	t.mem[name] = NewStringValue(value)
	// t.table.SetString(name, value)
}

func (t *TableStorage) Flush() {
	// t.table.Save(t.GetFileName())

	out := make(map[string]interface{})
	for k, v := range t.mem {
		switch v.GetType() {
		case IntegerType:
			out[k] = v.AsInt()
		case StringType:
			out[k] = v.AsString()
		default:
			panic(fmt.Sprintf("unknown type: %v", v.GetType()))
		}
	}

	bs, err := json.Marshal(out)
	if err != nil {
		panic(fmt.Errorf("marshal storage: %w", err))
	}

	err = os.WriteFile(t.GetFileName(), bs, 0664) //nolint:gofumpt
	if err != nil {
		panic(fmt.Errorf("write file storage: %w", err))
	}
}

func (t *TableStorage) parse() {
	bs, err := os.ReadFile(t.GetFileName())
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(fmt.Errorf("read storage file: %w", err))
	}

	var in map[string]interface{}
	err = json.Unmarshal(bs, &in)
	if err != nil {
		panic(fmt.Errorf("unmarshal storage file: %w", err))
	}

	for k, v := range in {
		switch v := v.(type) {
		case int:
			t.mem[k] = NewIntValue(v)
		case float64:
			t.mem[k] = NewIntValue(int(v))
		case string:
			t.mem[k] = NewStringValue(v)
		default:
			panic(fmt.Sprintf("unknown json type: %T", v))
		}
	}
}
