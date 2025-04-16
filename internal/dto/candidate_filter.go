package dto

type CandidateFilter struct {
	AIMatchMin int      `form:"ai_match"`
	Skills     []string `form:"skills"`
	City       string   `form:"city"`
	Position   string   `form:"position"`
}
