package apiduck

import "time"

type Method string

const (
	MethodGet     Method = "GET"
	MethodPost    Method = "POST"
	MethodPut     Method = "PUT"
	MethodDelete  Method = "DELETE"
	MethodPatch   Method = "PATCH"
	MethodHead    Method = "HEAD"
	MethodOptions Method = "OPTIONS"
)

type SecurityType string

const (
	SecurityTypeAPIKey SecurityType = "apiKey"
	SecurityTypeBearer SecurityType = "bearer"
	SecurityTypeBasic  SecurityType = "basic"
)

type SecurityLocation string

const (
	SecurityLocationHeader SecurityLocation = "header"
	SecurityLocationQuery  SecurityLocation = "query"
	SecurityLocationCookie SecurityLocation = "cookie"
)

type APIDocumentation struct {
	Info      Info       `json:"info"`
	Tags      []*Tag     `json:"tags"`
	Security  []Security `json:"security,omitempty"`
	Servers   []Server   `json:"servers,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

type Server struct {
	URL         string            `json:"url"`
	Description string            `json:"description"`
	Variables   map[string]string `json:"variables,omitempty"`
}

type Security struct {
	Type        SecurityType     `json:"type"`
	Name        string           `json:"name"`
	KeyName     string           `json:"keyName"`
	In          SecurityLocation `json:"in,omitempty"`
	Scheme      string           `json:"scheme,omitempty"`
	Description string           `json:"description,omitempty"`
}

type Tag struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Endpoints   []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Path        string       `json:"path"`
	Summary     string       `json:"summary"`
	Description string       `json:"description,omitempty"`
	Method      Method       `json:"method"`
	Query       []QueryParam `json:"queries,omitempty"`
	Headers     []Header     `json:"headers,omitempty"`
	Form        []FormField  `json:"formdata,omitempty"`
	Body        []BodyField  `json:"body,omitempty"`
	Responses   []Response   `json:"responses,omitempty"`
	Security    []string     `json:"security,omitempty"` // References to security schemes
	Tags        []string     `json:"tags,omitempty"`     // Can belong to multiple tags
	Deprecated  bool         `json:"deprecated"`
	OperationID string       `json:"operationId,omitempty"`
}

type Response struct {
	StatusCode  int         `json:"statusCode"`
	Description string      `json:"description"`
	Schema      []BodyField `json:"schema,omitempty"`
	Headers     []Header    `json:"headers,omitempty"`
	Examples    []Example   `json:"examples,omitempty"`
	ContentType string      `json:"contentType,omitempty"`
}

type Example struct {
	Name        string `json:"name"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
	Value       any    `json:"value"`
}

type Header struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Example     string `json:"example,omitempty"`
}

type BodyField struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	JSONName    string            `json:"jsonName,omitempty"`
	Description string            `json:"description,omitempty"`
	Required    bool              `json:"required"`
	Validation  map[string]string `json:"validation,omitempty"`
	Fields      []BodyField       `json:"fields,omitempty"`
	Example     any               `json:"example,omitempty"`
	Enum        []any             `json:"enum,omitempty"`
	Default     any               `json:"default,omitempty"`
}

type FormField struct {
	Name        string `json:"name"`
	Schema      any    `json:"schema,omitempty"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
	Example     any    `json:"example,omitempty"`
}

type QueryParam struct {
	Name        string          `json:"name"`
	Type        string          `json:"type"`
	Description string          `json:"description,omitempty"`
	Validation  ParamValidation `json:"validation"`
	Example     any             `json:"example,omitempty"`
	Enum        []any           `json:"enum,omitempty"`
}

type ParamValidation struct {
	Required bool   `json:"required"`
	Default  any    `json:"default,omitempty"`
	Min      int    `json:"min,omitempty"`
	Max      int    `json:"max,omitempty"`
	MinLen   int    `json:"minLength,omitempty"`
	MaxLen   int    `json:"maxLength,omitempty"`
	Pattern  string `json:"pattern,omitempty"`
}

type Info struct {
	Title          string      `json:"title"`
	Description    string      `json:"description,omitempty"`
	Contact        InfoContact `json:"contact"`
	License        InfoLicense `json:"license"`
	Version        string      `json:"version"`
	TermsOfService string      `json:"termsOfService,omitempty"`
}

type InfoContact struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}

type InfoLicense struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}
