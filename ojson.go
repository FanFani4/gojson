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
	Map   map[string]*JSON
	Array []*JSON
	Type  JSONType
}

// NewArray returns new array
func NewArray () *GoJSON {
	return &GoJSON{Type: JSONArray}
}

// NewObject returns new object
func NewObject () *GoJSON {
	return &GoJSON{Type: JSONObject}
}

// Get a node by string key or int index if object is an array
func (g *GoJSON) Get(key interface{}) *JSON {
	switch g.Type {
	case JSONObject:
		if val, ok := g.Map[key.(string)]; ok {
			return val
		}
	case JSONArray:
		if index, ok := key.(int); ok {
			if index > len(g.Array) {
				return &JSON{}
			}
			return g.Array[index]
		}
	}
	return &JSON{}
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
func (g *GoJSON) Pop(key interface{}) *JSON {
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
func (g *GoJSON) Values() (response []*JSON) {
	if g.Type == JSONObject {
		if g.Map == nil || len(g.Map) == 0 {return}
		response = make([]*JSON, len(g.Map))
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


// Set sets a pointer to JSON struct by key
func (g *GoJSON) Set (key interface{}, value *JSON) string {
	if g.Type == JSONInvalid {
		return "Invalid node"
	}
	switch g.Type {
	case JSONObject:
		if g.Map == nil {
			g.Map = make(map[string]*JSON)
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
			g.Array = append(g.Array, &JSON{})
			copy(g.Array[index + 1:], g.Array[index:])
			g.Array[index] = value
		}
	}

	return ""
}

// JSON represents a single node
type JSON struct {
    Type     JSONType
    Bytes    []byte
    Children *GoJSON
}

// Set just calls Children.Set if node is JSONArray or JSONObject
func (j *JSON) Set (key interface{}, value *JSON) string {
	if j.Type == JSONArray || j.Type == JSONObject {
		if j.Children == nil {
			j.Children = &GoJSON{}
		}
		j.Children.Set(key, value)
	}
	return "cannot set to non object / array"
}

// Get just calls Children.Get if node is JSONArray or JSONObject
func (j *JSON) Get(key interface{}) *JSON {
    if j.Type != JSONArray && j.Type != JSONObject {
        return &JSON{}
    }
    return j.Children.Get(key)
}

// Pop just calls Children.Pop if node is JSONArray or JSONObject
func (j *JSON) Pop(key interface{}) *JSON {
    if j.Type != JSONArray && j.Type != JSONObject {
        return &JSON{}
    }
    return j.Children.Pop(key)
}

// Update just calls Children.Update if node is JSONObject
func (j *JSON) Update(json *GoJSON) string {
    if j.Type != JSONObject {
		return "not an object"
	}
	if j.Children == nil {
		j.Children = &GoJSON{}
	}
	return j.Children.Update(json)
}

// Keys just calls Children.Keys if node is JSONObject
func (j *JSON) Keys() []string {
	if j.Type != JSONObject {
		return []string{}
	}
	if j.Children == nil {
		j.Children = &GoJSON{}
	}
	return j.Children.Keys()
}

// Values just calls Children.Values if node is JSONObject
func (j *JSON) Values() []*JSON {
	if j.Children == nil {
		j.Children = &GoJSON{}
	}
	return j.Children.Values()
}
