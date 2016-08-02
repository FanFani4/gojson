package gojson

import (
	"fmt"
	"testing"
)

var js *GoJSON

var data = []byte(`{
	"name": "Calarasanu Andrei",
	"format": {
		"type":       "rect",
		"width":      1920,
		"height":     1080,
		"interlace":  false,
		"frame rate": 24
	}
}`)

var data2 = []byte(`{
	"company": {
		"name": "Simpals"
	}
}`)

func TestMarshal(t *testing.T) {
	js = Marshal(data)

	fmt.Println(string(js.Unmarshal()))
}

func TestGoJSON_SetBytes(t *testing.T) {
	err := js.SetBytes("testJson", []byte(`{"hello": "world"}`), JSONObject)
	err = js.SetBytes("test_arr", []byte(`[12, 11, 10]`), JSONArray)
	if err != "" {
		panic(err)
	}

}

func TestGoJSON_Update(t *testing.T) {
	js2 := Marshal(data2)
	err := js.Update(js2)
	if err != "" {
		panic(err)
	}
	val, er := js.Get("company").Get("name").ValueString()
	if er != nil {
		panic(er)
	}
	if val != "Simpals" {
		panic("Wrong")
	}
}

func TestGoJSON_Delete(t *testing.T) {
	err := js.Delete("test_arr")
	if err != ""{
		panic(err)
	}
	if js.Get("test_arr").Type != JSONInvalid {
		panic("not deleted")
	}
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
