package locationiq

// Place represents a LocationIQ place
type Place struct {
	APIError
	PlaceID     string   `json:"place_id"`
	OsmType     string   `json:"osm_type"`
	OsmID       string   `json:"osm_id"`
	Licence     string   `json:"licence"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	DisplayName string   `json:"display_name"`
	BoundingBox []string `json:"boundingbox"`
	Importance  float64  `json:"importance"`
	Address     struct {
		HouseNumber  string `json:"house_number"`
		Road         string `json:"road"`
		Neighborhood string `json:"neighbourhood"`
		City         string `json:"city"`
		County       string `json:"county"`
		State        string `json:"state"`
		Country      string `json:"country"`
		CountryCode  string `json:"country_code"`
		Postcode     string `json:"postcode"`
	} `json:"address"`
}
