package credential

import (
	"encoding/json"
	"fmt"
	"strings"
)

const redactedValue = "[REDACTED]"

// redactFields walks a JSON string and replaces the values at the given
// dot-separated field paths with "[REDACTED]". If fieldPaths is empty the
// input is returned unchanged.
//
// Field paths use the proto/JSON field names emitted by protojson with
// UseProtoNames: true (snake_case). Nested paths are dot-separated, e.g.
// "spec.secret_access_key" or "spec.gcp_gke.service_account_key_base64".
//
// Only non-empty string values are redacted. Missing fields, empty strings,
// and non-string values are left untouched.
func redactFields(jsonStr string, fieldPaths []string) (string, error) {
	if len(fieldPaths) == 0 {
		return jsonStr, nil
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		return "", fmt.Errorf("redaction unmarshal: %w", err)
	}

	for _, path := range fieldPaths {
		redactPath(obj, strings.Split(path, "."))
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("redaction marshal: %w", err)
	}
	return string(out), nil
}

// redactPath traverses nested maps along parts and replaces the leaf value
// with "[REDACTED]" if it is a non-empty string.
func redactPath(obj map[string]any, parts []string) {
	if len(parts) == 0 || obj == nil {
		return
	}

	key := parts[0]

	if len(parts) == 1 {
		v, exists := obj[key]
		if !exists {
			return
		}
		if s, ok := v.(string); ok && s != "" {
			obj[key] = redactedValue
		}
		return
	}

	next, ok := obj[key].(map[string]any)
	if !ok {
		return
	}
	redactPath(next, parts[1:])
}
