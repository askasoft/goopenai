package models

import "github.com/askasoft/goopenai/openai/internal"

type Model struct {
	ID      string `json:"id,omitempty"`
	Object  string `json:"object,omitempty"`
	Created int64  `json:"created,omitempty"`
	OwnedBy string `json:"owned_by,omitempty"`
}

func (m *Model) String() string {
	return internal.ToString(m)
}

type ListModelsResponse struct {
	Object string   `json:"object"`
	Data   []*Model `json:"data"`
}
