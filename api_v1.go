package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/tjhorner/whereis/locationiq"
	"github.com/tjhorner/whereis/model"
)

// APIv1 implements v1 of the API
type APIv1 struct {
	DB *gorm.DB
}

// Prefix implements API.Prefix
func (api APIv1) Prefix() string {
	return "/api/v1"
}

// Route implements API.Route
func (api APIv1) Route(router *mux.Router) {
	router.HandleFunc("/location", api.getLatestLocation).Methods("GET")
	router.HandleFunc("/location", api.postLocation).Methods("POST")
	router.HandleFunc("/location/{key}", api.getLatestLocationWithKey).Methods("GET")
}

func (api *APIv1) contentTypeJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func (api *APIv1) getLatestLocation(w http.ResponseWriter, r *http.Request) {
	api.contentTypeJSON(w, r)

	var latestLocation model.Location
	var recordsFound int
	api.DB.Last(&latestLocation).Count(&recordsFound)

	if recordsFound == 0 {
		w.WriteHeader(404)
		w.Write([]byte("null"))
		return
	}

	// Censor some information
	latestLocation.Latitude = 0
	latestLocation.Longitude = 0
	latestLocation.Accuracy = 0

	j, err := json.Marshal(latestLocation)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(j)
}

func (api *APIv1) postLocation(w http.ResponseWriter, r *http.Request) {
	api.contentTypeJSON(w, r)

	if r.Header.Get("Authorization") != "Bearer "+envSharedKey {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	lat, _ := strconv.ParseFloat(r.Form.Get("lat"), 64)
	lon, _ := strconv.ParseFloat(r.Form.Get("lon"), 64)
	acc, _ := strconv.ParseFloat(r.Form.Get("acc"), 32)
	batt, _ := strconv.ParseFloat(r.Form.Get("batt"), 32)

	location := model.Location{
		Latitude:  lat,
		Longitude: lon,
		Accuracy:  float32(acc),
		Battery:   float32(batt),
	}

	// Reverse geocoding
	liq := locationiq.NewClient(envLocationIqKey)

	place, err := liq.ReverseGeocode(r.Form.Get("lat"), r.Form.Get("lon"))
	if err != nil || place.Error != "" {
		location.CoarseLocation = "Unknown Place"
	}

	if place.Address.Neighborhood != "" {
		location.CoarseLocation = place.Address.Neighborhood
	} else if place.Address.City != "" {
		location.CoarseLocation = place.Address.City
	} else {
		location.CoarseLocation = "Unknown Place"
	}

	api.DB.Create(&location)

	w.Write([]byte("true"))
}

func (api *APIv1) getLatestLocationWithKey(w http.ResponseWriter, r *http.Request) {
	api.contentTypeJSON(w, r)
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}
