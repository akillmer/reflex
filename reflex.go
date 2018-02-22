package reflex

import (
	"reflect"
	"strings"
)

// Model interface
type Model struct {
	name   string
	kind   reflect.Kind
	typ    reflect.Type
	val    reflect.Value
	fields map[string]*Model
}

// NewModel of the passed interface value
func NewModel(v interface{}) *Model {
	var m = &Model{
		typ:    reflect.TypeOf(v),
		val:    reflect.ValueOf(v),
		fields: make(map[string]*Model),
	}

	// dereference pointers to get an settable/addressable value
	for ; m.typ.Kind() == reflect.Ptr; m.typ = m.typ.Elem() {
		m.val = m.val.Elem()
	}

	if m.val.CanSet() == false {
		panic("NewModel's interface is not settable")
	} else if m.val.CanAddr() == false {
		panic("NewModel's interface is not addressable")
	}

	m.kind = m.typ.Kind()

	if m.kind == reflect.Slice {
		// m.kind is still `slice`, but now m.typ will indicate the slice's type
		m.typ = m.typ.Elem()
	} else if m.kind == reflect.Struct {
		// map all of the struct's fields to m.fields and set internal pointers
		m.addFields()
	}

	return m
}

// returns a new value of the model's type
func (m *Model) new() *Model {
	var v = reflect.New(m.typ).Interface()
	return NewModel(v)
}

// recursively adds struct fields to model.fields[...] and its nested structs
func (m *Model) addFields() {
	for i := 0; i < m.typ.NumField(); i++ {
		var (
			field  = m.typ.Field(i)
			v      = reflect.New(field.Type).Interface()
			target = NewModel(v)
		)

		target.name = field.Name
		m.fields[target.name] = target
		target.val = m.val.FieldByIndex(field.Index)

		if target.kind == reflect.Struct {
			target.addFields()
		}
	}
}

// get a nested model field (see above) by key
func (m *Model) getField(key string) *Model {
	var (
		target   = m
		splitKey = strings.Split(key, ".")
	)

	for _, k := range splitKey {
		if model, ok := target.fields[k]; ok {
			target = model
		}
	}

	return target
}

// Set a new value by key, where key can point to any nested member delimited by dots,
// e.g. "Movie.Name". If the key points to a slice then the new value is appended.
func (m *Model) Set(key string, v interface{}) {
	var (
		targetModel = m.getField(key)
		value       = reflect.ValueOf(v)
	)

	if targetModel.kind == reflect.Slice {
		if value.Kind() == reflect.Slice {
			value = reflect.AppendSlice(targetModel.val, value)
		} else {
			var model = targetModel.new()
			model.Set(key, v)
			value = reflect.Append(targetModel.val, model.val)
		}
	}

	targetModel.val.Set(value)
}

// Map a series of keys/values before adding it to the parent model.
func (m *Model) Map(kv map[string]interface{}) {
	var pendingModels = make(map[string]*Model)

	for k, v := range kv {
		var (
			model       = m.getField(k)
			pending, ok = pendingModels[model.name]
		)

		if !ok {
			pending = model.new()
			pendingModels[model.name] = pending
		}

		pending.Set(k, v)
	}

	for name, pending := range pendingModels {
		m.Set(name, pending.val.Interface())
	}
}
