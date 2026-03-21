package matcher

import (
	"regexp"
	"testing"
)

// MockMatchable is a test implementation of the Matchable interface
type MockMatchable struct {
	value string
}

func (m MockMatchable) Matchable() string {
	return m.value
}

func TestBuildPattern_EmptyArgs(t *testing.T) {
	pattern := BuildPattern([]string{})
	if pattern == nil {
		t.Error("expected non-nil pattern for empty args")
	}

	// Should match everything with ".*"
	if !pattern.MatchString("anything") {
		t.Error("empty args pattern should match any string")
	}
	if !pattern.MatchString("") {
		t.Error("empty args pattern should match empty string")
	}
}

func TestBuildPattern_SingleArg(t *testing.T) {
	pattern := BuildPattern([]string{"hello"})
	if pattern == nil {
		t.Error("expected non-nil pattern")
	}

	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", true},
		{"Hello", true},           // case insensitive
		{"HELLO", true},           // case insensitive
		{"hello world", true},     // partial match
		{"Say hello", true},       // partial match
		{"world", false},          // no match
		{"h e l l o", false},      // characters not in sequence
	}

	for _, test := range tests {
		if got := pattern.MatchString(test.input); got != test.expected {
			t.Errorf("BuildPattern([\"hello\"]).MatchString(%q) = %v, want %v", test.input, got, test.expected)
		}
	}
}

func TestBuildPattern_MultipleArgs(t *testing.T) {
	pattern := BuildPattern([]string{"hello", "world"})
	if pattern == nil {
		t.Error("expected non-nil pattern")
	}

	tests := []struct {
		input    string
		expected bool
	}{
		{"hello world", true},
		{"Hello World", true},           // case insensitive
		{"helloworld", true},            // no space between
		{"hello beautiful world", true}, // characters between
		{"world hello", false},          // wrong order
		{"hello", false},                // missing second part
		{"world", false},                // missing first part
	}

	for _, test := range tests {
		if got := pattern.MatchString(test.input); got != test.expected {
			t.Errorf("BuildPattern([\"hello\", \"world\"]).MatchString(%q) = %v, want %v", test.input, got, test.expected)
		}
	}
}

func TestBuildPattern_CaseInsensitive(t *testing.T) {
	pattern := BuildPattern([]string{"Test"})

	tests := []struct {
		input string
		want  bool
	}{
		{"test", true},
		{"TEST", true},
		{"Test", true},
		{"TeSt", true},
	}

	for _, test := range tests {
		if got := pattern.MatchString(test.input); got != test.want {
			t.Errorf("case insensitivity failed for input %q: got %v, want %v", test.input, got, test.want)
		}
	}
}

func TestBuildPattern_SpecialCharacters(t *testing.T) {
	// Test that the pattern handles literal strings (not special regex chars)
	// The current implementation joins with ".*" so we test basic joining
	pattern := BuildPattern([]string{"a", "b"})
	if pattern == nil {
		t.Error("expected non-nil pattern")
	}

	if !pattern.MatchString("a123b") {
		t.Error("pattern should match 'a' followed by any characters followed by 'b'")
	}
}

func TestMatchItems_EmptySlice(t *testing.T) {
	pattern := regexp.MustCompile(".*")
	items := []MockMatchable{}
	result := MatchItems(items, pattern)

	if len(result) != 0 {
		t.Errorf("MatchItems with empty slice should return empty result, got %d items", len(result))
	}
}

func TestMatchItems_NoMatches(t *testing.T) {
	pattern := regexp.MustCompile("xyz")
	items := []MockMatchable{
		{value: "abc"},
		{value: "def"},
		{value: "ghi"},
	}

	result := MatchItems(items, pattern)
	if len(result) != 0 {
		t.Errorf("MatchItems with no matches should return empty result, got %d items", len(result))
	}
}

func TestMatchItems_PartialMatches(t *testing.T) {
	pattern := regexp.MustCompile("(?i)test")
	items := []MockMatchable{
		{value: "testing"},
		{value: "contest"},
		{value: "hello"},
		{value: "Test"},
		{value: "notmatching"},
	}

	result := MatchItems(items, pattern)
	if len(result) != 3 {
		t.Errorf("MatchItems expected 3 matches, got %d", len(result))
	}

	expected := []string{"testing", "contest", "Test"}
	for i, item := range result {
		if item.value != expected[i] {
			t.Errorf("match %d: got %q, want %q", i, item.value, expected[i])
		}
	}
}

func TestMatchItems_AllMatch(t *testing.T) {
	pattern := regexp.MustCompile(".*")
	items := []MockMatchable{
		{value: "one"},
		{value: "two"},
		{value: "three"},
	}

	result := MatchItems(items, pattern)
	if len(result) != 3 {
		t.Errorf("MatchItems with wildcard pattern should match all items, got %d", len(result))
	}
}

func TestMatchItems_CaseSensitiveByDefault(t *testing.T) {
	// Pattern without (?i) flag should be case sensitive
	pattern := regexp.MustCompile("test")
	items := []MockMatchable{
		{value: "test"},
		{value: "Test"},
		{value: "TEST"},
	}

	result := MatchItems(items, pattern)
	if len(result) != 1 {
		t.Errorf("case sensitive pattern should match only 'test', got %d matches", len(result))
	}
	if result[0].value != "test" {
		t.Errorf("expected 'test', got %q", result[0].value)
	}
}

func TestMatchItems_PreservesOrder(t *testing.T) {
	pattern := regexp.MustCompile("(?i)^a")
	items := []MockMatchable{
		{value: "apple"},
		{value: "banana"},
		{value: "apricot"},
		{value: "carrot"},
		{value: "avocado"},
	}

	result := MatchItems(items, pattern)
	if len(result) != 3 {
		t.Errorf("expected 3 matches, got %d", len(result))
	}

	expected := []string{"apple", "apricot", "avocado"}
	for i, item := range result {
		if item.value != expected[i] {
			t.Errorf("item %d: got %q, want %q", i, item.value, expected[i])
		}
	}
}

func TestMatchItems_WithComplexPattern(t *testing.T) {
	// Test with a pattern created using BuildPattern
	buildPattern := BuildPattern([]string{"src", "test"})
	items := []MockMatchable{
		{value: "src/component/test.go"},
		{value: "src/test.go"},
		{value: "test/main.go"},
		{value: "src/util.go"},
	}

	result := MatchItems(items, buildPattern)
	if len(result) != 2 {
		t.Errorf("expected 2 matches for pattern from BuildPattern, got %d", len(result))
	}
}
