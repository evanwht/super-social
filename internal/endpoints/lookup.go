package endpoints

import (
	"net/http"
	"net/mail"
	"strconv"
	"time"

	"superhuman-social/internal/clearbit"
	"superhuman-social/internal/database"
	"superhuman-social/internal/models"

	"github.com/gin-gonic/gin"
)

func LookupPerson(clearbitAPI clearbit.API, db database.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.Query("email")
		if _, err := mail.ParseAddress(email); err != nil {
			c.String(http.StatusBadRequest, "email required")
			return
		}

		person, err := clearbitAPI.LookupPerson(email)
		if err != nil {
			if err == clearbit.RetryError {
				c.Header("Retry-After", strconv.Itoa(int(time.Minute.Seconds())))
				c.String(http.StatusAccepted, err.Error())
				return
			}
			c.String(http.StatusInternalServerError, "clearbit error: %v", err)
			return
		}

		dbPerson, err := db.GetPerson(c, email)
		if err != nil && err != database.NotFoundErr {
			c.String(http.StatusInternalServerError, "unable to increment lookup count: %v", err)
			return
		} else if err == database.NotFoundErr {
			dbPerson = models.Person{
				Name:  person.Name.FullName,
				Email: email,
			}
		}
		dbPerson.TimesLookedUp++
		if err = db.SavePerson(c, &dbPerson); err != nil {
			c.String(http.StatusInternalServerError, "unable to increment lookup count: %v", err)
			return
		}

		c.JSON(http.StatusOK, person)
	}
}
