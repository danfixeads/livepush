package models

import null "gopkg.in/guregu/null.v3"

// MultiplePush struct
type MultiplePush struct {
	ClientID null.Int               `json:"clientid"`
	Tokens   []null.String          `json:"tokens"`
	Payload  map[string]interface{} `json:"payload"`
}
