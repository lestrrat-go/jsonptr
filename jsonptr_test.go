package jsonptr_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/lestrrat-go/jsonptr"
)

func TestParse(t *testing.T) {
	const src = `{"a": 1, "b": 2, "c": {"d": 3, "e": 4}, "f": [5, 6, "g", {"h": 7}, ["i", 8, 9]]}`
	c, err := jsonptr.Parse([]byte(src))

	if err != nil {
		t.Fatalf("failed to parse JSON: %s", err)
	}

	testcases := []struct {
		Path     string
		Expected interface{}
		Source   []byte
		Error    bool
	}{
		{
			Path: "/a",
			// because we're doing Unmarshal(&interface{}), we need to use float64 here
			Expected: float64(1),
			Source:   []byte{'1'},
		},
		{
			Path: "/b",
			// because we're doing Unmarshal(&interface{}), we need to use float64 here
			Expected: float64(2),
			Source:   []byte{'2'},
		},
		{
			Path: "/c",
			Expected: map[string]interface{}{
				"d": float64(3),
				"e": float64(4),
			},
			Source: []byte(`{"d": 3, "e": 4}`),
		},
		{
			Path:     "/c/d",
			Expected: float64(3),
			Source:   []byte{'3'},
		},
		{
			Path:     "/c/e",
			Expected: float64(4),
			Source:   []byte{'4'},
		},
		{
			Path: "/f",
			Expected: []interface{}{
				float64(5),
				float64(6),
				"g",
				map[string]interface{}{
					"h": float64(7),
				},
				[]interface{}{
					"i",
					float64(8),
					float64(9),
				},
			},
			Source: []byte(`[5, 6, "g", {"h": 7}, ["i", 8, 9]]`),
		},
		{
			Path:     "/f/0",
			Expected: float64(5),
			Source:   []byte{'5'},
		},
		{
			Path:     "/f/1",
			Expected: float64(6),
			Source:   []byte{'6'},
		},
		{
			Path:     "/f/2",
			Expected: "g",
			Source:   []byte{'"', 'g', '"'},
		},
		{
			Path: "/f/3",
			Expected: map[string]interface{}{
				"h": float64(7),
			},
			Source: []byte(`{"h": 7}`),
		},
		{
			Path:     "/f/3/h",
			Expected: float64(7),
			Source:   []byte{'7'},
		},
		{
			Path: "/f/4",
			Expected: []interface{}{
				"i",
				float64(8),
				float64(9),
			},
			Source: []byte(`["i", 8, 9]`),
		},
		{
			Path:     "/f/4/0",
			Expected: "i",
			Source:   []byte{'"', 'i', '"'},
		},
		{
			Path:     "/f/4/1",
			Expected: float64(8),
			Source:   []byte{'8'},
		},
		{
			Path:     "/f/4/2",
			Expected: float64(9),
			Source:   []byte{'9'},
		},
		{
			Path:  "/g",
			Error: true,
		},
		{
			Path:  "/c/g",
			Error: true,
		},
		{
			Path:  "/f/5",
			Error: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Path, func(t *testing.T) {
			var v interface{}

			if tc.Error {
				if _, err := c.Get(tc.Path); err == nil {
					t.Errorf("c.Get should fail")
				}
				if err := c.Unmarshal(tc.Path, &v); err == nil {
					t.Errorf("c.Unmarshal should fail")
				}
			} else {
				fragment, err := c.Get(tc.Path)
				if err != nil {
					t.Fatalf("failed to find %s: %s", tc.Path, err)
				}
				if !bytes.Equal(fragment, tc.Source) {
					t.Fatalf("expected %s to be %s, but got %s", tc.Path, tc.Expected, fragment)
				}

				if err := c.Unmarshal(tc.Path, &v); err != nil {
					t.Fatalf("failed to find %s: %s", tc.Path, err)
				}

				if !reflect.DeepEqual(v, tc.Expected) {
					t.Fatalf("expected %s to be %#v, but got %#v", tc.Path, tc.Expected, v)
				}
				t.Logf("%s = %#v", tc.Path, v)
			}
		})
	}
}
