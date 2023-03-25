package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
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
		if limit > 1000 {
			limit = 1000
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
	type DailyStatsResults struct {
		Date  time.Time
		Count int64
	}

	type CachedStats struct {
		DailyStats   []*DailyStatsResults
		TotalCommits int64
	}

	var statsCache *CachedStats
	// Clear cache every 5 minutes
	go func() {
		for range time.Tick(5 * time.Minute) {
			statsCache = nil
		}
	}()

	return func(c echo.Context) error {
		// If cache is empty, perform data fetching
		if statsCache == nil {
			statsCache = &CachedStats{}
			dbSession := db.Session(&gorm.Session{PrepareStmt: true})

			err := dbSession.
				Raw(`
					SELECT COUNT(id) as 'count', DATE(created_at) as 'date'
					FROM commits
					GROUP BY date
					ORDER BY date DESC
					LIMIT 7;
				`).Scan(&statsCache.DailyStats).Error
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to retrieve stats")
			}

			err = dbSession.Raw(`SELECT COUNT(id) FROM commits;`).Scan(&statsCache.TotalCommits).Error
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to retrieve stats total")
			}
		}

		// Setup global chart options
		globalOpts := charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeWesteros,
		})

		// Prep dates stats chart
		datesChart := charts.NewBar()
		datesChart.SetGlobalOptions(globalOpts,
			charts.WithTitleOpts(opts.Title{
				Title: "Commit Stats by date",
				Subtitle: fmt.Sprintf(
					"🤓 total commits added each day over the last 7 days - out of %s in total",
					formatNumber(statsCache.TotalCommits),
				),
			}),
			charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		)

		// Prepare data to fit into chart
		dates := make([]string, 0)
		data := make([]opts.BarData, 0)
		// Loop through stats in reverse order since we want to display the date ascending from right to left
		for i := len(statsCache.DailyStats) - 1; i >= 0; i-- {
			row := statsCache.DailyStats[i]
			dates = append(dates, row.Date.Format("2006/01/02"))
			data = append(data, opts.BarData{Value: row.Count})
		}
		// Put dates data into chart
		datesChart.SetXAxis(dates).AddSeries("Commits on date", data)

		// Setup page to hold all charts
		page := components.NewPage()
		page.PageTitle = "Stats | LateNightCommits"
		page.AddCharts(datesChart)

		return page.Render(c.Response())
	}
}
