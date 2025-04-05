package entity

type HR struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	DepID     int `json:"dep_id"`
	CompanyID int `json:"company_id"`
}

type HRRegistration struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	DepID     int    `json:"dep_id"`
	CompanyID int    `json:"company_id"`
}
