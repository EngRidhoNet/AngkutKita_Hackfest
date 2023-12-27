// geolocation/utils.go
package geolocation

import (
   "encoding/json"
   "fmt"
   "net/http"
)

// GetCurrentLocation fetches the current location using Google Geolocation API
func GetCurrentLocation(w http.ResponseWriter, r *http.Request, apiKey string) {
   // Google Maps Geolocation API endpoint
   url := fmt.Sprintf("https://www.googleapis.com/geolocation/v1/geolocate?key=%s", apiKey)

   // HTTP POST request to Geolocation API
   resp, err := http.Post(url, "application/json", nil)
   if err != nil {
      http.Error(w, "Error sending request to Geolocation API", http.StatusInternalServerError)
      return
   }
   defer resp.Body.Close()

   // Decode the response JSON
   var result map[string]interface{}
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      http.Error(w, "Error decoding response from Geolocation API", http.StatusInternalServerError)
      return
   }

   // Extract and send the location information in the response
   if location, ok := result["location"].(map[string]interface{}); ok {
      latitude := location["lat"]
      longitude := location["lng"]
      response := fmt.Sprintf("Current Location: Lat %v, Lng %v", latitude, longitude)
      w.Write([]byte(response))
   } else {
      http.Error(w, "Unable to retrieve location information.", http.StatusInternalServerError)
   }
}
