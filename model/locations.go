package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/clbanning/mxj"
	"github.com/jinzhu/gorm"
)

type Location struct {
	Longitude float64 `json:"lng"`
	Latitude  float64 `json:"lat"`
}

type Marker struct {
	Location Location
	Infos    []string
}

const apiKey = "1zt39dn13glty5q8zjcbcsbs"

const xmlPost = `<?xml version="1.0" encoding="UTF-8"?>
<XLS version="1.2"
  xmlns="http://www.opengis.net/xls"
  xmlns:gml="http://www.opengis.net/gml"
  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xsi:schemaLocation="http://www.opengis.net/xls http://schemas.opengis.net/ols/1.2/olsAll.xsd">
    <RequestHeader/>
    %s
</XLS>`

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

func postGeoportail(xml string) (mxj.Map, error) {
	log.Println(xml)

	req, err := http.NewRequest("POST", "http://wxs.ign.fr/1zt39dn13glty5q8zjcbcsbs/geoportail/ols", strings.NewReader(xml))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Referer", "http://localhost")
	req.Header.Add("Content-Type", "application/xml")
	// log.Println(req)
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	mv, err := mxj.NewMapXmlReader(res.Body)
	if err != nil {
		return nil, err
	}
	data, _ := mv.Xml()
	log.Println(string(data))
	return mv, nil
}

func getAltitude(lat float64, lng float64) (int64, error) {
	url := fmt.Sprintf("http://wxs.ign.fr/1zt39dn13glty5q8zjcbcsbs/alti/rest/elevation.json?lat=%f&lon=%f&zonly=true", lat, lng)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Referer", "http://localhost")
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	type Altitude struct {
		Elevations []float64
	}
	var elevation Altitude
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&elevation)
	if err != nil {
		return 0, err
	}

	log.Println(elevation)

	altitude := int64(0)
	if len(elevation.Elevations) > 0 {
		altitude = int64(elevation.Elevations[0])
	}

	return altitude, nil
}

func FindLocation(lat float64, lng float64) (string, string, int64, error) {
	requestXML := fmt.Sprintf(`
    <Request methodName="ReverseGeocodeRequest" maximumResponses="1" requestID="1" version="1.2">
       <ReverseGeocodeRequest>
          <ReverseGeocodePreference>PositionOfInterest</ReverseGeocodePreference>
          <Position>
             <gml:Point>
                <gml:pos>%f %f</gml:pos>
			 </gml:Point>
			 <gml:CircleByCenterPoint>
                <gml:pos>%f %f</gml:pos>
                <gml:radius>500</gml:radius>
             </gml:CircleByCenterPoint>
          </Position>
       </ReverseGeocodeRequest>
  </Request>`, lat, lng, lat, lng)
	xml := fmt.Sprintf(xmlPost, requestXML)
	res, err := postGeoportail(xml)
	if err != nil {
		return "", "", 0, err
	}

	commune := ""
	code := ""

	communes, _ := res.ValuesForKey("Place", "-type:Commune")
	log.Print(communes)
	if len(communes) > 0 {
		commune = communes[0].(map[string]interface{})["#text"].(string)
		codes, _ := res.ValuesForKey("PostalCode")
		code = codes[0].(string)
	}
	codes, _ := res.ValuesForKey("Place", "-type:Departement")
	log.Print(codes)
	if len(codes) > 0 {
		code = codes[0].(map[string]interface{})["#text"].(string)
	}

	altitude, err := getAltitude(lat, lng)

	return commune, code, altitude, err
}

func FindLatLng(commune string) (float64, float64, error) {
	requestXML := fmt.Sprintf(`
    <Request methodName="LocationUtilityService" requestID="2" version="1.2">
       <GeocodeRequest returnFreeForm="false">
         <Address countryCode="PositionOfInterest">
                <freeFormAddress>%s</freeFormAddress>
         </Address>
       </GeocodeRequest>
	</Request>`, commune)
	xml := fmt.Sprintf(xmlPost, requestXML)
	res, err := postGeoportail(xml)
	if err != nil {
		return 0, 0, err
	}
	log.Println(res)
	positions, _ := res.ValuesForKey("pos")
	if len(positions) == 0 {
		return 0, 0, fmt.Errorf("%s not found", commune)
	}
	pos := strings.Split(positions[0].(string), " ")
	lat, err := strconv.ParseFloat(pos[0], 64)
	lng, err := strconv.ParseFloat(pos[1], 64)
	return lat, lng, err
}
