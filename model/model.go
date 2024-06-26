package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"time"
)

type Departement struct {
	Code string `gorm:"primaryKey"`
	Nom  string
}

type Classification struct {
	ID      uint
	Nom     string
	Ordre   sql.NullInt64
	Especes []Espece
}

type Espece struct {
	ID               uint
	ClassificationID uint
	Classification   Classification
	Genre            string
	SousGenre        string
	Espece           string
	SousEspece       string
	Descripteur      string
	Individus        []Individu
}

func (espece *Espece) NomEspece() string {
	buf := bytes.NewBufferString(espece.Genre)
	if espece.SousGenre != "" {
		buf.WriteString(" (")
		buf.WriteString(espece.SousGenre)
		buf.WriteString(")")
	}
	buf.WriteString(" ")
	buf.WriteString(espece.Espece)
	if espece.SousEspece != "" {
		buf.WriteString(" ")
		buf.WriteString(espece.SousEspece)
	}
	return buf.String()
}

type Individu struct {
	ID          uint
	Date        time.Time
	Sexe        string
	Altitude    sql.NullInt64
	Commentaire sql.NullString
	EspeceID    uint
	Espece      Espece
	Site        string
	Commune     string
	Code        string
	Departement Departement `gorm:"foreignKey:Code;references:Code"`
	Longitude   float64
	Latitude    float64
	Recolteur   string
	Selected    bool `gorm:"-"`
}

func (individu *Individu) Localization() string {
	return fmt.Sprintf("%f,%f", individu.Latitude, individu.Longitude)
}
