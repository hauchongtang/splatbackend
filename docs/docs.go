// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/cached/users/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get a User by id from cache",
                "parameters": [
                    {
                        "type": "string",
                        "description": "userId",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.userType"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controllers.errorResult"
                        }
                    }
                }
            }
        },
        "/users/login": {
            "post": {
                "description": "Responds with user details, including OAuth2 tokens.",
                "produces": [
                    "application/json"
                ],
                "summary": "User log in",
                "parameters": [
                    {
                        "description": "Sign in credentials",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.userLogin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.userType"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.errorResult"
                        }
                    }
                }
            }
        },
        "/users/signup": {
            "post": {
                "description": "Responds with userId",
                "produces": [
                    "application/json"
                ],
                "summary": "User sign up",
                "parameters": [
                    {
                        "description": "New user credentials",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.userSignUp"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.signUpResult"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controllers.errorResult"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.errorResult": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "controllers.signUpResult": {
            "type": "object",
            "properties": {
                "InsertedID": {
                    "type": "string"
                }
            }
        },
        "controllers.userLogin": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 6
                }
            }
        },
        "controllers.userSignUp": {
            "type": "object",
            "required": [
                "first_name",
                "last_name"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 1
                },
                "last_name": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 1
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "controllers.userType": {
            "type": "object",
            "required": [
                "Password",
                "email",
                "first_name",
                "last_name"
            ],
            "properties": {
                "Password": {
                    "type": "string",
                    "minLength": 6
                },
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 1
                },
                "id": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 1
                },
                "points": {
                    "type": "integer"
                },
                "refresh_token": {
                    "type": "string"
                },
                "timetable": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "SplatApp Backend API",
	Description:      "This is the backend service for splatapp at https://github.com/hauchongtang/splatbackend",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
