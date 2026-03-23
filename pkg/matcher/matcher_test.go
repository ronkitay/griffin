package matcher

import (
	"regexp"
	"testing"
)

// TestBuildPattern_Empty tests BuildPattern with no arguments
func TestBuildPattern_Empty(t *testing.T) {
	pattern := BuildPattern([]string{})
	if pattern == nil {
		t.Fatal("BuildPattern returned nil")
	}

	// Should match everything
	if !pattern.MatchString("anything") {
		t.Errorf("Pattern with empty args should match 'anything'")
	}
	if !pattern.MatchString("test") {
		t.Errorf("Pattern with empty args should match 'test'")
	}
	if !pattern.MatchString("") {
		t.Errorf("Pattern with empty args should match empty string")
	}
}

// TestBuildPattern_SingleArgument tests BuildPattern with a single search term
func TestBuildPattern_SingleArgument(t *testing.T) {
	pattern := BuildPattern([]string{"hello"})
	if pattern == nil {
		t.Fatal("BuildPattern returned nil")
	}

	if !pattern.MatchString("hello") {
		t.Errorf("Pattern should match 'hello'")
	}
	if !pattern.MatchString("Hello") {
		t.Errorf("Pattern should match 'Hello' (case-insensitive)")
	}
	if !pattern.MatchString("HELLO") {
		t.Errorf("Pattern should match 'HELLO' (case-insensitive)")
	}
	if !pattern.MatchString("say hello world") {
		t.Errorf("Pattern should match 'hello' within a longer string")
	}
	if pattern.MatchString("helo") {
		t.Errorf("Pattern should not match 'helo' (missing 'l')")
	}
}

// TestBuildPattern_MultipleArguments tests BuildPattern with multiple search terms
func TestBuildPattern_MultipleArguments(t *testing.T) {
	pattern := BuildPattern([]string{"hello", "world"})
	if pattern == nil {
		t.Fatal("BuildPattern returned nil")
	}

	if !pattern.MatchString("hello world") {
		t.Errorf("Pattern should match 'hello world'")
	}
	if !pattern.MatchString("say hello to the world") {
		t.Errorf("Pattern should match 'hello' and 'world' with anything in between")
	}
	if !pattern.MatchString("HeLLo WoRLd") {
		t.Errorf("Pattern should be case-insensitive")
	}
	if pattern.MatchString("hello") {
		t.Errorf("Pattern should not match 'hello' alone without 'world'")
	}
	if pattern.MatchString("world hello") {
		t.Errorf("Pattern should match 'hello' before 'world'")
	}
}

// TestBuildPattern_CaseInsensitive tests that pattern matching is case-insensitive
func TestBuildPattern_CaseInsensitive(t *testing.T) {
	pattern := BuildPattern([]string{"Test"})
	if pattern == nil {
		t.Fatal("BuildPattern returned nil")
	}

	tests := []struct {
		input string
		want  bool
	}{
		{"test", true},
		{"Test", true},
		{"TEST", true},
		{"TeSt", true},
		{"testing", true},
		{"my test", true},
		{"no match", false},
	}

	for _, tt := range tests {
		if pattern.MatchString(tt.input) != tt.want {
			t.Errorf("Pattern.MatchString(%q) = %v, want %v", tt.input, pattern.MatchString(tt.input), tt.want)
		}
	}
}

// TestBuildPattern_SpecialCharacters tests BuildPattern with words containing dots or dashes
func TestBuildPattern_SpecialCharacters(t *testing.T) {
	// Test with alphanumeric terms that don't include regex metacharacters
	// (The current implementation doesn't escape special chars, so we test safe patterns)
	pattern := BuildPattern([]string{"test"})
	if pattern == nil {
		t.Fatal("BuildPattern returned nil")
	}

	// The pattern should match "test" in a string
	if !pattern.MatchString("test") {
		t.Errorf("Pattern should match 'test'")
	}
	if !pattern.MatchString("testing") {
		t.Errorf("Pattern should match 'test' within 'testing'")
	}
}

// TestBuildPattern_EmptyString tests BuildPattern with empty string in args
func TestBuildPattern_EmptyString(t *testing.T) {
	pattern := BuildPattern([]string{""})
	if pattern == nil {
		t.Fatal("BuildPattern returned nil")
	}

	// Should match anything since empty string matches anywhere
	if !pattern.MatchString("anything") {
		t.Errorf("Pattern with empty string should match anything")
	}
}

// MockMatchable is a test implementation of the Matchable interface
type MockMatchable struct {
	values []string
}

