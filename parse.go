package gojson

import (
	"bytes"
	"reflect"
	"unsafe"
)

const (
	startArray  byte = '['
	stopArray   byte = ']'
	startObject byte = '{'
	stopObject  byte = '}'
	startString byte = '"'
	escape      byte = '\\'
)


// Unarshal parses input bytes and returns new json
func Unmarshal(value []byte) *GoJSON {
	json := &GoJSON{}
	var node *GoJSON
	parseValue(json, node, skip(value))
	return json
}


func skip(value []byte) []byte {
	i := 0
	for i < len(value) && value[i] <= 32 {
		i++
	}
	return value[i:]
}

func parseKey(json *GoJSON, node *GoJSON, value []byte) []byte {
	if value[0] != startString {
		syntaxError()
	}
	i := 1
	for i < len(value) {
		if value[i] == startString {
			if value[i-1] == escape {
				i++
				continue
			}
			break
		}
		i++
	}
	if json.Map == nil {
		json.Map = make(map[string]*GoJSON)
	}
	json.Map[bytesToStr(value[1:i])] = node
	return value[i+1:]
}

func parseValue(json *GoJSON, node *GoJSON, value []byte) []byte {
	if len(value) == 0 {
		return value
	}

	switch value[0] {
	case 'n': // n
		if len(value) >= 4 && value[1] == 'u' && value[2] == 'l' && value[3] == 'l' {
			node.Type = JSONNull
			node.Bytes = value[:4]
			return value[4:]
		}
	case 'f': // f
		if len(value) >= 5 && value[1] == 'a' && value[2] == 'l' && value[3] == 's' && value[4] == 'e' { // e
			node.Type = JSONBool
			node.Bytes = value[:5]
			return value[5:]
		}
	case 't': // t
		if len(value) >= 4 && value[1] == 'r' && value[2] == 'u' && value[3] == 'e' {
			node.Type = JSONBool
			node.Bytes = value[:4]
			return value[4:]
		}
	case startString: // "
		return parseString(node, value)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-': // - , 0-9
		return parseNumber(node, value)
	case startArray:
		return parseArray(json, node, value)
	case startObject:
		return parseObject(json, node, value)
	}
	return value
}

func parseString(node *GoJSON, value []byte) []byte {
	if value[0] != startString {
		syntaxError()
	}
	i := 1
	for i < len(value) {
		if value[i] == startString {
			if value[i-1] == escape {
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

func parseNumber(node *GoJSON, value []byte) []byte {
	i := 0
	nodeType := JSONInt
	if value[i] == 45 {
		i++
	} /* - */
	for i < len(value) {
		if value[i] >= 48 && value[i] <= 57 || value[i] == 'E' { /* 0 - 9 */
			i++
			continue
		} else {
			if value[i] == 46 { /* . */
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

func parseArray(json *GoJSON, node *GoJSON, value []byte) []byte {
	// Check for empty array
	var parent *GoJSON
	if json.Type != JSONInvalid {
		// is child
		if node != nil {
			node.Type = JSONArray
		} else {
			node = NewArray()
		}
		parent = node
	} else {
		parent = json
	}

	parent.Type = JSONArray
	value = skip(value[1:])
	if len(value) > 0 && value[0] == stopArray {
		return value[1:]
	}

	newNode := &GoJSON{}
	value = skip(parseValue(parent, newNode, skip(value)))
	if len(value) == 0 {
		return value
	}

	if parent.Array == nil {
		parent.Array = make([]*GoJSON, 1)
	}
	parent.Array[0] = newNode

	for value[0] == 44 { // ,
		newNode = &GoJSON{}
		value = skip(parseValue(parent, newNode, skip(value[1:])))
		if len(value) == 0 {
			return value
		}
		parent.Array = append(parent.Array, newNode)
	}

	if len(value) > 0 && value[0] == stopArray {
		return value[1:]
	}
	return value
}

func parseObject(json *GoJSON, node *GoJSON, value []byte) []byte {
	// check for empty object
	var parent *GoJSON

	if json.Type != JSONInvalid {
		// is child
		if node != nil {
			node.Type = JSONObject
		} else {
			node = NewObject()
		}
		parent = node
	} else {
		parent = json
	}

	parent.Type = JSONObject
	value = skip(value[1:])
	if len(value) > 0 && value[0] == stopObject {
		return value[1:]
	}

	newNode := &GoJSON{}
	value = skip(parseKey(parent, newNode, skip(value)))
	if len(value) == 0 || value[0] != ':' {
		syntaxError()
	}
	value = skip(parseValue(parent, newNode, skip(value[1:])))
	if len(value) == 0 {
		return value
	}

	for value[0] == ',' { // ,
		newNode = &GoJSON{}
		value = skip(parseKey(parent, newNode, skip(value[1:])))
		if len(value) == 0 || value[0] != ':' {
			syntaxError()
		}
		value = skip(parseValue(parent, newNode, skip(value[1:])))
		if len(value) == 0 {
			return value
		}
	}

	if len(value) > 0 && value[0] == stopObject {
		return value[1:]
	}
	return value
}

// Marshal transforms goJSON to []byte
func (g *GoJSON) Marshal(buf ...*bytes.Buffer) []byte {
	var bf *bytes.Buffer
	if len(buf) > 0 {
		bf = buf[0]
	} else {
		bf = &bytes.Buffer{}
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

func writeValue(value *GoJSON, bf *bytes.Buffer) {
	switch value.Type {
	case JSONString:
		bf.WriteByte(startString)
		var p int
		for i := 0; i < len(value.Bytes); i++ {
			c := value.Bytes[i]
			var e byte
			switch c {
			case '\t':
				e = 't'
			case '\r':
				e = 'r'
			case '\n':
				e = 'n'
			case '\\':
				e = '\\'
			case '"':
				e = '"'
			//case '<', '>':
			//	if !w.EscapeLtGt {
			//		continue
			//	}
			default:
				if c >= 0x20 {
					// no escaping is required
					continue
				}
			}
			if e != 0 {
				bf.Write(value.Bytes[p:i])
				bf.WriteByte(escape)
				bf.WriteByte(e)
			} else {
				bf.Write(value.Bytes[p:i])
			}
			p = i + 1
		}
		bf.Write(value.Bytes[p:])
		bf.WriteByte(startString)
	case JSONArray, JSONObject:
		value.Marshal(bf)
	default:
		bf.Write(value.Bytes)
	}
}


func bytesToStr(data []byte) string {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sHdr := reflect.StringHeader{h.Data, h.Len}
	return *(*string)(unsafe.Pointer(&sHdr))
}

func syntaxError() {
	panic("Invalid json")
}
