package main

import (
	"log"
	"net/http"

	ad "github.com/anvidev/apiduck"
)

type DummyBody struct {
	Name     string   `json:"name" validate:"min=3,required"`
	Email    string   `json:"email" validate:"email,required"`
	Password string   `json:"password" validate:"min=8,max=32,required"`
	Nested   []Nested `json:"nested" apiduck:"desc=this is a test"`
}

type Nested struct {
	Testing string       `json:"testing"`
	Deep    DeeplyNested `json:"deeplyNested"`
}

type DeeplyNested struct {
	Deep int64 `json:"deep"`
}

func main() {
	docs := ad.New(
		"Tidsregistrering API",
		"Internt værktøj for Skancode A/S til at dokumentere og overskue tid brugt på diverse projekter.",
		"v0.1.1",
		ad.WithContact(
			"Skancode Support",
			"support@skancode.dk",
			"www.skancode.dk",
		),
		ad.WithLicense(
			"MIT",
			"https://opensource.org/licenses/MIT",
		),
		ad.WithTerms("https://example.com/terms"),
	)

	docs.AddServer("http://localhost:9090", "Development Server")
	docs.AddServer("https://tid-api.anvi.dev", "Produktions Server")

	docs.AddSecurity(ad.BearerToken("Bruger authentication", "Alle brugere skal authenticates for at få adgang til ressourcer"))

	authResource := docs.AddResource("Auth", "Authentication og authorization endpoints")

	authResource.Post("/v1/auth/register", "Opret bruger", "Opretter en ny bruger til Tidsregistrerings Portal").
		Security("Bruger authentication").
		PathParams(
			ad.PathParam("Testing", "Test query parameter").
				Example("Foo"),
		).
		Queries(
			ad.QueryParam("adminsOnly", "Filter only admins").
				Required().
				Enum(true, false).
				Example(true).
				Example(false).
				Min(1).
				Max(10).
				MinLength(4).
				MaxLength(4),
		).
		Headers(
			ad.HeaderParam("X-Some-Key", "This is a very important header").
				Required().
				Example("Foo-bar-baz"),
			ad.HeaderParam("X-Some-Key2", "This is also a very important header").
				Required().
				Example("Foo-bar-baz2"),
		).
		Body(
			ad.JSONBody(DummyBody{}).
				Example(DummyBody{
					Name:     "John Doe",
					Email:    "john@doe.com",
					Password: "pa$$w0rd",
					Nested:   []Nested{},
				}),
		).
		Response(
			ad.JSONResponse(http.StatusOK, DummyBody{}).
				Example(DummyBody{
					Name:     "John Doe",
					Email:    "john@doe.com",
					Password: "pa$$w0rd",
					Nested: []Nested{
						{
							Testing: "",
							Deep: DeeplyNested{
								Deep: 0,
							},
						},
					},
				}),
		).
		Response(
			ad.JSONResponse(http.StatusConflict, DummyBody{}).
				Example(map[string]any{
					"error": "email er allerede i brug",
				}),
		)

	http.HandleFunc("/docs", docs.Serve)

	log.Println("Server running...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