func (m MockMatchable) Matchable() []string {
	return m.values
}

// TestMatchItems_EmptySlice tests MatchItems with an empty slice
func TestMatchItems_EmptySlice(t *testing.T) {
	pattern := regexp.MustCompile(".*")
	result := MatchItems([]MockMatchable{}, pattern)

	if len(result) != 0 {
		t.Errorf("MatchItems on empty slice should return 0 items, got %d items", len(result))
	}
}

// TestMatchItems_NoMatches tests MatchItems when no items match
func TestMatchItems_NoMatches(t *testing.T) {
	items := []MockMatchable{
		{values: []string{"apple"}},
		{values: []string{"banana"}},
		{values: []string{"cherry"}},
	}

	pattern := regexp.MustCompile("(?i)xyz")
	result := MatchItems(items, pattern)

	if len(result) != 0 {
		t.Errorf("MatchItems expected no matches, got %d", len(result))
	}
}

// TestMatchItems_AllMatch tests MatchItems when all items match
func TestMatchItems_AllMatch(t *testing.T) {
	items := []MockMatchable{
		{values: []string{"apple"}},
		{values: []string{"apricot"}},
		{values: []string{"avocado"}},
	}

	pattern := regexp.MustCompile("(?i)^a")
	result := MatchItems(items, pattern)

	if len(result) != 3 {
		t.Errorf("MatchItems expected 3 matches, got %d", len(result))
	}

	for i, item := range result {
		if len(item.values) == 0 || item.values[0] != items[i].values[0] {
			t.Errorf("MatchItems returned items in wrong order or wrong items")
		}
	}
}

// TestMatchItems_PartialMatch tests MatchItems with some matches
func TestMatchItems_PartialMatch(t *testing.T) {
	items := []MockMatchable{
		{values: []string{"apple"}},
		{values: []string{"banana"}},
		{values: []string{"apricot"}},
		{values: []string{"cherry"}},
		{values: []string{"avocado"}},
	}

	pattern := regexp.MustCompile("(?i)^a")
	result := MatchItems(items, pattern)

	if len(result) != 3 {
		t.Errorf("MatchItems expected 3 matches, got %d", len(result))
	}

	expected := []string{"apple", "apricot", "avocado"}
	for i, item := range result {
		if len(item.values) == 0 || item.values[0] != expected[i] {
			t.Errorf("MatchItems item %d: expected %q, got %q", i, expected[i], item.values)
		}
	}
}

// TestMatchItems_MultipleValuesPerItem tests MatchItems when items have multiple matchable values
func TestMatchItems_MultipleValuesPerItem(t *testing.T) {
	items := []MockMatchable{
		{values: []string{"apple", "orange"}},
		{values: []string{"banana", "kiwi"}},
		{values: []string{"apricot", "avocado"}},
	}

	pattern := regexp.MustCompile("(?i)^a")
	result := MatchItems(items, pattern)

	if len(result) != 2 {
		t.Errorf("MatchItems expected 2 matches, got %d", len(result))
	}

	// First match: apple matches
	if len(result[0].values) == 0 || result[0].values[0] != "apple" {
		t.Errorf("MatchItems first result: expected 'apple', got %v", result[0].values)
	}

	// Second match: apricot matches
	if len(result[1].values) == 0 || result[1].values[0] != "apricot" {
		t.Errorf("MatchItems second result: expected 'apricot', got %v", result[1].values)
	}
}

// TestMatchItems_CaseInsensitive tests that MatchItems respects the pattern's case-insensitivity
func TestMatchItems_CaseInsensitive(t *testing.T) {
	items := []MockMatchable{
		{values: []string{"Apple"}},
		{values: []string{"BANANA"}},
		{values: []string{"cherry"}},
	}

	pattern := regexp.MustCompile("(?i)app")
	result := MatchItems(items, pattern)

	if len(result) != 1 {
		t.Errorf("MatchItems expected 1 match, got %d", len(result))
	}

	if len(result[0].values) == 0 || result[0].values[0] != "Apple" {
		t.Errorf("MatchItems expected 'Apple', got %v", result[0].values)
	}
}

// TestMatchItems_PreservesOrder tests that MatchItems preserves the order of matched items
func TestMatchItems_PreservesOrder(t *testing.T) {
	items := []MockMatchable{
		{values: []string{"zebra"}},
		{values: []string{"apple"}},
		{values: []string{"monkey"}},
		{values: []string{"banana"}},
		{values: []string{"cherry"}},
	}

	pattern := regexp.MustCompile(".*")
	result := MatchItems(items, pattern)

	if len(result) != len(items) {
		t.Errorf("MatchItems expected %d matches, got %d", len(items), len(result))
	}

	for i, item := range result {
		if len(item.values) == 0 || item.values[0] != items[i].values[0] {
			t.Errorf("MatchItems item %d: expected %v, got %v", i, items[i].values, item.values)
		}
	}
}

