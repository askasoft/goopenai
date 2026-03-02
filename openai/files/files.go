package files

import (
	"bytes"
	"fmt"
	"io"

	"github.com/askasoft/pango/doc/jsonx"
	"github.com/askasoft/pango/net/httpx"
	"github.com/askasoft/pango/num"
)

func toString(o any) string {
	return jsonx.Prettify(o)
}

// The intended purpose of the file. Supported values are `assistants`,
// `assistants_output`, `batch`, `batch_output`, `fine-tune`, `fine-tune-results`,
// `vision`, and `user_data`.
type FileObjectPurpose string

const (
	FileObjectPurposeAssistants       FileObjectPurpose = "assistants"
	FileObjectPurposeAssistantsOutput FileObjectPurpose = "assistants_output"
	FileObjectPurposeBatch            FileObjectPurpose = "batch"
	FileObjectPurposeBatchOutput      FileObjectPurpose = "batch_output"
	FileObjectPurposeFineTune         FileObjectPurpose = "fine-tune"
	FileObjectPurposeFineTuneResults  FileObjectPurpose = "fine-tune-results"
	FileObjectPurposeVision           FileObjectPurpose = "vision"
	FileObjectPurposeUserData         FileObjectPurpose = "user_data"
)

// Deprecated. The current status of the file, which can be either `uploaded`,
// `processed`, or `error`.
type FileObjectStatus string

const (
	FileObjectStatusUploaded  FileObjectStatus = "uploaded"
	FileObjectStatusProcessed FileObjectStatus = "processed"
	FileObjectStatusError     FileObjectStatus = "error"
)

// The intended purpose of the uploaded file. One of:
//
// - `assistants`: Used in the Assistants API
// - `batch`: Used in the Batch API
// - `fine-tune`: Used for fine-tuning
// - `vision`: Images used for vision fine-tuning
// - `user_data`: Flexible file type for any purpose
// - `evals`: Used for eval data sets
type FilePurpose string

const (
	FilePurposeAssistants FilePurpose = "assistants"
	FilePurposeBatch      FilePurpose = "batch"
	FilePurposeFineTune   FilePurpose = "fine-tune"
	FilePurposeVision     FilePurpose = "vision"
	FilePurposeUserData   FilePurpose = "user_data"
	FilePurposeEvals      FilePurpose = "evals"
)

type CreateRequest struct {
	FileName string `json:"filename"`
	FileData any    `json:"-"` // string (filepath) or []byte (filedata) or io.Reader

	// The intended purpose of the uploaded file. One of:
	//
	// - `assistants`: Used in the Assistants API
	// - `batch`: Used in the Batch API
	// - `fine-tune`: Used for fine-tuning
	// - `vision`: Images used for vision fine-tuning
	// - `user_data`: Flexible file type for any purpose
	// - `evals`: Used for eval data sets
	//
	// Any of "assistants", "batch", "fine-tune", "vision", "user_data", "evals".
	Purpose FilePurpose `json:"purpose,omitempty"`

	// The number of seconds after the anchor time that the file will expire. Must be
	// between 3600 (1 hour) and 2592000 (30 days).
	// By default, files with `purpose=batch` expire
	// after 30 days and all other files are persisted until they are manually deleted.
	ExpiresAfter int `json:"expires_after,omitempty"`
}

func (cr CreateRequest) MarshalBody() (body io.Reader, contentType string, err error) {
	buf := &bytes.Buffer{}

	mw := httpx.NewMultipartWriter(buf)

	contentType = mw.FormDataContentType()

	if cr.Purpose != "" {
		if err = mw.WriteField("purpose", string(cr.Purpose)); err != nil {
			return
		}
	}

	if fd, ok := cr.FileData.([]byte); ok {
		if err = mw.WriteFileData("file", cr.FileName, fd); err != nil {
			return
		}
	} else if fr, ok := cr.FileData.(io.Reader); ok {
		if err = mw.WriteFileReader("file", cr.FileName, fr); err != nil {
			return
		}
	} else if fp, ok := cr.FileData.(string); ok {
		if err = mw.WriteFile("file", fp); err != nil {
			return
		}
	} else {
		err = fmt.Errorf("openai: invalid file data %T", cr.FileData)
		return
	}

	if cr.ExpiresAfter > 0 {
		if err = mw.WriteField("expires_after[anchor]", "created_at"); err != nil {
			return
		}
		if err = mw.WriteField("expires_after[seconds]", num.Itoa(cr.ExpiresAfter)); err != nil {
			return
		}
	}

	if err = mw.Close(); err != nil {
		return
	}

	body = buf
	return
}

// The `File` object represents a document that has been uploaded to OpenAI.
type FileObject struct {
	// The file identifier, which can be referenced in the API endpoints.
	ID string `json:"id"`

	// The name of the file.
	Filename string `json:"filename"`

	// The size of the file, in bytes.
	Bytes int64 `json:"bytes"`

	// The intended purpose of the file. Supported values are `assistants`,
	// `assistants_output`, `batch`, `batch_output`, `fine-tune`, `fine-tune-results`,
	// `vision`, and `user_data`.
	//
	// Any of "assistants", "assistants_output", "batch", "batch_output", "fine-tune",
	// "fine-tune-results", "vision", "user_data".
	Purpose FileObjectPurpose `json:"purpose"`

	// The Unix timestamp (in seconds) for when the file was created.
	CreatedAt int64 `json:"created_at"`

	// The Unix timestamp (in seconds) for when the file will expire.
	ExpiresAt int64 `json:"expires_at"`
}

func (fo *FileObject) String() string {
	return toString(fo)
}
