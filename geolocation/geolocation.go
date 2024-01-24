// geolocation/utils.go
package geolocation

import (
	"io"
	"net/http"
)

// GetCurrentLocation fetches the current location using Google Geolocation API
func getCurrentLocation(w http.ResponseWriter, r *http.Request) {
	url := "https://sanarora-get-ip-info.p.rapidapi.com/v3/ip-city/?ip=74.125.45.100&key=93b77a0add7dace51661cf559ef97326f3297ec27d6e5a9b903670e0246b8293&format=JSON"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	req.Header.Add("X-RapidAPI-Key", "6c9e9eb5b6mshda612731f65f1a1p1a470bjsna2a4357d9751")
	req.Header.Add("X-RapidAPI-Host", "sanarora-get-ip-info.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Error making request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}
	// Set the response content type
	w.Header().Set("Content-Type", "application/json")

	// Write the response body to the client
	w.WriteHeader(res.StatusCode)
	w.Write(body)
}