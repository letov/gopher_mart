package router

import (
	"github.com/go-andiamo/chioas"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/infrastructure/handler"
	"net/http"
	"os"
)

func NewOpenApi(
	hs *handler.List,
	config *config.Config,
) *chioas.Definition {
	tokenAuth := jwtauth.New("HS256", []byte(config.JwtKey), nil)
	ms := chi.Middlewares{
		jwtauth.Verifier(tokenAuth),
		jwtauth.Authenticator,
	}

	api := &chioas.Definition{
		AutoHeadMethods: true,
		DocOptions: chioas.DocOptions{
			ServeDocs:       true,
			HideHeadMethods: true,
		},
		Paths: chioas.Paths{
			"/api": {
				Paths: chioas.Paths{
					"/user": {
						Paths: chioas.Paths{
							"/register": {
								Tag: "auth",
								Methods: chioas.Methods{
									http.MethodPost: {
										Handler: hs.Get(handler.SaveUserName),
										Request: &chioas.Request{
											SchemaRef: "RequestSaveUser",
										},
										Responses: chioas.Responses{
											http.StatusOK:         {},
											http.StatusBadRequest: {},
										},
									},
								},
							},
							"/login": {
								Tag: "auth",
								Methods: chioas.Methods{
									http.MethodPost: {
										Handler: hs.Get(handler.LoginName),
										Request: &chioas.Request{
											SchemaRef: "RequestLogin",
										},
										Responses: chioas.Responses{
											http.StatusOK: {
												SchemaRef: "ResponsesLogin",
											},
											http.StatusBadRequest: {},
										},
									},
								},
							},
							"/orders": {
								Middlewares: ms,
								Tag:         "orders",
								Methods: chioas.Methods{
									http.MethodPost: {
										Security: chioas.SecuritySchemes{
											chioas.SecurityScheme{
												Name: "bearerAuth",
											},
										},
										Handler: hs.Get(handler.RequestAccrualName),
										Request: &chioas.Request{
											SchemaRef: "RequestRequestAccrual",
										},
										Responses: chioas.Responses{
											http.StatusOK:         {},
											http.StatusBadRequest: {},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Components: &chioas.Components{
			SecuritySchemes: chioas.SecuritySchemes{
				{
					Name:   "bearerAuth",
					Scheme: "bearer",
				},
			},
			Schemas: chioas.Schemas{
				{
					Name: "RequestSaveUser",
					RequiredProperties: []string{
						"login",
						"password",
						"name",
					},
					Properties: chioas.Properties{
						{
							Name: "login",
							Type: "string",
						},
						{
							Name: "password",
							Type: "string",
						},
						{
							Name: "name",
							Type: "string",
						},
					},
				},
				{
					Name: "RequestLogin",
					RequiredProperties: []string{
						"login",
						"password",
					},
					Properties: chioas.Properties{
						{
							Name: "login",
							Type: "string",
						},
						{
							Name: "password",
							Type: "string",
						},
					},
				},
				{
					Name: "ResponsesLogin",
					RequiredProperties: []string{
						"token",
					},
					Properties: chioas.Properties{
						{
							Name: "token",
							Type: "string",
						},
					},
				},
				{
					Name: "RequestRequestAccrual",
					RequiredProperties: []string{
						"orderId",
					},
					Properties: chioas.Properties{
						{
							Name: "orderId",
							Type: "number",
						},
					},
				},
			},
		},
	}

	data, _ := api.AsYaml()
	fp := "openapi.yaml"
	_ = os.WriteFile(fp, data, 0777)

	return api
}
