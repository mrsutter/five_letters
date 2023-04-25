package testutils

import (
  "bytes"
  _ "five_letters/routers"
  "github.com/beego/beego/v2/server/web"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
)

func HTTPGet(path string) (int, []byte, http.Header) {
  return HttpRequest("GET", path, nil, "")
}

func HTTPGetWithToken(path string, token string) (int, []byte, http.Header) {
  return HttpRequest("GET", path, nil, token)
}

func HTTPPost(path string, body []byte) (int, []byte, http.Header) {
  return HttpRequest("POST", path, body, "")
}

func HTTPPostWithToken(path string, body []byte, token string) (int, []byte, http.Header) {
  return HttpRequest("POST", path, body, token)
}

func HTTPPut(path string, body []byte) (int, []byte, http.Header) {
  return HttpRequest("PUT", path, body, "")
}

func HTTPPutWithToken(path string, body []byte, token string) (int, []byte, http.Header) {
  return HttpRequest("PUT", path, body, token)
}

func HTTPPatch(path string, body []byte) (int, []byte, http.Header) {
  return HttpRequest("PATCH", path, body, "")
}

func HTTPPatchWithToken(path string, body []byte, token string) (int, []byte, http.Header) {
  return HttpRequest("PATCH", path, body, token)
}

func HTTPDelete(path string) (int, []byte, http.Header) {
  return HttpRequest("DELETE", path, nil, "")
}

func HTTPDeleteWithToken(path string, token string) (int, []byte, http.Header) {
  return HttpRequest("DELETE", path, nil, token)
}

func HttpRequest(method string, path string, body []byte, token string) (
  int, []byte, http.Header) {

  request, _ := http.NewRequest(method, path, bytes.NewBuffer(body))

  if token != "" {
    request.Header.Set("Authorization", "Bearer "+token)
  }

  response := httptest.NewRecorder()

  web.BeeApp.Handlers.ServeHTTP(response, request)

  bodyRaw, _ := ioutil.ReadAll(response.Body)
  return response.Code, bodyRaw, response.Header()
}
