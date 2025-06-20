package main

import (
	"log"
	"net/http"
	"time"

	"github.com/anvidev/apiduck"
)

type User struct {
	ID       int       `json:"id" description:"Unique user identifier" example:"123"`
	Name     string    `json:"name" description:"User's full name" validate:"required,min=2,max=100" example:"John Doe"`
	Email    string    `json:"email" description:"User's email address" validate:"required,email" example:"john@example.com"`
	Age      int       `json:"age" description:"User's age" validate:"min=13,max=120" example:"25"`
	Role     string    `json:"role" description:"User's role" enum:"admin,user,moderator" example:"user"`
	Active   bool      `json:"active" description:"Whether the user is active" default:"true"`
	Profile  *Profile  `json:"profile,omitempty" description:"User's profile information"`
	Tags     []string  `json:"tags,omitempty" description:"User tags" example:"developer,golang"`
	JoinedAt time.Time `json:"joined_at" description:"When the user joined" example:"2023-01-15T10:30:00Z"`
}

type Profile struct {
	Bio       string `json:"bio" description:"User's biography" validate:"max=500" example:"Software developer with 5 years of experience"`
	Website   string `json:"website,omitempty" description:"User's website URL" example:"https://johndoe.dev"`
	Location  string `json:"location,omitempty" description:"User's location" example:"San Francisco, CA"`
	AvatarURL string `json:"avatar_url,omitempty" description:"URL to user's avatar image" example:"https://example.com/avatar.jpg"`
}

type CreateUserRequest struct {
	Name    string   `json:"name" description:"User's full name" validate:"required,min=2,max=100"`
	Email   string   `json:"email" description:"User's email address" validate:"required,email"`
	Age     int      `json:"age" description:"User's age" validate:"min=13,max=120"`
	Role    string   `json:"role" description:"User's role" enum:"admin,user,moderator" default:"user"`
	Tags    []string `json:"tags,omitempty" description:"User tags"`
	Profile *Profile `json:"profile,omitempty" description:"User's profile information"`
}

type UpdateUserRequest struct {
	Name    *string  `json:"name,omitempty" description:"User's full name" validate:"min=2,max=100"`
	Email   *string  `json:"email,omitempty" description:"User's email address" validate:"email"`
	Age     *int     `json:"age,omitempty" description:"User's age" validate:"min=13,max=120"`
	Role    *string  `json:"role,omitempty" description:"User's role" enum:"admin,user,moderator"`
	Active  *bool    `json:"active,omitempty" description:"Whether the user is active"`
	Profile *Profile `json:"profile,omitempty" description:"User's profile information"`
}

type ErrorResponse struct {
	Error   string            `json:"error" description:"Error message" example:"Validation failed"`
	Code    string            `json:"code" description:"Error code" example:"VALIDATION_ERROR"`
	Details map[string]string `json:"details,omitempty" description:"Additional error details"`
}

type PaginatedResponse struct {
	Data       []User `json:"data" description:"List of users"`
	Total      int    `json:"total" description:"Total number of items" example:"150"`
	Page       int    `json:"page" description:"Current page number" example:"1"`
	PageSize   int    `json:"page_size" description:"Number of items per page" example:"20"`
	TotalPages int    `json:"total_pages" description:"Total number of pages" example:"8"`
}

