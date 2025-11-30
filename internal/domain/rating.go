package domain

// RatingSubmit 评分提交请求
type RatingSubmit struct {
	Rating float64 `json:"rating" binding:"required"`
}

// RatingResult 评分结果
type RatingResult struct {
	MovieTitle string  `json:"movieTitle"`
	RaterID    string  `json:"raterId"`
	Rating     float64 `json:"rating"`
}

// RatingAggregate 评分聚合结果
type RatingAggregate struct {
	Average float64 `json:"average"`
	Count   int     `json:"count"`
}
