package endpoints

import (
	"net/http"

	"superhuman-social/internal/database"

	"github.com/gin-gonic/gin"
)

func ListPopular(db database.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		persons, err := db.ListAllPersons(c)
		if err != nil {
			c.String(http.StatusInternalServerError, "error looking up users: %v", err)
			return
		}
		emails := make([]string, len(persons))
		for i := range persons {
			emails[i] = persons[i].Email
		}
		c.JSON(http.StatusOK, emails)
	}
}
