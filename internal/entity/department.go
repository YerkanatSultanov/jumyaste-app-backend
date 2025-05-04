package entity

type Department struct {
	ID        int    `json:"id"`
	CompanyId int    `json:"company_id"`
	Name      string `json:"name"`
	HrCount   int    `json:"hr_count"`
}
