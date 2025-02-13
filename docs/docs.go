// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": [[ marshal .Schemes ]],
    "swagger": "2.0",
    "info": {
        "description": "[[escape .Description]]",
        "title": "[[.Title]]",
        "contact": {
            "name": "API Support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "[[.Version]]"
    },
    "host": "[[.Host]]",
    "basePath": "[[.BasePath]]",
    "paths": {
        "/api/v1": {
            "get": {
                "description": "Display the main page with redirect button",
                "consumes": [
                    "text/html"
                ],
                "produces": [
                    "text/html"
                ],
                "tags": [
                    "Index"
                ],
                "summary": "Show index page",
                "responses": {
                    "200": {
                        "description": "HTML content",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/auth/login": {
            "post": {
                "description": "Authenticate user with username and password to obtain an API token for protected endpoints",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "Login credentials",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully logged in",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.AuthResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Invalid request or validation error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/response.APIError"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "Invalid credentials",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/response.APIError"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/auth/register": {
            "post": {
                "description": "Register a new user account with username, password and email. The password must be at least 8 characters long.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "Registration details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Successfully registered",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.AuthResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Invalid request or validation error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/response.APIError"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "409": {
                        "description": "Username or email already exists",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/response.APIError"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/response.APIError"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/ffmpeg": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Submit a video processing job using FFMPEG. The command should use placeholders like {{in1}} for input files and {{out1}} for output files.\nThese placeholders will be replaced with actual file paths during processing.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "FFMPEG"
                ],
                "summary": "Process video with FFMPEG",
                "parameters": [
                    {
                        "description": "FFMPEG processing details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.FFMPEGRequest"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Job accepted for processing",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.FFMPEGResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Invalid request or validation error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/response.APIError"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "Missing or invalid API token",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/response.APIError"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/response.APIError"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/ffmpeg/progress/{uuid}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get the current status and progress of a video processing job. Returns details about output files when the job is completed.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "FFMPEG"
                ],
                "summary": "Get job progress",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Job UUID returned from the process endpoint",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Job status retrieved successfully",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.JobStatus"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Invalid UUID format",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/response.APIError"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "Missing or invalid API token",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/response.APIError"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "Job not found",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/response.APIError"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.OutputFileMetadata": {
            "type": "object",
            "properties": {
                "file_format": {
                    "type": "string"
                },
                "file_id": {
                    "type": "string"
                },
                "file_type": {
                    "type": "string"
                },
                "height": {
                    "type": "integer"
                },
                "size_mbytes": {
                    "type": "number"
                },
                "storage_url": {
                    "type": "string"
                },
                "width": {
                    "type": "integer"
                }
            }
        },
        "dto.AuthResponse": {
            "type": "object",
            "properties": {
                "api_token": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "dto.FFMPEGRequest": {
            "type": "object",
            "required": [
                "ffmpeg_command",
                "input_files",
                "output_files"
            ],
            "properties": {
                "ffmpeg_command": {
                    "type": "string",
                    "example": "-i {{in1}} {{out1}}"
                },
                "input_files": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    },
                    "example": {
                        "{\"in1\"": " \"https://storage.googleapis.com/ffmpeg-api-test-bucket/user_1/input/test.mp4\"}"
                    }
                },
                "output_files": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    },
                    "example": {
                        "{\"out1\"": " \"string.mp4\"}"
                    }
                }
            }
        },
        "dto.FFMPEGResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "dto.JobStatus": {
            "type": "object",
            "required": [
                "status"
            ],
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                },
                "output_files": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/domain.OutputFileMetadata"
                    }
                },
                "progress": {
                    "type": "integer"
                },
                "result": {
                    "type": "string"
                },
                "status": {
                    "type": "string",
                    "enum": [
                        "pending",
                        "processing",
                        "completed",
                        "failed"
                    ]
                },
                "updated_at": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "dto.LoginRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "dto.RegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 8
                },
                "username": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 3
                }
            }
        },
        "response.APIError": {
            "type": "object",
            "properties": {
                "message": {},
                "type": {
                    "type": "string"
                }
            }
        },
        "response.Response": {
            "type": "object",
            "properties": {
                "data": {},
                "error": {
                    "$ref": "#/definitions/response.APIError"
                },
                "success": {
                    "type": "boolean"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "API token obtained after login. Required for all protected endpoints.",
            "type": "apiKey",
            "name": "X-API-Token",
            "in": "header"
        }
    },
    "tags": [
        {
            "description": "Authentication endpoints for user registration and login",
            "name": "Auth"
        },
        {
            "description": "Video processing endpoints using FFMPEG",
            "name": "FFMPEG"
        },
        {
            "description": "Main page and general information",
            "name": "Index"
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "/api/v1",
	Schemes:          []string{"http"},
	Title:            "FFMPEG Serverless API",
	Description:      "A serverless API for processing videos using FFMPEG. This API allows you to submit video processing jobs, monitor their progress, and manage user authentication.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "[[",
	RightDelim:       "]]",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
