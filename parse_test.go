package gojson

import (
	"fmt"
	"testing"
)

var js *GoJSON

var data []byte = []byte(`{
        "name": "Calarasanu Andrei",
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
	js.Get("person").Get("name").SetString("fullName", "test")
	js.SetString("key", "value")
	fmt.Println(js.Get("key").ValueString())
	js.SetBytes("testJson", []byte(`{"hello": "world"}`), JSONObject)
	js.SetBytes("test_arr", []byte(`[12, 11, 10]`), JSONArray)
	fmt.Println(js.Get("person").Get("name").Get("fullName").ValueString())
	fmt.Println(js.Get("testJson").Get("hello").ValueString())
	js.Delete("testJson")
	fmt.Println(js.Get("testJson").ValueString("dft value"))

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
