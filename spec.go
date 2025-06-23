package apiduck

type httpMethod string

const (
	GET     httpMethod = "GET"
	POST    httpMethod = "POST"
	PUT     httpMethod = "PUT"
	DELETE  httpMethod = "DELETE"
	PATCH   httpMethod = "PATCH"
	HEAD    httpMethod = "HEAD"
	OPTIONS httpMethod = "OPTIONS"
)

type contentType string

const (
	ContentTypeJSON contentType = "application/json"
	ContentTypeXML  contentType = "application/xml"
	ContentTypeForm contentType = "application/x-www-form-urlencoded"
	ContentTypeText contentType = "text/plain"
	ContentTypeHTML contentType = "text/html"
)

type securityType string

const (
	BearerType securityType = "bearerToken"
	ApiKeyType              = "apiKey"
	BasicType               = "basic"
)

type securityLocation string

const (
	Headers securityLocation = "header"
	Cookies                  = "cookies"
	Query                    = "query"
)

type Documentation struct {
	Info      Info
	Servers   []Server
	Security  []SecurityScheme
	Resources []Resource
}

type Info struct {
	Title   string
	Desc    string
	Version string
	Contact Contact
	License License
	Terms   string
}

type Contact struct {
	Name  string
	Email string
	URL   string
}

type License struct {
	Name string
	URL  string
}

type Server struct {
	URL  string
	Desc string
}

type SecurityScheme struct {
	Name    string
	Desc    string
	KeyName string
	Type    securityType
	In      securityLocation
}

type Resource struct {
	Name      string
	Desc      string
	Endpoints []Endpoint
}

type Endpoint struct {
	Method     httpMethod
	Path       string
	Summary    string
	Desc       string
	Deprecated bool

	Parameters Parameters
	Request    *Request
	Responses  []Response

	Securities []string
}

type PathParameter Parameter
type QueryParameter Parameter
type HeaderParameter Parameter
type CookieParameter Parameter

type Parameter struct {
	Name         string
	Type         string
	Desc         string
	Req          bool
	DefaultValue any

	Ex any

	Enums   []any
	MinLen  *int
	MaxLen  *int
	Minimum *float64
	Maximum *float64
}

type Parameters struct {
	Path   []*PathParameter
	Query  []*QueryParameter
	Header []*HeaderParameter
	Cookie []*CookieParameter
}

type Request struct {
	Desc        string
	ContentType contentType
	Fields      []Field
	Ex          any
	JSON        string
}

type Response struct {
	StatusCode  int
	Desc        string
	Headers     []any
	ContentType contentType
	Fields      []Field
	Ex          any
	JSON        string
}

type Field struct {
	Fields []Field

	Name         string
	Desc         string
	Type         string
	Req          bool
	DefaultValue any

	Ex any

	Enums   []any
	MinLen  *string
	MaxLen  *string
	Minimum *string
	Maximum *string
}

type M map[string]any
