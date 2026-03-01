package openai

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/askasoft/goopenai/openai/chat/completions"
	"github.com/askasoft/goopenai/openai/embeddings"
	"github.com/askasoft/pango/log"
)

func testNewAzureOpenAI(t *testing.T, deploy string) *Client {
	apikey := os.Getenv("AZURE_OPENAI_APIKEY")
	if apikey == "" {
		t.Skip("AZURE_OPENAI_APIKEY not set")
		return nil
	}

	domain := os.Getenv("AZURE_OPENAI_DOMAIN")
	if domain == "" {
		t.Skip("AZURE_OPENAI_DOMAIN not set")
		return nil
	}

	deployment := os.Getenv(deploy)

	logs := log.NewLog()
	logs.SetLevel(log.LevelDebug)
	aoai := &Client{
		BaseURL:     AzureOpenAIBaseURL(domain, deployment),
		APIKey:      apikey,
		Logger:      logs.GetLogger("AZUREOPENAI"),
		MaxRetries:  1,
		RetryAfter:  time.Second * 3,
		ServicePath: AzureOpenAIServicePath("2024-06-01"),
	}

	return aoai
}

func TestAzureOpenAICreateChatCompletion(t *testing.T) {
	aoai := testNewAzureOpenAI(t, "AZURE_OPENAI_CHAT_DEPLOYMENT")
	if aoai == nil {
		return
	}

	req := &completions.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []completions.ChatMessage{
			{Role: RoleUser, Content: "あなたはだれですか？"},
		},
	}

	res, err := aoai.CreateChatCompletion(context.TODO(), req)
	if err != nil {
		t.Fatalf("AzureOpenAI.CreateChatCompletion(): %v", err)
	}

	fmt.Println(res)
}

func TestAzureOpenAICreateTextEmbeddingsAda002(t *testing.T) {
	aoai := testNewAzureOpenAI(t, "AZURE_OPENAI_TEMB_DEPLOYMENT")
	if aoai == nil {
		return
	}

	req := &embeddings.TextEmbeddingsRequest{
		Model: "text-embedding-ada-002",
		Input: []string{"あなたはだれですか？"},
	}

	res, err := aoai.CreateTextEmbeddings(context.TODO(), req)
	if err != nil {
		t.Fatalf("AzureOpenAI.CreateTextEmbeddings(): %v", err)
	} else {
		fmt.Println(len(res.Embedding()), res.Usage)
	}
}

func TestAzureOpenAICreateTextEmbeddings3Small(t *testing.T) {
	aoai := testNewAzureOpenAI(t, "AZURE_OPENAI_TEMB_DEPLOYMENT")
	if aoai == nil {
		return
	}

	req := &embeddings.TextEmbeddingsRequest{
		Model: "text-embedding-3-small",
		Input: []string{"あなたはだれですか？"},
	}

	res, err := aoai.CreateTextEmbeddings(context.TODO(), req)
	if err != nil {
		t.Fatalf("AzureOpenAI.CreateTextEmbeddings(): %v", err)
	} else {
		fmt.Println(len(res.Embedding()), res.Usage)
	}
}

func TestAzureCreateTextEmbeddings3LargeWithDimensions(t *testing.T) {
	aoai := testNewAzureOpenAI(t, "AZURE_OPENAI_TEMB_DEPLOYMENT")
	if aoai == nil {
		return
	}

	req := &embeddings.TextEmbeddingsRequest{
		Model:      "text-embedding-3-large",
		Input:      []string{"あなたはだれですか？"},
		Dimensions: 1536,
	}

	res, err := aoai.CreateTextEmbeddings(context.TODO(), req)
	if err != nil {
		t.Fatalf("AzureOpenAI.CreateTextEmbeddings(): %v", err)
	} else {
		fmt.Println(len(res.Embedding()), res.Usage)
	}
}
