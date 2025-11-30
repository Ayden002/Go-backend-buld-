package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"interview/internal/domain"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// MovieCreateRequest 创建电影的请求结构
type MovieCreateRequest struct {
	Title       string  `json:"title" binding:"required"`
	Genre       string  `json:"genre" binding:"required"`
	ReleaseDate string  `json:"releaseDate" binding:"required"`
	Distributor *string `json:"distributor"`
	Budget      *int64  `json:"budget"`
	MpaRating   *string `json:"mpaRating"`
}

// POST /movies
func (h *HandlerSet) CreateMovie(c *gin.Context) {
	var req MovieCreateRequest
	// ShouldBindJSON 不会自动写入响应，但会添加错误到 c.Errors
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(422, gin.H{"code": "VALIDATION_ERROR", "message": "Invalid request body"})
		return
	}

	// Step 1 — minimal DB insert
	var movieID int
	err := h.DB.QueryRow(`
		INSERT INTO movies (title, genre, release_date)
		VALUES ($1, $2, $3)
		RETURNING id
	`, req.Title, req.Genre, req.ReleaseDate).Scan(&movieID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			c.JSON(422, gin.H{"code": "VALIDATION_ERROR", "message": "Movie already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": "Failed to insert movie"})
		return
	}

	// Step 2 — call boxoffice (忽略错误，降级处理)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	boxData, _ := h.BoxOffice.FetchFull(ctx, req.Title)

	// Step 3 — merge fields (用户提供的值 > 上游返回的值)
	var finalDistributor *string
	var finalBudget *int64
	var finalMpaRating *string
	var box *domain.BoxOffice

	if boxData != nil {
		box = boxData.BoxOffice
	}

	// Distributor: 用户提供 > boxoffice.distributor
	if req.Distributor != nil {
		finalDistributor = req.Distributor
	} else if boxData != nil && boxData.Distributor != nil {
		finalDistributor = boxData.Distributor
	}

	// Budget: 用户提供 > boxoffice.budget
	if req.Budget != nil {
		finalBudget = req.Budget
	} else if boxData != nil && boxData.Budget != nil {
		finalBudget = boxData.Budget
	}

	// MpaRating: 用户提供 > boxoffice.mpaRating
	if req.MpaRating != nil {
		finalMpaRating = req.MpaRating
	} else if boxData != nil && boxData.MpaRating != nil {
		finalMpaRating = boxData.MpaRating
	}

	// Step 4 — update DB with merged results
	_, err = h.DB.Exec(`
		UPDATE movies
		SET distributor = $1,
		    budget = $2,
		    mpa_rating = $3,
		    boxoffice_revenue_worldwide = $4,
		    boxoffice_revenue_opening_weekend_usa = $5,
		    boxoffice_currency = $6,
		    boxoffice_source = $7,
		    boxoffice_last_updated = $8
		WHERE id = $9
	`, finalDistributor, finalBudget, finalMpaRating,
		func() *int64 {
			if box == nil {
				return nil
			}
			return &box.Revenue.Worldwide
		}(),
		func() *int64 {
			if box == nil {
				return nil
			}
			return &box.Revenue.OpeningWeekendUSA
		}(),
		func() *string {
			if box == nil {
				return nil
			}
			return &box.Currency
		}(),
		func() *string {
			if box == nil {
				return nil
			}
			return &box.Source
		}(),
		func() *time.Time {
			if box == nil {
				return nil
			}
			return &box.LastUpdated
		}(),
		movieID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": "Failed to update movie"})
		return
	}

	// 构建响应
	resp := domain.Movie{
		ID:          fmt.Sprintf("%d", movieID),
		Title:       req.Title,
		Genre:       req.Genre,
		ReleaseDate: req.ReleaseDate,
		Distributor: finalDistributor,
		Budget:      finalBudget,
		MpaRating:   finalMpaRating,
		BoxOffice:   box,
	}

	// 201 + Location header
	c.Header("Location", "/movies/"+url.PathEscape(req.Title))
	c.JSON(http.StatusCreated, resp)
}

