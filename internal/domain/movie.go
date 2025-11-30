package domain

import "time"

// BoxOffice Revenue
type BoxOfficeRevenue struct {
	Worldwide         int64 `json:"worldwide"`
	OpeningWeekendUSA int64 `json:"openingWeekendUSA"`
}

// BoxOffice information
type BoxOffice struct {
	Revenue     BoxOfficeRevenue `json:"revenue"`
	Currency    string           `json:"currency"`
	Source      string           `json:"source"`
	LastUpdated time.Time        `json:"lastUpdated"`
}

// Movie entity 返回到客户端
type Movie struct {
	ID          string     `json:"id"` // 注意 openapi 要求 string ID
	Title       string     `json:"title"`
	ReleaseDate string     `json:"releaseDate"`
	Genre       string     `json:"genre"`
	Distributor *string    `json:"distributor"` // 可为空，但始终返回（即使为null）
	Budget      *int64     `json:"budget"`      // 可为空，但始终返回（即使为null）
	MpaRating   *string    `json:"mpaRating"`   // 可为空，但始终返回（即使为null）
	BoxOffice   *BoxOffice `json:"boxOffice"`   // 可为空，但始终返回（即使为null）
}
