package properties

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type PropertiesMap map[string]string

func LoadProperties(path string) (PropertiesMap, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	props := make(PropertiesMap)
	scanner := bufio.NewScanner(file)

	var prevLine string
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}

		if strings.HasSuffix(line, "\\") {
			prevLine += strings.TrimSuffix(line, "\\")
			continue
		}
		line = prevLine + line
		prevLine = ""

		var key, value string
		if idx := strings.IndexAny(line, "=:"); idx != -1 {
			key = strings.TrimSpace(line[:idx])
			value = strings.TrimSpace(line[idx+1:])
		} else {
			key = strings.TrimSpace(line)
			value = ""
		}

		value = unescapeValue(value)

		props[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return props, nil
}

func unescapeValue(s string) string {
	replacer := strings.NewReplacer(
		"\\t", "\t",
		"\\n", "\n",
		"\\r", "\r",
		"\\\\", "\\",
	)
	s = replacer.Replace(s)

	for {
		idx := strings.Index(s, `\u`)
		if idx == -1 || idx+6 > len(s) {
			break
		}
		code, err := strconv.ParseInt(s[idx+2:idx+6], 16, 32)
		if err != nil {
			break
		}
		r := rune(code)
		s = s[:idx] + string(r) + s[idx+6:]
	}

	return s
}
