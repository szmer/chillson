package chillson

import (
	"encoding/json"
	"fmt"
    "reflect"
	"testing"
)

func TestGet(t *testing.T) {
	sampleJSON := "{\"foo\":\"bar\",\"arr\":[\"joe\", \"mary\"]}"
	var data interface{}
	err := json.Unmarshal([]byte(sampleJSON), &data)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	chill := Son{data}
	all, err := chill.Get("")
	if _, ok := all.(map[string]interface{}); err != nil || !ok {
		t.Errorf("Empty path doesn't return the entire JSON object.")
	}
	c1, err := chill.Get("[foo]")
	if err != nil || c1 != "bar" {
		t.Errorf(fmt.Sprintf("[foo] doesn't return \"bar\" (%v)."), c1)
	}
	c2, err := chill.Get("[arr][1]")
	if err != nil || c2 != "mary" {
		t.Errorf(fmt.Sprintf("[arr][1] doesn't return \"mary\" (%v)."), c2)
	}
}

func TestGetTyped(t *testing.T) {
        sampleJSON := "{\"foo\":\"bar\",\"arr\":[\"joe\", \"mary\", 42, false, {\"apples\": \"oranges\"}]}"
	var data interface{}
	err := json.Unmarshal([]byte(sampleJSON), &data)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	chill := Son{data}
    s, err := chill.GetStr("[foo]")
    if err != nil || s != "bar" || reflect.ValueOf(s).Kind() != reflect.String {
        t.Errorf(fmt.Sprintf("[foo] doesn't return string bar (%v of kind %v).", s, reflect.ValueOf(s).Kind()))
    }
    i, err := chill.GetInt("[arr][2]")
    if err != nil || i != 42 || reflect.ValueOf(i).Kind() != reflect.Int {
        t.Errorf(fmt.Sprintf("[arr][2] doesn't return integer 42 when requested (%v of kind %v).", i, reflect.ValueOf(i).Kind()))
    }
    f, err := chill.GetFloat("[arr][2]")
    if err != nil || f != 42.0 || reflect.ValueOf(f).Kind() != reflect.Float64 {
        t.Errorf(fmt.Sprintf("[arr][2] doesn't return float64 42.0 when requested (%v of kind %v).", f, reflect.ValueOf(f).Kind()))
    }
    b, err := chill.GetBool("[arr][3]")
    if err != nil || b != false || reflect.ValueOf(b).Kind() != reflect.Bool {
        t.Errorf(fmt.Sprintf("[arr][3] doesn't return bool false when requested (%v of kind %v).", b, reflect.ValueOf(b).Kind()))
    }
    // Test access to the values via an JSON-object Chilllson object.
    o, err := chill.GetObj("[arr][4]")
    if err != nil || reflect.ValueOf(o).Kind() != reflect.Map {
        t.Errorf(fmt.Sprintf("[arr] doesn't return an array (%v).", o, reflect.ValueOf(o).Kind()))
    }
    chillo := Son{o}
    s, err = chillo.GetStr("[apples]")
    if err != nil || s != "oranges" || reflect.ValueOf(s).Kind() != reflect.String {
        t.Errorf(fmt.Sprintf("[apples] on JSON-object object doesn't return string oranges (%v of kind %v).", s, reflect.ValueOf(s).Kind()))
    }
    // Test access to the values via an JSON-array Chillson object.
    a, err := chill.GetArr("[arr]")
    if err != nil || reflect.ValueOf(a).Kind() != reflect.Slice {
        t.Errorf(fmt.Sprintf("[arr] doesn't return an array (%v).", a, reflect.ValueOf(a).Kind()))
    }
    chilla := Son{a}
    s, err = chilla.GetStr("[0]")
    if err != nil || s != "joe" || reflect.ValueOf(s).Kind() != reflect.String {
        t.Errorf(fmt.Sprintf("[0] on array object doesn't return string joe (%v of kind %v).", s, reflect.ValueOf(s).Kind()))
    }
    i, err = chilla.GetInt("[2]")
    if err != nil || i != 42 || reflect.ValueOf(i).Kind() != reflect.Int {
        t.Errorf(fmt.Sprintf("[2] on array object doesn't return integer 42 when requested (%v of kind %v).", i, reflect.ValueOf(i).Kind()))
    }
    f, err = chilla.GetFloat("[2]")
    if err != nil || f != 42.0 || reflect.ValueOf(f).Kind() != reflect.Float64 {
        t.Errorf(fmt.Sprintf("[2] on array object doesn't return float64 42.0 when requested (%v of kind %v).", f, reflect.ValueOf(f).Kind()))
    }
    b, err = chilla.GetBool("[3]")
    if err != nil || b != false || reflect.ValueOf(b).Kind() != reflect.Bool {
        t.Errorf(fmt.Sprintf("[3] on array object doesn't return bool false when requested (%v of kind %v).", b, reflect.ValueOf(b).Kind()))
    }
}