// GET /movies - 列表和搜索
func (h *HandlerSet) ListMovies(c *gin.Context) {
	// 解析查询参数
	q := c.Query("q")                     // 关键词搜索
	yearStr := c.Query("year")            // 年份过滤
	genre := c.Query("genre")             // 类型过滤
	distributor := c.Query("distributor") // 发行商过滤
	budgetStr := c.Query("budget")        // 预算过滤
	mpaRating := c.Query("mpaRating")     // MPA评级过滤
	limitStr := c.Query("limit")
	cursor := c.Query("cursor")

	// 默认limit
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// 解析cursor（简单实现：offset）
	offset := 0
	if cursor != "" {
		if o, err := strconv.Atoi(cursor); err == nil && o > 0 {
			offset = o
		}
	}

	// 构建SQL查询
	query := "SELECT id, title, genre, release_date, distributor, budget, mpa_rating, " +
		"boxoffice_revenue_worldwide, boxoffice_revenue_opening_weekend_usa, " +
		"boxoffice_currency, boxoffice_source, boxoffice_last_updated " +
		"FROM movies WHERE 1=1"

	args := []interface{}{}
	argIdx := 1

	// 关键词搜索（模糊匹配title）
	if q != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", argIdx)
		args = append(args, "%"+q+"%")
		argIdx++
	}

	// 年份过滤
	if yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			query += fmt.Sprintf(" AND EXTRACT(YEAR FROM release_date) = $%d", argIdx)
			args = append(args, year)
			argIdx++
		}
	}

	// 类型过滤
	if genre != "" {
		query += fmt.Sprintf(" AND LOWER(genre) = LOWER($%d)", argIdx)
		args = append(args, genre)
		argIdx++
	}

	// 发行商过滤
	if distributor != "" {
		query += fmt.Sprintf(" AND LOWER(distributor) = LOWER($%d)", argIdx)
		args = append(args, distributor)
		argIdx++
	}

	// 预算过滤
	if budgetStr != "" {
		if budget, err := strconv.ParseInt(budgetStr, 10, 64); err == nil {
			query += fmt.Sprintf(" AND budget <= $%d", argIdx)
			args = append(args, budget)
			argIdx++
		}
	}

	// MPA评级过滤
	if mpaRating != "" {
		query += fmt.Sprintf(" AND mpa_rating = $%d", argIdx)
		args = append(args, mpaRating)
		argIdx++
	}

	// 排序和分页
	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit+1, offset) // 多查一条判断是否有下一页

	// 执行查询
	rows, err := h.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": "Failed to query movies"})
		return
	}
	defer rows.Close()

	movies := []domain.Movie{}
	for rows.Next() {
		var m domain.Movie
		var id int
		var distributor, mpaRating, boxCurrency, boxSource sql.NullString
		var budget, boxWorldwide, boxOpeningUSA sql.NullInt64
		var boxLastUpdated sql.NullTime

		err := rows.Scan(&id, &m.Title, &m.Genre, &m.ReleaseDate,
			&distributor, &budget, &mpaRating,
			&boxWorldwide, &boxOpeningUSA, &boxCurrency, &boxSource, &boxLastUpdated)
		if err != nil {
			continue
		}

		m.ID = fmt.Sprintf("%d", id)

		if distributor.Valid {
			m.Distributor = &distributor.String
		}
		if budget.Valid {
			m.Budget = &budget.Int64
		}
		if mpaRating.Valid {
			m.MpaRating = &mpaRating.String
		}

		// 构建BoxOffice对象
		if boxWorldwide.Valid && boxCurrency.Valid && boxSource.Valid && boxLastUpdated.Valid {
			m.BoxOffice = &domain.BoxOffice{
				Revenue: domain.BoxOfficeRevenue{
					Worldwide:         boxWorldwide.Int64,
					OpeningWeekendUSA: boxOpeningUSA.Int64,
				},
				Currency:    boxCurrency.String,
				Source:      boxSource.String,
				LastUpdated: boxLastUpdated.Time,
			}
		}

		movies = append(movies, m)
	}

	// 判断是否有下一页
	var nextCursor *string
	if len(movies) > limit {
		movies = movies[:limit]
		next := fmt.Sprintf("%d", offset+limit)
		nextCursor = &next
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      movies,
		"nextCursor": nextCursor,
	})
}

// POST /movies/{title}/ratings - 提交评分
func (h *HandlerSet) SubmitRating(c *gin.Context) {
	title := c.Param("title")
	raterID, _ := c.Get("rater_id")

	var req struct {
		Rating float64 `json:"rating" binding:"required"`
	}
	// ShouldBindJSON 不会自动写入响应
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(422, gin.H{"code": "VALIDATION_ERROR", "message": "Invalid request body"})
		return
	}

	// 验证rating值（必须是0.5的倍数，范围0.5-5.0）
	validRatings := []float64{0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5.0}
	isValid := false
	for _, v := range validRatings {
		if req.Rating == v {
			isValid = true
			break
		}
	}
	if !isValid {
		c.JSON(422, gin.H{"code": "INVALID_RATING", "message": "Rating must be one of: 0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5.0"})
		return
	}

	// 检查电影是否存在
	var exists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM movies WHERE title = $1)", title).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "Movie not found"})
		return
	}

	// 先检查评分是否已存在
	var existingRating sql.NullFloat64
	err = h.DB.QueryRow("SELECT rating FROM ratings WHERE movie_title = $1 AND rater_id = $2",
		title, raterID).Scan(&existingRating)

	isNew := err == sql.ErrNoRows

	// Upsert评分
	_, err = h.DB.Exec(`
		INSERT INTO ratings (movie_title, rater_id, rating, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (movie_title, rater_id)
		DO UPDATE SET rating = EXCLUDED.rating, updated_at = NOW()
	`, title, raterID, req.Rating)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": "Failed to save rating"})
		return
	}

	resp := gin.H{
		"movieTitle": title,
		"raterId":    raterID,
		"rating":     req.Rating,
	}

	// 根据是否是新记录返回不同的状态码
	if isNew {
		c.Header("Location", fmt.Sprintf("/movies/%s/ratings", url.PathEscape(title)))
		c.JSON(http.StatusCreated, resp)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

// GET /movies/{title}/rating - 获取评分聚合
func (h *HandlerSet) GetRating(c *gin.Context) {
	title := c.Param("title")

	// 检查电影是否存在
	var exists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM movies WHERE title = $1)", title).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "Movie not found"})
		return
	}

	// 计算平均分和数量
	var avg sql.NullFloat64
	var count int
	err = h.DB.QueryRow(`
		SELECT AVG(rating), COUNT(*)
		FROM ratings
		WHERE movie_title = $1
	`, title).Scan(&avg, &count)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": "Failed to calculate rating"})
		return
	}

	// 如果没有评分，返回0和0
	avgValue := 0.0
	if avg.Valid {
		// 四舍五入到1位小数
		avgValue = math.Round(avg.Float64*10) / 10
	}

	c.JSON(http.StatusOK, gin.H{
		"average": avgValue,
		"count":   count,
	})
}
