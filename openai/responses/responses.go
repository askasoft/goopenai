package responses

import (
	"strings"

	"github.com/askasoft/goopenai/openai/shared"
	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/doc/jsonx"
	"github.com/askasoft/pango/net/dataurl"
	"github.com/askasoft/pango/net/mimex"
)

// The error code for the response.
const (
	ResponseErrorCodeServerError                 = "server_error"
	ResponseErrorCodeRateLimitExceeded           = "rate_limit_exceeded"
	ResponseErrorCodeInvalidPrompt               = "invalid_prompt"
	ResponseErrorCodeVectorStoreTimeout          = "vector_store_timeout"
	ResponseErrorCodeInvalidImage                = "invalid_image"
	ResponseErrorCodeInvalidImageFormat          = "invalid_image_format"
	ResponseErrorCodeInvalidBase64Image          = "invalid_base64_image"
	ResponseErrorCodeInvalidImageURL             = "invalid_image_url"
	ResponseErrorCodeImageTooLarge               = "image_too_large"
	ResponseErrorCodeImageTooSmall               = "image_too_small"
	ResponseErrorCodeImageParseError             = "image_parse_error"
	ResponseErrorCodeImageContentPolicyViolation = "image_content_policy_violation"
	ResponseErrorCodeInvalidImageMode            = "invalid_image_mode"
	ResponseErrorCodeImageFileTooLarge           = "image_file_too_large"
	ResponseErrorCodeUnsupportedImageMediaType   = "unsupported_image_media_type"
	ResponseErrorCodeEmptyImageFile              = "empty_image_file"
	ResponseErrorCodeFailedToDownloadImage       = "failed_to_download_image"
	ResponseErrorCodeImageFileNotFound           = "image_file_not_found"
)

type Metadata = map[string]string
type Reasoning = shared.Reasoning

// Specifies the processing type used for serving the request.
//
//   - If set to 'auto', then the request will be processed with the service tier
//     configured in the Project settings. Unless otherwise configured, the Project
//     will use 'default'.
//   - If set to 'default', then the request will be processed with the standard
//     pricing and performance for the selected model.
//   - If set to '[flex](https://platform.openai.com/docs/guides/flex-processing)' or
//     '[priority](https://openai.com/api-priority-processing/)', then the request
//     will be processed with the corresponding service tier.
//   - When not set, the default behavior is 'auto'.
//
// When the `service_tier` parameter is set, the response body will include the
// `service_tier` value based on the processing mode actually used to serve the
// request. This response value may be different from the value set in the
// parameter.
type ResponseServiceTier string

const (
	ResponseServiceTierAuto     ResponseServiceTier = "auto"
	ResponseServiceTierDefault  ResponseServiceTier = "default"
	ResponseServiceTierFlex     ResponseServiceTier = "flex"
	ResponseServiceTierScale    ResponseServiceTier = "scale"
	ResponseServiceTierPriority ResponseServiceTier = "priority"
)

// The status of the response generation. One of `completed`, `failed`,
// `in_progress`, `cancelled`, `queued`, or `incomplete`.
type ResponseStatus string

const (
	ResponseStatusCompleted  ResponseStatus = "completed"
	ResponseStatusFailed     ResponseStatus = "failed"
	ResponseStatusInProgress ResponseStatus = "in_progress"
	ResponseStatusCancelled  ResponseStatus = "cancelled"
	ResponseStatusQueued     ResponseStatus = "queued"
	ResponseStatusIncomplete ResponseStatus = "incomplete"
)

func toString(o any) string {
	return jsonx.Prettify(o)
}

// Options for streaming responses. Only set this when you set `stream: true`.
type ResponseStreamOptions struct {
	// When true, stream obfuscation will be enabled. Stream obfuscation adds random
	// characters to an `obfuscation` field on streaming delta events to normalize
	// payload sizes as a mitigation to certain side-channel attacks. These obfuscation
	// fields are included by default, but add a small amount of overhead to the data
	// stream. You can set `include_obfuscation` to false to optimize for bandwidth if
	// you trust the network links between your application and the OpenAI API.
	IncludeObfuscation bool `json:"include_obfuscation"`
}

// The property Type is required.
type ResponseContextManagement struct {
	// The context management entry type. Currently only 'compaction' is supported.
	Type string `json:"type"`

	// Token threshold at which compaction should be triggered for this entry.
	CompactThreshold int64 `json:"compact_threshold,omitempty"`
}

// Constrains the verbosity of the model's response. Lower values will result in
// more concise responses, while higher values will result in more verbose
// responses. Currently supported values are `low`, `medium`, and `high`.
type ResponseTextConfigVerbosity string

const (
	ResponseTextConfigVerbosityLow    ResponseTextConfigVerbosity = "low"
	ResponseTextConfigVerbosityMedium ResponseTextConfigVerbosity = "medium"
	ResponseTextConfigVerbosityHigh   ResponseTextConfigVerbosity = "high"
)

type ResponseTextConfigFormat struct {
	Type string `json:"type"`
}

// Configuration options for a text response from the model. Can be plain text or
// structured JSON data. Learn more:
//
// - [Text inputs and outputs](https://platform.openai.com/docs/guides/text)
// - [Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs)
type ResponseTextConfigParam struct {
	// Constrains the verbosity of the model's response. Lower values will result in
	// more concise responses, while higher values will result in more verbose
	// responses. Currently supported values are `low`, `medium`, and `high`.
	//
	// Any of "low", "medium", "high".
	Verbosity ResponseTextConfigVerbosity `json:"verbosity,omitempty"`

	// An object specifying the format that the model must output.
	//
	// Configuring `{ "type": "json_schema" }` enables Structured Outputs, which
	// ensures the model will match your supplied JSON schema. Learn more in the
	// [Structured Outputs guide](https://platform.openai.com/docs/guides/structured-outputs).
	//
	// The default format is `{ "type": "text" }` with no additional options.
	//
	// **Not recommended for gpt-4o and newer models:**
	//
	// Setting to `{ "type": "json_object" }` enables the older JSON mode, which
	// ensures the message the model generates is valid JSON. Using `json_schema` is
	// preferred for models that support it.
	Format *ResponseTextConfigFormat `json:"format,omitempty"`
}

