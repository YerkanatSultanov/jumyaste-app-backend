package helper

import (
	"encoding/json"
	"regexp"
	"strings"
)

func StructToMap(data interface{}) map[string]interface{} {
	var result map[string]interface{}
	tmp, _ := json.Marshal(data)
	json.Unmarshal(tmp, &result)
	return result
}

func ExtractJSON(text string) string {
	re := regexp.MustCompile(`(?s)\{.*?\}`) // (?s) — чтобы . матчило переносы строк, .*? — ленивое

	jsonText := re.FindString(text)

	if jsonText == "" || !strings.HasPrefix(jsonText, "{") || !strings.HasSuffix(jsonText, "}") {
		return ""
	}

	return jsonText
}
