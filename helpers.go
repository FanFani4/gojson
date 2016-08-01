package gojson

import (
	"strconv"
	"errors"
	"fmt"
)

// region Getters

// Value returns bytes and type of the current node
func (j *GoJSON) Value() ([]byte, JSONType) {
	return j.Bytes, j.Type
}

// ValueInt returns int representation of the node if its Type is JSONInt or JSONFloat
// if node is empty and dft was specified if will be returned otherwise 0 and error
func (j *GoJSON) ValueInt(dft ...int) (result int, err error) {
	if j.Type != JSONInt && j.Type != JSONFloat {
		err = errors.New("Type missmatch")
	} else {
		result, err = strconv.Atoi(string(j.Bytes))
	}
	if err != nil {
		if len(dft) > 0 {
			return dft[0], nil
		}
	}
	return
}

// ValueFloat returns float representation of the node if its Type is JSONFloat or JSONInt
// if node is empty and dft was specified if will be returned otherwise 0 and error
func (j *GoJSON) ValueFloat(dft ...float64) (result float64, err error) {
	if j.Type != JSONFloat && j.Type != JSONInt {
		err = errors.New("Type missmatch")
	} else {
		result, err = strconv.ParseFloat(string(j.Bytes), 64)
	}
	if err != nil {
		if len(dft) > 0 {
			return dft[0], nil
		}
	}
	return
}

// ValueString returns string representation of the node if its Type is JSONString
// if node is empty and dft was specified if will be returned otherwise "" and error
func (j *GoJSON) ValueString(dft ...string) (result string, err error) {
	if j.Type != JSONString {
		err = errors.New("Type missmatch")
	}
	if len(j.Bytes) > 0 {
		result = string(j.Bytes)
	}
	if len(dft) > 0 {
		return dft[0], nil
	}
	return
}

// ValueBool returns string representation of the node if its Type is JSONBool
// if node is empty and dft was specified if will be returned otherwise false and error
func (j *GoJSON) ValueBool(dft ...bool) (result bool, err error) {
	if j.Type != JSONBool {
		err = errors.New("Type missmatch")
	} else {
		result, err = strconv.ParseBool(string(j.Bytes))
	}
	if err != nil {
		if len(dft) > 0 {
			return dft[0], nil
		}
	}
	return
}

// endregion

// region Setters

/*
SetBytes is a universal method to add a node
js.SetBytes("test_obj", []byte(`{"test": "best"}`), JSONObject)
js.SetBytes("test_arr", []byte(`[12, 11, 10]`), JSONObject)
js.SetBytes("yes", []byte("true"), JSONBool)
*/
func (g *GoJSON) SetBytes(key interface{}, value []byte, Type JSONType) string {
	if Type != JSONInvalid {
		if node := g.Get(key); node.Type != JSONInvalid {
			return node.setBytes(value, Type)
		}
		node := &GoJSON{}
		err := node.setBytes(value, Type)
		if err != "" {
			return err
		}
		return g.Set(key, node)
	}
	return "trying to assign invalid node"
}

func (g *GoJSON) setBytes(value []byte, Type JSONType) string {
	var err error
	switch Type {
	case JSONInt:
		_, err = strconv.Atoi(string(value))
		if err != nil {
			return "invalid int"
		}
	case JSONFloat:
		_, err = strconv.ParseFloat(string(value), 64)
		if err != nil {
			return "invalid float"
		}
	case JSONBool:
		_, err = strconv.ParseBool(string(value))
		if err != nil {
			fmt.Println(err)
			return "invalid float"
		}
	case JSONArray, JSONObject:
		child := Marshal(value)
		if child == nil {
			return "array or object expected"
		}
		g.Type = Type
		g.Children = child
		return ""
	}
	g.Type = Type
	g.Bytes = value
	return ""
}

// SetInt is a helper for setting a int node
func (g *GoJSON) SetInt(key interface{}, value int) string {
	node := &GoJSON{Type: JSONInt, Bytes: []byte(strconv.Itoa(value))}
	err := g.Set(key, node)
	return err
}

// SetString is a helper for setting string
func (g *GoJSON) SetString(key interface{}, value string) string {
	node := &GoJSON{Type: JSONString, Bytes: []byte(value)}
	err := g.Set(key, node)
	return err
}

// SetFloat is a helper for setting float
func (g *GoJSON) SetFloat(key interface{}, value float64) string {
	node := &GoJSON{Type: JSONFloat, Bytes: []byte(strconv.FormatFloat(value, 'f', -1, 64))}
	err := g.Set(key, node)
	return err
}

// SetBool is a helper for setting bool
func (g *GoJSON) SetBool(key interface{}, value bool) string {
	node := &GoJSON{Type: JSONBool, Bytes: []byte(strconv.FormatBool(value))}
	err := g.Set(key, node)
	return err
}

// SetNull sets a json null
func (g *GoJSON) SetNull(key interface{}) string {
	node := &GoJSON{Type: JSONNull, Bytes: []byte("null")}
	err := g.Set(key, node)
	return err
}

// endregion
