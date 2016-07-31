package gojson

import "bytes"

const (
	startArray byte = 91
	stopArray byte = 93
	startObject byte = 123
	stopObject byte = 125
	startString byte = 34
	escape byte = 92
)

// Loads parses input bytes and returns new json
func Marshal(value []byte) *GoJSON {
	json := &GoJSON{}
	var node *JSON
	parseValue(json, node, skip(value));
	return json
}


func skip (value []byte) []byte {
	i := 0
	for i < len(value) && value[i] <= 32 {
 		i++
	}
	return value[i:]
}

func parseKey(json *GoJSON, node *JSON, value []byte) []byte {
	if value[0] != startString { return []byte{}}
	i := 1
	for i < len(value) {
		if value[i] == startString {
			if value[i - 1] == escape {
				i++
				continue
			}
			break
		}
		i++
	}
	if json.Map == nil {
		json.Map = make(map[string]*JSON)
	}
	json.Map[string(value[1:i])] = node
	return value[i+1:]
}

func parseValue(json *GoJSON, node *JSON, value []byte) []byte {
	if len(value) == 0 { return []byte{} }

	switch value[0] {
	case 110: // n
		if len(value) >= 4 && value[3] == 108 {
			node.Type = JSONNull; node.Bytes = value[:4]
			return value[4:]
		}
		return []byte{}
	case 102: // f
		if len(value) >= 5 && value[4] == 101 {  // e
			node.Type = JSONBool; node.Bytes = value[:5]
			return value[5:]
		}
		return []byte{}
	case 116: // t
		if len(value) >= 4 && value[3] == 101 {   // e
			node.Type = JSONBool; node.Bytes = value[:4]
			return value[4:]
		}
		return []byte{}
	case startString:  // "
		return parseString(node, value);
	case 45, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57: // - , 0-9
		return parseNumber(node, value);
	case startArray:
		return parseArray(json, node, value);
	case startObject:
		return parseObject(json, node, value);

	}
	return []byte{}
}

func parseString(node *JSON, value []byte) []byte {
	if value[0] != startString { return []byte{} }
	i := 1
	for i < len(value) {
		if value[i] == startString {
			if value[i - 1] == escape {
				i++
				continue
			}
			break
		}
		i++
	}

	node.Bytes = value[1:i]
	node.Type = JSONString
	return value[i+1:]
}

func parseNumber(node *JSON, value []byte) []byte {
	i := 0
	nodeType := JSONInt
	if value[i] == 45 { i++ }    /* - */
	for i < len(value) {
		if value[i] >= 48 && value[i] <= 57 {  /* 0 - 9 */
			i++
			continue
		} else {
			if value[i] == 46 {  /* . */
				nodeType = JSONFloat
				i++
				continue
			} else {
				break
			}
		}
	}
	node.Type = nodeType
	node.Bytes = value[:i]
	return value[i:]
}


func parseArray(json *GoJSON, node *JSON, value []byte) []byte {
	// Check for empty array
	var parent *GoJSON
	if json.Type != JSONInvalid {
		// is child
		if node != nil { node.Type = JSONArray }
		if node.Children == nil {
			node.Children = NewArray()
		}
		parent = node.Children
	} else {
		parent = json
	}

	parent.Type = JSONArray
	value = skip(value[1:])
	if len(value) > 0 && value[0] == stopArray { return value[1:] }

	newNode := &JSON{}
	value=skip(parseValue(parent, newNode, skip(value)))
	if len(value) == 0 {return []byte{}}

	if parent.Array == nil { parent.Array = make([]*JSON, 1) }
	parent.Array[0] = newNode

	for value[0] == 44 { // ,
		newNode = &JSON{}
		value=skip(parseValue(parent, newNode, skip(value[1:])))
		if len(value) == 0 {return []byte{}}
		parent.Array = append(parent.Array, newNode)
	}

	if len(value) > 0 && value[0] == stopArray {
		return value[1:]
	}
	return []byte{}
}


func parseObject(json *GoJSON, node *JSON, value []byte) []byte {
	// check for empty object
	var parent *GoJSON

	if json.Type != JSONInvalid {
		// is child
		if node != nil { node.Type = JSONObject }
		if node.Children == nil {
			node.Children = NewObject()
		}
		parent = node.Children
	} else {
		parent = json
	}

	parent.Type = JSONObject
	value = skip(value[1:])
	if len(value) > 0 && value[0] == stopObject { return value[1:] }

	newNode := &JSON{}
	value = skip(parseKey(parent, newNode, skip(value)))
	if len(value) == 0 || value[0] != 58 { return []byte{} } // :
	value=skip(parseValue(parent, newNode, skip(value[1:])))
	if len(value) == 0 {return []byte{}}

	for value[0] == 44 { // ,
		newNode = &JSON{}
		value = skip(parseKey(parent, newNode, skip(value[1:])))
		if len(value) == 0 || value[0] != 58 { return []byte{} } // :
		value=skip(parseValue(parent, newNode, skip(value[1:])))
		if len(value) == 0 {return []byte{}}
	}

	if len(value) > 0 && value[0] == stopObject {
		return value[1:]
	}
	return []byte{}
}

// Unmarshal transforms goJSON to []byte
func (g *GoJSON) Unmarshal(buf ...*bytes.Buffer) []byte {
	var bf *bytes.Buffer
	if len(buf) > 0 {
		bf = buf[0]
	} else {
		bf = &bytes.Buffer{}
		bf.Grow(1024)
	}
	if g.Type == JSONObject {
		bf.WriteByte(startObject)
		idx := 0
		for key, value := range g.Map {
			if idx > 0 {
				bf.WriteByte(44)
			} else {
				idx++
			}
			bf.WriteByte(startString)
			bf.WriteString(key)
			bf.WriteByte(startString)
			bf.WriteByte(58)
			writeValue(value, bf)
		}
		bf.WriteByte(stopObject)
	} else {
		bf.WriteByte(startArray)
		for idx, value := range g.Array {
			if idx > 0 {
				bf.WriteByte(44)
			}
			writeValue(value, bf)
		}
		bf.WriteByte(stopArray)
	}
	return bf.Bytes()
}


func writeValue (value *JSON, bf *bytes.Buffer) {
	switch value.Type {
	case JSONString:
		bf.WriteByte(startString)
		bf.Write(value.Bytes)
		bf.WriteByte(startString)
	case JSONArray, JSONObject:
		value.Children.Unmarshal(bf)
	default:
		bf.Write(value.Bytes)
	}
}
