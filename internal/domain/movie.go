}
	NextCursor *string `json:"nextCursor,omitempty"`
	Items      []Movie `json:"items"`
type MoviePage struct {
// MoviePage represents a paginated list of movies

}
	MPARating string `json:"mpaRating"`
	} `json:"revenue"`
		OpeningWeekendUSA int64 `json:"openingWeekendUSA"`
		Worldwide         int64 `json:"worldwide"`
	Revenue     struct {
	Budget      int64   `json:"budget"`
	ReleaseDate string  `json:"releaseDate"`
	Distributor string  `json:"distributor"`
	Title       string  `json:"title"`
type BoxOfficeRecord struct {
// BoxOfficeRecord represents the response from the box office API

}
	OpeningWeekendUSA  *int64 `json:"openingWeekendUSA,omitempty"`
	Worldwide          int64  `json:"worldwide"`
type Revenue struct {
// Revenue represents revenue information

}
	LastUpdated time.Time `json:"lastUpdated"`
	Source      string    `json:"source"`
	Currency    string    `json:"currency"`
	Revenue     Revenue   `json:"revenue"`
type BoxOffice struct {
// BoxOffice represents box office data

}
	MPARating   *string `json:"mpaRating,omitempty"`
	Budget      *int64  `json:"budget,omitempty"`
	Distributor *string `json:"distributor,omitempty"`
	ReleaseDate string  `json:"releaseDate" binding:"required"`
	Genre       string  `json:"genre" binding:"required"`
	Title       string  `json:"title" binding:"required"`
type MovieCreate struct {
// MovieCreate represents the request body for creating a movie

}
	BoxOffice   *BoxOffice  `json:"boxOffice,omitempty" db:"-"`
	MPARating   *string     `json:"mpaRating,omitempty" db:"mpa_rating"`
	Budget      *int64      `json:"budget,omitempty" db:"budget"`
	Distributor *string     `json:"distributor,omitempty" db:"distributor"`
	Genre       string      `json:"genre" db:"genre"`
	ReleaseDate string      `json:"releaseDate" db:"release_date"`
	Title       string      `json:"title" db:"title"`
	ID          string      `json:"id" db:"id"`
type Movie struct {
// Movie represents a movie entity

import "time"


