package main

import (
	"AngkutKita/RegisterLogin/handlers"
	"AngkutKita/RegisterLogin/models"
	"encoding/json"
	"io/ioutil"
	// "errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error
func initialMigration() {
	db.AutoMigrate(&models.User{})
}
type Trayek struct {
	Jalur             string `json:"Jalur"`
	Trayek            string `json:"Trayek"`
	DariTerminalArjosari string `json:"Dari Terminal Arjosari"`
	DariTerminalLandungsari string `json:"Dari Terminal Landungsari"`
}
func getTrayeks(c *gin.Context) {
	// Mengambil data dari URL yang diberikan
	url := "https://script.googleusercontent.com/macros/echo?user_content_key=EoLYkhs3mY7o8hi7zz2briQo2hLUVpFnNSg88mZeJR8gc0ChCHiOYU14BhAvKrLNpH-zSEB9gyWP5MoVu_JctkQcmHO4zT1Mm5_BxDlH2jW0nuo2oDemN9CCS2h10ox_1xSncGQajx_ryfhECjZEnJILmmzQhVB_6nmRSh4U-9w2BAm8YzC8yeeMjWbCX3mriMyxg6Uva8MjBJ_UNSjdurF3DIeZ__rV8b2IhxNXwhi7qVj0lxuan9z9Jw9Md8uu&lib=MMCczYl0a-TDyef2SR8nJub3GWsEnYdVl"
	response, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from the URL"})
		return
	}
	defer response.Body.Close()

	// Membaca data dari response
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response body"})
		return
	}

	// Parsing JSON ke dalam struktur data
	var trayeks []Trayek
	err = json.Unmarshal(data, &trayeks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing JSON"})
		return
	}

	// Mengirimkan data JSON sebagai response
	c.JSON(http.StatusOK, trayeks)
}

const apiKey = "AIzaSyDX_f2ODZCxR2vKIn7kzo5UjBuig2EEYLA"
func main() {
	dsn := "root:@tcp(localhost:3306)/angkutkita?parseTime=true"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	initialMigration()

	router := gin.Default()
	// Register handlers with the DB instance
	router.POST("/register", handlers.RegisterHandler(db))
	router.GET("/login/:username/:password", handlers.LoginHandler(db))
	router.GET("/get-user-location", GetUserLocation)
	router.GET("/trayeks", getTrayeks) // Tambahkan endpoint untuk trayeks
	router.GET("/trayeks/:jalur", getTrayeksByJalur)

	fmt.Println("Server is running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}

func getTrayeksByJalur(c *gin.Context) {
	// Mengambil nilai parameter jalur dari URL
	jalur := c.Param("jalur")

	// Mengambil data dari URL yang diberikan
	url := "https://script.googleusercontent.com/macros/echo?user_content_key=EoLYkhs3mY7o8hi7zz2briQo2hLUVpFnNSg88mZeJR8gc0ChCHiOYU14BhAvKrLNpH-zSEB9gyWP5MoVu_JctkQcmHO4zT1Mm5_BxDlH2jW0nuo2oDemN9CCS2h10ox_1xSncGQajx_ryfhECjZEnJILmmzQhVB_6nmRSh4U-9w2BAm8YzC8yeeMjWbCX3mriMyxg6Uva8MjBJ_UNSjdurF3DIeZ__rV8b2IhxNXwhi7qVj0lxuan9z9Jw9Md8uu&lib=MMCczYl0a-TDyef2SR8nJub3GWsEnYdVl"
	response, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from the URL"})
		return
	}
	defer response.Body.Close()

	// Membaca data dari response
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response body"})
		return
	}

	// Parsing JSON ke dalam struktur data
	var trayeks []Trayek
	err = json.Unmarshal(data, &trayeks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing JSON"})
		return
	}

	// Menyaring data trayek berdasarkan jalur
	var trayeksByJalur []Trayek
	for _, trayek := range trayeks {
		if trayek.Jalur == jalur {
			trayeksByJalur = append(trayeksByJalur, trayek)
		}
	}

	// Mengirimkan data JSON sebagai response
	c.JSON(http.StatusOK, trayeksByJalur)
}
type GeocodingResponse struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
	} `json:"results"`
	Status string `json:"status"`
}

func Geocode(lat, lng float64) (string, error) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&key=%s", lat, lng, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var geocodingResponse GeocodingResponse
	err = json.Unmarshal(body, &geocodingResponse)
	if err != nil {
		return "", err
	}
	// if geocodingResponse.Status != "OK" || len(geocodingResponse.Results) == 0 {
	// 	return "", errors.New("Geocoding failed");
	// }
	return geocodingResponse.Results[0].FormattedAddress, nil
}

func GetUserLocation(c *gin.Context) {
	requestData := map[string]interface{}{
		"considerIp": true,
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestData).
		Post("https://www.googleapis.com/geolocation/v1/geolocate?key=" + apiKey)

	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	var locationData struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}
	err = json.Unmarshal(resp.Body(), &locationData)
	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	address, err := Geocode(locationData.Lat, locationData.Lng)
	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, gin.H{"location": address})
}
