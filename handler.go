package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

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
