package validate

import "testing"

func TestEventHubClusterName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// basic example
			input:    "ab-12-c",
			expected: true,
		},
		{
			// can't start with a number
			input:    "1ab-12-c",
			expected: false,
		},
		{
			// can't contain underscore
			input:    "hello_world",
			expected: false,
		},
		{
			// can't end with hyphen
			input:    "example-",
			expected: false,
		},
		{
			// can not short than 6 characters
			input:    "hello",
			expected: false,
		},
		{
			// 50 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqr",
			expected: true,
		},
		{
			// 51 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrs",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := EventHubClusterName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
