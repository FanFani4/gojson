package gojson

import (
	"testing"
	"fmt"
)


var js *GoJSON

var data []byte = []byte(`{
    "name": "Jack (\"Bee\") Nimble",
    "format": {
        "type":       "rect",
        "width":      1920,
        "height":     1080,
        "interlace":  false,
        "frame rate": 24
    }
}`)



func TestMarshal(t *testing.T) {
	js = Marshal(data)
	//js.Get("person").Get("gravatar").SetString("urls", "yoyo")
	//fmt.Println(js.Get("person").Get("gravatar").Get("urls").ValueString())
	//js.SetBytes("test_obj", []byte(`{"test": "best"}`), JSONObject)
	//js.SetBytes("test_arr", []byte(`[12, 11, 10]`), JSONObject)
	//js.SetBytes("yes", []byte("true"), JSONBool)
	//

	js.SetInt("name", 1)

	fmt.Println(js.Get("format").Get("width").ValueFloat())

	value, _ := js.Get("format").Pop("type").ValueString()
	fmt.Println(value)

	fmt.Println(js.Get("name").ValueString())
	fmt.Println(js.Keys())
	fmt.Println(js.Values())
	fmt.Println(js.Delete("name"))



	data2 := []byte(`{
		"foo": "bar"
	}`)

	json2 := Marshal(data2)

	js.Update(json2)

	//js.SetString("company", "Yoyoyo")
	//fmt.Println(js.Get("not_existent_key").ValueString("default value"))
	//fmt.Println(js.Get("not_existent_key").ValueInt(100))
	//fmt.Println(js.Get("not_existent_key").ValueInt())
	//fmt.Println(js.Get("not_existent_key").ValueBool(true))
	//fmt.Println(js.Get("not_existent_key").ValueFloat(100.100))
	//fmt.Println(js.Get("person").Get("gravatar").Get("avatars").Get(0).Get("url").ValueString())
	fmt.Println(string(js.Unmarshal()))
}

func BenchmarkMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Marshal(data)
	}
}


func BenchmarkGoJSON_Unmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		js.Unmarshal()
	}
}
