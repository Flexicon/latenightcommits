package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func commitsLogHandler(db *gorm.DB) echo.HandlerFunc {
	type CommitsLogResponse struct {
		Log         []*Commit `json:"log"`
		HasNextPage bool      `json:"has_next_page"`
	}

	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page < 1 {
			page = 1
		}

		limit := 10
		offset := (page - 1) * limit

		var commits []*Commit
		db.Order("created_at desc").Limit(limit).Offset(offset).Find(&commits)

		return c.JSON(http.StatusOK, CommitsLogResponse{
			Log:         commits,
			HasNextPage: len(commits) == limit,
		})
	}
}
