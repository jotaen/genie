# Genie 🧞‍♂️

Your friendly .ini config file parser, written in Go.

## Usage

```go
package main
import ("fmt"; "github.com/jotaen/genie")

const iniText = `
# Comment.
key = value

[section1]
foo = bar
test = 123
`

func main() {
	data, err := genie.Parse(iniText)
	if err != nil {
		panic(err)
	}

	fmt.Println(data.Get("key")) // Prints `value`.
	fmt.Println(data.GetFromSection("section1", "foo")) // Prints `bar`.
	fmt.Println(data.GetFromSection("section1", "test")) // Prints `123`.

	fmt.Println(data.Get("asdfasdf")) // Prints empty string.
	fmt.Println(data.GetFromSection("section26487", "bla")) // Prints empty string.
}
```

## Rules

### Whitespace
Whitespace is a space ` ` or a tab `\t`.
Leading whitespace is insignificant and always skipped.

### Entries
An entry is a key/value pair.
The delimiter between key and value is the sequence of a `=` character surrounded by a space ` ` character.
The key itself cannot contain whitespace.
All what follows behind the delimiter is value-land, including any whitespace or `#`’s.
The value might be absent, in which case the space behind the `=` may be absent too.
There is no distinction between “absent” value or “empty” value.
Both the key and the value are case-sensitive.

### Comments
If a line starts with `#`, it’s treated as comment and is ignored. There can’t be trailing comments.

### Sections
A label wrapped in square brackets denotes the section for all following entries.
All entries until the first explicit section belong to the top-level section.
Section names cannot contain whitespace; they cannot be empty; they are case-sensitive.

## License

[MIT](License.txt)
