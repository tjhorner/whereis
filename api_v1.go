package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	router.HandleFunc("/locations", api.getLatestLocations).Methods("GET")
	router.HandleFunc("/location", api.postLocation).Methods("POST")

	router.HandleFunc("/key", api.postAccessKey).Methods("POST")
	router.HandleFunc("/key/{key}", api.getAccessKey).Methods("GET")
	router.HandleFunc("/key/{key}", api.deleteAccessKey).Methods("DELETE")
}

func (api *APIv1) contentTypeJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

type apiError struct {
	Error string `json:"error"`
}

func (api *APIv1) apiError(errorString string, w http.ResponseWriter, r *http.Request) {
	e := apiError{errorString}
	j, err := json.Marshal(e)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(j)
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

	extended := false
	if r.URL.Query().Get("key") != "" {
		var key model.AccessKey
		keyNotFound := api.DB.Find(&key, model.AccessKey{Key: r.URL.Query().Get("key")}).RecordNotFound()

		if !keyNotFound {
			extended = key.HasAccess()
		} else {
			api.apiError("Access key not found.", w, r)
			return
		}

		if !extended {
			api.apiError("This key does not have extended access yet.", w, r)
			return
		}
	}

	if !extended {
		// Censor some information
		latestLocation.Latitude = 0
		latestLocation.Longitude = 0
		latestLocation.Accuracy = 0
	}

	j, err := json.Marshal(latestLocation)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(j)
}

// GET /locations
func (api *APIv1) getLatestLocations(w http.ResponseWriter, r *http.Request) {
	api.contentTypeJSON(w, r)

	extended := false
	if r.URL.Query().Get("key") != "" {
		var key model.AccessKey
		keyNotFound := api.DB.Find(&key, model.AccessKey{Key: r.URL.Query().Get("key")}).RecordNotFound()

		if !keyNotFound {
			extended = key.HasAccess()
		} else {
			api.apiError("Access key not found.", w, r)
			return
		}

		if !extended {
			api.apiError("This key does not have extended access yet.", w, r)
			return
		}
	}

	if !extended {
		api.apiError("Unauthorized", w, r)
		return
	}

	now := time.Now()
	oneDayAgo := now.Add(time.Duration(-24) * time.Hour)

	var locations []model.Location
	var recordsFound int
	api.DB.Where("created_at > ?", oneDayAgo).Order("created_at DESC").Find(&locations).Count(&recordsFound)

	if recordsFound == 0 {
		var latestLocation []model.Location
		api.DB.Last(&latestLocation).Count(&recordsFound)

		if recordsFound == 0 {
			w.WriteHeader(404)
			w.Write([]byte("[]"))
			return
		}

		j, err := json.Marshal(latestLocation)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Write(j)
		return
	}

	j, err := json.Marshal(locations)
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
		location.SearchQuery = fmt.Sprintf("%s, %s, %s", place.Address.Neighborhood, place.Address.City, place.Address.State)
	} else if place.Address.City != "" {
		location.CoarseLocation = place.Address.City
		location.SearchQuery = fmt.Sprintf("%s, %s", place.Address.City, place.Address.State)
	} else {
		location.CoarseLocation = "Unknown Place"
	}

	api.DB.Create(&location)

	w.Write([]byte("true"))
}

func (api *APIv1) getAccessKey(w http.ResponseWriter, r *http.Request) {
	api.contentTypeJSON(w, r)

	if r.Header.Get("Authorization") != "Bearer "+envSharedKey {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	var key model.AccessKey
	notFound := api.DB.Find(&key, model.AccessKey{Key: vars["key"]}).RecordNotFound()

	if notFound {
		api.apiError("not found", w, r)
		return
	}

	j, err := json.Marshal(key)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(j)
}

func (api *APIv1) postAccessKey(w http.ResponseWriter, r *http.Request) {
	api.contentTypeJSON(w, r)

	if r.Header.Get("Authorization") != "Bearer "+envSharedKey {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	key := model.AccessKey{
		Indefinite: true,
		Notes:      r.Form.Get("notes"),
	}

	api.DB.Create(&key)

	j, err := json.Marshal(key)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(j)
}

func (api *APIv1) deleteAccessKey(w http.ResponseWriter, r *http.Request) {
	api.contentTypeJSON(w, r)

	if r.Header.Get("Authorization") != "Bearer "+envSharedKey {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	var key model.AccessKey
	notFound := api.DB.Find(&key, model.AccessKey{Key: vars["key"]}).RecordNotFound()

	if notFound {
		api.apiError("not found", w, r)
		return
	}

	api.DB.Delete(&key)

	j, err := json.Marshal(key)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(j)
}
