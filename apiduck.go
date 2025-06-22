package apiduck

import (
	"embed"
	"encoding/json"
	"html/template"
	"net/http"
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

type documentationOption func(*Documentation)

func New(title, description, version string, options ...documentationOption) *Documentation {
	docs := &Documentation{
		Info: Info{
			Title:   title,
			Desc:    description,
			Version: version,
		},
	}

	for _, opt := range options {
		opt(docs)
	}

	return docs
}

func WithContact(name, email, url string) documentationOption {
	return func(d *Documentation) {
		d.Info.Contact = Contact{
			Name:  name,
			Email: email,
			URL:   url,
		}
	}
}

func WithLicense(name, url string) documentationOption {
	return func(d *Documentation) {
		d.Info.License = License{
			Name: name,
			URL:  url,
		}
	}
}

func WithTerms(terms string) documentationOption {
	return func(d *Documentation) {
		d.Info.Terms = terms
	}
}

func (d *Documentation) AddServer(url, description string) *Documentation {
	d.Servers = append(d.Servers, Server{url, description})
	return d
}

func (d *Documentation) AddSecurity(security SecurityScheme) *Documentation {
	d.Security = append(d.Security, security)
	return d
}

func BearerToken(name, description string) SecurityScheme {
	return SecurityScheme{
		Name:    name,
		Desc:    description,
		KeyName: "Authorization",
		Type:    BearerType,
		In:      Headers,
	}
}

func ApiKey(name, description, keyname string, in securityLocation) SecurityScheme {
	return SecurityScheme{
		Name:    name,
		Desc:    description,
		KeyName: keyname,
		Type:    ApiKeyType,
		In:      in,
	}
}

func Basic(name, description string) SecurityScheme {
	return SecurityScheme{
		Name:    name,
		Desc:    description,
		KeyName: "Authorization",
		Type:    BasicType,
		In:      Headers,
	}
}

func NewSecurityScheme(name, description, keyname string, typ securityType, in securityLocation) SecurityScheme {
	return SecurityScheme{
		Name:    name,
		Desc:    description,
		KeyName: keyname,
		Type:    typ,
		In:      in,
	}
}

func (d *Documentation) AddResource(name, description string) *Resource {
	resource := Resource{
		Name:      name,
		Desc:      description,
		Endpoints: make([]Endpoint, 0),
	}
	d.Resources = append(d.Resources, resource)
	return &d.Resources[len(d.Resources)-1]
}

func (r *Resource) Post(path, summary, description string) *Endpoint {
	endpoint := Endpoint{
		Method:  POST,
		Path:    path,
		Summary: summary,
		Desc:    description,
	}
	r.Endpoints = append(r.Endpoints, endpoint)
	return &r.Endpoints[len(r.Endpoints)-1]
}

func (r *Resource) Get(path, summary, description string) *Endpoint {
	endpoint := Endpoint{
		Method:  GET,
		Path:    path,
		Summary: summary,
		Desc:    description,
	}
	r.Endpoints = append(r.Endpoints, endpoint)
	return &r.Endpoints[len(r.Endpoints)-1]
}

func (r *Resource) Put(path, summary, description string) *Endpoint {
	endpoint := Endpoint{
		Method:  PUT,
		Path:    path,
		Summary: summary,
		Desc:    description,
	}
	r.Endpoints = append(r.Endpoints, endpoint)
	return &r.Endpoints[len(r.Endpoints)-1]
}

func (r *Resource) Delete(path, summary, description string) *Endpoint {
	endpoint := Endpoint{
		Method:  DELETE,
		Path:    path,
		Summary: summary,
		Desc:    description,
	}
	r.Endpoints = append(r.Endpoints, endpoint)
	return &r.Endpoints[len(r.Endpoints)-1]
}

func (r *Resource) Patch(path, summary, description string) *Endpoint {
	endpoint := Endpoint{
		Method:  PATCH,
		Path:    path,
		Summary: summary,
		Desc:    description,
	}
	r.Endpoints = append(r.Endpoints, endpoint)
	return &r.Endpoints[len(r.Endpoints)-1]
}

func (r *Resource) Options(path, summary, description string) *Endpoint {
	endpoint := Endpoint{
		Method:  OPTIONS,
		Path:    path,
		Summary: summary,
		Desc:    description,
	}
	r.Endpoints = append(r.Endpoints, endpoint)
	return &r.Endpoints[len(r.Endpoints)-1]
}

func (e *Endpoint) Security(names ...string) *Endpoint {
	e.Securities = append(e.Securities, names...)
	return e
}

func (e *Endpoint) PathParams(params ...*PathParameter) *Endpoint {
	e.Parameters.Path = append(e.Parameters.Path, params...)
	return e
}

func (e *Endpoint) Queries(queries ...*QueryParameter) *Endpoint {
	e.Parameters.Query = append(e.Parameters.Query, queries...)
	return e
}

func (e *Endpoint) Headers(headers ...*HeaderParameter) *Endpoint {
	e.Parameters.Header = append(e.Parameters.Header, headers...)
	return e
}

func (e *Endpoint) Body(body *Request) *Endpoint {
	e.Request = body
	return e
}

func (e *Endpoint) Response(response *Response) *Endpoint {
	e.Responses = append(e.Responses, *response)
	return e
}

func PathParam(name, description string) *PathParameter {
	return &PathParameter{
		Name: name,
		Type: "string", // TODO: make type and inclide in function parameters
		Desc: description,
		Req:  true,
		Ex:   nil,
	}
}

func (p *PathParameter) Example(v any) *PathParameter {
	p.Ex = v
	return p
}

func QueryParam(name, description string) *QueryParameter {
	return &QueryParameter{
		Name: name,
		Type: "string", // TODO: make type and inclide in function parameters
		Desc: description,
		Req:  false,
		Ex:   nil,
	}
}

func (p *QueryParameter) Required() *QueryParameter {
	p.Req = true
	return p
}

func (p *QueryParameter) Example(v any) *QueryParameter {
	p.Ex = v
	return p
}

func (p *QueryParameter) Enum(v ...any) *QueryParameter {
	p.Enums = append(p.Enums, v...)
	return p
}

func (p *QueryParameter) Min(min float64) *QueryParameter {
	p.Minimum = &min
	return p
}

func (p *QueryParameter) Max(max float64) *QueryParameter {
	p.Maximum = &max
	return p
}

func (p *QueryParameter) MinLength(min int) *QueryParameter {
	p.MinLen = &min
	return p
}

func (p *QueryParameter) MaxLength(max int) *QueryParameter {
	p.MaxLen = &max
	return p
}

func HeaderParam(key, description string) *HeaderParameter {
	return &HeaderParameter{
		Name: key,
		Type: "string", // TODO: make type and inclide in function parameters
		Desc: description,
		Req:  false,
		Ex:   nil,
	}
}

func (p *HeaderParameter) Required() *HeaderParameter {
	p.Req = true
	return p
}

func (p *HeaderParameter) Example(v any) *HeaderParameter {
	p.Ex = v
	return p
}

func JSONBody(v any) *Request {
	return &Request{
		ContentType: ContentTypeJSON,
		Fields:      parseStruct(v),
		Ex:          nil,
	}
}

func (b *Request) Example(v any) *Request {
	jsonBytes, _ := json.Marshal(v)
	b.Ex = v
	b.JSON = string(jsonBytes)
	return b
}

func JSONResponse(statusCode int, v any) *Response {
	statusMessage, ok := statusCodeMessages[statusCode]
	if !ok {
		statusMessage = ""
	}
	response := &Response{
		StatusCode:  statusCode,
		Desc:        statusMessage,
		ContentType: ContentTypeJSON,
		Fields:      parseStruct(v),
		Ex:          nil,
	}
	return response
}

func (r *Response) Description(description string) *Response {
	r.Desc = description
	return r
}

func (r *Response) Example(v any) *Response {
	jsonBytes, _ := json.Marshal(v)
	r.Ex = v
	r.JSON = string(jsonBytes)
	return r
}

func (d *Documentation) Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base.html", d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
