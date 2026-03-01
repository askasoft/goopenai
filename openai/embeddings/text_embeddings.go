package embeddings

import (
	"fmt"

	"github.com/askasoft/pango/doc/jsonx"
	"github.com/askasoft/pango/str"
)

type TextEmbeddingsRequest struct {
	// Input Input text to embed (required)
	Input []string `json:"input,omitempty"`

	// Model ID of the model to use (required)
	Model string `json:"model,omitempty"`

	// Dimensions (optional) The number of dimensions the resulting output embeddings should have. Only supported in text-embedding-3 and later models.
	Dimensions int `json:"dimensions,omitempty"`

	// EncodingFormat (optional) "float" or "base64"
	EncodingFormat string `json:"encoding_format,omitempty"`

	// User (optional) A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse.
	User string `json:"user,omitempty"`
}

func (te *TextEmbeddingsRequest) String() string {
	return toString(te)
}

func (te *TextEmbeddingsRequest) InputRuneCount() int {
	cnt := 0
	for _, s := range te.Input {
		cnt += str.RuneCount(s)
	}
	return cnt
}

type EmbeddingData struct {
	// The index of the embedding in the list of embeddings.
	Index int `json:"index"`

	// The object type, which is always "embedding".
	Object string `json:"object,omitempty"`

	// The embedding vector, which is a list of floats.
	Embedding []float64 `json:"embedding"`
}

type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

func (u *Usage) Add(a *Usage) {
	u.PromptTokens += a.PromptTokens
	u.TotalTokens += a.TotalTokens
}

func (u *Usage) String() string {
	return fmt.Sprintf("P: %d, T: %d", u.PromptTokens, u.TotalTokens)
}

type TextEmbeddingsResponse struct {
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Object string          `json:"object"`
	Usage  Usage           `json:"usage"`
}

func (te *TextEmbeddingsResponse) String() string {
	return toString(te)
}

func (te *TextEmbeddingsResponse) Embedding() []float64 {
	if len(te.Data) > 0 {
		return te.Data[0].Embedding
	}
	return nil
}

func toString(o any) string {
	return jsonx.Prettify(o)
}
