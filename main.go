package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"
  "os"
  "path/filepath"
  "github.com/drone/drone-plugin-go/plugin"
)

type Cowpoke struct {
  Url string `json:"cowpoke_url"`
  Port int `json:"cowpoke_port"`
}

type ImageJson struct {
  Image string `json:"image"`
}

func main() {
  fmt.Println("starting drone-cowpoke...")

  workspace := plugin.Workspace{}
  vargs := Cowpoke{}

  plugin.Param("workspace", &workspace)
  plugin.Param("vargs", &vargs)
  plugin.MustParse()

  if len(vargs.Url) == 0 {
    fmt.Println("no cowpoke url was specified")
    os.Exit(1)
  }

  if vargs.Port == 0 {
    fmt.Println("no cowpoke port was specified")
    os.Exit(1)
  }

  fmt.Println("loading image data from", filepath.Join(workspace.Path, ".docker.json"))
  image := GetImageName(filepath.Join(workspace.Path, ".docker.json"))

  if(len(image) <= 0) {
    fmt.Println("image load failed from .docker.json")
    os.Exit(1)
  }

  var cowpokeUrl = fmt.Sprintf("%s:%d/api/environment/", vargs.Url, vargs.Port)
  fmt.Println("cowpoke url set to:", cowpokeUrl)
  fmt.Println(".docker.json value being posted:", image)
  ExecutePut(cowpokeUrl + url.QueryEscape(image));

  fmt.Println("finished drone-cowpoke.")
}

func ExecutePut(putUrl string) {
  fmt.Println("executing a PUT request for:", putUrl)

  client := &http.Client{}
  request, err := http.NewRequest("PUT", putUrl, nil)

  response, err := client.Do(request)
  if err != nil {
    fmt.Println("error executing request:", err)
    os.Exit(1)
  }

  defer response.Body.Close()
    contents, err := ioutil.ReadAll(response.Body)
    if err != nil {
      fmt.Println("error reading response:", err)
    }

    fmt.Println("response status code:", response.StatusCode)
    fmt.Println("content:", string(contents))
}

func GetImageName(path string) string {
  file, err := ioutil.ReadFile(path)

  if err != nil {
    fmt.Println("error opening json file", err)
  }

  var jsonobject ImageJson
  json.Unmarshal(file, &jsonobject)

  return jsonobject.Image
}
