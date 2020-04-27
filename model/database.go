package model

import (
	"log"

	"github.com/jinzhu/gorm"
)

func LoadIndividus(db *gorm.DB, individus *[]Individu, index int, size int, search string, espece_id uint, order string) int {
	db = db.Joins("join especes on espece_id=especes.id")
	db = db.Joins("join classifications on especes.classification_id=classifications.id")

	if espece_id != 0 {
		db = db.Where("especes.id = ?", espece_id)
	} else {
		if search != "" {
			search = "%" + search + "%"
			whereClause := `especes.genre like ? or especes.sous_genre like ? or especes.espece like ? or especes.sous_espece like ? or site like ? or commune like ? or code like ? or classifications.nom like ?`
			db = db.Where(whereClause, search, search, search, search, search, search, search, search)
		}
	}

	var count int
	if err := db.Model(&Individu{}).Count(&count).
		Offset(index).Limit(size).
		Order(order).
		Preload("Espece.Classification").
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

func CountAllEspeces(db *gorm.DB) int {
	nb := 0
	db.Model(&Espece{}).Count(&nb)
	return nb
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
