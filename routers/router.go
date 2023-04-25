package routers

import (
  "five_letters/controllers"
  "five_letters/middlewares"
  "github.com/beego/beego/v2/server/web"
  "github.com/beego/beego/v2/server/web/filter/cors"
)

type route struct {
  path       string
  controller web.ControllerInterface
  methods    string
  auth       bool
}

func getRoutes() []route {
  return []route{
    {
      path:       "/api/v1/languages",
      controller: &controllers.LanguagesController{},
      methods:    "get:Index",
      auth:       false,
    },
    {
      path:       "/api/v1/profile",
      controller: &controllers.UsersController{},
      methods:    "get:Get;put,patch:Update",
      auth:       true,
    },
    {
      path:       "/api/v1/games",
      controller: &controllers.GamesController{},
      methods:    "get:Index;post:Create",
      auth:       true,
    },
    {
      path:       "/api/v1/games/active",
      controller: &controllers.GamesController{},
      methods:    "get:Active",
      auth:       true,
    },
    {
      path:       "/api/v1/games/:id",
      controller: &controllers.GamesController{},
      methods:    "get:Read",
      auth:       true,
    },
    {
      path:       "/api/v1/games/active/attempts",
      controller: &controllers.GamesController{},
      methods:    "post:CreateAttempt",
      auth:       true,
    },
    {
      path:       "/api/v1/auth/register",
      controller: &controllers.AuthController{},
      methods:    "post:Register",
      auth:       false,
    },
    {
      path:       "/api/v1/auth/login",
      controller: &controllers.AuthController{},
      methods:    "post:Login",
      auth:       false,
    },
    {
      path:       "/api/v1/auth/refresh",
      controller: &controllers.AuthController{},
      methods:    "post:Refresh",
      auth:       false,
    },
    {
      path:       "/api/v1/auth/logout",
      controller: &controllers.AuthController{},
      methods:    "post:Logout",
      auth:       true,
    },
  }
}

func init() {
  initRoutes()
  initErrorHandlers()
  initCors()
  initDocs()
}

func initRoutes() {
  for _, r := range getRoutes() {
    web.Router(r.path, r.controller, r.methods)
    if r.auth {
      web.InsertFilter(r.path, web.BeforeRouter, middlewares.Authenticate)
    }
  }
}

func initErrorHandlers() {
  web.ErrorController(&controllers.ErrorsController{})
  web.BConfig.RecoverFunc = middlewares.RecoverFromPanic
}

func initCors() {
  web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
    AllowOrigins:     []string{},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "User-Agent"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
  }))
}

func initDocs() {
  if web.BConfig.RunMode != "prod" {
    web.BConfig.WebConfig.DirectoryIndex = true
    web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
  }
}
