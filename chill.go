/*
Chillson package is a convenience tool, providing way to process unmarshaled JSON data in schema-agnostic way. Thus you don't have to
specify type structure conforming to the expected JSON structure, or write countless Golang type assertions. The latter is automated by
chillson's Son type methods, like in the example below.

    import (
	    "chillson"
	    "encoding/json"
    )

    var jsonData interface{}
    json.Unmarshal([]byte(rawJson), &jsonData)
    chill := chillson.Son{jsonData}

    // now you can use Son-type variable like this:
    strField, err := chill.GetStr("[gophers][0][name]")
    intField, err := chill.GetInt("[gophers][0][weight]")

    // you can also spawn more specific Son{}'s to save some type assertions:
    gophersTable, err := chill.GetArr("[gophers]")
    for i := 0; i < len(gophersTable); i++ {
	    gophersRow := chillson.Son{gophersTable[i]}
	    strField, err = chill.GetStr("[name]")
	    intField, err = chill.GetInt("[weight]")
    }

Chillson is MIT-licensed (see LICENSE). Pull requests, general suggestions (also regarding quality of documentation) and filing issues
are welcome.
*/
package chillson

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/* Son wraps an unmarshaled JSON document.*/
type Son struct {
	Data (interface{})
}

/* Get() returns value from given location in Son data. Object keys and array indices should be both enclosed in
[square brackets], WITHOUT "quotation marks". String indices (object keys) can be arbitrary, but they shouldn't
contain square brackets ([, ]). */
func (c *Son) Get(path string) (*(interface{}), error) {
	format := regexp.MustCompile("(?:\\[([^\\[\\]]+)\\])+?")
	matches := format.FindAllString(path, -1)
	if len(matches) == 0 && len(path) != 0 {
		return nil, errors.New(fmt.Sprintf("No indices recognized in %s, did you forget about square brackets?", path))
	}
	var currLeaf *(interface{}) = &(*c).Data
	for _, label := range matches {
		label = strings.Trim(label, "[]")
		// If label is parse'able to integer, try to convert the parent to JSON array (= go slice).
		if numIndex, err := strconv.Atoi(label); err == nil {
			leafArr, ok := (*currLeaf).([]interface{})
			if ok {
				if numIndex < len(leafArr) {
					currLeaf = &(leafArr[numIndex])
					continue
				}
				return nil, errors.New(fmt.Sprintf("Chillson: %s is out of range of array %v", label, currLeaf))
			}
			// If leaf isn't an array, try to parse it as JSON object...
		}
		leafObj, ok := (*currLeaf).(map[string]interface{})
		if !ok {
			return nil, errors.New(fmt.Sprintf("Chillson: %v (parent node of %s) cannot be parsed as JSON object.", currLeaf, label))
		}
		val, ok := leafObj[label]
		if !ok {
			return nil, errors.New(fmt.Sprintf("Chillson: parent object %v doesn't contain entry labeled as %s", currLeaf, label))
		}
		currLeaf = &val
	}
	return currLeaf, nil
}

func (c *Son) GetArr(path string) ([]interface{}, error) {
	val, err := (*c).Get(path)
	if err != nil {
		return nil, err
	}
	arr, ok := (*val).([]interface{})
	if !ok {
		return nil, errors.New(fmt.Sprintf("Son: value cannot be converted to a []interface{}: %v", val))
	}
	return arr, nil
}

func (c *Son) GetFloat(path string) (float64, error) {
	val, err := (*c).Get(path)
	if err != nil {
		return -1, err
	}
	num, ok := (*val).(float64)
	if !ok {
		return -1, errors.New(fmt.Sprintf("Son: value cannot be converted to a float64: %v", val))
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
	str, ok := (*val).(string)
	if !ok {
		return "", errors.New(fmt.Sprintf("Son: value cannot be converted to a string: %v", val))
	}
	return str, nil
}

func (c *Son) GetObj(path string) (map[string]interface{}, error) {
	val, err := (*c).Get(path)
	if err != nil {
		return nil, err
	}
	obj, ok := (*val).(map[string]interface{})
	if !ok {
		return nil, errors.New(fmt.Sprintf("Son: value cannot be converted to a map[string]interface{}: %v", val))
	}
	return obj, nil
}
