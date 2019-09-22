package locationiq

import (
	"encoding/json"
	"net/http"
)

// Client is a LocationIQ client
type Client struct {
	APIKey string
}

func (c *Client) httpGet(endpoint string, qs map[string]string, resp interface{}) error {
	req, _ := http.NewRequest("GET", "https://us1.locationiq.com/v1/"+endpoint, nil)

	q := req.URL.Query()
	for k, v := range qs {
		q.Set(k, v)
	}

	q.Set("key", c.APIKey)
	q.Set("format", "json")

	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return err
	}

	return nil
}

// ReverseGeocode reverse geocodes
func (c *Client) ReverseGeocode(lat, lon string) (*Place, error) {
	var resp Place
	err := c.httpGet("reverse.php", map[string]string{
		"lat": lat,
		"lon": lon,
	}, &resp)

	return &resp, err
}
