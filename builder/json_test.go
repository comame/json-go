package builder_test

import (
	"fmt"

	"github.com/comame/json-go/builder"
)

func Example() {
	obj := builder.Object(
		builder.KV("string", builder.String("hoge")),
		builder.KV("int", builder.Int64(12)),
		builder.KV("float", builder.Float64(1.0)),
		builder.KV("bool", builder.Bool(true)),
		builder.KV("null", builder.Null()),
		builder.KV("nested", builder.Object(
			builder.KV("foo", builder.String("bar")),
		)),
		builder.KV("array", builder.Array(
			builder.String("foo"),
			builder.String("bar"),
			builder.Null(),
		)),
	)

	arr := builder.Array(builder.String("hi"))
	arr.MustPush(builder.String("hello"))

	obj.MustSet("greetings", arr)

	str := obj.MustBuild()

	fmt.Println(string(str))
	// Output: {"array":["foo","bar",null],"bool":true,"float":1,"greetings":["hi","hello"],"int":12,"nested":{"foo":"bar"},"null":null,"string":"hoge"}
}
