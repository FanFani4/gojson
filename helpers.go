package gojson

import (
	"strconv"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"math"
	"time"
)

func (g *GoJSON) String() string {
	if g.Type == JSONObject || g.Type == JSONArray {
		return bytesToStr(g.Unmarshal())
	} else {
		val, _ := g.ValueString()
		return val
	}
}

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
		result, err = strconv.Atoi(bytesToStr(j.Bytes))
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
		result, err = strconv.ParseFloat(bytesToStr(j.Bytes), 64)
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
		result = bytesToStr(j.Bytes)
		return result, nil
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
		result, err = strconv.ParseBool(bytesToStr(j.Bytes))
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
		_, err = strconv.Atoi(bytesToStr(value))
		if err != nil {
			return "invalid int"
		}
	case JSONFloat:
		_, err = strconv.ParseFloat(bytesToStr(value), 64)
		if err != nil {
			return "invalid float"
		}
	case JSONBool:
		_, err = strconv.ParseBool(bytesToStr(value))
		if err != nil {
			return "invalid float"
		}
	case JSONArray, JSONObject:
		child := Marshal(value)
		if child == nil {
			return "array or object expected"
		}
		g.Type = Type
		g.Array = child.Array
		g.Map = child.Map
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

// region bson
func (g *GoJSON) GetBSON() (interface{}, error) {
	return g.ToMap(), nil
}

func (g *GoJSON) parseObject (d *decoder, obj *GoJSON) {
	if g.Type == JSONInvalid {
		g.Type = JSONObject
	}
	end := int(d.readInt32())
	end += d.i - 4
	if end <= d.i || end > len(d.in) || d.in[end-1] != '\x00' {
		panic("bson corupted")
	}
	for d.in[d.i] != '\x00' {
		kind := d.readByte()
		name := d.readCStr()
		if d.i >= len(d.in) {
			return
		}
		obj := &GoJSON{}
		g.setBSON(d, kind, obj)
		if g.Type == JSONObject {
			g.Set(name, obj)
		}
	}
	d.i += 1
}

func (g *GoJSON) parseSlice(d *decoder, obj *GoJSON) {
	if g.Type == JSONInvalid {
		g.Type = JSONArray
		if g.Array == nil {
			g.Array = make([]*GoJSON, 0)
		}
	}

	end := int(d.readInt32())
	end += d.i - 4
	if end <= d.i || end > len(d.in) || d.in[end-1] != '\x00' {
		panic("bson corupted")
	}
	for d.in[d.i] != '\x00' {
		kind := d.readByte()
		for d.i < end && d.in[d.i] != '\x00' {
			d.i++
		}
		if d.i >= end {
			panic("corupted")
		}
		d.i++
		obj := &GoJSON{}
		g.setBSON(d, kind, obj)
		if g.Type == JSONArray {
			g.Array = append(g.Array, obj)
		}
		if d.i >= end {
			panic("corupted")
		}
	}
	d.i++ // '\x00'
	if d.i != end {
		panic("corupted")
	}
}


