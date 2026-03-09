package openai

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/askasoft/goopenai/openai/chat/completions"
	"github.com/askasoft/goopenai/openai/embeddings"
	"github.com/askasoft/goopenai/openai/files"
	"github.com/askasoft/goopenai/openai/responses"
	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/log/httplog"
)

func testFilename(name string) string {
	return filepath.Join("testdata", name)
}

func testReadFile(t *testing.T, name string) []byte {
	fn := testFilename(name)
	bs, err := fsu.ReadFile(fn)
	if err != nil {
		t.Fatalf("Failed to read file %q: %v", fn, err)
	}
	return bs
}

func testNewOpenAI(t *testing.T) *Client {
	apikey := os.Getenv("OPENAI_APIKEY")
	if apikey == "" {
		t.Skip("OPENAI_APIKEY not set")
		return nil
	}

	logs := log.NewLog()
	logs.SetLevel(log.LevelDebug)
	logger := logs.GetLogger("OPENAI")

	oai := &Client{
		BaseURL:   OpenAIBaseURL,
		APIKey:    apikey,
		Transport: httplog.LoggingRoundTripper(logger),
		Retryer:   NewRetryer(logger, 1, time.Second*3),
	}

	return oai
}

func TestOpenAICreateTextEmbeddingsAda002(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	req := &embeddings.TextEmbeddingsRequest{
		Model: "text-embedding-ada-002",
		Input: []string{"あなたはだれですか？"},
	}

	res, err := oai.CreateTextEmbeddings(context.Background(), req)
	if err != nil {
		t.Fatalf("OpenAI.CreateTextEmbeddings(): %v", err)
	} else {
		fmt.Println(len(res.Embedding()), res.Usage)
	}
}

func TestOpenAICreateTextEmbeddings3Small(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	req := &embeddings.TextEmbeddingsRequest{
		Model: "text-embedding-3-small",
		Input: []string{"あなたはだれですか？"},
	}

	res, err := oai.CreateTextEmbeddings(context.Background(), req)
	if err != nil {
		t.Fatalf("OpenAI.CreateTextEmbeddings(): %v", err)
	} else {
		fmt.Println(len(res.Embedding()), res.Usage)
	}
}

func TestOpenAICreateTextEmbeddings3LargeWithDimensions(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	req := &embeddings.TextEmbeddingsRequest{
		Model:      "text-embedding-3-large",
		Input:      []string{"あなたはだれですか？"},
		Dimensions: 1536,
	}

	res, err := oai.CreateTextEmbeddings(context.Background(), req)
	if err != nil {
		t.Fatalf("OpenAI.CreateTextEmbeddings(): %v", err)
	} else {
		fmt.Println(len(res.Embedding()), res.Usage)
	}
}

func TestOpenAICreateChatCompletion(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	req := &completions.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []completions.ChatMessage{
			{Role: RoleUser, Content: "あなたはだれですか？"},
		},
	}

	res, err := oai.CreateChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("OpenAI.CreateChatCompletion(): %v", err)
	}

	fmt.Println("-------------------------------------------")
	fmt.Println(res)
	fmt.Println(res.Usage.String())
}

func TestOpenAIWebSearchTool(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	req := &completions.ChatCompletionRequest{
		Model: "gpt-4o-search-preview",
		Messages: []completions.ChatMessage{
			{Role: RoleUser, Content: "今年春アニメのおすすめは？"},
		},
		WebSearchOptions: &completions.WebSearchOptions{
			SearchContextSize: "medium",
			UserLocation: &completions.UserLocation{
				Type: "approximate",
				Approximate: &completions.Approximate{
					Country: "JP",
				},
			},
		},
	}

	res, err := oai.CreateChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("OpenAI.CreateChatCompletion(): %v", err)
	}

	fmt.Println("-------------------------------------------")
	fmt.Println(res)
	fmt.Println(res.Usage.String())
}

func TestOpenAIImageAnalyze(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	req := &completions.ChatCompletionRequest{
		Model: "gpt-4.1",
		Messages: []completions.ChatMessage{
			{
				Role: RoleUser,
				Content: []completions.MessageContent{
					completions.TextContent("画像の中に「個人情報が含まれているかどうか」を判定してください。"),
					completions.ImageURLContent("https://s3.amazonaws.com/cdn.freshdesk.com/data/helpdesk/attachments/production/50012396079/original/j3UQrTiD9AcapYi98QjFjTKXptsLq4TSBA.png?1720516588", ""),
				},
			},
		},
	}

	res, err := oai.CreateChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("OpenAI.CreateChatCompletion(): %v", err)
	}

	fmt.Println("-------------------------------------------")
	fmt.Println(res)
	fmt.Println(res.Usage.String())
}

