// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/activities": {
            "get": {
                "description": "Извлачење листе активности са опционим филтрирањем и страничењем",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "activities"
                ],
                "summary": "Листа активности",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Page size",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Activity type filter",
                        "name": "type",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Activity level filter",
                        "name": "level",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityListResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "post": {
                "description": "Креирање нове активности у систему",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "activities"
                ],
                "summary": "Креирање нове активности",
                "parameters": [
                    {
                        "description": "Activity data",
                        "name": "activity",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/activities/reset": {
            "delete": {
                "description": "Брисање свих активности из система",
                "tags": [
                    "activities"
                ],
                "summary": "Ресетовање свих података о активностима",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/activities/stats": {
            "get": {
                "description": "Преузимање свеобухватних статистика активности",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "activities"
                ],
                "summary": "Статистике активности",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityStatsResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/activities/{id}": {
            "get": {
                "description": "Преузимање одређене активности по њеном ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "activities"
                ],
                "summary": "Преузимање активности по ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Activity ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "delete": {
                "description": "Брисање одређене активности по њеном ID",
                "tags": [
                    "activities"
                ],
                "summary": "Брисање активности",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Activity ID",
                        "name": "id",
                        "in": "path",
                        "required": true
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
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityCreateRequest": {
            "type": "object",
            "required": [
                "description",
                "level",
                "title",
                "type"
            ],
            "properties": {
                "actorId": {
                    "type": "integer"
                },
                "actorName": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "level": {
                    "$ref": "#/definitions/github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityLevel"
                },
                "metadata": {
                    "type": "string"
                },
                "targetId": {
                    "type": "integer"
                },
                "targetType": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "type": {
                    "$ref": "#/definitions/github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityType"
                }
            }
        },
        "github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityLevel": {
            "type": "string",
            "enum": [
                "info",
                "warning",
                "error",
                "critical"
            ],
            "x-enum-varnames": [
                "ActivityLevelInfo",
                "ActivityLevelWarning",
                "ActivityLevelError",
                "ActivityLevelCritical"
            ]
        },
        "github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityListResponse": {
            "type": "object",
            "properties": {
                "activities": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityResponse"
                    }
                },
                "page": {
                    "type": "integer"
                },
                "pageSize": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                },
                "totalPages": {
                    "type": "integer"
                }
            }
        },
        "github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityResponse": {
            "type": "object",
            "properties": {
                "actorId": {
                    "description": "ID of the user who performed the action",
                    "type": "integer"
                },
                "actorName": {
                    "description": "Name of the user who performed the action",
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "level": {
                    "$ref": "#/definitions/github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityLevel"
                },
                "metadata": {
                    "description": "JSON string with additional data",
                    "type": "string"
                },
                "targetId": {
                    "description": "ID of the target entity",
                    "type": "integer"
                },
                "targetType": {
                    "description": "Type of the target entity (employee, urgency, etc.)",
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "type": {
                    "$ref": "#/definitions/github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityType"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityStatsResponse": {
            "type": "object",
            "properties": {
                "activitiesByLevel": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "integer",
                        "format": "int64"
                    }
                },
                "activitiesByType": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "integer",
                        "format": "int64"
                    }
                },
                "activitiesLast24h": {
                    "type": "integer"
                },
                "activitiesLast30Days": {
                    "type": "integer"
                },
                "activitiesLast7Days": {
                    "type": "integer"
                },
                "recentActivities": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityResponse"
                    }
                },
                "totalActivities": {
                    "type": "integer"
                }
            }
        },
        "github_com_pd120424d_mountain-service_api_contracts_activity_v1.ActivityType": {
            "type": "string",
            "enum": [
                "employee_created",
                "employee_updated",
                "employee_deleted",
                "employee_login",
                "shift_assigned",
                "shift_removed",
                "urgency_created",
                "urgency_updated",
                "urgency_deleted",
                "emergency_assigned",
                "emergency_accepted",
                "emergency_declined",
                "notification_sent",
                "notification_failed",
                "system_reset"
            ],
            "x-enum-varnames": [
                "ActivityEmployeeCreated",
                "ActivityEmployeeUpdated",
                "ActivityEmployeeDeleted",
                "ActivityEmployeeLogin",
                "ActivityShiftAssigned",
                "ActivityShiftRemoved",
                "ActivityUrgencyCreated",
                "ActivityUrgencyUpdated",
                "ActivityUrgencyDeleted",
                "ActivityEmergencyAssigned",
                "ActivityEmergencyAccepted",
                "ActivityEmergencyDeclined",
                "ActivityNotificationSent",
                "ActivityNotificationFailed",
                "ActivitySystemReset"
            ]
        }
    },
    "securityDefinitions": {
        "OAuth2Password": {
            "type": "oauth2",
            "flow": "password",
            "tokenUrl": "/api/v1/oauth/token",
            "scopes": {
                "read": "Grants read access",
                "write": "Grants write access"
            }
        }
    },
    "security": [
        {
            "OAuth2Password": []
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8084",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Activity Service API",
	Description:      "Activity tracking and audit service for the Mountain Emergency Management System",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
