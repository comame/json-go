// 動的に JSON を組み立てる。
// 使用するかどうかは慎重に検討すること。乱用するとメンテナンス性の低下を招きます。
package builder

import (
	"encoding/json"
	"errors"
)

var (
	ErrInvalidReceiverType = errors.New("invalid receiver type")
)

type Json struct {
	t   jsonType
	v   interface{}
	err error
}

type jsonType int

const (
	tobject jsonType = iota
	tarray
	tint
	tfloat
	tstring
	tbool
	tnull
	traw
)

type KV struct {
	key   string
	value Json
}

func Object(values ...KV) Json {
	var j Json
	j.t = tobject
	m := make(map[string]json.RawMessage)

	for _, kv := range values {
		v, err := kv.value.toRaw()
		if err != nil {
			return Json{
				err: err,
			}
		}
		m[kv.key] = v
	}

	j.v = m

	return j
}

func Entry(key string, value Json) KV {
	return KV{
		key:   key,
		value: value,
	}
}

func Array(values ...Json) Json {
	var j Json
	j.t = tarray
	var arr []interface{}

	for _, value := range values {
		v, err := value.toRaw()
		if err != nil {
			return Json{
				err: err,
			}
		}
		arr = append(arr, v)
	}

	j.v = arr

	return j
}

func String(v string) Json {
	return Json{
		t: tstring,
		v: v,
	}
}

func Int64(v int64) Json {
	return Json{
		t: tint,
		v: v,
	}
}

func Float64(v float64) Json {
	return Json{
		t: tfloat,
		v: v,
	}
}

func Bool(v bool) Json {
	return Json{
		t: tbool,
		v: v,
	}
}

func Null() Json {
	return Json{
		t: tnull,
		v: nil,
	}
}

func Raw(v []byte) Json {
	return Json{
		t: traw,
		v: v,
	}
}

func (j *Json) Set(key string, value Json) error {
	if j.t != tobject {
		return ErrInvalidReceiverType
	}

	m, ok := (j.v).(map[string]json.RawMessage)
	if !ok {
		panic("cast error")
	}

	v, err := value.toRaw()
	if err != nil {
		return err
	}

	m[key] = v
	return nil
}

func (j *Json) Push(value Json) error {
	if j.t != tarray {
		return ErrInvalidReceiverType
	}

	arr, ok := (j.v).([]interface{})
	if !ok {
		panic("cast error")
	}

	v, err := value.toRaw()
	if err != nil {
		return err
	}

	arr = append(arr, v)
	j.v = arr
	return nil
}

func (j *Json) Build() ([]byte, error) {
	if j.err != nil {
		return nil, j.err
	}

	b, err := j.toRaw()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (j *Json) MustSet(key string, value Json) {
	if err := j.Set(key, value); err != nil {
		panic(err)
	}
}

func (j *Json) MustPush(value Json) {
	if err := j.Push(value); err != nil {
		panic(err)
	}
}

func (j *Json) MustBuild() []byte {
	s, err := j.Build()
	if err != nil {
		panic(err)
	}
	return s
}

func (j *Json) toRaw() (json.RawMessage, error) {
	if j.err != nil {
		return nil, j.err
	}

	if j.t == traw {
		bytes := j.v.([]byte)
		return bytes, nil
	}

	v, err := json.Marshal(j.v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
