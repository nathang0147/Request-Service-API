package walt

import "net/http"

func applyAPIKey(req *http.Request, apiKey string) {
	if apiKey == "" {
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
}
