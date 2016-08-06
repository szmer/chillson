# chillson
--
    import "chillson"

Chillson package is a convenience tool, providing way to process unmarshaled
JSON data in schema-agnostic way. Thus you don't have to specify type structure
conforming to the expected JSON structure, or write countless Golang type
assertions. The latter is automated by chillson's Son type methods, like in the
example below.

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
        }

Chillson is MIT-licensed (see LICENSE).

## Usage

```go
const (
	InvalidPath chillsonErr = iota
	OutOfRange
	ParentNotObject
	FieldNotFound
	BadValueType
)
```

#### type Son

```go
type Son struct {
	Data (interface{})
}
```

Son wraps an unmarshaled JSON document.

#### func (*Son) Get

```go
func (c *Son) Get(path string) (interface{}, error)
```

    Get() returns value from given location in Son data. Object keys and array indices should be both enclosed in
[square brackets], WITHOUT "quotation marks". String indices (= object keys) can
be arbitrary JSON strings as in JSON source, but they shouldn't contain square
brackets [ ].

#### func (*Son) GetArr

```go
func (c *Son) GetArr(path string) ([]interface{}, error)
```

    GetArr returns a JSON array as Golang slice.

#### func (*Son) GetFloat

```go
func (c *Son) GetFloat(path string) (float64, error)
```

    GetFloat returns JSON number as Golang float64.

#### func (*Son) GetInt

```go
func (c *Son) GetInt(path string) (int, error)
```

#### func (*Son) GetObj

```go
func (c *Son) GetObj(path string) (map[string]interface{}, error)
```

#### func (*Son) GetStr

```go
func (c *Son) GetStr(path string) (string, error)
```
