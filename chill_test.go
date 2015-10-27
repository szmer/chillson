package chillson

import (
	"encoding/json"
	"fmt"
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
