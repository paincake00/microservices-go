package health

import "net/http"

func CheckHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	//nolint:gosec
	_, _ = w.Write([]byte("."))
}