func (g *GoJSON) setBSON(d *decoder, kind byte, obj *GoJSON) {
	switch kind {
	case 0x01: // Float64
		in := d.readFloat64()
		obj.Type = JSONFloat
		obj.Bytes = []byte(strconv.FormatFloat(in, 'f', -1, 64))
	case 0x02: // UTF-8 string
		obj.Type = JSONString
		obj.Bytes = d.readStr()
	case 0x03: // Document
		newObj := NewObject()
		obj.Type = JSONObject
		obj.Map = make(map[string]*GoJSON)
		obj.parseObject(d, newObj)
	case 0x04: // Array
		newArr := NewArray()
		obj.Type = JSONArray
		obj.Array = make([]*GoJSON, 0)
		obj.parseSlice(d, newArr)
	case 0x05: // Binary
		b := d.readBinary()
		obj.Type = JSONString
		obj.Bytes = b.Data
	case 0x07: // ObjectId
		obj.Type = JSONString
		obj.Bytes = d.readBytes(12)
	case 0x08: // Bool
		obj.Type = JSONBool
		obj.Bytes = d.readBool()
	case 0x09: // Timestamp
		// MongoDB handles timestamps as milliseconds.
		obj.Type = JSONFloat
		i := d.readInt64()
		if i == -62135596800000 {
			obj.Bytes = []byte{0} // In UTC for convenience.
		} else {
			obj.Bytes = []byte(time.Unix(i / 1e3, i % 1e3 * 1e6).String())
		}
	case 0x0A: // Nil
		obj.Type = JSONNull
		obj.Bytes = []byte("null")
	case 0x0D, 0x0E: // JavaScript without scope
		obj.Type = JSONString
		obj.Bytes = d.readStr()
	case 0x10: // Int32
		obj.Type = JSONInt
		obj.Bytes = []byte(strconv.Itoa(int(d.readInt32())))
	case 0x11: // Mongo-specific timestamp
		obj.Type = JSONInt
		obj.Bytes = []byte(strconv.Itoa(int(d.readInt64())))
	case 0x12: // Int64
		obj.Type = JSONInt
		obj.Bytes = []byte(strconv.Itoa(int(d.readInt64())))
	default:
		panic(fmt.Sprintf("Unknown element kind (0x%02X)", kind))
	}
	return
}

func (g *GoJSON) SetBSON(raw bson.Raw) error {
	d := decoder{in: raw.Data}
	g.parseObject(&d, g)
	return nil
}

func (d *decoder) readRegEx() bson.RegEx {
	re := bson.RegEx{}
	re.Pattern = d.readCStr()
	re.Options = d.readCStr()
	return re
}

func (d *decoder) readBinary() bson.Binary {
	l := d.readInt32()
	b := bson.Binary{}
	b.Kind = d.readByte()
	b.Data = d.readBytes(l)
	if b.Kind == 0x02 && len(b.Data) >= 4 {
		// Weird obsolete format with redundant length.
		b.Data = b.Data[4:]
	}
	return b
}

func (d *decoder) readStr() []byte {
	l := d.readInt32()
	b := d.readBytes(l - 1)
	if d.readByte() != '\x00' {
		panic("bad")
	}
	return b
}

type decoder struct {
	in      []byte
	i       int
}


func (d *decoder) readCStr() string {
	start := d.i
	end := start
	l := len(d.in)
	for ; end != l; end++ {
		if d.in[end] == '\x00' {
			break
		}
	}
	d.i = end + 1
	if d.i > l {
		return "bad"
	}
	return bytesToStr(d.in[start:end])
}

func (d *decoder) readBool() []byte {
	b := d.readByte()
	if b == 0 {
		return []byte("false")
	}
	if b == 1 {
		return []byte("true")
	}
	panic(fmt.Sprintf("encoded boolean must be 1 or 0, found %d", b))
}

func (d *decoder) readFloat64() float64 {
	return math.Float64frombits(uint64(d.readInt64()))
}

func (d *decoder) readInt32() int32 {
	b := d.readBytes(4)
	return int32((uint32(b[0]) << 0) |
		(uint32(b[1]) << 8) |
		(uint32(b[2]) << 16) |
		(uint32(b[3]) << 24))
}

func (d *decoder) readInt64() int64 {
	b := d.readBytes(8)
	return int64((uint64(b[0]) << 0) |
		(uint64(b[1]) << 8) |
		(uint64(b[2]) << 16) |
		(uint64(b[3]) << 24) |
		(uint64(b[4]) << 32) |
		(uint64(b[5]) << 40) |
		(uint64(b[6]) << 48) |
		(uint64(b[7]) << 56))
}

func (d *decoder) readByte() byte {
	i := d.i
	d.i++
	if d.i > len(d.in) {
		return 0
	}
	return d.in[i]
}

func (d *decoder) readBytes(length int32) []byte {
	if length < 0 {
		return []byte("aaaa")
	}
	start := d.i
	d.i += int(length)
	if d.i < start || d.i > len(d.in) {
		return []byte("aaaa")
	}
	return d.in[start : start+int(length)]
}
// endregion