func main() {
	docs := apiduck.NewDocumentation(apiduck.Info{
		Title:       "User Management API",
		Description: "A comprehensive API for managing users in the system",
		Version:     "v0.2.3",
		Contact: apiduck.InfoContact{
			Name:  "API Support",
			Email: "support@example.com",
			URL:   "https://example.com/support",
		},
		License: apiduck.InfoLicense{
			Name: "MIT",
			URL:  "https://opensource.org/licenses/MIT",
		},
		TermsOfService: "https://example.com/terms",
	})

	docs.AddServer("https://api.example.com", "Production server", nil)
	docs.AddServer("https://staging-api.example.com", "Staging server", nil)
	docs.AddServer("http://localhost:8080", "Development server", nil)

	// Add security schemes
	docs.AddSecurity(apiduck.NewBearerSecurity("bearerAuth", "JWT Bearer token authentication"))
	docs.AddSecurity(apiduck.NewAPIKeySecurity("apiKey", "X-API-Key", "header", apiduck.SecurityLocationHeader))
	docs.AddSecurity(apiduck.NewBasicSecurity("admin", "Admin authentication"))

	// Create tags
	userTag := docs.AddTag("Users", "Operations related to user management")
	authTag := docs.AddTag("Authentication", "Authentication and authorization endpoints")

	// Add user endpoints
	userTag.AddEndpoint(
		apiduck.MethodGet,
		"/api/v1/users",
		"List all users",
		apiduck.WithDescription("Retrieve a paginated list of all users in the system"),
		apiduck.WithOperationID("listUsers"),
		apiduck.WithSecurity("bearerAuth"),
		apiduck.WithQuery("page", "integer", "Page number for pagination",
			apiduck.WithDefault(1),
			apiduck.WithMin(1),
			apiduck.WithParamExample(1),
		),
		apiduck.WithQuery("page_size", "integer", "Number of items per page",
			apiduck.WithDefault(20),
			apiduck.WithMin(1),
			apiduck.WithMax(100),
			apiduck.WithParamExample(20),
		),
		apiduck.WithQuery("search", "string", "Search users by name or email",
			apiduck.WithMinLength(2),
			apiduck.WithMaxLength(100),
			apiduck.WithParamExample("john"),
		),
		apiduck.WithQuery("role", "string", "Filter users by role",
			apiduck.WithEnum("admin", "user", "moderator"),
			apiduck.WithParamExample("user"),
		),
		apiduck.WithQuery("active", "boolean", "Filter users by active status",
			apiduck.WithParamExample(true),
		),
		apiduck.WithHeader("X-Request-ID", "string", "Unique request identifier", false, "req-123456"),
		apiduck.WithResponse(200, "Successful response", PaginatedResponse{},
			apiduck.NewExample("success", "Successful response", "List of users with pagination", map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":    1,
						"name":  "John Doe",
						"email": "john@example.com",
						"age":   25,
						"role":  "user",
					},
				},
				"total":       1,
				"page":        1,
				"page_size":   20,
				"total_pages": 1,
			}),
		),
		apiduck.WithResponse(400, "Bad request", ErrorResponse{}),
		apiduck.WithResponse(401, "Unauthorized", ErrorResponse{}),
		apiduck.WithResponse(500, "Internal server error", ErrorResponse{}),
	)

	userTag.AddEndpoint(
		apiduck.MethodGet,
		"/api/v1/users/{id}",
		"Get user by ID",
		apiduck.WithDescription("Retrieve a specific user by their unique identifier"),
		apiduck.WithOperationID("getUserById"),
		apiduck.WithSecurity("bearerAuth"),
		apiduck.WithQuery("include", "string", "Include additional user data",
			apiduck.WithEnum("profile", "tags", "profile,tags"),
			apiduck.WithParamExample("profile"),
		),
		apiduck.WithResponse(200, "User found", User{}),
		apiduck.WithResponse(404, "User not found", ErrorResponse{}),
		apiduck.WithResponse(401, "Unauthorized", ErrorResponse{}),
	)

	userTag.AddEndpoint(
		apiduck.MethodPost,
		"/api/v1/users",
		"Create a new user",
		apiduck.WithDescription("Create a new user in the system"),
		apiduck.WithOperationID("createUser"),
		apiduck.WithSecurity("bearerAuth"),
		apiduck.WithBody(CreateUserRequest{}),
		apiduck.WithHeader("Content-Type", "string", "Request content type", true, "application/json"),
		apiduck.WithResponse(201, "User created successfully", User{},
			apiduck.NewExample("created", "User created", "Newly created user", map[string]interface{}{
				"id":    123,
				"name":  "John Doe",
				"email": "john@example.com",
				"age":   25,
				"role":  "user",
			}),
		),
		apiduck.WithResponse(400, "Validation error", ErrorResponse{},
			apiduck.NewExample("validation_error", "Validation failed", "Request body validation failed", map[string]interface{}{
				"error": "Validation failed",
				"code":  "VALIDATION_ERROR",
				"details": map[string]string{
					"email": "Invalid email format",
					"age":   "Age must be between 13 and 120",
				},
			}),
		),
		apiduck.WithResponse(401, "Unauthorized", ErrorResponse{}),
		apiduck.WithResponse(409, "User already exists", ErrorResponse{}),
	)

	userTag.AddEndpoint(
		apiduck.MethodPut,
		"/api/v1/users/{id}",
		"Update user",
		apiduck.WithDescription("Update an existing user's information"),
		apiduck.WithOperationID("updateUser"),
		apiduck.WithSecurity("bearerAuth"),
		apiduck.WithBody(UpdateUserRequest{}),
		apiduck.WithResponse(200, "User updated successfully", User{}),
		apiduck.WithResponse(400, "Validation error", ErrorResponse{}),
		apiduck.WithResponse(401, "Unauthorized", ErrorResponse{}),
		apiduck.WithResponse(404, "User not found", ErrorResponse{}),
	)

	userTag.AddEndpoint(
		apiduck.MethodDelete,
		"/api/v1/users/{id}",
		"Delete user",
		apiduck.WithDescription("Delete a user from the system"),
		apiduck.WithOperationID("deleteUser"),
		apiduck.WithSecurity("bearerAuth"),
		apiduck.WithResponse(204, "User deleted successfully", nil),
		apiduck.WithResponse(401, "Unauthorized", ErrorResponse{}),
		apiduck.WithResponse(404, "User not found", ErrorResponse{}),
	)

	// Add authentication endpoints
	authTag.AddEndpoint(
		apiduck.MethodPost,
		"/api/v1/auth/login",
		"User login",
		apiduck.WithDescription("Authenticate a user and return an access token"),
		apiduck.WithOperationID("login"),
		apiduck.WithFormField("email", "string", "User's email address", true, "john@example.com"),
		apiduck.WithFormField("password", "string", "User's password", true, "password123"),
		apiduck.WithResponse(200, "Login successful", map[string]interface{}{
			"access_token":  "string",
			"refresh_token": "string",
			"expires_in":    "integer",
			"token_type":    "string",
		},
			apiduck.NewExample("login_success", "Successful login", "User authenticated successfully", map[string]interface{}{
				"access_token":  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				"refresh_token": "def502004a8b9c...",
				"expires_in":    3600,
				"token_type":    "Bearer",
			}),
		),
		apiduck.WithResponse(401, "Invalid credentials", ErrorResponse{},
			apiduck.NewExample("invalid_credentials", "Login failed", "Invalid email or password", map[string]interface{}{
				"error": "Invalid credentials",
				"code":  "INVALID_CREDENTIALS",
			}),
		),
		apiduck.WithResponse(400, "Bad request", ErrorResponse{}),
	)

	authTag.AddEndpoint(
		apiduck.MethodPost,
		"/api/v1/auth/refresh",
		"Refresh access token",
		apiduck.WithDescription("Refresh an expired access token using a refresh token"),
		apiduck.WithOperationID("refreshToken"),
		apiduck.WithFormField("refresh_token", "string", "Valid refresh token", true, "def502004a8b9c..."),
		apiduck.WithResponse(200, "Token refreshed successfully", map[string]interface{}{
			"access_token": "string",
			"expires_in":   "integer",
			"token_type":   "string",
		}),
		apiduck.WithResponse(401, "Invalid refresh token", ErrorResponse{}),
	)

	authTag.AddEndpoint(
		apiduck.MethodPost,
		"/api/v1/auth/logout",
		"User logout",
		apiduck.WithDescription("Logout the current user and invalidate their token"),
		apiduck.WithOperationID("logout"),
		apiduck.WithSecurity("bearerAuth"),
		apiduck.WithResponse(200, "Logout successful", map[string]interface{}{
			"message": "string",
		}),
		apiduck.WithResponse(401, "Unauthorized", ErrorResponse{}),
	)

	// Add a deprecated endpoint as an example
	userTag.AddEndpoint(
		apiduck.MethodGet,
		"/api/v1/users/profile",
		"Get current user profile",
		apiduck.WithDescription("This endpoint is deprecated. Use GET /api/v1/users/{id} instead."),
		apiduck.WithOperationID("getCurrentUserProfile"),
		apiduck.WithSecurity("bearerAuth"),
		apiduck.WithDeprecated(true),
		apiduck.WithResponse(200, "User profile", User{}),
		apiduck.WithResponse(401, "Unauthorized", ErrorResponse{}),
	)

	http.HandleFunc("/docs", docs.Serve)
	http.HandleFunc("/docs.json", docs.ServeJSON)

	// Optional: Save documentation to file
	if jsonData, err := docs.Export(); err == nil {
		log.Println("Documentation exported successfully")
		// You could write this to a file:
		// os.WriteFile("api-docs.json", jsonData, 0644)
		_ = jsonData
	}

	log.Println("API Documentation server starting on :8080")
	log.Println("Visit http://localhost:8080/docs for HTML documentation")
	log.Println("Visit http://localhost:8080/docs.json for JSON documentation")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
