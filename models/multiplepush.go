package models

import null "gopkg.in/guregu/null.v3"

// MultiplePush struct
type MultiplePush struct {
	ClientID null.String            `json:"clientid"`
	Tokens   []null.String          `json:"tokens"`
	Payload  map[string]interface{} `json:"payload"`
}
