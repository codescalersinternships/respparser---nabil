package respparser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      []interface{}
		expectedError bool
	}{
		{
			name:  "Simple String",
			input: "+OK\r\n",
			expected: []interface{}{
				"OK",
			},
			expectedError: false,
		},
		{
			name:  "Error Message",
			input: "-Error message\r\n",
			expected: []interface{}{
				fmt.Errorf("Error message"),
			},
			expectedError: false,
		},
		{
			name:  "Integer",
			input: ":1000\r\n",
			expected: []interface{}{
				int64(1000),
			},
			expectedError: false,
		},
		{
			name:  "Bulk String",
			input: "$6\r\nfoobar\r\n",
			expected: []interface{}{
				"foobar",
			},
			expectedError: false,
		},
		{
			name:  "Null Bulk String",
			input: "$-1\r\n",
			expected: []interface{}{
				nil,
			},
			expectedError: false,
		},
		{
			name:  "Array of Bulk Strings",
			input: "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			expected: []interface{}{
				[]interface{}{"foo", "bar"},
			},
			expectedError: false,
		},
		{
			name:  "Array of Integers",
			input: "*3\r\n:1\r\n:2\r\n:3\r\n",
			expected: []interface{}{
				[]interface{}{int64(1), int64(2), int64(3)},
			},
			expectedError: false,
		},
		{
			name:  "Null Element",
			input: "_\r\n",
			expected: []interface{}{
				nil,
			},
			expectedError: false,
		},
		{
			name:  "Boolean True",
			input: "#t\r\n",
			expected: []interface{}{
				true,
			},
			expectedError: false,
		},
		{
			name:  "Boolean False",
			input: "#f\r\n",
			expected: []interface{}{
				false,
			},
			expectedError: false,
		},
		{
			name:  "Floating Point Number",
			input: ",123.45\r\n",
			expected: []interface{}{
				float64(123.45),
			},
			expectedError: false,
		},
		{
			name:  "Bulk Error Message",
			input: "!21\r\nSYNTAX invalid syntax\r\n",
			expected: []interface{}{
				fmt.Errorf("SYNTAX invalid syntax"),
			},
			expectedError: false,
		},
		{
			name:  "Map with Key-Value Pairs",
			input: "%2\r\n+first\r\n:1\r\n+second\r\n:2\r\n",
			expected: []interface{}{
				map[interface{}]interface{}{
					"first":  int64(1),
					"second": int64(2),
				},
			},
			expectedError: false,
		},
		{
			name:  "nested arrays",
			input: "*3\r\n*3\r\n*3\r\n:1\r\n:2\r\n:3\r\n\r\n:2\r\n:3\r\n\r\n:2\r\n:3\r\n",
			expected : []interface {}{[]interface {}{[]interface {}{[]interface {}{int64(1), int64(2), int64(3)}, int64(2), int64(3)}, int64(2), int64(3)}},
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Parser(test.input)
			if test.expectedError {
				assert.Error(t, err)
			}

			assert.Equal(t, test.expected, result)
		})
	}
}
