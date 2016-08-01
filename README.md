
Welcome to GoJSON.

GoJSON is a very simple json parser, its approach is similar to https://github.com/DaveGamble/cJSON .

Some JSON:
    
    {
        {
        "name": "Calarasanu Andrei", 
        "format": {
            "type":       "rect", 
            "width":      1920, 
            "height":     1080, 
            "interlace":  false, 
            "frame rate": 24
        }
    }

Get it parsed:

    json := gojson.Marshal(jsonBytes)
    
Get a value:

    value, err := json.Get("format").Get("type").ValueString()
    if err != "" {
        fmt.Println(err)
    }
    fmt.Println(value)
    
if such key exists then its value will be returned , otherwise the zero value of the type , if
such behaviour is not acceptable , you can pass default value and get it in case of error:

    value, _ := json.Get("not_existent_key").ValueString("default value"))
    fmt.Println(value)

if you want to get a value and remove it from json you can use Pop:

    value, err := json.Get("format").Pop("type").ValueString()
    
change or add a value:

    json.SetInt("name", 1)
    value, err := json.Get("name").ValueInt()
    
delete a key:

    json.Delete("name")
    
get slice of keys:

    json.Keys()

get slice of value nodes:

    json.Values()
    
update one json with another:

    data := []byte(`{
		"foo": "bar"
	}`)
	json2 := Marshal(data2)
	json.Update(json2)

back to []byte:

    b := json.Unmarshal()

Here's the structure:

    type GoJSON struct {
        Type     JSONType
        Bytes    []byte
        Map      map[string]*GoJSON
        Array    []*GoJSON
        Children *GoJSON
    }
    


medium size json benchmark:

    Marshal                 50000             29142 ns/op           13842 B/op        201 allocs/op
    Unmarshal              100000             12199 ns/op            4979 B/op          5 allocs/op
