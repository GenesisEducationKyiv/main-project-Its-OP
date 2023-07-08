// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/rate": {
            "get": {
                "description": "Get the current rate of BTC to UAH using any third-party service with public API",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "rate"
                ],
                "summary": "Get current BTC to UAH rate",
                "responses": {
                    "200": {
                        "description": "Successful operation",
                        "schema": {
                            "type": "number"
                        }
                    },
                    "400": {
                        "description": "Invalid status value",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/sendEmails": {
            "post": {
                "description": "Send the current BTC to UAH rate to all subscribed emails",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscription"
                ],
                "summary": "Send email with BTC rate",
                "responses": {
                    "200": {
                        "description": "E-mails sent",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/subscribe": {
            "post": {
                "description": "Add an email to the database if it does not exist already",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscription"
                ],
                "summary": "Subscribe email to get BTC rate",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email to be subscribed",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "E-mail added",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "409": {
                        "description": "E-mail already exists in the database",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "GSES2 BTC application API",
	Description:      "This is a sample server for a BTC to UAH rate application.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}