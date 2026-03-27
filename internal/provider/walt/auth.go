package walt

import "net/http"

func applyBearerToken(req *http.Request, bearerToken string) {
	if bearerToken == "" {
		return
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)
}
