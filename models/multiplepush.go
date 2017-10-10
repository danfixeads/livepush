package models

import null "gopkg.in/guregu/null.v3"

// MultiplePush struct
type MultiplePush struct {
	Tokens  []null.String          `json:"tokens"`
	Payload map[string]interface{} `json:"payload"`
}
