package gojson


/*
package GoJSON provides methods for Marshaling/Unmarshaling and editing JSON
*/

// JSONType represents type of the field
type JSONType int

const (
	JSONInvalid JSONType = iota
	JSONBool
	JSONNull
	JSONInt
	JSONFloat
	JSONString
	JSONArray
	JSONObject
)

/*
GoJSON is base json type
expected types are JSONObject and JSONArray
*/
type GoJSON struct {
	Type     JSONType
	Bytes    []byte
	Map      map[string]*GoJSON
	Array    []*GoJSON
}

func (g *GoJSON) ToMap() interface{} {
	if g.Type == JSONObject {
		m := make(map[string]interface{})
		for key, value := range g.Map {
			m[key] = toMap(value)
		}
		return m
	} else {
		if g.Type == JSONArray {
			m := make([]interface{}, 0)
			for _, value := range g.Array {
				m = append(m, toMap(value))
			}
			return m
		}
	}
	return nil
}

func toMap(value *GoJSON) interface{} {
	var r interface{}
	switch value.Type {
	case JSONObject, JSONArray:
		r = value.ToMap()
	case JSONString:
		r, _ = value.ValueString()
	case JSONInt:
		r, _ = value.ValueInt()
	case JSONBool:
		r, _ = value.ValueBool()
	case JSONNull:
		r = nil
	}
	return r
}

func (g *GoJSON) GetBSON() (interface{}, error) {
	return g.ToMap(), nil
}

// NewArray returns new array
func NewArray() *GoJSON {
	return &GoJSON{Type: JSONArray}
}

// NewObject returns new object
func NewObject() *GoJSON {
	return &GoJSON{Type: JSONObject}
}

// Get a node by string key or int index if object is an array
func (g *GoJSON) Get(key interface{}) *GoJSON {
	switch g.Type {
	case JSONObject:
		if val, ok := g.Map[key.(string)]; ok {
			return val
		}
	case JSONArray:
		if index, ok := key.(int); ok {
			if index > len(g.Array) {
				return &GoJSON{}
			}
			return g.Array[index]
		}
	}
	return &GoJSON{}
}

// Delete a key from map or item from array by index
func (g *GoJSON) Delete(key interface{}) string {
	switch g.Type {
	case JSONObject:
		if strKey, ok := key.(string); ok {
			delete(g.Map, strKey)
		} else {
			return "You can delete from object just by string key"
		}
	case JSONArray:
		if index, ok := key.(int); ok {
			g.Array = append(g.Array[:index], g.Array[index+1:]...)
		} else {
			return "You can delete from array just by index"
		}
	default:
		return "cannot delete from non object/array"
	}
	return ""
}

// Pop - same as get but removes node from json
func (g *GoJSON) Pop(key interface{}) *GoJSON {
	value := g.Get(key)
	if value.Type != JSONInvalid {
		g.Delete(key)
	}
	return value
}

// Update merges two objects - available only for objects
func (g *GoJSON) Update(json *GoJSON) string {
	if g.Type != JSONObject {
		return "json is not an object"
	}

	if json.Map == nil || len(json.Map) == 0 {
		// update empty json
		return ""
	}

	for key, value := range json.Map {
		err := g.Set(key, value)
		if err != "" {
			return err
		}
	}
	return ""
}

// Keys returns keys of json object
func (g *GoJSON) Keys() []string {
	if g.Type != JSONObject || g.Map == nil || len(g.Map) == 0 {
		return []string{}
	}
	response := make([]string, len(g.Map))

	i := 0
	for key := range g.Map {
		response[i] = key
		i++
	}
	return response
}

// Values returns values of object or array
func (g *GoJSON) Values() (response []*GoJSON) {
	if g.Type == JSONObject {
		if g.Map == nil || len(g.Map) == 0 {
			return
		}
		response = make([]*GoJSON, len(g.Map))
		i := 0
		for _, value := range g.Map {
			response[i] = value
			i++
		}
	}
	if g.Type == JSONArray {
		return g.Array
	}
	return
}

// Len of Array Object or string
func (g *GoJSON) Len() int {
	switch g.Type {
	case JSONString:
		return len(g.Bytes)
	case JSONObject:
		return len(g.Map)
	case JSONArray:
		return len(g.Array)
	default:
		return 0
	}
}

// Set sets a pointer to JSON struct by key
func (g *GoJSON) Set(key interface{}, value *GoJSON) string {
	if g.Type == JSONInvalid {
		return "Invalid node"
	}
	switch g.Type {
	case JSONObject:
		if g.Map == nil {
			g.Map = make(map[string]*GoJSON)
		}
		g.Map[key.(string)] = value
	case JSONArray:
		var index int
		var ok bool
		if index, ok = key.(int); !ok {
			index = -1
		}
		if index == -1 || index > len(g.Array) {
			//	append
			g.Array = append(g.Array, value)
		} else {
			//	insert
			g.Array = append(g.Array, &GoJSON{})
			copy(g.Array[index+1:], g.Array[index:])
			g.Array[index] = value
		}
	}

	return ""
}
