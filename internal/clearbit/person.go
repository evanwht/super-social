package clearbit

type (
	Person struct {
		Name     Name   `json:"name"`
		Location string `json:"location"` // possibly break down into parts
		Avatar   string `json:"avatar"`   // image TODO what is empty? pull from other sources?
	}
	Name struct {
		FullName   string `json:"fullName"`
		GivenName  string `json:"givenName"`
		FamilyName string `json:"familyName"`
	}
)

func (n Name) String() string {
	return n.FullName
}
