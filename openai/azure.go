package openai

import (
	"fmt"
	"net/http"
)

func AzureOpenAIBaseURL(domain, deployment string) string {
	return "https://" + domain + "/openai/deployments/" + deployment
}

func AzureOpenAIServicePath(apiver string) func(string, ...any) string {
	return func(format string, args ...any) string {
		return fmt.Sprintf(format, args...) + "?api-version=" + apiver
	}
}

func AzureOpenAIAuthenticate(req *http.Request, apikey string) {
	req.Header.Set("Api-Key", apikey)
}

func NewAzureClient(domain, deployment, apiver string) *Client {
	return &Client{
		BaseURL:      AzureOpenAIBaseURL(domain, deployment),
		Authenticate: AzureOpenAIAuthenticate,
		ServicePath:  AzureOpenAIServicePath(apiver),
	}
}
