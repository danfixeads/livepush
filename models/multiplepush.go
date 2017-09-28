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
	Type     null.Int      `json:"type"`
	TheID    null.String   `json:"the_id"`
	UserID   null.Int      `json:"user_id"`
	Track    null.String   `json:"track"`
	Main     null.String   `json:"main"`
	Options  []PushOption  `json:"options"`
}
