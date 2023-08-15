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

type tjson struct {
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
)

type kv struct {
	key   string
	value tjson
}

func Object(values ...kv) tjson {
	var j tjson
	j.t = tobject
	m := make(map[string]json.RawMessage)

	for _, kv := range values {
		v, err := kv.value.toRaw()
		if err != nil {
			return tjson{
				err: err,
			}
		}
		m[kv.key] = v
	}

	j.v = m

	return j
}

func KV(key string, value tjson) kv {
	return kv{
		key:   key,
		value: value,
	}
}

func Array(values ...tjson) tjson {
	var j tjson
	j.t = tarray
	var arr []interface{}

	for _, value := range values {
		v, err := value.toRaw()
		if err != nil {
			return tjson{
				err: err,
			}
		}
		arr = append(arr, v)
	}

	j.v = arr

	return j
}

func String(v string) tjson {
	return tjson{
		t: tstring,
		v: v,
	}
}

func Int64(v int64) tjson {
	return tjson{
		t: tint,
		v: v,
	}
}

func Float64(v float64) tjson {
	return tjson{
		t: tfloat,
		v: v,
	}
}

func Bool(v bool) tjson {
	return tjson{
		t: tbool,
		v: v,
	}
}

func Null() tjson {
	return tjson{
		t: tnull,
		v: nil,
	}
}

func (j *tjson) Set(key string, value tjson) error {
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

func (j *tjson) Push(value tjson) error {
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

func (j *tjson) Build() ([]byte, error) {
	if j.err != nil {
		return nil, j.err
	}

	b, err := j.toRaw()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (j *tjson) MustSet(key string, value tjson) {
	if err := j.Set(key, value); err != nil {
		panic(err)
	}
}

func (j *tjson) MustPush(value tjson) {
	if err := j.Push(value); err != nil {
		panic(err)
	}
}

func (j *tjson) MustBuild() []byte {
	s, err := j.Build()
	if err != nil {
		panic(err)
	}
	return s
}

func (j *tjson) toRaw() (json.RawMessage, error) {
	if j.err != nil {
		return nil, j.err
	}

	v, err := json.Marshal(j.v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
