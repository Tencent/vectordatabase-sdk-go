package model

type Embedding struct {
	Field       string         `json:"field,omitempty"`
	VectorField string         `json:"vectorField,omitempty"`
	Model       EmbeddingModel `json:"model,omitempty"`
	Enabled     bool           `json:"enabled,omitempty"`
}
