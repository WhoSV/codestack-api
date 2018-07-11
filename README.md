# Codestack-Api

<!-- v1.2.0 -->

[![GitHub version](https://badge.fury.io/gh/WhoSV%2Fcodestack-api.svg)](https://badge.fury.io/gh/WhoSV%2Fcodestack-api)
[![TeamCity CodeBetter](https://img.shields.io/teamcity/codebetter/bt428.svg)](https://github.com/WhoSV/codestack-api)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](ttps://github.com/WhoSV/codestack-api)

A REST-ful API for [CodeStack](https://github.com/WhoSV/codestack) application with Go (golang)

![alternativetext](screenshot.png)

## Quick Start

**Get dependencies**

`$ cd codestack-api/`

`$ go get`

**Run**

`$ go run main.go`

**Browse**

`http://localhost:8000`

## Structure

```
├── data
├── database              // Database
│   └──db.go
├── endpoints             // Endpoints
│   ├── handlers          // API core handlers
│   │   ├── course.go
│   │   ├── favorite.go
│   │   ├── survey.go
│   │   └── user.go
│   ├── auth.go
│   └── middleware.go
├── errors                // Errors
│   └──error.go
├── model                 // Models for our application
│   ├── course.go
│   ├── favorite.go
│   ├── survey.go
│   └── user.go
├── repository            // Repository
│   ├── user.go
├── router                // Routes
│   ├── router.go
├── server                // Server
│   ├── server.go
└── main.go
```

## API

#### /people

- `GET` : Get all users
- `POST` : Create a new user

#### /people/{id}

- `GET` : Get user
- `PUT` `PATCH` : Update user
- `DELETE` : Delete user

#### /people/{id}/update

- `PUT` `PATCH` : Update user password

#### /people/{id}/reset

- `POST` : Reset user password

#### /favorite

- `GET` : Get favorite
- `POST` : Create favorite

#### /favorite{id}

- `DELETE` : Detele favorite

##### /courses

- `GET` : Get all courses
- `POST` : Create a new course

#### /courses/{id}

- `GET` : Get course
- `PUT` `PATCH` : Update course
- `DELETE` : Delete course

#### /courses/{id}/status

- `PUT` `PATCH` : Update course status

#### /courses/{id}/open

- `GET` : Open course

#### /survey

- `GET` : Get surveys
- `POST` : Create a new survey

## Frontend for this Application

[CodeStack](https://github.com/WhoSV/codestack)
