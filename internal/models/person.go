package models

type Person struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	TimesLookedUp int64  `json:"times_looked_up"`
}
