package model

type Embedding struct {
	TextField   string         `json:"textField,omitempty"`
	VectorField string         `json:"vectorField,omitempty"`
	Model       EmbeddingModel `json:"model,omitempty"`
	Enabled     bool           `json:"enabled,omitempty"`
}
