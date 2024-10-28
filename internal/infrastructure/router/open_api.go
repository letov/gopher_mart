package router

import (
	"github.com/go-andiamo/chioas"
	"gopher_mart/internal/infrastructure/handler"
	"net/http"
)

func NewOpenApi(
	hs *handler.List,
) *chioas.Definition {
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
						},
					},
				},
			},
		},
		Components: &chioas.Components{
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
					Name:               "RequestLogin",
					RequiredProperties: []string{"token"},
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
					Name:               "ResponsesLogin",
					RequiredProperties: []string{"token"},
					Properties: chioas.Properties{
						{
							Name: "token",
							Type: "string",
						},
					},
				},
			},
		},
	}

	//data, _ := api.AsYaml()
	//fp := "openapi.yaml"
	//_ = os.WriteFile(fp, data, 0777)

	return api
}
