package models

import null "gopkg.in/guregu/null.v3"

// MultiplePush struct
type MultiplePush struct {
	ClientID null.Int      `json:"clientid"`
	Tokens   []null.String `json:"tokens"`
	Title    null.String   `json:"title"`
	Subtitle null.String   `json:"subtitle"`
	Body     null.String   `json:"body"`
	Badge    null.Int      `json:"badge"`
	Image    null.String   `json:"image"`
	Sound    null.String   `json:"sound"`
}