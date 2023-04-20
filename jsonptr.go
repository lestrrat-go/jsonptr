package jsonptr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"unicode"
)

type Context struct {
	data  []byte
	paths map[string]int64
}

func (c *Context) Register(path string, i int64) {
	c.paths[path] = i
}

// Get retrieves the JSON value at the specified path. The value is
// returned as a byte slice, which you can then pass to json.Unmarshal
func (c *Context) Get(path string) ([]byte, error) {
	v, ok := c.paths[path]
	if !ok {
		return nil, fmt.Errorf("jsonptr.Context.Get: path %s not found", path)
	}

	var msg json.RawMessage
	if err := json.NewDecoder(bytes.NewReader(c.data[v:])).Decode(&msg); err != nil {
		return nil, err
	}
	return []byte(msg), nil
}

// Unmarshal retrieves the value at the specified path. It follows the
// rules of json.Unmarshal, so you can pass in a pointer to a struct
// or a pointer to a slice, etc, and it will be populated with the value.
//
// The operation will fail if the value at the specified path is not
// a valid JSON value.
//
// The `path` argument must start with a `/` character. If you are using
// path notations that include a leading `#`, you must strip it off before
// passing it to this function.
func (c *Context) Unmarshal(path string, value interface{}) error {
	v, ok := c.paths[path]
	if !ok {
		return fmt.Errorf("jsonptr.Context.Get: path %s not found", path)
	}

	if err := json.NewDecoder(bytes.NewReader(c.data[v:])).Decode(value); err != nil {
		return err
	}
	return nil
}

func Parse(data []byte) (*Context, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	c := &Context{
		data:  data,
		paths: make(map[string]int64),
	}
	if err := parse(c, dec, ""); err != nil {
		return nil, err
	}
	return c, nil
}

func parse(c *Context, dec *json.Decoder, path string) error {
	pos := dec.InputOffset()
	// pos is at ":", so we need to move the position to
	// the location of the value
	pos++
	for _, r := range c.data[pos:] {
		if unicode.IsSpace(rune(r)) {
			pos += int64(1)
		}
		break
	}

	tok, err := dec.Token()
	if err != nil {
		return err
	}

	switch tok := tok.(type) {
	case json.Delim:
		switch tok {
		case '{':
			c.Register(path, pos)
			return parseObject(c, dec, path)
		case '[':
			c.Register(path, pos)
			return parseArray(c, dec, path)
		default:
			return fmt.Errorf("jsonptr.parse: unexpected token %v", tok)
		}
	default:
		c.Register(path, pos)
		return nil
	}
}

func parseObject(c *Context, dec *json.Decoder, path string) error {
	for {
		tok, err := dec.Token()
		if err != nil {
			return err
		}

		switch tok := tok.(type) {
		case json.Delim:
			switch tok {
			case '}':
				return nil
			default:
				return fmt.Errorf("unexpected token %v", tok)
			}
		case string:
			if err := parse(c, dec, path+"/"+tok); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unexpected token %v", tok)
		}
	}
}

func parseArray(c *Context, dec *json.Decoder, path string) error {
	for i := 0; ; i++ {
		pos := dec.InputOffset()
		// pos is at "," (iff this is the second element or later)
		// so move the position to the beginning of the next token
		if i > 0 {
			pos++
			for _, r := range c.data[pos:] {
				if unicode.IsSpace(rune(r)) {
					pos += int64(1)
				}
				break
			}
		}

		tok, err := dec.Token()
		if err != nil {
			return err
		}

		switch tok := tok.(type) {
		case json.Delim:
			switch tok {
			case ']':
				return nil
			case '[':
				c.Register(fmt.Sprintf("%s/%d", path, i), pos)
				if err := parseArray(c, dec, fmt.Sprintf("%s/%d", path, i)); err != nil {
					return fmt.Errorf("jsonptr.parseArray: failed to parse array: %s", err)
				}
			case '{':
				c.Register(fmt.Sprintf("%s/%d", path, i), pos)
				if err := parseObject(c, dec, fmt.Sprintf("%s/%d", path, i)); err != nil {
					return fmt.Errorf("jsonptr.parseArray: failed to parse object: %s", err)
				}
			default:
				return fmt.Errorf("jsonptr.parseArray: unexpected token %v", tok)
			}
		default:
			c.Register(fmt.Sprintf("%s/%d", path, i), pos)
		}
	}
}
