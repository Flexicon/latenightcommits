package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func indexHandler(e *echo.Echo) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"api":       "latenightcommits",
			"resources": e.Routes(),
		})
	}
}

func commitLogHandler(db *gorm.DB) echo.HandlerFunc {
	type CommitLogResponse struct {
		Log         []*Commit `json:"log"`
		HasNextPage bool      `json:"has_next_page"`
	}

	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page < 1 {
			page = 1
		}

		limit, _ := strconv.Atoi(c.QueryParam("per_page"))
		if limit < 1 {
			limit = 20
		}

		var commits []*Commit
		offset := (page - 1) * limit
		db.Order("created_at desc").Limit(limit).Offset(offset).Find(&commits)

		return c.JSON(http.StatusOK, CommitLogResponse{
			Log:         commits,
			HasNextPage: len(commits) == limit,
		})
	}
}

func statsHandler(db *gorm.DB) echo.HandlerFunc {
	type Results struct {
		Count int64
		Date  time.Time
	}
	var StatsCache []*Results
	// Clear cache every 5 minutes
	go func() {
		for range time.Tick(5 * time.Minute) {
			StatsCache = nil
		}
	}()

	return func(c echo.Context) error {
		if StatsCache == nil {
			err := db.Session(&gorm.Session{PrepareStmt: true}).
				Raw(`
					SELECT COUNT(id) as 'count', DATE(created_at) as 'date'
					FROM commits
					GROUP BY date
					ORDER BY date DESC
					LIMIT 7;
				`).Scan(&StatsCache).Error

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to retrieve stats")
			}
		}

		// Prep chart
		bar := charts.NewBar()
		bar.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{
				Title:    "Commit Stats",
				Subtitle: "🤓 total commits added each day over the last 7 days",
			}),
			charts.WithInitializationOpts(opts.Initialization{
				PageTitle: "Stats | LateNightCommits",
			}),
			charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		)

		// Prepare data to fit into chart
		dates := make([]string, 0)
		data := make([]opts.BarData, 0)

		for _, row := range StatsCache {
			dates = append(dates, row.Date.Format("2006/01/02"))
			data = append(data, opts.BarData{Value: row.Count})
		}
		// Put data into chart
		bar.SetXAxis(dates).AddSeries("Commits on date", data)

		return bar.Render(c.Response())
	}
}
