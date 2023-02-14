package genie

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseFileWithSingleEntry(t *testing.T) {
	result, err := Parse(`
hello = world
`)
	require.Nil(t, err)
	assert.Equal(t, result.CountAllEntries(), 1)
	assert.Equal(t, "world", result.Get("hello"))
}

func TestParseFileWithMultipleEntries(t *testing.T) {
	result, err := Parse(`
hello = world
foo = bar
`)
	require.Nil(t, err)
	assert.Equal(t, result.CountAllEntries(), 2)
	assert.Equal(t, "world", result.Get("hello"))
	assert.Equal(t, "bar", result.Get("foo"))
}

func TestParseFileWithEntryWithEmptyValue(t *testing.T) {
	result, err := Parse("" +
		"key1 =" + "\n" +
		"key2 = " + "\n" +
		"")
	require.Nil(t, err)
	assert.Equal(t, "", result.Get("key1"))
	assert.Equal(t, "", result.Get("key2"))
}

func TestTreatsValueLiterally(t *testing.T) {
	result, err := Parse("" +
		"key1 =     world     " + "\n" +
		"key2 = \t world \t \t" + "\n" +
		"key3 = world # foo" + "\n" +
		"key4 = world = foo" + "\n" +
		`key5 = "world = foo"` + "\n" +
		"")
	require.Nil(t, err)
	assert.Equal(t, "    world     ", result.Get("key1"))
	assert.Equal(t, "\t world \t \t", result.Get("key2"))
	assert.Equal(t, "world # foo", result.Get("key3"))
	assert.Equal(t, "world = foo", result.Get("key4"))
	assert.Equal(t, `"world = foo"`, result.Get("key5"))
}

func TestParseFileWithComments(t *testing.T) {
	result, err := Parse(`
# Hello World.
hello = world

foo = bar

# Test
# 123
test = 123
`)
	require.Nil(t, err)
	assert.Equal(t, result.CountAllEntries(), 3)
	assert.Equal(t, "world", result.Get("hello"))
	assert.Equal(t, "bar", result.Get("foo"))
	assert.Equal(t, "123", result.Get("test"))
}

func TestParseFileWithWindowsLineEndings(t *testing.T) {
	result, err := Parse("hello = world\r\nfoo = bar\r\n\r\n")
	require.Nil(t, err)
	assert.Equal(t, result.CountAllEntries(), 2)
	assert.Equal(t, "world", result.Get("hello"))
	assert.Equal(t, "bar", result.Get("foo"))
}

func TestParseFileWithSections(t *testing.T) {
	result, err := Parse(`
key = value

[section1]
hello = world
foo = bar

[section 2]
test = 123

[section3]
`)
	require.Nil(t, err)
	assert.Equal(t, result.CountAllEntries(), 4)
	assert.Equal(t, "value", result.GetFromSection("", "key"))
	assert.Equal(t, "world", result.GetFromSection("section1", "hello"))
	assert.Equal(t, "bar", result.GetFromSection("section1", "foo"))
	assert.Equal(t, "123", result.GetFromSection("section 2", "test"))
}

func TestAlignmentOnKeySide(t *testing.T) {
	result, err := Parse(`
key         = 123
another_key = 124
`)
	require.Nil(t, err)
	assert.Equal(t, "123", result.Get("key"))
	assert.Equal(t, "124", result.Get("another_key"))
}

func TestSectionWithTrailingWhitespace(t *testing.T) {
	// This is only allowed because it doesn’t create any harm.
	result, err := Parse("" +
		"[section]  \t " + "\n" +
		"key = value" +
		"")
	require.Nil(t, err)
	assert.Equal(t, "value", result.GetFromSection("section", "key"))
}

func TestHandlesAbsentOrEmptyKeysAndSectionsGracefully(t *testing.T) {
	result, err := Parse("")
	require.Nil(t, err)
	assert.Equal(t, "", result.Get(""))
	assert.Equal(t, "", result.Get("  "))
	assert.Equal(t, "", result.Get("hello"))
	assert.Equal(t, "", result.GetFromSection("", "hello"))
	assert.Equal(t, "", result.GetFromSection("  ", "hello"))
	assert.Equal(t, "", result.GetFromSection("foo", ""))
	assert.Equal(t, "", result.GetFromSection("foo", "  "))
	assert.Equal(t, "", result.GetFromSection("foo", "hello"))
}

func TestParseEmptyOrBlankFile(t *testing.T) {
	for _, text := range []string{
		"",
		"  ",
		"\n\n",
		"  \n   \n    \n ",
		"\r\n\n\n\r\n",
	} {
		result, err := Parse(text)
		require.Nil(t, err)
		assert.Equal(t, 0, result.CountAllEntries())
	}
}

func TestErrorCases(t *testing.T) {
	for _, text := range []string{
		// General
		"foo",   // Key without value
		"  foo", // Key without value with leading whitespace
		"\tfoo", // Key without value with leading whitespace

		// Sections
		"[section",             // Invalid section
		"section]",             // Invalid section
		"[]",                   // Empty section name
		"[ ]",                  // Blank section name
		"[\t]",                 // Blank section name
		"[section] # Comment?", // Trailing comment
		"[section] Text?",      // Trailing text
		" [section]",           // Leading whitespace
		"\t[section]",          // Leading whitespace
		"[[section]]",          // Section name itself contains square bracket
		"[se]ct[ion]",          // Section name itself contains square bracket
		"[key] = 123",          // Key cannot look like section
		"[key = 123",           // Key cannot start with square bracket

		// Entries
		"key= 123",      // Missing space delimiter after key
		"key =123",      // Missing space delimiter before value
		"key == 123",    // Invalid delimiter
		"key != 123",    // Invalid delimiter
		" key = 123",    // Leading whitespace
		"\tkey = 123",   // Leading whitespace
		"k e y = 123",   // Key with whitespace
		"k\te\ty = 123", // Key with whitespace

		// Comments
		" # Comment",  // Leading whitespace
		"\t# Comment", // Leading whitespace
	} {
		_, err := Parse(text)
		require.Error(t, err, text)
	}
}
