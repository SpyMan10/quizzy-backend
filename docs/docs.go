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
        "/ping": {
            "get": {
                "description": "Vérifie si Firebase et Redis sont accessibles et retourne leur état",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "HealthCheck"
                ],
                "summary": "Vérifier la disponibilité des services",
                "responses": {
                    "200": {
                        "description": "Statut des services",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/quiz": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retourne la liste des quiz créés par l'utilisateur authentifié",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Quizzes"
                ],
                "summary": "Récupérer tous mes quiz",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cvotre_token\u003e",
                        "description": "Token d'authentification Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Liste des quiz de l'utilisateur",
                        "schema": {
                            "$ref": "#/definitions/quizzes.UserQuizzesResponse"
                        }
                    },
                    "401": {
                        "description": "Utilisateur non authentifié",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Erreur interne du serveur",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Permet à l'utilisateur authentifié de créer un quiz",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Quizzes"
                ],
                "summary": "Créer un quiz",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cvotre_token\u003e",
                        "description": "Token d'authentification Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Informations du quiz à créer",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/quizzes.CreateQuizRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Quiz créé avec succès",
                        "schema": {
                            "$ref": "#/definitions/quizzes.Quiz"
                        }
                    },
                    "400": {
                        "description": "Requête invalide",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Utilisateur non authentifié",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Erreur interne du serveur",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/quiz/{quiz-id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retourne les informations d'un quiz appartenant à l'utilisateur authentifié",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Quizzes"
                ],
                "summary": "Récupérer un quiz",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cvotre_token\u003e",
                        "description": "Token d'authentification Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID du quiz",
                        "name": "quiz-id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Détails du quiz",
                        "schema": {
                            "$ref": "#/definitions/quizzes.Quiz"
                        }
                    },
                    "401": {
                        "description": "Utilisateur non authentifié",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Quiz non trouvé",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Erreur interne du serveur",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Met à jour un quiz existant en fonction des champs envoyés",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Quizzes"
                ],
                "summary": "Modifier un quiz",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cvotre_token\u003e",
                        "description": "Token d'authentification Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID du quiz",
                        "name": "quiz-id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Champs à modifier",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/quizzes.FieldPatchOp"
                            }
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Quiz mis à jour avec succès",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Requête invalide",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Utilisateur non authentifié",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Erreur interne du serveur",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/quiz/{quiz-id}/questions": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retourne toutes les questions du quiz spécifié par son ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Quizzes"
                ],
                "summary": "Récupérer les questions d'un quiz",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cvotre_token\u003e",
                        "description": "Token d'authentification Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID du quiz",
                        "name": "quiz-id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Liste des questions du quiz",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/quizzes.Question"
                            }
                        }
                    },
                    "401": {
                        "description": "Utilisateur non authentifié",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Quiz non trouvé",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Erreur interne du serveur",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Permet d'ajouter une question à un quiz existant",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Quizzes"
                ],
                "summary": "Ajouter une question",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cvotre_token\u003e",
                        "description": "Token d'authentification Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID du quiz",
                        "name": "quiz-id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Détails de la question",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/quizzes.CreateQuestionRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Question ajoutée avec succès",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Requête invalide",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Utilisateur non authentifié",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Erreur interne du serveur",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/quiz/{quiz-id}/questions/{question-id}": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Met à jour une question spécifique d'un quiz",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Quizzes"
                ],
                "summary": "Modifier une question",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cvotre_token\u003e",
                        "description": "Token d'authentification Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID du quiz",
                        "name": "quiz-id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID de la question",
                        "name": "question-id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Mise à jour de la question",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/quizzes.UpdateQuestionRequest"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Question mise à jour avec succès",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Requête invalide",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Utilisateur non authentifié",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Quiz ou question non trouvée",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Erreur interne du serveur",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/quiz/{quiz-id}/start": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Démarre un quiz et retourne son code d'exécution",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Quizzes"
                ],
                "summary": "Démarrer un quiz",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cvotre_token\u003e",
                        "description": "Token d'authentification Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID du quiz",
                        "name": "quiz-id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Quiz démarré avec succès",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Quiz non prêt à être démarré",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Utilisateur non authentifié",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Quiz non trouvé",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Erreur interne du serveur",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Cette route permet de créer un nouvel utilisateur à partir d'un username et de l'email récupéré via l'authentification",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Créer un utilisateur",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cvotre_token\u003e",
                        "description": "Token d'authentification Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Informations de l'utilisateur à créer",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.createUserRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Utilisateur créé avec succès",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Requête invalide",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Erreur interne du serveur",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users/me": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Cette route permet d'obtenir les informations du compte actuellement authentifié",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Récupérer les informations de l'utilisateur connecté",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cvotre_token\u003e",
                        "description": "Token d'authentification Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Informations de l'utilisateur",
                        "schema": {
                            "$ref": "#/definitions/users.User"
                        }
                    },
                    "401": {
                        "description": "Utilisateur non authentifié",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Erreur interne du serveur",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "quizzes.Answer": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "isCorrect": {
                    "type": "boolean"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "quizzes.CreateQuestionRequest": {
            "type": "object",
            "properties": {
                "answers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/quizzes.Answer"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "quizzes.CreateQuizRequest": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "quizzes.FieldPatchOp": {
            "type": "object",
            "properties": {
                "op": {
                    "type": "string"
                },
                "path": {
                    "type": "string"
                },
                "value": {}
            }
        },
        "quizzes.Links": {
            "type": "object",
            "properties": {
                "create": {
                    "type": "string"
                },
                "start": {
                    "type": "string"
                }
            }
        },
        "quizzes.Question": {
            "type": "object",
            "properties": {
                "answers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/quizzes.Answer"
                    }
                },
                "id": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "quizzes.Quiz": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "questions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/quizzes.Question"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "quizzes.QuizWithLinks": {
            "type": "object",
            "properties": {
                "_links": {
                    "$ref": "#/definitions/quizzes.Links"
                },
                "code": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "questions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/quizzes.Question"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "quizzes.UnidentifiedAnswer": {
            "type": "object",
            "properties": {
                "isCorrect": {
                    "type": "boolean"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "quizzes.UpdateQuestionRequest": {
            "type": "object",
            "properties": {
                "answers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/quizzes.UnidentifiedAnswer"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "quizzes.UserQuizzesResponse": {
            "type": "object",
            "properties": {
                "_links": {
                    "$ref": "#/definitions/quizzes.Links"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/quizzes.QuizWithLinks"
                    }
                }
            }
        },
        "users.User": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "uid": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "users.createUserRequest": {
            "type": "object",
            "properties": {
                "username": {
                    "type": "string"
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
