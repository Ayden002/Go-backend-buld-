package boxoffice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"interview/internal/domain"
)

// Client BoxOffice API客户端
type Client struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

// 创建新的BoxOffice客户端
func New(baseURL, key string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  key,
		Client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// responseDTO上游API响应结构
type responseDTO struct {
	Title       string `json:"title"`
	Distributor string `json:"distributor"`
	ReleaseDate string `json:"releaseDate"`
	Budget      int64  `json:"budget"`
	Revenue     struct {
		Worldwide         int64 `json:"worldwide"`
		OpeningWeekendUSA int64 `json:"openingWeekendUSA"`
	} `json:"revenue"`
	MpaRating string `json:"mpaRating"`
}

// BoxOfficeData 扩展的票房数据，包含额外字段
type BoxOfficeData struct {
	BoxOffice   *domain.BoxOffice
	Distributor *string
	Budget      *int64
	MpaRating   *string
}

// Fetch 查询票房信息
func (c *Client) Fetch(ctx context.Context, title string) (*domain.BoxOffice, error) {
	data, err := c.FetchFull(ctx, title)
	if err != nil || data == nil {
		return nil, err
	}
	return data.BoxOffice, nil
}

// FetchFull 查询完整的票房信息
func (c *Client) FetchFull(ctx context.Context, title string) (*BoxOfficeData, error) {
	u, _ := url.Parse(c.BaseURL)
	u.Path = "/boxoffice"

	q := u.Query()
	q.Set("title", title)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	req.Header.Set("X-API-Key", c.APIKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upstream_error: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, nil //如果上游未找到，返回nil
	case http.StatusOK:

	default:
		return nil, fmt.Errorf("upstream_error: status %d", resp.StatusCode)
	}

	var dto responseDTO
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return nil, fmt.Errorf("decode_error: %w", err)
	}

	t, _ := time.Parse(time.RFC3339, dto.ReleaseDate)

	result := &BoxOfficeData{
		BoxOffice: &domain.BoxOffice{
			Revenue: domain.BoxOfficeRevenue{
				Worldwide:         dto.Revenue.Worldwide,
				OpeningWeekendUSA: dto.Revenue.OpeningWeekendUSA,
			},
			Currency:    "USD",
			Source:      "ExampleBoxOfficeAPI",
			LastUpdated: t,
		},
	}

	// 添加额外字段
	if dto.Distributor != "" {
		result.Distributor = &dto.Distributor
	}
	if dto.Budget > 0 {
		result.Budget = &dto.Budget
	}
	if dto.MpaRating != "" {
		result.MpaRating = &dto.MpaRating
	}

	return result, nil
}
