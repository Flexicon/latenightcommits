package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func indexHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page < 1 {
			page = 1
		}

		limit := 10
		offset := (page - 1) * limit

		var commits []*Commit
		db.Order("created_at desc").Limit(limit).Offset(offset).Find(&commits)

		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"commits": commits,
		})
	}
}