// Reference to a prompt template and its variables.
// [Learn more](https://platform.openai.com/docs/guides/text?api-mode=responses#reusable-prompts).
//
// The property ID is required.
type ResponsePrompt struct {
	// The unique identifier of the prompt template to use.
	ID string `json:"id"`

	// Optional version of the prompt template.
	Version string `json:"version,omitempty"`

	// Optional map of values to substitute in for variables in your prompt. The
	// substitution values can either be strings, or other Response input types like
	// images or files.
	Variables map[string]any `json:"variables,omitempty"`
}

// Defines a function in your own code the model can choose to call. Learn more
// about
// [function calling](https://platform.openai.com/docs/guides/function-calling).
//
// The properties Name, Parameters, Strict, Type are required.
type FunctionTool struct {
	// Whether to enforce strict parameter validation. Default `true`.
	Strict bool `json:"strict"`

	// A JSON schema object describing the parameters of the function.
	Parameters map[string]any `json:"parameters"`

	// The name of the function to call.
	Name string `json:"name"`

	// A description of the function. Used by the model to determine whether or not to
	// call the function.
	Description string `json:"description,omitempty"`

	// The type of the function tool. Always `function`.
	//
	// This field can be elided, and will marshal its zero value as "function".
	Type string `json:"type"`
}

// Weights that control how reciprocal rank fusion balances semantic embedding
// matches versus sparse keyword matches when hybrid search is enabled.
//
// The properties EmbeddingWeight, TextWeight are required.
type FileSearchToolRankingOptionsHybridSearch struct {
	// The weight of the embedding in the reciprocal ranking fusion.
	EmbeddingWeight float64 `json:"embedding_weight"`

	// The weight of the text in the reciprocal ranking fusion.
	TextWeight float64 `json:"text_weight"`
}

// Ranking options for search.
type FileSearchToolRankingOptions struct {
	// The score threshold for the file search, a number between 0 and 1. Numbers
	// closer to 1 will attempt to return only the most relevant results, but may
	// return fewer results.
	ScoreThreshold float64 `json:"score_threshold,omitempty"`

	// Weights that control how reciprocal rank fusion balances semantic embedding
	// matches versus sparse keyword matches when hybrid search is enabled.
	HybridSearch FileSearchToolRankingOptionsHybridSearch `json:"hybrid_search,omitzero"`

	// The ranker to use for the file search.
	//
	// Any of "auto", "default-2024-11-15".
	Ranker string `json:"ranker,omitempty"`
}

// A tool that searches for relevant content from uploaded files. Learn more about
// the
// [file search tool](https://platform.openai.com/docs/guides/tools-file-search).
//
// The properties Type, VectorStoreIDs are required.
type FileSearchTool struct {
	// The IDs of the vector stores to search.
	VectorStoreIDs []string `json:"vector_store_ids"`

	// The maximum number of results to return. This number should be between 1 and 50
	// inclusive.
	MaxNumResults int `json:"max_num_results,omitempty"`

	// A filter to apply.
	Filters any `json:"filters,omitempty"`

	// Ranking options for search.
	RankingOptions FileSearchToolRankingOptions `json:"ranking_options,omitzero"`

	// The type of the file search tool. Always `file_search`.
	//
	// This field can be elided, and will marshal its zero value as "file_search".
	Type string `json:"type"`
}

