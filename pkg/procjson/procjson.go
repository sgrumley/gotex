package procjson

import (
	"encoding/json"
	"html"
	"strings"

	"github.com/hokaccha/go-prettyjson"
)

type TestEvent struct {
	Time    string `json:"Time"`
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test,omitempty"`
	Output  string `json:"Output,omitempty"`
}

func Format(jason string) string {
	jason = strings.ReplaceAll(jason, `\n`, "")
	lines := strings.Split(jason, "\n")

	var out string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var event TestEvent
		err := json.Unmarshal([]byte(line), &event)
		if err != nil {
			return "JSON unmarshal error: " + err.Error() + "\n" + line
		}

		colored, err := prettyjson.Marshal(event)
		if err != nil {
			return "prettyjson marshal error: " + err.Error()
		}

		out += "\n" + html.UnescapeString(string(colored))
	}

	return out
}
