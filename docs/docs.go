// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/expand": {
            "get": {
                "description": "Преобразует короткую ссылку в исходную длинную ссылку.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Расширение URL"
                ],
                "summary": "Расширить короткую ссылку до её оригинальной формы",
                "parameters": [
                    {
                        "description": "Короткая ссылка",
                        "name": "shortUrl",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ShortURL"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Расширенная длинная ссылка",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Неверный ввод",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/shorten": {
            "post": {
                "description": "Преобразует длинную ссылку в компактную форму.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Сокращение URL"
                ],
                "summary": "Сократить длинную ссылку",
                "parameters": [
                    {
                        "description": "Длинная ссылка",
                        "name": "longUrl",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.LongURL"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Сокращённая ссылка",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Неверный ввод",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.ErrorResponse": {
            "description": "Формат ответа об ошибке",
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "model.LongURL": {
            "type": "object",
            "required": [
                "long_url"
            ],
            "properties": {
                "long_url": {
                    "type": "string"
                }
            }
        },
        "model.ShortURL": {
            "type": "object",
            "required": [
                "short_url"
            ],
            "properties": {
                "short_url": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "URL Shortener API",
	Description:      "This is a sample API for a URL shortener with Swagger documentation.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
