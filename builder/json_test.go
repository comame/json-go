package builder_test

import (
	"fmt"

	"github.com/comame/json-go/builder"
)

func Example() {
	obj := builder.Object(
		builder.Entry("string", builder.String("hoge")),
		builder.Entry("int", builder.Int64(12)),
		builder.Entry("float", builder.Float64(1.0)),
		builder.Entry("bool", builder.Bool(true)),
		builder.Entry("null", builder.Null()),
		builder.Entry("nested", builder.Object(
			builder.Entry("foo", builder.String("bar")),
		)),
		builder.Entry("array", builder.Array(
			builder.String("foo"),
			builder.String("bar"),
			builder.Null(),
		)),
		builder.Entry("raw", builder.Raw([]byte(`{"raw":0}`))),
	)

	arr := builder.Array(builder.String("hi"))
	arr.MustPush(builder.String("hello"))

	obj.MustSet("greetings", arr)

	str := obj.MustBuild()

	fmt.Println(string(str))
	// Output: {"array":["foo","bar",null],"bool":true,"float":1,"greetings":["hi","hello"],"int":12,"nested":{"foo":"bar"},"null":null,"raw":{"raw":0},"string":"hoge"}
}
