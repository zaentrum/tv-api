package http

// sampleItems returns hardcoded items shaped for the tv product. Replaced once
// the shared catalog backend is online.
func sampleItems() []map[string]any {
	return []map[string]any{
		{"id": "ch-1", "kind": "channel", "name": "Channel One", "live": true},
		{"id": "ch-2", "kind": "channel", "name": "Channel Two", "live": true},
	}
}
