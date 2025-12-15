package cmdutils

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	clientSDK "github.com/MagaluCloud/mgc-sdk-go/client"
)

type CliError struct {
	Message string
	Details string
}

func (e *CliError) Error() string {
	if e == nil {
		return "nil CLI error"
	}

	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

func NewCliError(message string) *CliError {
	return &CliError{
		Message: message,
	}
}

func NewCliErrorWithDetails(message, details string) *CliError {
	return &CliError{
		Message: message,
		Details: details,
	}
}

const (
	simpleHttpError       = "API request failed with HTTP error"
	simpleValidationError = "Request validation failed"
	simpleGenericError    = "An unexpected error occurred"
	simpleMaxRetriesError = "Max HTTP retries exceeded"

	MgcTraceIDKey = "x-mgc-trace-id"
)

type HttpErrorResponse struct {
	Status     string
	Body       string
	URL        string
	RequestID  string
	MgcTraceID string
	Payload    []byte
	Message    string // MGC reports this in the json body
	Slug       string // MGC reports this in the json body

}

type IdentifiableHttpError struct {
	*HttpErrorResponse
	RequestID string `json:"requestID"`
}

type BaseApiError struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

func (e HttpErrorResponse) String() string {
	return fmt.Sprintf("\nHTTP Error:\n  Status: %s\n  Body: %s\n  URL: %s\n  Request ID: %s\n  MGC Trace ID: %s\n  Message: %s\n  Slug: %s\n  Payload: %s",
		e.Status,
		e.Body,
		e.URL,
		e.RequestID,
		e.MgcTraceID,
		e.Message,
		e.Slug,
		e.Payload,
	)
}

func (e HttpErrorResponse) Error() string {
	return e.String()
}

func buildFromSDKError(err *clientSDK.HTTPError) (HttpErrorResponse, error) {
	if err == nil {
		return HttpErrorResponse{}, fmt.Errorf("cannot build error response from nil error")
	}

	e := HttpErrorResponse{
		Status: err.Status,
		Body:   string(err.Body),
	}

	if err.Response != nil && err.Response.Request != nil {
		e.URL = err.Response.Request.URL.String()
		for key, value := range err.Response.Header {
			if strings.ToLower(key) == string(clientSDK.RequestIDKey) {
				e.RequestID = value[0]
			}
			if strings.ToLower(key) == string(MgcTraceIDKey) {
				e.MgcTraceID = value[0]
			}
		}
	}

	return e, nil
}

func ParseSDKError(err error) (msg, detail string) {
	if err == nil {
		return simpleGenericError, "nil error provided"
	}

	switch e := err.(type) {
	case *CliError:

		if e.Details != "" {
			detail = fmt.Sprintf("%s\n%s", detail, e.Details)
		}
		return e.Message, detail

	case *clientSDK.HTTPError:
		errorResponse, buildErr := buildFromSDKError(e)
		if buildErr != nil {
			return simpleGenericError, buildErr.Error()
		}
		return simpleHttpError, errorResponse.String()

	case *clientSDK.ValidationError:
		if e == nil {
			return simpleValidationError, "nil validation error"
		}
		return simpleValidationError, fmt.Sprintf("Field: %s - Message: %s", e.Field, e.Message)

	case *clientSDK.RetryError:
		if e == nil {
			return simpleMaxRetriesError, ""
		}
		if e.LastError == nil {
			return simpleMaxRetriesError, "unexpected last retry error"
		}
		if he, ok := e.LastError.(*clientSDK.HTTPError); ok {
			errorResponse, buildErr := buildFromSDKError(he)
			if buildErr != nil {
				return simpleMaxRetriesError, buildErr.Error()
			}
			return simpleMaxRetriesError, fmt.Sprintf("Max HTTP retries exceeded at %d retries.\nLast error:\n %s", e.Retries, errorResponse.String())
		}
		return simpleMaxRetriesError, fmt.Sprintf("Max HTTP retries exceeded at %d retries.\nLast error: %s", e.Retries, e.LastError.Error())

	default:
		return simpleGenericError, err.Error()
	}
}

func NewHttpErrorFromResponse(resp *http.Response, req *http.Request) *IdentifiableHttpError {
	slug := "unknown"
	message := resp.Status

	defer resp.Body.Close()
	payload, _ := io.ReadAll(resp.Body)

	contentType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		fmt.Println(
			"ignored invalid response",
			"Content-Type:", resp.Header.Get("Content-Type"),
			"error:", err.Error(),
		)
	}
	if contentType == "application/json" {
		data := BaseApiError{}
		if err := json.Unmarshal(payload, &data); err == nil {
			if data.Message != "" {
				message = data.Message
			}
			if data.Slug != "" {
				slug = data.Slug
			}
		}
	}

	body, _ := io.ReadAll(resp.Body)

	httpError := &HttpErrorResponse{
		Body:    string(body),
		URL:     req.URL.String(),
		Status:  resp.Status,
		Payload: payload,
		Message: message,
		Slug:    slug,
	}

	return NewIdentifiableHttpError(httpError, req, resp)
}

func NewIdentifiableHttpError(httpError *HttpErrorResponse, request *http.Request, response *http.Response) *IdentifiableHttpError {
	a := IdentifiableHttpError{
		HttpErrorResponse: httpError,
	}
	if response != nil {
		if id := response.Header.Get("X-Request-Id"); id != "" {
			a.RequestID = id
		}
		if id := response.Header.Get("X-Mgc-Trace-Id"); id != "" {
			a.MgcTraceID = id
		}
	}

	return &a
}
