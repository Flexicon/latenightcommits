package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func indexHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var commits []*Commit
		db.Order("created_at desc").Find(&commits)

		return c.Render(http.StatusOK, "index.html", map[string]interface{}{"commits": commits})
	}
}
