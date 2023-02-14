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

func TestIgnoresRedundantWhiteSpaceOnKeySide(t *testing.T) {
	result, err := Parse("" +
		"     hello1      = world\n" +
		"\t\thello2 = world\n" +
		"")
	require.Nil(t, err)
	assert.Equal(t, "world", result.Get("hello1"))
	assert.Equal(t, "world", result.Get("hello2"))
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

[section2]
test = 123

[section3]
`)
	require.Nil(t, err)
	assert.Equal(t, result.CountAllEntries(), 4)
	assert.Equal(t, "value", result.GetFromSection("", "key"))
	assert.Equal(t, "world", result.GetFromSection("section1", "hello"))
	assert.Equal(t, "bar", result.GetFromSection("section1", "foo"))
	assert.Equal(t, "123", result.GetFromSection("section2", "test"))
}

func TestSectionWithSurroundingWhitespace(t *testing.T) {
	result, err := Parse("" +
		"   [section]  \t " + "\n" +
		"key = value" +
		"")
	require.Nil(t, err)
	assert.Equal(t, "value", result.GetFromSection("section", "key"))
}

func TestIgnoresLeadingWhitespace(t *testing.T) {
	result, err := Parse(`
          key = value
      #FooBAR
   [section1]
         hello = world
 foo = bar

  # Test 123
    [section2]
        test = 123
`)
	require.Nil(t, err)
	assert.Equal(t, result.CountAllEntries(), 4)
	assert.Equal(t, "value", result.GetFromSection("", "key"))
	assert.Equal(t, "world", result.GetFromSection("section1", "hello"))
	assert.Equal(t, "bar", result.GetFromSection("section1", "foo"))
	assert.Equal(t, "123", result.GetFromSection("section2", "test"))
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

func TestParseEmptyFile(t *testing.T) {
	for _, text := range []string{
		"",
		"  ",
		"\n\n",
		"\r\n\n\n\r\n",
	} {
		result, err := Parse(text)
		require.Nil(t, err)
		assert.Equal(t, 0, result.CountAllEntries())
	}
}

func TestErrorCases(t *testing.T) {
	for _, text := range []string{
		"foo",                  // Key without value
		"[section",             // Invalid section
		"section]",             // Invalid section
		"[section] # Comment?", // Illegal trailing comment
		"[[section]]",          // Section name itself contains square bracket
		"[se]ct[ion]",          // Section name itself contains square bracket
		"[key] = 123",          // Key cannot look like section
		"key= 123",             // Missing space delimiter after key
		"key =123",             // Missing space delimiter before value
		"k e y = 123",          // Key with whitespace
		"k\te\ty = 123",        // Key with whitespace
	} {
		_, err := Parse(text)
		require.Error(t, err)
	}
}