func TestOpenAICompeletionsFileData(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	testfiles := []string{"earth.pdf"} // pdf only

	for i, testfile := range testfiles {
		filedata := testReadFile(t, testfile)

		req := &completions.ChatCompletionRequest{
			Model: "gpt-5.2",
			Messages: []completions.ChatMessage{
				{
					Role: RoleUser,
					Content: []completions.MessageContent{
						completions.TextContent("ファイルの中に「個人情報が含まれているかどうか」を判定してください。"),
						completions.FileDataContent(testfile, filedata),
					},
				},
			},
		}

		res, err := oai.CreateChatCompletion(context.Background(), req)
		if err != nil {
			t.Errorf("#%d OpenAI.CreateChatCompletion(): %v", i, err)
			continue
		}

		fmt.Println("-------------------------------------------")
		fmt.Println(res)
	}
}

func TestOpenAIResponsesFileData(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	testfiles := []string{
		"earth.pdf",
		"earth.docx",
		// "earth.pptx",
		// "earth.xlsx",
		// "earth.txt",
		// "earth.csv",
		// "earth.tsv", // unsupport
	}

	for i, testfile := range testfiles {
		filedata := testReadFile(t, testfile)

		req := &responses.CreateRequest{
			Model: "gpt-5.2",
			Input: []responses.ResponseMessage{
				{
					Role: RoleUser,
					Content: []responses.ResponseMessageContent{
						responses.TextContent("ファイルの中に「個人情報が含まれているかどうか」を判定してください。"),
						responses.FileDataContent(testfile, filedata),
					},
				},
			},
		}

		res, err := oai.CreateResponse(context.Background(), req)
		if err != nil {
			t.Errorf("#%d OpenAI.CreateResponse(): %v", i, err)
			continue
		}

		fmt.Println("-------------------------------------------")
		fmt.Println(res)
		fmt.Println("-------------------------------------------")
		fmt.Println(res.OutputText())
	}
}

func TestOpenAICompeletionsFileID(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	testfiles := []string{"earth.pdf", "earth.docx"}

	for i, testfile := range testfiles {
		filedata := testReadFile(t, testfile)

		freq := &files.CreateRequest{
			FileName:     testfile,
			FileData:     filedata,
			Purpose:      files.FilePurposeAssistants,
			ExpiresAfter: 3600,
		}

		fres, err := oai.CreateFile(context.Background(), freq)
		if err != nil {
			t.Errorf("#%d OpenAI.CreateFile(): %v", i, err)
			continue
		}

		creq := &completions.ChatCompletionRequest{
			Model: "gpt-5.2",
			Messages: []completions.ChatMessage{
				{
					Role: RoleUser,
					Content: []completions.MessageContent{
						completions.TextContent("ファイルの中に「個人情報が含まれているかどうか」を判定してください。"),
						completions.FileIDContent(fres.ID),
					},
				},
			},
		}

		cres, err := oai.CreateChatCompletion(context.Background(), creq)
		if err != nil {
			t.Errorf("#%d OpenAI.CreateChatCompletion(): %v", i, err)
			continue
		}

		fmt.Println("-------------------------------------------")
		fmt.Println(cres)
	}
}

func TestOpenAIResponsesFileID(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	testfiles := []string{
		"earth.pdf",
		"earth.docx",
		// "earth.pptx",
		// "earth.xlsx",
		// "earth.txt",
		// "earth.csv",
		// "earth.tsv", // unsupport
	}

	for i, testfile := range testfiles {
		filedata := testReadFile(t, testfile)

		freq := &files.CreateRequest{
			FileName:     testfile,
			FileData:     filedata,
			Purpose:      files.FilePurposeAssistants,
			ExpiresAfter: 3600,
		}

		fres, err := oai.CreateFile(context.Background(), freq)
		if err != nil {
			t.Errorf("#%d OpenAI.CreateFile(): %v", i, err)
			continue
		}

		creq := &responses.CreateRequest{
			Model: "gpt-5.2",
			Input: []responses.ResponseMessage{
				{
					Role: RoleUser,
					Content: []responses.ResponseMessageContent{
						responses.TextContent("ファイルの中に「個人情報が含まれているかどうか」を判定してください。"),
						responses.FileIDContent(fres.ID),
					},
				},
			},
		}

		cres, err := oai.CreateResponse(context.Background(), creq)
		if err != nil {
			t.Errorf("#%d OpenAI.CreateResponse(): %v", i, err)
			continue
		}

		fmt.Println("-------------------------------------------")
		fmt.Println(cres)
		fmt.Println("-------------------------------------------")
		fmt.Println(cres.OutputText())
	}
}

func TestOpenAICreateFile(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	cs := []string{
		"earth.docx",
		// "earth.tsv", // unsupport
	}

	for i, file := range cs {
		data := testReadFile(t, file)

		req := &files.CreateRequest{
			FileName:     file,
			FileData:     data,
			Purpose:      files.FilePurposeAssistants,
			ExpiresAfter: 3600,
		}

		res, err := oai.CreateFile(context.Background(), req)
		if err != nil {
			t.Errorf("#%d OpenAI.CreateFile(): %v", i, err)
			continue
		}

		fmt.Println("-------------------------------------------")
		fmt.Println(res)
	}
}
