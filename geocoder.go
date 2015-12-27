package geocoder

/* Initialize with a key from https://developers.google.com/maps/documentation/geocoding/get-api-key */
//  geo := Geocoder{key: "......"}

/* This will return a JSON reponse */
//  json := geo.getJSON("Tokyo, Japan")

/* This will return a Response struct as seen below */
//  response := geo.decodeJSON("zipcode M2H 2G6")

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Geocoder struct
type Geocoder struct {
	url string
	key string
}

// Response json
type Response struct {
	Results []Result `json:"results"`
	Status  string   `json:"status"`
}

// Result json
type Result struct {
	AddressComponents []Component `json:"address_components"`
	FormattedAddress  string      `json:"formatted_address"`
	Geometry          Geometry
	PartialMatch      bool     `json:"partial_match"`
	PlaceID           string   `json:"place_id"`
	Types             []string `json:"types"`
}

// Component json
type Component struct {
	LongName  string `json:"long_name"`
	ShortName string `json:"short_name"`
	Types     []string
}

// Geometry json
type Geometry struct {
	Bounds       Square     `json:"bounds"`
	Location     Coordinate `json:"location"`
	LocationType string     `json:"location_type"`
	Viewport     Square     `json:"viewport"`
}

// Square json
type Square struct {
	Northeast Coordinate `json:"northeast"`
	Southwest Coordinate `json:"southwest"`
}

// Coordinate json
type Coordinate struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

func (g *Geocoder) setURL(geo string) {
	geo = strings.Replace(geo, " ", "+", -1)
	g.url = fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", geo, g.key)
}

// DecodeJSON returns the Google result as a struct
func (g *Geocoder) DecodeJSON(geo string) Response {
	g.setURL(geo)
	response := new(Response)
	r, err := http.Get(g.url)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(response)
	if err != nil {
		panic(err)
	}
	return *response
}

// GetJSON returs the Google result as a JSON string
func (g *Geocoder) GetJSON(geo string) string {
	g.setURL(geo)
	fmt.Println(geo)
	fmt.Println(g.url)
	r, err := http.Get(g.url)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	return string(contents)
}

// GetComponent fetches an item by name from the matched places
func GetComponent(components []Component, name string) string {
	for _, c := range components {
		if c.Types[0] == name {
			return c.LongName
		}
	}
	return ""
}
