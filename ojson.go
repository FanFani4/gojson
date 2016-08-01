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
	Children *GoJSON
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
	node := g
	if node.Children != nil {
		node = node.Children
	}
	switch node.Type {
	case JSONObject:
		if val, ok := node.Map[key.(string)]; ok {
			return val
		}
	case JSONArray:
		if index, ok := key.(int); ok {
			if index > len(node.Array) {
				return &GoJSON{}
			}
			return node.Array[index]
		}
	}
	return &GoJSON{}
}

// Delete a key from map or item from array by index
func (g *GoJSON) Delete(key interface{}) string {
	node := g
	if node.Children != nil {
		node = node.Children
	}
	switch node.Type {
	case JSONObject:
		if strKey, ok := key.(string); ok {
			delete(node.Map, strKey)
		} else {
			return "You can delete from object just by string key"
		}
	case JSONArray:
		if index, ok := key.(int); ok {
			node.Array = append(node.Array[:index], node.Array[index+1:]...)
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
	node := g
	if node.Children != nil {
		node = node.Children
	}
	value := node.Get(key)
	if value.Type != JSONInvalid {
		node.Delete(key)
	}
	return value
}

// Update merges two objects - available only for objects
func (g *GoJSON) Update(json *GoJSON) string {
	node := g
	if node.Children != nil {
		node = node.Children
	}
	if node.Type != JSONObject {
		return "json is not an object"
	}

	if json.Map == nil || len(json.Map) == 0 {
		// update empty json
		return ""
	}

	for key, value := range json.Map {
		err := node.Set(key, value)
		if err != "" {
			return err
		}
	}
	return ""
}

// Keys returns keys of json object
func (g *GoJSON) Keys() []string {
	node := g
	if node.Children != nil {
		node = node.Children
	}
	if node.Type != JSONObject || node.Map == nil || len(node.Map) == 0 {
		return []string{}
	}
	response := make([]string, len(node.Map))

	i := 0
	for key := range node.Map {
		response[i] = key
		i++
	}
	return response
}

// Values returns values of object or array
func (g *GoJSON) Values() (response []*GoJSON) {
	node := g
	if node.Children != nil {
		node = node.Children
	}
	if node.Type == JSONObject {
		if node.Map == nil || len(node.Map) == 0 {
			return
		}
		response = make([]*GoJSON, len(node.Map))
		i := 0
		for _, value := range node.Map {
			response[i] = value
			i++
		}
	}
	if node.Type == JSONArray {
		return node.Array
	}
	return
}

// Set sets a pointer to JSON struct by key
func (g *GoJSON) Set(key interface{}, value *GoJSON) string {
	if g.Type == JSONInvalid {
		return "Invalid node"
	}
	node := g
	if node.Children != nil {
		node = node.Children
	}
	switch node.Type {
	case JSONObject:
		if node.Map == nil {
			node.Map = make(map[string]*GoJSON)
		}
		node.Map[key.(string)] = value
	case JSONArray:
		var index int
		var ok bool
		if index, ok = key.(int); !ok {
			index = -1
		}
		if index == -1 || index > len(node.Array) {
			//	append
			node.Array = append(node.Array, value)
		} else {
			//	insert
			node.Array = append(node.Array, &GoJSON{})
			copy(node.Array[index+1:], node.Array[index:])
			node.Array[index] = value
		}
	}

	return ""
}
