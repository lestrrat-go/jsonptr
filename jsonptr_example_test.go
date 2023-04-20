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

	var i int
	if err := c.Unmarshal("/a", &i); err != nil {
		fmt.Printf("failed to find /a: %s", err)
		return
	}
	fmt.Printf("/a = %d\n", i)

	if err := c.Unmarshal("/b", &i); err != nil {
		fmt.Printf("failed to find /b: %s", err)
		return
	}
	fmt.Printf("/b = %d\n", i)

	var iface interface{}
	if err := c.Unmarshal("/c", &iface); err != nil {
		fmt.Printf("failed to find /c: %s", err)
		return
	}
	fmt.Printf("/c = %#v\n", iface)

	if err := c.Unmarshal("/c/d", &i); err != nil {
		fmt.Printf("failed to find /c/d: %s", err)
		return
	}
	fmt.Printf("/c/d = %d\n", i)

	if err := c.Unmarshal("/c/e", &i); err != nil {
		fmt.Printf("failed to find /c/e: %s", err)
		return
	}
	fmt.Printf("/c/e = %d\n", i)

	if err := c.Unmarshal("/f", &iface); err != nil {
		fmt.Printf("failed to find /f: %s", err)
		return
	}
	fmt.Printf("/f = %#v\n", iface)

	if err := c.Unmarshal("/f/0", &i); err != nil {
		fmt.Printf("failed to find /f/0: %s", err)
		return
	}
	fmt.Printf("/f/0 = %d\n", i)

	if err := c.Unmarshal("/f/1", &i); err != nil {
		fmt.Printf("failed to find /f/1: %s", err)
		return
	}
	fmt.Printf("/f/1 = %d\n", i)

	var s string
	if err := c.Unmarshal("/f/2", &s); err != nil {
		fmt.Printf("failed to find /f/2: %s", err)
		return
	}
	fmt.Printf("/f/2 = %q\n", s)

	if err := c.Unmarshal("/f/3", &iface); err != nil {
		fmt.Printf("failed to find /f/3: %s", err)
		return
	}
	fmt.Printf("/f/3 = %#v\n", iface)

	if err := c.Unmarshal("/f/3/h", &i); err != nil {
		fmt.Printf("failed to find /f/3/h: %s", err)
		return
	}
	fmt.Printf("/f/3/h = %d\n", i)

	if err := c.Unmarshal("/f/4", &iface); err != nil {
		fmt.Printf("failed to find /f/4: %s", err)
		return
	}
	fmt.Printf("/f/4 = %#v\n", iface)

	if err := c.Unmarshal("/f/4/0", &s); err != nil {
		fmt.Printf("failed to find /f/4/0: %s", err)
		return
	}
	fmt.Printf("/f/4/0 = %q\n", s)

	if err := c.Unmarshal("/f/4/1", &i); err != nil {
		fmt.Printf("failed to find /f/4/1: %s", err)
		return
	}
	fmt.Printf("/f/4/1 = %d\n", i)

	if err := c.Unmarshal("/f/4/2", &i); err != nil {
		fmt.Printf("failed to find /f/4/2: %s", err)
		return
	}
	fmt.Printf("/f/4/2 = %d\n", i)
	// OUTPUT:
	// /a = 1
	// /b = 2
	// /c = map[string]interface {}{"d":3, "e":4}
	// /c/d = 3
	// /c/e = 4
	// /f = []interface {}{5, 6, "g", map[string]interface {}{"h":7}, []interface {}{"i", 8, 9}}
	// /f/0 = 5
	// /f/1 = 6
	// /f/2 = "g"
	// /f/3 = map[string]interface {}{"h":7}
	// /f/3/h = 7
	// /f/4 = []interface {}{"i", 8, 9}
	// /f/4/0 = "i"
	// /f/4/1 = 8
	// /f/4/2 = 9
}
