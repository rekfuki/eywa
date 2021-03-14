package pagination

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Validate validates pagination and sets the alues if they are missing
func Validate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == http.MethodGet || c.Request().Method == http.MethodPost {

				var err error
				page := 1
				perPage := 20

				qPage := c.QueryParam("page")
				if qPage != "" {
					page, err = strconv.Atoi(qPage)
					if err != nil {
						return c.JSON(http.StatusBadRequest, map[string]interface{}{
							"message": "Validation error",
							"details": map[string]interface{}{
								"page": err,
							},
						})
					}
				}

				qPerPage := c.QueryParam("per_page")
				if qPerPage != "" {
					perPage, err = strconv.Atoi(qPerPage)
					if err != nil {
						return c.JSON(http.StatusBadRequest, map[string]interface{}{
							"message": "Validation error",
							"details": map[string]interface{}{
								"per_page": err,
							},
						})
					}
				}

				c.Set("page_number", page)
				c.Set("per_page", perPage)
			}
			return next(c)
		}
	}
}
