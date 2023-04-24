# jsonptr

`github.com/lestrrat-go/jsonptr` is a tool that allows users to query and unmarshal JSON data using JSON Pointers (RFC6901). This modules does not support updating existing JSON data, and only supports querying / unmarshaling.

# SYNOPSIS

```go
package jsonptr_test

import (
	"fmt"

	"github.com/lestrrat-go/jsonptr"
)

func Example() {
	const src = `
	{
		"a": 1,
		"b": 2,
		"c": {"d": 3, "e": 4},
		"f": [
			5,
			6,
			"g",
			{
				"h": 7
			},
			[
				"i",
				8,
				9
			]
		]
	}`
	c, err := jsonptr.Parse([]byte(src))
	if err != nil {
		fmt.Printf("failed to parse JSON: %s", err)
		return
	}

    // If you know the type of object beforehand, you can specify
    // a strongly typed variable for a bit of added efficiency.
    //
    // This example only uses a simple `int`, but you can also use
    // this method to intentionally trigger `UnmarshalJSON`
    // to be called for that particular type, if it implements one.
    var i int
	if err := c.Unmarshal("/a", &i); err != nil {
		fmt.Printf("failed to find /a: %s", err)
		return
	}
	fmt.Printf("/a = %d\n", i)

    // If you do not know the type of value before hand, you can
    // just pass it an `interface{}`. In this case, it will populate
    // `iface` with a `map[string]interface{}`
    var iface interface{}
   	if err := c.Unmarshal("/f", &iface); err != nil {
		fmt.Printf("failed to find /f: %s", err)
		return
	}
	fmt.Printf("/f = %#v\n", iface)
}
```

# Motivation

There are many implementations of RFC6901, but the fact that most of them immediately decoded the JSON value pointed by the pointer specification to a `map[string]interface{}` was problematic.

This is because in my use case I would like to unmarshal these JSON objects into specific Go structs with no exported fields. If we could simply leverage the `encoding/json` mechanics, the following pseudocode is all we would need:

```go
type Foo struct {
    // unexported fields
}

func (f *Foo) UnmarshalJSON(data []byte) error {
    // use json.Decoder or some other trick to properly assign to unexported fields
}

func (f *Foo) MarshalJSON() ([]byte, error) {
    // construct JSON using bytes.Buffer and printing out the values stored in the unexported fields
}
```

But by introducing other modules that involve returning `map`s, we would need to add an extra way to construct the objects from the `map`:

```go
func (f *Foo) ConstructFromMap(m map[string]interface{}) error {
    // do the equivalent of UnmarshalJSON but using the map
}
```

This is doable, but feels a bit awkward. 

Instead of this above, this module first parses the JSON data once, only to record the available paths in the JSON data, as well as the location of the data pointed by the paths in the original JSON data.

So for example, given the following JSON data:

```go
data := []byte(`{"a":{"b":"c"}}`)
```

This module makes a single pass to record the following information:

```
"/a": 5
"/a/b": 10
```

The above signifies that the data pointed to by the pointer "/a" is located at byte 5 (note: 0-based counting), and the data pointed by "/a/b" is located at byte 10.

Given this information, we can get to the raw JSON data pointed by the JSON pointer, and use the normal `encoding/json` unmarshal mechanics against that data.

For example, in the above case, we can do the following to get at the JSON data `"c"` pointed by "/a/b":

```go
c, _ := jsonptr.Parse(data)
target, _ := c.Get("/a/b")
```

Then `target` is a piece of `[]byte` data that contains `"c"`, which you can then safely pass to `json.Unmarshal` and the like.

This module provides one more utility, which allows you to skip having to call `c.Get()` and `json.Unmarshal` yourself:

```go
var s string
_ := c.Unmarshal("/a/b", &s)
```

Notice that we defined `s` as a string, because we already knew that the data pointed by `"/a/b"` was a string value. If we do not know this before hand, we could use `interface{}` instead to let `encoding/json` handle it.

Going back to our original requirement, which is to use JSON pointer with Go structs with custom `UnmarshalJSON` methods, we can simply pass an instance of said object in the `Unmarshal` method to take full advantage of the custom `UnmarshalJSON` method:

```go
f := &Foo{}
_ := c.Unmarshal("/a", f) // would call (*Foo).UnmarshalJSON
```

The beauty of this approach is that `Foo` does not need to be aware of anything about JSON pointers. It can merrily be a plain old Go object, and it us users will be able to query this object out of the JSON data by simply parsing the data one extra time

# TODO

* Test against object keys that resemble array indices
* Test against UTF-8 keys and values outside of the ASCII range