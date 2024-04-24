module github.com/rverpillot/coleoptera

go 1.22

require (
	github.com/go-resty/resty/v2 v2.12.0
	github.com/jung-kurt/gofpdf v1.16.2
	github.com/rverpillot/ihui v1.3.0
	gorm.io/driver/sqlite v1.5.5
	gorm.io/gorm v1.25.9
)

require (
	github.com/PuerkitoBio/goquery v1.9.1 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	golang.org/x/net v0.24.0 // indirect
)

replace github.com/rverpillot/ihui v1.3.0 => ../ihui
