package model

import (
	"log"

	"github.com/jinzhu/gorm"
)

func LoadIndividus(db *gorm.DB, individus *[]Individu, index int, size int, search string, espece_id uint) int {
	if espece_id != 0 {
		db = db.
			Joins("join especes on espece_id=especes.id").
			Where("especes.id = ?", espece_id)
	} else {
		if search != "" {
			db = db.
				Joins("join especes on espece_id=especes.id").
				Where("especes.genre like ? or especes.sous_genre like ? or especes.espece like ? or especes.sous_espece like ? or site like ? or commune like ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
		}
	}

	var count int
	if err := db.Model(&Individu{}).Count(&count).
		Offset(index).Limit(size).
		Order("date desc").
		Preload("Espece").
		Find(individus).Error; err != nil {
		log.Println(err)
	}
	return count
}

func AllClassifications(db *gorm.DB) []Classification {
	var classifications []Classification
	if err := db.Order("nom").
		Preload("Especes", func(db *gorm.DB) *gorm.DB {
			return db.Order("genre,sous_genre,espece,sous_espece")
		}).
		Find(&classifications).Error; err != nil {
		log.Println(err)
	}
	return classifications
}

func AllEspeces(db *gorm.DB) []Espece {
	var especes []Espece
	if err := db.Order("genre,sous_genre,espece,sous_espece").Find(&especes).Error; err != nil {
		log.Println(err)
	}
	return especes
}

func AllDepartements(db *gorm.DB) []Departement {
	var departements []Departement
	if err := db.Order("code").Find(&departements).Error; err != nil {
		log.Println(err)
	}
	return departements
}

func AllGenres(db *gorm.DB) []string {
	var result []string
	rows, err := db.Raw("select distinct(genre) from especes order by 1").Rows()
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		rows.Scan(&name)
		result = append(result, name)
	}
	return result
}

func AllSousGenres(db *gorm.DB) []string {
	var result []string
	rows, err := db.Raw("select distinct(sous_genre) from especes order by 1").Rows()
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		rows.Scan(&name)
		result = append(result, name)
	}
	return result
}

func AllNomEspeces(db *gorm.DB) []string {
	var result []string
	rows, err := db.Raw("select distinct(espece) from especes order by 1").Rows()
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		rows.Scan(&name)
		result = append(result, name)
	}
	return result
}

func AllSousEspeces(db *gorm.DB) []string {
	var result []string
	rows, err := db.Raw("select distinct(sous_espece) from especes order by 1").Rows()
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		rows.Scan(&name)
		result = append(result, name)
	}
	return result
}

func AllRecolteurs(db *gorm.DB) []string {
	var recolteurs []string
	rows, err := db.Raw("select distinct(recolteur) from individus order by recolteur").Rows()
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		rows.Scan(&name)
		recolteurs = append(recolteurs, name)
	}
	return recolteurs
}

func AllSites(db *gorm.DB) []string {
	var sites []string
	rows, err := db.Raw("select distinct(site) from individus order by site").Rows()
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		rows.Scan(&name)
		sites = append(sites, name)
	}
	return sites
}

func AllCommunes(db *gorm.DB) []string {
	var communes []string
	rows, err := db.Raw("select distinct(commune) from individus order by commune").Rows()
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		rows.Scan(&name)
		communes = append(communes, name)
	}
	return communes
}

/*
func UpdateLocations(db *gorm.DB) {
	var individus []Individu
	if err := db.Where("longitude is null or latitude is null").Find(&individus).Error; err != nil {
		log.Println(err)
		return
	}
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyCP7koK86PlGmkHnO6EciwaUQI43laSmZo"))
	if err != nil {
		log.Println(err)
		return
	}

	for _, i := range individus {
		t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
		//		address, _, _ := transform.String(t, fmt.Sprintf("%s. %s. %s", i.Site, i.Code, i.Commune))
		address, _, _ := transform.String(t, fmt.Sprintf("%s. %s", i.Code, i.Commune))
		log.Println(address)
		rs, err := c.Geocode(context.Background(), &maps.GeocodingRequest{Region: "fr", Address: address})
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(rs)
		for _, r := range rs {
			i.Longitude = r.Geometry.Location.Lng
			i.Latitude = r.Geometry.Location.Lat
			db.Set("gorm:save_associations", false).Save(i)
			break
		}
	}
}
*/
