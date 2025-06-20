package apiduck

import (
	"embed"
	"encoding/json"
	"html/template"
	"net/http"
	"time"
)

//go:embed templates/*.html
var templateFS embed.FS

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		panic("failed to parse api documentation templates: " + err.Error())
	}
}

// NewDocumentation creates a new API documentation instance
func NewDocumentation(info Info) *APIDocumentation {
	return &APIDocumentation{
		Info:      info,
		Tags:      []*Tag{},
		Security:  []Security{},
		Servers:   []Server{},
		CreatedAt: time.Now(),
	}
}

// AddServer adds a server to the documentation
func (docs *APIDocumentation) AddServer(url, description string, variables map[string]string) *APIDocumentation {
	server := Server{
		URL:         url,
		Description: description,
		Variables:   variables,
	}
	docs.Servers = append(docs.Servers, server)
	return docs
}

// AddSecurity adds a security scheme to the documentation
func (docs *APIDocumentation) AddSecurity(security Security) *APIDocumentation {
	docs.Security = append(docs.Security, security)
	return docs
}

// AddTag adds a new tag to the documentation
func (docs *APIDocumentation) AddTag(name, description string) *Tag {
	tag := &Tag{
		Name:        name,
		Description: description,
		Endpoints:   []Endpoint{},
	}
	docs.Tags = append(docs.Tags, tag)
	return tag
}

// GetTag retrieves an existing tag by name or creates a new one
func (docs *APIDocumentation) GetTag(name string) *Tag {
	for _, tag := range docs.Tags {
		if tag.Name == name {
			return tag
		}
	}
	return docs.AddTag(name, "")
}

// EndpointOption represents a function option for configuring endpoints
type EndpointOption func(e *Endpoint)

// AddEndpoint adds a new endpoint to the tag
func (tag *Tag) AddEndpoint(method Method, path, summary string, opts ...EndpointOption) *Endpoint {
	endpoint := &Endpoint{
		Path:      path,
		Summary:   summary,
		Method:    method,
		Query:     []QueryParam{},
		Headers:   []Header{},
		Form:      []FormField{},
		Body:      []BodyField{},
		Responses: []Response{},
	}

	for _, opt := range opts {
		opt(endpoint)
	}

	tag.Endpoints = append(tag.Endpoints, *endpoint)
	return endpoint
}

// Endpoint configuration options
func WithDescription(description string) EndpointOption {
	return func(e *Endpoint) {
		e.Description = description
	}
}

func WithOperationID(operationID string) EndpointOption {
	return func(e *Endpoint) {
		e.OperationID = operationID
	}
}

func WithBody(data any) EndpointOption {
	return func(e *Endpoint) {
		e.Body = structToSlice(data)
	}
}

func WithDeprecated(deprecated bool) EndpointOption {
	return func(e *Endpoint) {
		e.Deprecated = deprecated
	}
}

func WithSecurity(securityNames ...string) EndpointOption {
	return func(e *Endpoint) {
		e.Security = securityNames
	}
}

func WithTags(tags ...string) EndpointOption {
	return func(e *Endpoint) {
		e.Tags = tags
	}
}

// Query parameter options
type ParamOption func(q *QueryParam)

func WithRequired(required bool) ParamOption {
	return func(q *QueryParam) {
		q.Validation.Required = required
	}
}

func WithDefault(defaultValue any) ParamOption {
	return func(q *QueryParam) {
		q.Validation.Default = defaultValue
	}
}

func WithMin(min int) ParamOption {
	return func(q *QueryParam) {
		q.Validation.Min = min
	}
}

func WithMax(max int) ParamOption {
	return func(q *QueryParam) {
		q.Validation.Max = max
	}
}

func WithMinLength(minLen int) ParamOption {
	return func(q *QueryParam) {
		q.Validation.MinLen = minLen
	}
}

func WithMaxLength(maxLen int) ParamOption {
	return func(q *QueryParam) {
		q.Validation.MaxLen = maxLen
	}
}

func WithPattern(pattern string) ParamOption {
	return func(q *QueryParam) {
		q.Validation.Pattern = pattern
	}
}

func WithParamExample(example any) ParamOption {
	return func(q *QueryParam) {
		q.Example = example
	}
}

func WithEnum(values ...any) ParamOption {
	return func(q *QueryParam) {
		q.Enum = values
	}
}

func WithQuery(name, typ, description string, opts ...ParamOption) EndpointOption {
	return func(e *Endpoint) {
		query := QueryParam{
			Name:        name,
			Type:        typ,
			Description: description,
		}
		for _, opt := range opts {
			opt(&query)
		}
		e.Query = append(e.Query, query)
	}
}

// Header configuration
func WithHeader(name, typ, description string, required bool, example string) EndpointOption {
	return func(e *Endpoint) {
		header := Header{
			Name:        name,
			Type:        typ,
			Description: description,
			Required:    required,
			Example:     example,
		}
		e.Headers = append(e.Headers, header)
	}
}

// Form field configuration
func WithFormField(name, typ, description string, required bool, example any) EndpointOption {
	return func(e *Endpoint) {
		field := FormField{
			Name:        name,
			Type:        typ,
			Description: description,
			Required:    required,
			Example:     example,
		}
		e.Form = append(e.Form, field)
	}
}

// Response configuration
func WithResponse(statusCode int, description string, schema any, examples ...Example) EndpointOption {
	return func(e *Endpoint) {
		response := Response{
			StatusCode:  statusCode,
			Description: description,
			Examples:    examples,
		}
		if schema != nil {
			response.Schema = structToSlice(schema)
		}
		e.Responses = append(e.Responses, response)
	}
}

// Serve renders the documentation as HTML
func (d *APIDocumentation) Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base.html", d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ServeJSON serves the documentation as JSON
func (d *APIDocumentation) ServeJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Export returns the documentation as JSON bytes
func (d *APIDocumentation) Export() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}

// Security scheme helpers
func NewAPIKeySecurity(name, keyName, description string, in SecurityLocation) Security {
	return Security{
		Type:        SecurityTypeAPIKey,
		Name:        name,
		In:          in,
		Description: description,
		KeyName:     keyName,
	}
}

func NewBearerSecurity(name, description string) Security {
	return Security{
		Type:        SecurityTypeBearer,
		Name:        name,
		Scheme:      "bearer",
		In:          SecurityLocationHeader,
		Description: description,
	}
}

func NewBasicSecurity(name, description string) Security {
	return Security{
		Type:        SecurityTypeBasic,
		Name:        name,
		In:          SecurityLocationHeader,
		Scheme:      "basic",
		Description: description,
	}
}

// Helper for creating examples
func NewExample(name, summary, description string, value any) Example {
	return Example{
		Name:        name,
		Summary:     summary,
		Description: description,
		Value:       value,
	}
}
