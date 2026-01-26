package core

type NewChat struct {
	ExternalId string `json:"external_id"`
}

type Chat struct {
	ExternalId string `json:"external_id"`
	ID         string `json:"id"`
}
