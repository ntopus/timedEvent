package data_types

type ArangoCloudEvent struct {
	ArangoKey string `json:"_key"`
	ArangoId  string `json:"_id"`
	ArangoRev string `json:"_rev"`
	CloudEvent
}
