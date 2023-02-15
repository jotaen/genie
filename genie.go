package genie

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type Data struct {
	entries map[string]map[string]string
	count   int
}

func (d Data) Get(key string) string {
	return d.GetFromSection("", key)
}

func (d Data) GetFromSection(section string, key string) string {
	return d.entries[section][key]
}

func (d Data) CountAllEntries() int {
	return d.count
}

var (
	whitespacePattern = regexp.MustCompile(`^[ \t]*$`)
)

func Parse(iniText string) (Data, error) {
	data := Data{
		entries: make(map[string]map[string]string),
		count:   0,
	}

	// Convert Windows line endings.
	iniText = strings.Replace(iniText, "\r\n", "\n", -1)

	section := ""
	// Parse text line by line.
	for i, l := range strings.Split(iniText, "\n") {

		// Skip empty, blank, or comment lines.
		if l == "" || strings.HasPrefix(l, "#") || whitespacePattern.MatchString(l) {
			continue
		}

		// Parse section header.
		if strings.HasPrefix(l, "[") {
			l = strings.TrimRight(l, " \t")
			if !strings.HasSuffix(l, "]") {
				return Data{}, newError(i, "Invalid section declaration")
			}
			section = l[1 : len(l)-1]
			if section == "" || whitespacePattern.MatchString(section) {
				return Data{}, newError(i, "Invalid section name")
			}
			if strings.Contains(section, "[") || strings.Contains(section, "]") {
				return Data{}, newError(i, "Invalid section name")
			}
			continue
		}

		// Parse key<>value pairs.
		parts := strings.SplitN(l, "=", 2)
		if len(parts) != 2 {
			return Data{}, newError(i, "Invalid key")
		}

		// Process key.
		key := parts[0]
		if !strings.HasSuffix(key, " ") {
			return Data{}, newError(i, "Invalid delimiter sequence")
		}
		key = strings.TrimRight(key, " ")
		if strings.Contains(key, " ") || strings.Contains(key, "\t") {
			return Data{}, newError(i, "Invalid key")
		}

		// Process value.
		value := parts[1]
		if (value != "" && value != " ") && !strings.HasPrefix(value, " ") {
			return Data{}, newError(i, "Invalid delimiter sequence")
		}
		if strings.HasPrefix(value, " ") {
			value = strings.TrimPrefix(value, " ")
		}

		// Save.
		if data.entries[section] == nil {
			data.entries[section] = make(map[string]string)
		}
		data.entries[section][key] = value
		data.count++
	}
	return data, nil
}

func newError(lineIndex int, msg string) error {
	nr := strconv.Itoa(lineIndex + 1)
	return errors.New("Malformed syntax in line " + nr + ": " + msg)
}