type ResponseMessageContent struct {
	Type     string `json:"type,omitempty"`
	Text     string `json:"text,omitempty"`
	FileID   string `json:"file_id,omitempty"`
	Filename string `json:"filename,omitempty"`
	FileData string `json:"file_data,omitempty"`
	FileURL  string `json:"file_url,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
	Detail   string `json:"detail,omitempty"`
}

type ResponseMessage struct {
	Role    string                   `json:"role"`
	Content []ResponseMessageContent `json:"content"`
}

const (
	ResponseInputTypeText  = "input_text"
	ResponseInputTypeImage = "input_image"
	ResponseInputTypeFile  = "input_file"
)

func TextContent(text string) ResponseMessageContent {
	return ResponseMessageContent{Type: ResponseInputTypeText, Text: text}
}

func ImageDataContent(name string, data []byte, detail ...string) ResponseMessageContent {
	mediaType := mimex.MediaTypeByFilename(name, "image/jpeg")
	dataURL := dataurl.Encode(mediaType, data)
	return ResponseMessageContent{Type: ResponseInputTypeImage, ImageURL: dataURL, Detail: asg.First(detail)}
}

func ImageURLContent(url string, detail ...string) ResponseMessageContent {
	return ResponseMessageContent{Type: ResponseInputTypeImage, ImageURL: url, Detail: asg.First(detail)}
}

// https://developers.openai.com/api/docs/assistants/tools/file-search#supported-files
func FileDataContent(filename string, data []byte) ResponseMessageContent {
	mediaType := mimex.MediaTypeByFilename(filename, "text/plain")
	dataURL := dataurl.Encode(mediaType, data)
	return ResponseMessageContent{Type: ResponseInputTypeFile, Filename: filename, FileData: dataURL}
}

func FileIDContent(fileid, filename string) ResponseMessageContent {
	return ResponseMessageContent{Type: ResponseInputTypeFile, FileID: fileid, Filename: filename}
}

func FileURLContent(fileurl, filename string) ResponseMessageContent {
	return ResponseMessageContent{Type: ResponseInputTypeFile, FileURL: fileurl, Filename: filename}
}

type CreateRequest struct {
	// Whether to run the model response in the background.
	// [Learn more](https://platform.openai.com/docs/guides/background).
	Background bool `json:"background,omitempty"`

	// A system (or developer) message inserted into the model's context.
	//
	// When using along with `previous_response_id`, the instructions from a previous
	// response will not be carried over to the next response. This makes it simple to
	// swap out system (or developer) messages in new responses.
	Instructions string `json:"instructions,omitempty"`

	// An upper bound for the number of tokens that can be generated for a response,
	// including visible output tokens and
	// [reasoning tokens](https://platform.openai.com/docs/guides/reasoning).
	MaxOutputTokens int64 `json:"max_output_tokens,omitempty"`

	// The maximum number of total calls to built-in tools that can be processed in a
	// response. This maximum number applies across all built-in tool calls, not per
	// individual tool. Any further attempts to call a tool by the model will be
	// ignored.
	MaxToolCalls int64 `json:"max_tool_calls,omitempty"`

	// Whether to allow the model to run tool calls in parallel.
	ParallelToolCalls bool `json:"parallel_tool_calls"`

	// The unique ID of the previous response to the model. Use this to create
	// multi-turn conversations. Learn more about
	// [conversation state](https://platform.openai.com/docs/guides/conversation-state).
	// Cannot be used in conjunction with `conversation`.
	PreviousResponseID string `json:"previous_response_id,omitempty"`

	// Whether to store the generated model response for later retrieval via API.
	Store bool `json:"store"`

	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it more
	// focused and deterministic. We generally recommend altering this or `top_p` but
	// not both.
	Temperature float64 `json:"temperature,omitempty"`

	// An integer between 0 and 20 specifying the number of most likely tokens to
	// return at each token position, each with an associated log probability.
	TopLogprobs int64 `json:"top_logprobs,omitempty"`

	// An alternative to sampling with temperature, called nucleus sampling, where the
	// model considers the results of the tokens with top_p probability mass. So 0.1
	// means only the tokens comprising the top 10% probability mass are considered.
	//
	// We generally recommend altering this or `temperature` but not both.
	TopP float64 `json:"top_p,omitempty"`

	// Used by OpenAI to cache responses for similar requests to optimize your cache
	// hit rates. Replaces the `user` field.
	// [Learn more](https://platform.openai.com/docs/guides/prompt-caching).
	PromptCacheKey string `json:"prompt_cache_key,omitempty"`

	// A stable identifier used to help detect users of your application that may be
	// violating OpenAI's usage policies. The IDs should be a string that uniquely
	// identifies each user, with a maximum length of 64 characters. We recommend
	// hashing their username or email address, in order to avoid sending us any
	// identifying information.
	// [Learn more](https://platform.openai.com/docs/guides/safety-best-practices#safety-identifiers).
	SafetyIdentifier string `json:"safety_identifier,omitempty"`

	// This field is being replaced by `safety_identifier` and `prompt_cache_key`. Use
	// `prompt_cache_key` instead to maintain caching optimizations. A stable
	// identifier for your end-users. Used to boost cache hit rates by better bucketing
	// similar requests and to help OpenAI detect and prevent abuse.
	// [Learn more](https://platform.openai.com/docs/guides/safety-best-practices#safety-identifiers).
	User string `json:"user,omitempty"`

	// Context management configuration for this request.
	ContextManagement []ResponseContextManagement `json:"context_management,omitempty"`

	// The conversation that this response belongs to. Items from this conversation are
	// prepended to `input_items` for this response request. Input items and output
	// items from this response are automatically added to this conversation after this
	// response completes.
	Conversation any `json:"conversation,omitempty"`

	// Specify additional output data to include in the model response. Currently
	// supported values are:
	//
	//   - `web_search_call.action.sources`: Include the sources of the web search tool
	//     call.
	//   - `code_interpreter_call.outputs`: Includes the outputs of python code execution
	//     in code interpreter tool call items.
	//   - `computer_call_output.output.image_url`: Include image urls from the computer
	//     call output.
	//   - `file_search_call.results`: Include the search results of the file search tool
	//     call.
	//   - `message.input_image.image_url`: Include image urls from the input message.
	//   - `message.output_text.logprobs`: Include logprobs with assistant messages.
	//   - `reasoning.encrypted_content`: Includes an encrypted version of reasoning
	//     tokens in reasoning item outputs. This enables reasoning items to be used in
	//     multi-turn conversations when using the Responses API statelessly (like when
	//     the `store` parameter is set to `false`, or when an organization is enrolled
	//     in the zero data retention program).
	Include []string `json:"include,omitempty"`

	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata Metadata `json:"metadata,omitempty"`

	// Reference to a prompt template and its variables.
	// [Learn more](https://platform.openai.com/docs/guides/text?api-mode=responses#reusable-prompts).
	Prompt ResponsePrompt `json:"prompt,omitzero"`

	// The retention policy for the prompt cache. Set to `24h` to enable extended
	// prompt caching, which keeps cached prefixes active for longer, up to a maximum
	// of 24 hours.
	// [Learn more](https://platform.openai.com/docs/guides/prompt-caching#prompt-cache-retention).
	//
	// Any of "in-memory", "24h".
	PromptCacheRetention ResponsePromptCacheRetention `json:"prompt_cache_retention,omitempty"`

	// Specifies the processing type used for serving the request.
	//
	//   - If set to 'auto', then the request will be processed with the service tier
	//     configured in the Project settings. Unless otherwise configured, the Project
	//     will use 'default'.
	//   - If set to 'default', then the request will be processed with the standard
	//     pricing and performance for the selected model.
	//   - If set to '[flex](https://platform.openai.com/docs/guides/flex-processing)' or
	//     '[priority](https://openai.com/api-priority-processing/)', then the request
	//     will be processed with the corresponding service tier.
	//   - When not set, the default behavior is 'auto'.
	//
	// When the `service_tier` parameter is set, the response body will include the
	// `service_tier` value based on the processing mode actually used to serve the
	// request. This response value may be different from the value set in the
	// parameter.
	//
	// Any of "auto", "default", "flex", "scale", "priority".
	ServiceTier ResponseServiceTier `json:"service_tier,omitempty"`

	// Options for streaming responses.
	Stream bool `json:"stream,omitempty"`

	// Options for streaming responses. Only set this when you set `stream: true`.
	StreamOptions ResponseStreamOptions `json:"stream_options,omitzero"`

	// The truncation strategy to use for the model response.
	//
	//   - `auto`: If the input to this Response exceeds the model's context window size,
	//     the model will truncate the response to fit the context window by dropping
	//     items from the beginning of the conversation.
	//   - `disabled` (default): If the input size will exceed the context window size
	//     for a model, the request will fail with a 400 error.
	//
	// Any of "auto", "disabled".
	Truncation string `json:"truncation,omitempty"`

	// Text, image, or file inputs to the model, used to generate a response.
	//
	// Learn more:
	//
	// - [Text inputs and outputs](https://platform.openai.com/docs/guides/text)
	// - [Image inputs](https://platform.openai.com/docs/guides/images)
	// - [File inputs](https://platform.openai.com/docs/guides/pdf-files)
	// - [Conversation state](https://platform.openai.com/docs/guides/conversation-state)
	// - [Function calling](https://platform.openai.com/docs/guides/function-calling)
	Input []ResponseMessage `json:"input,omitempty"`

	// Model ID used to generate the response, like `gpt-4o` or `o3`. OpenAI offers a
	// wide range of models with different capabilities, performance characteristics,
	// and price points. Refer to the
	// [model guide](https://platform.openai.com/docs/models) to browse and compare
	// available models.
	Model string `json:"model,omitempty"`

	// **gpt-5 and o-series models only**
	//
	// Configuration options for
	// [reasoning models](https://platform.openai.com/docs/guides/reasoning).
	Reasoning Reasoning `json:"reasoning,omitzero"`

	// Configuration options for a text response from the model. Can be plain text or
	// structured JSON data. Learn more:
	//
	// - [Text inputs and outputs](https://platform.openai.com/docs/guides/text)
	// - [Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs)
	Text ResponseTextConfigParam `json:"text,omitzero"`

	// How the model should select which tool (or tools) to use when generating a
	// response. See the `tools` parameter to see how to specify which tools the model
	// can call.
	ToolChoice any `json:"tool_choice,omitempty"`

	// An array of tools the model may call while generating a response. You can
	// specify which tool to use by setting the `tool_choice` parameter.
	//
	// We support the following categories of tools:
	//
	//   - **Built-in tools**: Tools that are provided by OpenAI that extend the model's
	//     capabilities, like
	//     [web search](https://platform.openai.com/docs/guides/tools-web-search) or
	//     [file search](https://platform.openai.com/docs/guides/tools-file-search).
	//     Learn more about
	//     [built-in tools](https://platform.openai.com/docs/guides/tools).
	//   - **MCP Tools**: Integrations with third-party systems via custom MCP servers or
	//     predefined connectors such as Google Drive and SharePoint. Learn more about
	//     [MCP Tools](https://platform.openai.com/docs/guides/tools-connectors-mcp).
	//   - **Function calls (custom tools)**: Functions that are defined by you, enabling
	//     the model to call your own code with strongly typed arguments and outputs.
	//     Learn more about
	//     [function calling](https://platform.openai.com/docs/guides/function-calling).
	//     You can also use custom tools to call your own code.
	Tools []any `json:"tools,omitempty"`
}

func (r *CreateRequest) String() string {
	return toString(r)
}

// The phase of an assistant message.
//
// Use `commentary` for an intermediate assistant message and `final_answer` for
// the final assistant message. For follow-up requests with models like
// `gpt-5.3-codex` and later, preserve and resend phase on all assistant messages.
// Omitting it can degrade performance. Not used for user messages.
type ResponseOutputMessagePhase string

const (
	ResponseOutputMessagePhaseCommentary  ResponseOutputMessagePhase = "commentary"
	ResponseOutputMessagePhaseFinalAnswer ResponseOutputMessagePhase = "final_answer"
)

// Details about why the response is incomplete.
type ResponseIncompleteDetails struct {
	// The reason why the response is incomplete.
	//
	// Any of "max_output_tokens", "content_filter".
	Reason string `json:"reason"`
}

// The top log probability of a token.
type ResponseOutputTextLogprobTopLogprob struct {
	Token   string  `json:"token"`
	Bytes   []int64 `json:"bytes"`
	Logprob float64 `json:"logprob"`
}

// The log probability of a token.
type ResponseOutputTextLogprob struct {
	Token       string                                `json:"token"`
	Bytes       []int64                               `json:"bytes"`
	Logprob     float64                               `json:"logprob"`
	TopLogprobs []ResponseOutputTextLogprobTopLogprob `json:"top_logprobs"`
}

type ResponseOutputMessageContent struct {
	// This field is from variant [ResponseOutputText].
	Annotations []any `json:"annotations"`

	// This field is from variant [ResponseOutputText].
	Text string `json:"text"`

	// Any of "output_text", "refusal".
	Type string `json:"type"`

	// This field is from variant [ResponseOutputText].
	Logprobs []ResponseOutputTextLogprob `json:"logprobs"`

	// This field is from variant [ResponseOutputRefusal].
	Refusal string `json:"refusal"`
}

// A source used in the search.
type ResponseFunctionWebSearchActionSearchSource struct {
	// The type of source. Always `url`.
	Type string `json:"type"`
	// The URL of the source.
	URL string `json:"url"`
}

// An x/y coordinate pair, e.g. `{ x: 100, y: 200 }`.
type ResponseComputerToolCallActionDragPath struct {
	// The x-coordinate.
	X int64 `json:"x"`
	// The y-coordinate.
	Y int64 `json:"y"`
}

// ResponseOutputItemUnionAction is an implicit subunion of
// [ResponseOutputItemUnion]. ResponseOutputItemUnionAction provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ResponseOutputItemUnion].
type ResponseOutputItemUnionAction struct {
	// This field is from variant [ResponseFunctionWebSearchActionUnion].
	Query string `json:"query"`
	Type  string `json:"type"`
	// This field is from variant [ResponseFunctionWebSearchActionUnion].
	Queries []string `json:"queries"`
	// This field is from variant [ResponseFunctionWebSearchActionUnion].
	Sources []ResponseFunctionWebSearchActionSearchSource `json:"sources"`
	URL     string                                        `json:"url"`
	// This field is from variant [ResponseFunctionWebSearchActionUnion].
	Pattern string `json:"pattern"`
	// This field is from variant [ResponseComputerToolCallActionUnion].
	Button string `json:"button"`
	X      int64  `json:"x"`
	Y      int64  `json:"y"`
	// This field is from variant [ResponseComputerToolCallActionUnion].
	Path []ResponseComputerToolCallActionDragPath `json:"path"`
	// This field is from variant [ResponseComputerToolCallActionUnion].
	Keys []string `json:"keys"`
	// This field is from variant [ResponseComputerToolCallActionUnion].
	ScrollX int64 `json:"scroll_x"`
	// This field is from variant [ResponseComputerToolCallActionUnion].
	ScrollY int64 `json:"scroll_y"`
	// This field is from variant [ResponseComputerToolCallActionUnion].
	Text string `json:"text"`
	// This field is from variant [ResponseOutputItemLocalShellCallAction].
	Command []string `json:"command"`
	// This field is from variant [ResponseOutputItemLocalShellCallAction].
	Env       map[string]string `json:"env"`
	TimeoutMs int64             `json:"timeout_ms"`
	// This field is from variant [ResponseOutputItemLocalShellCallAction].
	User string `json:"user"`
	// This field is from variant [ResponseOutputItemLocalShellCallAction].
	WorkingDirectory string `json:"working_directory"`
	// This field is from variant [ResponseFunctionShellToolCallAction].
	Commands []string `json:"commands"`
	// This field is from variant [ResponseFunctionShellToolCallAction].
	MaxOutputLength int64 `json:"max_output_length"`
}

// A pending safety check for the computer call.
type ResponseComputerToolCallPendingSafetyCheck struct {
	// The ID of the pending safety check.
	ID string `json:"id"`
	// The type of the pending safety check.
	Code string `json:"code" api:"nullable"`
	// Details about the pending safety check.
	Message string `json:"message" api:"nullable"`
}

// A summary text from the model.
type ResponseReasoningItemSummary struct {
	// A summary of the reasoning output from the model so far.
	Text string `json:"text"`
	// The type of the object. Always `summary_text`.
	Type string `json:"type"`
}

// ResponseCodeInterpreterToolCallOutputUnion contains all possible properties and
// values from [ResponseCodeInterpreterToolCallOutputLogs],
// [ResponseCodeInterpreterToolCallOutputImage].
//
// Use the [ResponseCodeInterpreterToolCallOutputUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ResponseCodeInterpreterToolCallOutputUnion struct {
	// This field is from variant [ResponseCodeInterpreterToolCallOutputLogs].
	Logs string `json:"logs"`
	// Any of "logs", "image".
	Type string `json:"type"`
	// This field is from variant [ResponseCodeInterpreterToolCallOutputImage].
	URL string `json:"url"`
}

// ResponseFunctionShellToolCallEnvironmentUnion contains all possible properties
// and values from [ResponseLocalEnvironment], [ResponseContainerReference].
//
// Use the [ResponseFunctionShellToolCallEnvironmentUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ResponseFunctionShellToolCallEnvironmentUnion struct {
	// Any of "local", "container_reference".
	Type string `json:"type"`
	// This field is from variant [ResponseContainerReference].
	ContainerID string `json:"container_id"`
}

// ResponseApplyPatchToolCallOperationUnion contains all possible properties and
// values from [ResponseApplyPatchToolCallOperationCreateFile],
// [ResponseApplyPatchToolCallOperationDeleteFile],
// [ResponseApplyPatchToolCallOperationUpdateFile].
//
// Use the [ResponseApplyPatchToolCallOperationUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ResponseApplyPatchToolCallOperationUnion struct {
	Diff string `json:"diff"`
	Path string `json:"path"`
	// Any of "create_file", "delete_file", "update_file".
	Type string `json:"type"`
}

// A tool available on an MCP server.
type ResponseOutputItemMcpListToolsTool struct {
	// The JSON schema describing the tool's input.
	InputSchema any `json:"input_schema"`
	// The name of the tool.
	Name string `json:"name"`
	// Additional annotations about the tool.
	Annotations any `json:"annotations" api:"nullable"`
	// The description of the tool.
	Description string `json:"description" api:"nullable"`
}

// ResponseOutputItemUnion contains all possible properties and values from
// [ResponseOutputMessage], [ResponseFileSearchToolCall],
// [ResponseFunctionToolCall], [ResponseFunctionWebSearch],
// [ResponseComputerToolCall], [ResponseReasoningItem], [ResponseCompactionItem],
// [ResponseOutputItemImageGenerationCall], [ResponseCodeInterpreterToolCall],
// [ResponseOutputItemLocalShellCall], [ResponseFunctionShellToolCall],
// [ResponseFunctionShellToolCallOutput], [ResponseApplyPatchToolCall],
// [ResponseApplyPatchToolCallOutput], [ResponseOutputItemMcpCall],
// [ResponseOutputItemMcpListTools], [ResponseOutputItemMcpApprovalRequest],
// [ResponseCustomToolCall].
//
// Use the [ResponseOutputItemUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ResponseOutputItemUnion struct {
	ID      string                         `json:"id"`
	Content []ResponseOutputMessageContent `json:"content"`

	// This field is from variant [ResponseOutputMessage].
	Role   string `json:"role"`
	Status string `json:"status"`

	// Any of "message", "file_search_call", "function_call", "web_search_call",
	// "computer_call", "reasoning", "compaction", "image_generation_call",
	// "code_interpreter_call", "local_shell_call", "shell_call", "shell_call_output",
	// "apply_patch_call", "apply_patch_call_output", "mcp_call", "mcp_list_tools",
	// "mcp_approval_request", "custom_tool_call".
	Type string `json:"type"`

	// This field is from variant [ResponseOutputMessage].
	Phase ResponseOutputMessagePhase `json:"phase"`

	// This field is from variant [ResponseFileSearchToolCall].
	Queries []string `json:"queries"`

	// This field is from variant [ResponseFileSearchToolCall].
	Results []ResponseFileSearchToolCallResult `json:"results"`

	Arguments string `json:"arguments"`
	CallID    string `json:"call_id"`
	Name      string `json:"name"`

	// This field is a union of [ResponseFunctionWebSearchActionUnion],
	// [ResponseComputerToolCallActionUnion], [ResponseOutputItemLocalShellCallAction],
	// [ResponseFunctionShellToolCallAction]
	Action ResponseOutputItemUnionAction `json:"action"`

	// This field is from variant [ResponseComputerToolCall].
	PendingSafetyChecks []ResponseComputerToolCallPendingSafetyCheck `json:"pending_safety_checks"`

	// This field is from variant [ResponseReasoningItem].
	Summary []ResponseReasoningItemSummary `json:"summary"`

	EncryptedContent string `json:"encrypted_content"`
	CreatedBy        string `json:"created_by"`

	// This field is from variant [ResponseOutputItemImageGenerationCall].
	Result string `json:"result"`

	// This field is from variant [ResponseCodeInterpreterToolCall].
	Code string `json:"code"`

	// This field is from variant [ResponseCodeInterpreterToolCall].
	ContainerID string `json:"container_id"`

	// This field is from variant [ResponseCodeInterpreterToolCall].
	Outputs []ResponseCodeInterpreterToolCallOutputUnion `json:"outputs"`

	// This field is from variant [ResponseFunctionShellToolCall].
	Environment ResponseFunctionShellToolCallEnvironmentUnion `json:"environment"`

	// This field is from variant [ResponseFunctionShellToolCallOutput].
	MaxOutputLength int64 `json:"max_output_length"`

	// This field is a union of [[]ResponseFunctionShellToolCallOutputOutput],
	// [string], [string]
	Output any `json:"output"`

	// This field is from variant [ResponseApplyPatchToolCall].
	Operation ResponseApplyPatchToolCallOperationUnion `json:"operation"`

	ServerLabel string `json:"server_label"`

	// This field is from variant [ResponseOutputItemMcpCall].
	ApprovalRequestID string `json:"approval_request_id"`

	Error string `json:"error"`

	// This field is from variant [ResponseOutputItemMcpListTools].
	Tools []ResponseOutputItemMcpListToolsTool `json:"tools"`

	// This field is from variant [ResponseCustomToolCall].
	Input string `json:"input"`
}

type ResponseFileSearchToolCallResult struct {
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard. Keys are strings with a maximum
	// length of 64 characters. Values are strings with a maximum length of 512
	// characters, booleans, or numbers.
	Attributes map[string]any `json:"attributes" api:"nullable"`

	// The unique ID of the file.
	FileID string `json:"file_id"`

	// The name of the file.
	Filename string `json:"filename"`

	// The relevance score of the file - a value between 0 and 1.
	Score float64 `json:"score"`

	// The text that was retrieved from the file.
	Text string `json:"text"`
}

// The conversation that this response belonged to. Input items and output items
// from this response were automatically added to this conversation.
type ResponseConversation struct {
	// The unique ID of the conversation that this response was associated with.
	ID string `json:"id"`
}

// The retention policy for the prompt cache. Set to `24h` to enable extended
// prompt caching, which keeps cached prefixes active for longer, up to a maximum
// of 24 hours.
// [Learn more](https://platform.openai.com/docs/guides/prompt-caching#prompt-cache-retention).
type ResponsePromptCacheRetention string

const (
	ResponsePromptCacheRetentionInMemory ResponsePromptCacheRetention = "in-memory"
	ResponsePromptCacheRetention24h      ResponsePromptCacheRetention = "24h"
)

// Configuration options for a text response from the model. Can be plain text or
// structured JSON data. Learn more:
//
// - [Text inputs and outputs](https://platform.openai.com/docs/guides/text)
// - [Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs)
type ResponseTextConfig struct {
	// An object specifying the format that the model must output.
	//
	// Configuring `{ "type": "json_schema" }` enables Structured Outputs, which
	// ensures the model will match your supplied JSON schema. Learn more in the
	// [Structured Outputs guide](https://platform.openai.com/docs/guides/structured-outputs).
	//
	// The default format is `{ "type": "text" }` with no additional options.
	//
	// **Not recommended for gpt-4o and newer models:**
	//
	// Setting to `{ "type": "json_object" }` enables the older JSON mode, which
	// ensures the message the model generates is valid JSON. Using `json_schema` is
	// preferred for models that support it.
	Format *ResponseTextConfigFormat `json:"format"`

	// Constrains the verbosity of the model's response. Lower values will result in
	// more concise responses, while higher values will result in more verbose
	// responses. Currently supported values are `low`, `medium`, and `high`.
	//
	// Any of "low", "medium", "high".
	Verbosity ResponseTextConfigVerbosity `json:"verbosity,omitempty"`
}

// A detailed breakdown of the input tokens.
type ResponseUsageInputTokensDetails struct {
	// The number of tokens that were retrieved from the cache.
	// [More on prompt caching](https://platform.openai.com/docs/guides/prompt-caching).
	CachedTokens int64 `json:"cached_tokens"`
}

// Represents token usage details including input tokens, output tokens, a
// breakdown of output tokens, and the total tokens used.
type ResponseUsage struct {
	// The number of input tokens.
	InputTokens int64 `json:"input_tokens"`

	// A detailed breakdown of the input tokens.
	InputTokensDetails ResponseUsageInputTokensDetails `json:"input_tokens_details"`

	// The number of output tokens.
	OutputTokens int64 `json:"output_tokens"`

	// A detailed breakdown of the output tokens.
	OutputTokensDetails ResponseUsageOutputTokensDetails `json:"output_tokens_details"`

	// The total number of tokens used.
	TotalTokens int64 `json:"total_tokens"`
}

// A detailed breakdown of the output tokens.
type ResponseUsageOutputTokensDetails struct {
	// The number of reasoning tokens.
	ReasoningTokens int64 `json:"reasoning_tokens"`
}

type CreateResponse struct {
	// Unique identifier for this Response.
	ID string `json:"id"`

	// Unix timestamp (in seconds) of when this Response was created.
	CreatedAt float64 `json:"created_at"`

	// Details about why the response is incomplete.
	IncompleteDetails *ResponseIncompleteDetails `json:"incomplete_details,omitempty"`

	// A system (or developer) message inserted into the model's context.
	//
	// When using along with `previous_response_id`, the instructions from a previous
	// response will not be carried over to the next response. This makes it simple to
	// swap out system (or developer) messages in new responses.
	Instructions any `json:"instructions,omitempty"`

	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata Metadata `json:"metadata,omitempty"`

	// Model ID used to generate the response, like `gpt-4o` or `o3`. OpenAI offers a
	// wide range of models with different capabilities, performance characteristics,
	// and price points. Refer to the
	// [model guide](https://platform.openai.com/docs/models) to browse and compare
	// available models.
	Model string `json:"model"`

	// The object type of this resource - always set to `response`.
	Object string `json:"object"`

	// An array of content items generated by the model.
	//
	//   - The length and order of items in the `output` array is dependent on the
	//     model's response.
	//   - Rather than accessing the first item in the `output` array and assuming it's
	//     an `assistant` message with the content generated by the model, you might
	//     consider using the `output_text` property where supported in SDKs.
	Output []ResponseOutputItemUnion `json:"output"`

	// Whether to allow the model to run tool calls in parallel.
	ParallelToolCalls bool `json:"parallel_tool_calls"`

	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it more
	// focused and deterministic. We generally recommend altering this or `top_p` but
	// not both.
	Temperature float64 `json:"temperature"`

	// How the model should select which tool (or tools) to use when generating a
	// response. See the `tools` parameter to see how to specify which tools the model
	// can call.
	ToolChoice any `json:"tool_choice,omitempty"`

	// An array of tools the model may call while generating a response. You can
	// specify which tool to use by setting the `tool_choice` parameter.
	//
	// We support the following categories of tools:
	//
	//   - **Built-in tools**: Tools that are provided by OpenAI that extend the model's
	//     capabilities, like
	//     [web search](https://platform.openai.com/docs/guides/tools-web-search) or
	//     [file search](https://platform.openai.com/docs/guides/tools-file-search).
	//     Learn more about
	//     [built-in tools](https://platform.openai.com/docs/guides/tools).
	//   - **MCP Tools**: Integrations with third-party systems via custom MCP servers or
	//     predefined connectors such as Google Drive and SharePoint. Learn more about
	//     [MCP Tools](https://platform.openai.com/docs/guides/tools-connectors-mcp).
	//   - **Function calls (custom tools)**: Functions that are defined by you, enabling
	//     the model to call your own code with strongly typed arguments and outputs.
	//     Learn more about
	//     [function calling](https://platform.openai.com/docs/guides/function-calling).
	//     You can also use custom tools to call your own code.
	Tools []any `json:"tools,omitempty"`

	// An alternative to sampling with temperature, called nucleus sampling, where the
	// model considers the results of the tokens with top_p probability mass. So 0.1
	// means only the tokens comprising the top 10% probability mass are considered.
	//
	// We generally recommend altering this or `temperature` but not both.
	TopP float64 `json:"top_p,omitempty"`

	// Whether to run the model response in the background.
	// [Learn more](https://platform.openai.com/docs/guides/background).
	Background bool `json:"background,omitempty"`

	// Unix timestamp (in seconds) of when this Response was completed. Only present
	// when the status is `completed`.
	CompletedAt float64 `json:"completed_at,omitempty"`

	// The conversation that this response belonged to. Input items and output items
	// from this response were automatically added to this conversation.
	Conversation ResponseConversation `json:"conversation,omitzero"`

	// An upper bound for the number of tokens that can be generated for a response,
	// including visible output tokens and
	// [reasoning tokens](https://platform.openai.com/docs/guides/reasoning).
	MaxOutputTokens int64 `json:"max_output_tokens,omitempty"`

	// The maximum number of total calls to built-in tools that can be processed in a
	// response. This maximum number applies across all built-in tool calls, not per
	// individual tool. Any further attempts to call a tool by the model will be
	// ignored.
	MaxToolCalls int64 `json:"max_tool_calls,omitempty"`

	// The unique ID of the previous response to the model. Use this to create
	// multi-turn conversations. Learn more about
	// [conversation state](https://platform.openai.com/docs/guides/conversation-state).
	// Cannot be used in conjunction with `conversation`.
	PreviousResponseID string `json:"previous_response_id,omitempty"`

	// Reference to a prompt template and its variables.
	// [Learn more](https://platform.openai.com/docs/guides/text?api-mode=responses#reusable-prompts).
	Prompt *ResponsePrompt `json:"prompt,omitempty"`

	// Used by OpenAI to cache responses for similar requests to optimize your cache
	// hit rates. Replaces the `user` field.
	// [Learn more](https://platform.openai.com/docs/guides/prompt-caching).
	PromptCacheKey string `json:"prompt_cache_key,omitempty"`

	// The retention policy for the prompt cache. Set to `24h` to enable extended
	// prompt caching, which keeps cached prefixes active for longer, up to a maximum
	// of 24 hours.
	// [Learn more](https://platform.openai.com/docs/guides/prompt-caching#prompt-cache-retention).
	//
	// Any of "in-memory", "24h".
	PromptCacheRetention ResponsePromptCacheRetention `json:"prompt_cache_retention,omitempty"`

	// **gpt-5 and o-series models only**
	//
	// Configuration options for
	// [reasoning models](https://platform.openai.com/docs/guides/reasoning).
	Reasoning *Reasoning `json:"reasoning,omitempty"`

	// A stable identifier used to help detect users of your application that may be
	// violating OpenAI's usage policies. The IDs should be a string that uniquely
	// identifies each user, with a maximum length of 64 characters. We recommend
	// hashing their username or email address, in order to avoid sending us any
	// identifying information.
	// [Learn more](https://platform.openai.com/docs/guides/safety-best-practices#safety-identifiers).
	SafetyIdentifier string `json:"safety_identifier,omitempty"`

	// Specifies the processing type used for serving the request.
	//
	//   - If set to 'auto', then the request will be processed with the service tier
	//     configured in the Project settings. Unless otherwise configured, the Project
	//     will use 'default'.
	//   - If set to 'default', then the request will be processed with the standard
	//     pricing and performance for the selected model.
	//   - If set to '[flex](https://platform.openai.com/docs/guides/flex-processing)' or
	//     '[priority](https://openai.com/api-priority-processing/)', then the request
	//     will be processed with the corresponding service tier.
	//   - When not set, the default behavior is 'auto'.
	//
	// When the `service_tier` parameter is set, the response body will include the
	// `service_tier` value based on the processing mode actually used to serve the
	// request. This response value may be different from the value set in the
	// parameter.
	//
	// Any of "auto", "default", "flex", "scale", "priority".
	ServiceTier ResponseServiceTier `json:"service_tier,omitempty"`

	// The status of the response generation. One of `completed`, `failed`,
	// `in_progress`, `cancelled`, `queued`, or `incomplete`.
	//
	// Any of "completed", "failed", "in_progress", "cancelled", "queued",
	// "incomplete".
	Status ResponseStatus `json:"status"`

	// Configuration options for a text response from the model. Can be plain text or
	// structured JSON data. Learn more:
	//
	// - [Text inputs and outputs](https://platform.openai.com/docs/guides/text)
	// - [Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs)
	Text ResponseTextConfig `json:"text"`

	// An integer between 0 and 20 specifying the number of most likely tokens to
	// return at each token position, each with an associated log probability.
	TopLogprobs int64 `json:"top_logprobs,omitempty"`

	// The truncation strategy to use for the model response.
	//
	//   - `auto`: If the input to this Response exceeds the model's context window size,
	//     the model will truncate the response to fit the context window by dropping
	//     items from the beginning of the conversation.
	//   - `disabled` (default): If the input size will exceed the context window size
	//     for a model, the request will fail with a 400 error.
	//
	// Any of "auto", "disabled".
	Truncation string `json:"truncation,omitempty"`

	// Represents token usage details including input tokens, output tokens, a
	// breakdown of output tokens, and the total tokens used.
	Usage ResponseUsage `json:"usage"`

	// This field is being replaced by `safety_identifier` and `prompt_cache_key`. Use
	// `prompt_cache_key` instead to maintain caching optimizations. A stable
	// identifier for your end-users. Used to boost cache hit rates by better bucketing
	// similar requests and to help OpenAI detect and prevent abuse.
	// [Learn more](https://platform.openai.com/docs/guides/safety-best-practices#safety-identifiers).
	//
	// Deprecated: deprecated
	User string `json:"user,omitempty"`
}

func (r *CreateResponse) String() string {
	return toString(r)
}

func (r *CreateResponse) OutputText() string {
	var outputText strings.Builder
	for _, item := range r.Output {
		for _, content := range item.Content {
			if content.Type == "output_text" {
				outputText.WriteString(content.Text)
			}
		}
	}
	return outputText.String()
}
