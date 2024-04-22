package model

import (
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type Location struct {
	Longitude float64 `json:"lng"`
	Latitude  float64 `json:"lat"`
}

type Marker struct {
	Location Location
	Infos    []string
}

func Markers(db *gorm.DB, search string, espece_id uint) ([]Marker, error) {
	locations := make(map[Location][]string)

	db = db.Table("individus").Joins("join especes on espece_id=especes.id")

	if espece_id != 0 {
		db = db.
			Where("especes.id = ?", espece_id)
	} else {
		if search != "" {
			search = "%" + search + "%"
			db = db.
				Where("especes.genre like ? or especes.sous_genre like ? or especes.espece like ? or especes.sous_espece like ? or site like ? or commune like ? or code like ?",
					search, search, search, search, search, search, search)
		}
	}

	rows, err := db.Select("sexe,genre,sous_genre,espece,sous_espece,date,latitude,longitude").Rows()
	if err != nil {
		return nil, nil
	}
	defer rows.Close()
	for rows.Next() {
		var sexe1, genre, sous_genre, espece, sous_espece string
		var date time.Time
		var latitude, longitude float64
		if err := rows.Scan(&sexe1, &genre, &sous_genre, &espece, &sous_espece, &date, &latitude, &longitude); err != nil {
			log.Println(err)
			return nil, err
		}
		loc := Location{longitude, latitude}
		sexe := "♀"
		if sexe1 == "M" {
			sexe = "♂"
		}
		infos := append(locations[loc], fmt.Sprintf("%s %s (%s) %s %s - (%s)", sexe, genre, sous_genre, espece, sous_espece, date.Format("02/01/2006")))
		locations[loc] = infos
	}

	var result []Marker
	for k, v := range locations {
		result = append(result, Marker{Location: k, Infos: v})
	}
	return result, nil
}

func getAltitude(lng, lat float64) (int64, error) {
	url := "https://data.geopf.fr/altimetrie/1.0/calcul/alti/rest/elevation.json"
	resMap := make(map[string]interface{})
	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetQueryParams(map[string]string{
			"lat":      fmt.Sprintf("%f", lat),
			"lon":      fmt.Sprintf("%f", lng),
			"resource": "ign_rge_alti_wld",
			"zonly":    "true",
		}).
		SetResult(&resMap).
		Get(url)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode() != 200 {
		return 0, fmt.Errorf("error %d", resp.StatusCode())
	}
	elevations := resMap["elevations"].([]interface{})
	if len(elevations) == 0 {
		return 0, fmt.Errorf("no elevation found")
	}
	elevation := elevations[0].(float64)
	return int64(elevation), nil
}

func FindLocation(lng, lat float64) (string, string, int64, error) {
	url := "https://data.geopf.fr/geocodage/reverse"
	resMap := make(map[string]interface{})
	client := resty.New()
	// client.AllowGetMethodPayload = true
	// client.SetDebug(true)
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetQueryParam("searchgeom", fmt.Sprintf(`{"type": "Circle", "coordinates": [%f,%f], "radius": 500}`, lng, lat)).
		SetResult(&resMap).
		Get(url)
	if err != nil {
		return "", "", 0, err
	}
	if resp.StatusCode() != 200 {
		return "", "", 0, fmt.Errorf("error %d", resp.StatusCode())
	}
	log.Println("findLocation", resMap)
	features := resMap["features"].([]interface{})
	if len(features) == 0 {
		return "", "", 0, fmt.Errorf("no location found")
	}
	properties := features[0].(map[string]interface{})["properties"].(map[string]interface{})

	commune := ""
	code := ""
	if properties["city"] != nil {
		commune = properties["city"].(string)
	}
	if properties["citycode"] != nil {
		code = properties["citycode"].(string)[0.:2]
	}
	altitude, err := getAltitude(lng, lat)
	return commune, code, altitude, err
}
