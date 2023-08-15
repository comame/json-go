// 動的な JSON をパースする。
// 使用するかどうかは慎重に検討すること。乱用するとメンテナンス性の低下を招きます。
package parser

import (
	"encoding/json"
	"errors"
	"strings"
)

type Json struct {
	raw []byte
	err error
}

var (
	ErrOutOfRange = errors.New("index out of range")
	ErrNoKey      = errors.New("no key")
)

func New(v []byte) *Json {
	return &Json{
		raw: v,
		err: nil,
	}
}

func (j *Json) Keys() []string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(j.raw, &obj); err != nil {
		return nil
	}

	var keys []string
	for key := range obj {
		keys = append(keys, key)
	}

	return keys
}

func (j *Json) Key(key string) *Json {
	if j.err != nil {
		return &Json{
			raw: nil,
			err: j.err,
		}
	}

	q := strings.Split(key, ".")

	raw := j.raw

	for _, query := range q {
		if query == "" {
			continue
		}

		var m map[string]json.RawMessage
		if err := json.Unmarshal(raw, &m); err != nil {
			return &Json{
				raw: nil,
				err: j.err,
			}
		}
		v, ok := m[query]
		if !ok {
			return &Json{
				raw: nil,
				err: ErrNoKey,
			}
		}
		raw = v
	}

	return &Json{
		raw: raw,
		err: nil,
	}
}

func (v *Json) Int64() (int64, error) {
	if v.err != nil {
		return 0, v.err
	}

	var num int64
	if err := json.Unmarshal(v.raw, &num); err != nil {
		return 0, err
	}
	return num, nil
}

func (v *Json) Float64() (float64, error) {
	if v.err != nil {
		return 0, v.err
	}

	var num float64
	if err := json.Unmarshal(v.raw, &num); err != nil {
		return 0, err
	}
	return num, nil
}

func (v *Json) String() (string, error) {
	if v.err != nil {
		return "", v.err
	}

	var str string
	if err := json.Unmarshal(v.raw, &str); err != nil {
		return "", err
	}
	return str, nil
}

func (v *Json) IsNull() bool {
	if v.err != nil {
		return false
	}

	var o *int
	if err := json.Unmarshal(v.raw, &o); err != nil {
		return false
	}

	if o != nil {
		return false
	}

	return true
}

func (v *Json) Index(i int) *Json {
	var obj []json.RawMessage
	if err := json.Unmarshal(v.raw, &obj); err != nil {
		return &Json{
			raw: nil,
			err: err,
		}
	}

	if i > len(obj) {
		return &Json{
			raw: nil,
			err: ErrOutOfRange,
		}
	}

	item, err := json.Marshal(obj[i])
	if err != nil {
		return &Json{
			raw: nil,
			err: err,
		}
	}

	return &Json{
		raw: item,
		err: err,
	}
}

// レシーバが配列ではないとき、0 を返す。
func (v *Json) Len() int {
	var obj []json.RawMessage
	if err := json.Unmarshal(v.raw, &obj); err != nil {
		return 0
	}
	return len(obj)
}

func (v *Json) Raw() ([]byte, error) {
	if v.err != nil {
		return nil, v.err
	}
	return v.raw, nil
}