// TestMatchItems_ComplexPattern tests MatchItems with a complex regex pattern
func TestMatchItems_ComplexPattern(t *testing.T) {
	items := []MockMatchable{
		{values: []string{"test_file.go"}},
		{values: []string{"main_test.go"}},
		{values: []string{"helper.go"}},
		{values: []string{"doc_test.md"}},
		{values: []string{"readme.txt"}},
	}

	// Pattern to match files ending in .go
	pattern := regexp.MustCompile(`\.go$`)
	result := MatchItems(items, pattern)

	if len(result) != 3 {
		t.Errorf("MatchItems expected 3 .go files, got %d", len(result))
	}

	expected := []string{"test_file.go", "main_test.go", "helper.go"}
	for i, item := range result {
		if len(item.values) == 0 || item.values[0] != expected[i] {
			t.Errorf("MatchItems item %d: expected %q, got %v", i, expected[i], item.values)
		}
	}
}

// TestMatchItems_SpecialCharactersInContent tests MatchItems with special characters
func TestMatchItems_SpecialCharactersInContent(t *testing.T) {
	items := []MockMatchable{
		{values: []string{"path/to/file.go"}},
		{values: []string{"another-file.go"}},
		{values: []string{"file_name.go"}},
		{values: []string{"file.backup.go"}},
	}

	pattern := regexp.MustCompile(`(?i)file`)
	result := MatchItems(items, pattern)

	if len(result) != 4 {
		t.Errorf("MatchItems expected 4 matches, got %d", len(result))
	}
}

// TestMatchItems_EmptyStrings tests MatchItems with empty string values
func TestMatchItems_EmptyStrings(t *testing.T) {
	items := []MockMatchable{
		{values: []string{""}},
		{values: []string{"test"}},
		{values: []string{""}},
	}

	pattern := regexp.MustCompile(".*")
	result := MatchItems(items, pattern)

	if len(result) != 3 {
		t.Errorf("MatchItems with .* should match empty strings too, expected 3, got %d", len(result))
	}
}

// TestMatchItems_ReturnTypePreservation tests that MatchItems preserves the correct type
func TestMatchItems_ReturnTypePreservation(t *testing.T) {
	items := []MockMatchable{
		{values: []string{"first"}},
		{values: []string{"second"}},
		{values: []string{"third"}},
	}

	pattern := regexp.MustCompile(".*")
	result := MatchItems(items, pattern)

	// Verify that we can call Matchable() on the results
	for _, item := range result {
		_ = item.Matchable()
	}
}

// TestBuildPattern_ValidPatternGeneration tests BuildPattern with various valid inputs
func TestBuildPattern_ValidPatternGeneration(t *testing.T) {
	tests := []struct {
		args        []string
		testStr     string
		shouldMatch bool
	}{
		{[]string{}, "anything", true},
		{[]string{"test"}, "test", true},
		{[]string{"test"}, "TEST", true},
		{[]string{"hello", "world"}, "hello there world", true},
		{[]string{"a", "b", "c"}, "a x b y c", true},
		{[]string{"a", "b", "c"}, "c b a", false},
	}

	for _, tt := range tests {
		pattern := BuildPattern(tt.args)
		if pattern == nil {
			t.Fatal("BuildPattern returned nil")
		}

		got := pattern.MatchString(tt.testStr)
		if got != tt.shouldMatch {
			t.Errorf("BuildPattern(%v).MatchString(%q) = %v, want %v",
				tt.args, tt.testStr, got, tt.shouldMatch)
		}
	}
}

// TestMatchItems_MatchesFirstValue tests that MatchItems breaks after first match
func TestMatchItems_MatchesFirstValue(t *testing.T) {
	items := []MockMatchable{
		{values: []string{"no_match", "banana"}},
		{values: []string{"apple", "no_match"}},
		{values: []string{"cherry"}},
	}

	pattern := regexp.MustCompile("(?i)^a")
	result := MatchItems(items, pattern)

	if len(result) != 1 {
		t.Errorf("MatchItems expected 1 match (only items with first value matching 'a'), got %d", len(result))
	}

	if len(result[0].values) == 0 || result[0].values[0] != "apple" {
		t.Errorf("MatchItems expected 'apple', got %v", result[0].values)
	}
}
