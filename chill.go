/*
Chillson package is a convenience tool, providing way to process unmarshaled JSON data in schema-agnostic way. Thus you don't have to
specify type structure conforming to the expected JSON structure, or write countless Golang type assertions. The latter is automated by
chillson's Son type methods, like in the example below.

    import (
	    "chillson"
	    "encoding/json"
	    "fmt"
    )

    var jsonData interface{}
    json.Unmarshal([]byte(rawJson), &jsonData)
    chill := chillson.Son{jsonData}

    // now you can use Son-type variable like this:
    strField, err := chill.GetStr("[gophers][0][name]")
    fmt.Println(strField)
    intField, err := chill.GetInt("[gophers][0][weight]")
    fmt.Println(intField)

    // you can also spawn "smaller" Son{}'s to avoid some underlying type assertions:
    gophersTable, err := chill.GetArr("[gophers]")
    for i := 0; i < len(gophersTable); i++ {
	    gophersRow := chillson.Son{gophersTable[i]}
	    strField, err = chill.GetStr("[name]")
            fmt.Println(strField)
	    intField, err = chill.GetInt("[weight]")
            fmt.Println(intField)
    }

    // testing for specific error values, if you happen to need them:
    switch err {
    case nil:
        break
    case InvalidPath:     // "Chillson: cannot parse nonempty value path, did you forget about square brackets?"
        ...
    case OutOfRange:      // "Chillson: value's parent seems to be a JSON array, but the index is out of range."
        ...
    case ParentNotObject: // "Chillson: value's parent is neither JSON object nor array."
        ...
    case FieldNotFound:    // "Chillson: value's parent seems to be a JSON object, but the field cannot be found."
        ...
    case BadValueType:    // "Chillson: retrieved value cannot be converted to the requested type."
        ...
    case NullLeaf:        // "Chillson: null leaf encountered in the structure"
        ...
    }

Chillson is MIT-licensed (see LICENSE). Pull requests, general suggestions (also regarding quality of documentation) and filing issues
are welcome.
*/
package chillson

import (
	"regexp"
	"strconv"
	"strings"
)

type chillsonErr int

const (
	InvalidPath chillsonErr = iota
	OutOfRange
	ParentNotObject
	FieldNotFound
	BadValueType
	NullLeaf
)

func (err chillsonErr) Error() string {
	switch err {
	case InvalidPath:
		return "Chillson: cannot parse nonempty value path, did you forget about square brackets?"
	case OutOfRange:
		return "Chillson: value's parent seems to be a JSON array, but the index is out of range."
	case ParentNotObject:
		return "Chillson: value's parent is neither JSON object nor array."
	case FieldNotFound:
		return "Chillson: value's parent seems to be a JSON object, but the field cannot be found."
	case BadValueType:
		return "Chillson: retrieved value cannot be converted to the requested type."
	case NullLeaf:
		return "Chillson: null leaf encountered in the structure"
	}
	return "Undefined Chillson error."
}

/* Son wraps an unmarshaled JSON document.*/
type Son struct {
	Data (interface{})
}

/* Get() returns value from given location in Son data. Object keys and array indices should be both enclosed in
[square brackets], WITHOUT "quotation marks". String indices (= object keys) can be arbitrary JSON strings as in
JSON source, but they shouldn't contain square brackets [ ]. */
func (c *Son) Get(path string) (interface{}, error) {
	format := regexp.MustCompile("(?:\\[([^\\[\\]]+)\\])+?")
	matches := format.FindAllString(path, -1)
	if len(matches) == 0 && len(path) != 0 {
		return nil, InvalidPath
	}
	var currLeaf *(interface{}) = &(*c).Data
	for _, label := range matches {
		if currLeaf == nil {
			return nil, NullLeaf
		}
		label = strings.Trim(label, "[]")
		// If label is parse'able to integer, try to convert the parent to JSON array (= go slice).
		if numIndex, err := strconv.Atoi(label); err == nil {
			leafArr, ok := (*currLeaf).([]interface{})
			if ok {
				if numIndex < len(leafArr) {
					currLeaf = &(leafArr[numIndex])
					continue
				}
				return nil, OutOfRange
			}
			// If leaf isn't an array, try to parse it as JSON object...
		}
		leafObj, ok := (*currLeaf).(map[string]interface{})
		if !ok {
			return nil, ParentNotObject
		}
		val, ok := leafObj[label]
		if !ok {
			return nil, FieldNotFound
		}
		currLeaf = &val
	}
	return *currLeaf, nil
}

func (c *Son) GetArr(path string) ([]interface{}, error) {
	val, err := (*c).Get(path)
	if err != nil {
		return nil, err
	}
	arr, ok := val.([]interface{})
	if !ok {
		return nil, BadValueType
	}
	return arr, nil
}

func (c *Son) GetBool(path string) (bool, error) {
	val, err := (*c).Get(path)
	if err != nil {
		return false, err
	}
	truth, ok := val.(bool)
	if !ok {
		return false, BadValueType
	}
	return truth, nil
}

func (c *Son) GetFloat(path string) (float64, error) {
	val, err := (*c).Get(path)
	if err != nil {
		return -1, err
	}
	num, ok := val.(float64)
	if !ok {
		return -1, BadValueType
	}
	return num, nil
}

func (c *Son) GetInt(path string) (int, error) {
	num, err := (*c).GetFloat(path)
	if err != nil {
		return -1, err
	}
	return int(num), nil
}

func (c *Son) GetStr(path string) (string, error) {
	val, err := (*c).Get(path)
	if err != nil {
		return "", err
	}
	str, ok := val.(string)
	if !ok {
		return "", BadValueType
	}
	return str, nil
}

func (c *Son) GetObj(path string) (map[string]interface{}, error) {
	val, err := (*c).Get(path)
	if err != nil {
		return nil, err
	}
	obj, ok := val.(map[string]interface{})
	if !ok {
		return nil, BadValueType
	}
	return obj, nil
}

func (c *Son) Require(path string) interface{} {
	ret, err := c.Get(path)
	if err != nil {
		panic(err.Error())
	}
	return ret
}

func (c *Son) RequireArr(path string) []interface{} {
	ret, err := c.GetArr(path)
	if err != nil {
		panic(err.Error())
	}
	return ret
}

func (c *Son) RequireBool(path string) bool {
	ret, err := c.GetBool(path)
	if err != nil {
		panic(err.Error())
	}
	return ret
}

func (c *Son) RequireFloat(path string) float64 {
	ret, err := c.GetFloat(path)
	if err != nil {
		panic(err.Error())
	}
	return ret
}

func (c *Son) RequireInt(path string) int {
	ret, err := c.GetInt(path)
	if err != nil {
		panic(err.Error())
	}
	return ret
}

func (c *Son) RequireObj(path string) map[string]interface{} {
	ret, err := c.GetObj(path)
	if err != nil {
		panic(err.Error())
	}
	return ret
}

func (c *Son) RequireStr(path string) string {
	ret, err := c.GetStr(path)
	if err != nil {
		panic(err.Error())
	}
	return ret
}
