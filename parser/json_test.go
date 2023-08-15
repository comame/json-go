package parser_test

import (
	"fmt"
	"testing"

	"github.com/comame/json-go/parser"
)

func TestInt64(t *testing.T) {
	input := []byte(`32`)
	o := parser.New(input)

	v, err := o.Int64()
	if err != nil {
		t.Error(err)
	}

	if v != 32 {
		t.Errorf("Expected 32, Got %d", v)
	}
}

func TestInt64_Fail(t *testing.T) {
	input := []byte(`1.0`)
	o := parser.New(input)

	_, err := o.Int64()
	if err == nil {
		t.Error(err)
	}
}

func TestIsNull(t *testing.T) {
	input := []byte("null")
	v := parser.New(input).IsNull()
	if !v {
		t.Fail()
	}
}

func TestIsNull_InMap(t *testing.T) {
	input := []byte(`{ "foo": null }`)
	v := parser.New(input).Key(".foo").IsNull()
	if !v {
		t.Fail()
	}
}

func TestIndex(t *testing.T) {
	input := []byte(`[0, 1, 2]`)
	o := parser.New(input)

	v0, err := o.Index(0).Int64()
	if err != nil {
		t.Error(err)
	}
	if v0 != 0 {
		t.Errorf("Expected 0, Got %d", v0)
	}

	v2, err := o.Index(2).Int64()
	if err != nil {
		t.Error(err)
	}
	if v2 != 2 {
		t.Errorf("Expected 2, Got %d", v0)
	}
}

func TestIndex_OutOfRange(t *testing.T) {
	input := []byte(`[0, 1, 2]`)
	_, err := parser.New(input).Index(10).Int64()
	if err == nil {
		t.Fail()
	}

	if err != parser.ErrOutOfRange {
		t.Fail()
	}
}

func TestLen(t *testing.T) {
	input := []byte(`[0, 1, 2]`)
	got := parser.New(input).Len()

	if got != 3 {
		t.Errorf("Expect 3, Got %d", got)
	}
}

func TestLen_zero(t *testing.T) {
	input := []byte(`fooo`)
	got := parser.New(input).Len()

	if got != 0 {
		t.Errorf("Expect 0, Got %d", got)
	}
}

func TestKeys(t *testing.T) {
	input := []byte(`{
		"foo": 0,
		"bar": 1
	}`)
	got := parser.New(input).Keys()

	if !compSlice(got, []string{"foo", "bar"}) {
		t.Fail()
	}
}

func TestKeys_FailArray(t *testing.T) {
	input := []byte(`[0, 1, 2]`)
	got := parser.New(input).Keys()

	if !compSlice(got, []string{}) {
		t.Fail()
	}
}

func TestKeys_FailPrimitive(t *testing.T) {
	input := []byte(`0`)
	got := parser.New(input).Keys()

	if !compSlice(got, []string{}) {
		t.Fail()
	}
}

func TestKey(t *testing.T) {
	input := []byte(`{ "foo": 0 }`)
	o := parser.New(input)

	got, err := o.Key(".foo").Int64()
	if err != nil {
		t.Error(err)
	}

	if got != 0 {
		t.Fail()
	}
}

func TestKey_Nested(t *testing.T) {
	input := []byte(`{ "foo": { "bar": 0 } }`)
	got, err := parser.New(input).Key("foo.bar").Int64()

	if err != nil {
		t.Error(err)
	}

	if got != 0 {
		t.Fail()
	}
}

func TestKey_NoKey(t *testing.T) {
	input := []byte(`{ "foo": { "bar": 0 } }`)
	_, err := parser.New(input).Key("foo.baz").Int64()

	if err != parser.ErrNoKey {
		t.Fail()
	}
}

func Example() {
	input := []byte(`{"foo":{"bar":[{"baz":"baz"},0,1,"array"]}}`)
	obj := parser.New(input)

	baz, err := obj.Key("foo.bar").Index(0).Key("baz").String()
	if err != nil {
		panic(err)
	}
	if baz != "baz" {
		panic(fmt.Sprintf("Expect baz, got %s", baz))
	}

	zero, err := obj.Key("foo.bar").Index(1).Int64()
	if err != nil {
		panic(err)
	}
	if zero != 0 {
		panic("wrong")
	}
}

func compSlice[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
