package middlewares

import (
  "five_letters/controllers"
  "fmt"
  "github.com/beego/beego/v2/core/logs"
  "github.com/beego/beego/v2/server/web"
  "github.com/beego/beego/v2/server/web/context"
  "reflect"
  "runtime"
  "strconv"
)

func RecoverFromPanic(ctx *context.Context, cfg *web.Config) {
  if err := recover(); err != nil {

    ec := &controllers.ErrorsController{}

    if err == web.ErrAbort {
      return
    }

    if !cfg.RecoverPanic {
      panic(err)
    }

    if cfg.EnableErrorsShow {
      if _, ok := web.ErrorMaps[fmt.Sprint(err)]; ok {
        statusCode, _ := strconv.Atoi(fmt.Sprint(err))

        methodName := fmt.Sprintf("Error%d", statusCode)
        method := reflect.ValueOf(ec).MethodByName(methodName)

        if method.IsValid() {
          if ec.Data == nil {
            ec.Data = make(map[interface{}]interface{})
          }
          ec.Ctx = ctx
          method.Call([]reflect.Value{})
          return
        }
      }
    }

    logs.Critical("the request url is ", ctx.Input.URL())
    logs.Critical("Handler crashed with error", err)

    var stack string

    for i := 1; ; i++ {
      _, file, line, ok := runtime.Caller(i)
      if !ok {
        break
      }
      logs.Critical(fmt.Sprintf("%s:%d", file, line))
      stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
    }

    if ec.Data == nil {
      ec.Data = make(map[interface{}]interface{})
    }
    ec.Ctx = ctx
    ec.Error500()
  }
}
