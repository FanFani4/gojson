
Welcome to GoJSON.

GoJSON is a very simple json parser, its approach is similar to https://github.com/DaveGamble/cJSON .
it was created to work in more dynamical way with json. its not finished yet and not tested.

Some JSON:

{

    "name": "Jack (\"Bee\") Nimble", 
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