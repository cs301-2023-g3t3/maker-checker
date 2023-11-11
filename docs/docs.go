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
        "/permission": {
            "get": {
                "description": "Retrieves a list of makerchecker permission",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "permission"
                ],
                "summary": "Get all Makerchecker Permission",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Permission"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            },
            "post": {
                "description": "Add a Makerchecker Permission",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "permission"
                ],
                "summary": "Add a Makerchecker Permission",
                "parameters": [
                    {
                        "description": "Permission Body",
                        "name": "requestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Permission"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Permission"
                        }
                    },
                    "400": {
                        "description": "Invalid permission object or endpoint already exists",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            }
        },
        "/permission/by-endpoint": {
            "get": {
                "description": "Retrieve a Makerchecker Permission by Endpoint route",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "permission"
                ],
                "summary": "Get Makerchecker Permission by Endpoint route",
                "parameters": [
                    {
                        "description": "endpoint",
                        "name": "endpoint",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Permission"
                            }
                        }
                    },
                    "400": {
                        "description": "Endpoint cannot be empty",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "404": {
                        "description": "Makerchecker permission cannot be found with endpoint route",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            }
        },
        "/permission/{id}": {
            "get": {
                "description": "Retrieve a Makerchecker permission by Id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "permission"
                ],
                "summary": "Get Makerchecker permission by Id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Permission"
                        }
                    },
                    "400": {
                        "description": "Id cannot be empty",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "404": {
                        "description": "No permission can be found with Id",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            },
            "put": {
                "description": "Update Makerchecker Permission by Id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "permission"
                ],
                "summary": "Update Makerchecker Permission by Id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request Body",
                        "name": "requestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/permission.UpdatePermission"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Permission"
                        }
                    },
                    "400": {
                        "description": "Id cannot be empty and permission object is invalid",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "404": {
                        "description": "No permission found",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a Makerchecker Permission By ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "permission"
                ],
                "summary": "Delete a Makerchecker Permission by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Permission"
                        }
                    },
                    "400": {
                        "description": "Id cannot be empty",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "404": {
                        "description": "No permission found with Id",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            }
        },
        "/record": {
            "get": {
                "description": "Retrieves a list of Makerchecker records",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "makerchecker"
                ],
                "summary": "Get all Makerchecker Record",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Makerchecker"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            },
            "put": {
                "description": "Update Makerchecker by approving or rejecting request",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "makerchecker"
                ],
                "summary": "Update Makerchecker by approving or rejecting request",
                "parameters": [
                    {
                        "description": "Request Body",
                        "name": "requestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/makerchecker.UpdateMakerchecker"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Makerchecker"
                        }
                    },
                    "400": {
                        "description": "Bad request due to invalid JSON body",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "403": {
                        "description": "User is not authorize to approve the request",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "404": {
                        "description": "No makerchecker record found with makercheckerId",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a Makerchecker",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "makerchecker"
                ],
                "summary": "Create a Makerchecker",
                "parameters": [
                    {
                        "description": "makerchecker",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Makerchecker"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Makerchecker"
                        }
                    },
                    "400": {
                        "description": "Bad request due to invalid JSON body",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "403": {
                        "description": "Not enough permission to do makerchecker",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "404": {
                        "description": "Unable to find resources",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            }
        },
        "/record/pending-approve": {
            "get": {
                "description": "Retrieves a list of pending approval records",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "makerchecker"
                ],
                "summary": "Get all Pending Approval as a Maker using MakerID",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Makerchecker"
                            }
                        }
                    },
                    "400": {
                        "description": "Maker Id cannot be found in the header provided",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "404": {
                        "description": "No pending requests found",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            }
        },
        "/record/to-approve": {
            "get": {
                "description": "Retrieves a list of pending approval records",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "makerchecker"
                ],
                "summary": "Get all Pending Approval as a Checker using CheckerID",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Makerchecker"
                            }
                        }
                    },
                    "400": {
                        "description": "Checker Id cannot be found in the header provided",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "404": {
                        "description": "No pending requests found",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            }
        },
        "/record/{id}": {
            "get": {
                "description": "Retrieve a Makerchecker By ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "makerchecker"
                ],
                "summary": "Get Makerchecker recordby Id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Makerchecker"
                        }
                    },
                    "400": {
                        "description": "Id cannot be empty",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "404": {
                        "description": "Record not found with Id",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            }
        },
        "/verify": {
            "post": {
                "description": "Verify if a User can do Makerchecker",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "makerchecker"
                ],
                "summary": "Verify if a User can do Makerchecker",
                "parameters": [
                    {
                        "description": "Request Body",
                        "name": "requestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/makerchecker.VerifyMaker"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad request due to invalid JSON body",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "403": {
                        "description": "Not enough permission to do makerchecker",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "404": {
                        "description": "Unable to find makerchecker permission",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "makerchecker.UpdateMakerchecker": {
            "type": "object",
            "required": [
                "id",
                "status"
            ],
            "properties": {
                "id": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "makerchecker.VerifyMaker": {
            "type": "object",
            "required": [
                "endpoint"
            ],
            "properties": {
                "endpoint": {
                    "type": "string"
                }
            }
        },
        "models.HttpError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "data": {
                    "type": "object",
                    "additionalProperties": true
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        },
        "models.Makerchecker": {
            "type": "object",
            "required": [
                "checkerId",
                "data",
                "endpoint"
            ],
            "properties": {
                "_id": {
                    "type": "string"
                },
                "checkerEmail": {
                    "type": "string"
                },
                "checkerId": {
                    "type": "string"
                },
                "data": {
                    "type": "object",
                    "additionalProperties": true
                },
                "endpoint": {
                    "type": "string"
                },
                "makerEmail": {
                    "type": "string"
                },
                "makerId": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "models.Permission": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string"
                },
                "checker": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "endpoint": {
                    "type": "string"
                },
                "maker": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                }
            }
        },
        "permission.UpdatePermission": {
            "type": "object",
            "properties": {
                "checker": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "endpoint": {
                    "type": "string"
                },
                "maker": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
