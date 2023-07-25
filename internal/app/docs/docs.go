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
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/courses": {
            "get": {
                "description": "This API for getting all courses",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Course"
                ],
                "summary": "GetAllCourses",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.AllCourses"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    }
                }
            },
            "put": {
                "description": "This API for updating Course",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Course"
                ],
                "summary": "UpdateCourse",
                "parameters": [
                    {
                        "description": "CourseBody",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Course"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Course"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    }
                }
            },
            "post": {
                "description": "This API for creating course",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Course"
                ],
                "summary": "CreateCourse",
                "parameters": [
                    {
                        "description": "CourseBody",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Course"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Course"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    }
                }
            }
        },
        "/api/courses/{id}": {
            "get": {
                "description": "This API for getting Course",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Course"
                ],
                "summary": "GetCourse",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Course_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Course"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    }
                }
            },
            "delete": {
                "description": "This API for deleting Course",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Course"
                ],
                "summary": "DeleteCourse",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Course_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Empty"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    }
                }
            }
        },
        "/api/educenters": {
            "get": {
                "description": "This API for getting EduCenter",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "EduCenter"
                ],
                "summary": "Get Edu Center",
                "parameters": [
                    {
                        "type": "string",
                        "description": "EduCenter_ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.EduCenter"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    }
                }
            },
            "put": {
                "description": "This API for updating eduCenter",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "EduCenter"
                ],
                "summary": "Update EduCenter",
                "parameters": [
                    {
                        "description": "EduCenter",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.EduCenter"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.EduCenter"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    }
                }
            },
            "post": {
                "description": "This API for creating Edu Center",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "EduCenter"
                ],
                "summary": "Create Edu Center",
                "parameters": [
                    {
                        "description": "CourseBody",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.EduCenter"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.EduCenter"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    }
                }
            },
            "delete": {
                "description": "This API for deleting EduCenter",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "EduCenter"
                ],
                "summary": "Delete EduCenter",
                "parameters": [
                    {
                        "type": "string",
                        "description": "EduCenter_ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Empty"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.CustomError"
                        }
                    }
                }
            }
        },
        "/example/helloworld": {
            "get": {
                "description": "do ping",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "example"
                ],
                "summary": "ping example",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.AllCourses": {
            "type": "object",
            "properties": {
                "courses": {
                    "$ref": "#/definitions/models.Course"
                }
            }
        },
        "models.AllEduCenters": {
            "type": "object",
            "properties": {
                "edu_centers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.EduCenter"
                    }
                }
            }
        },
        "models.Course": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "edu_center_id": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "teacher": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.CustomError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "statusCode": {
                    "type": "integer"
                }
            }
        },
        "models.EduCenter": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "address": {
                    "type": "string"
                },
                "cover_image": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "html_description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "location": {
                    "$ref": "#/definitions/models.Point"
                },
                "name": {
                    "type": "string"
                },
                "owner_id": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.Empty": {
            "type": "object"
        },
        "models.Point": {
            "type": "object",
            "properties": {
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
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
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}