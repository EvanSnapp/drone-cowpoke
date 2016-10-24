package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/LeanKit-Labs/drone-cowpoke/types"
	dronePlugin "github.com/drone/drone-go/plugin"
	yaml "gopkg.in/yaml.v2"
)

var catalogsFile = "/drone/.CatalogData.yml"

func main() {
	fmt.Println("starting drone-cowpoke")
	/*
	   Drone pkg types are abstracted into "plugin" in order
	   to make the migration to Drone's 0.5 way of getting
	   plugin args easier (i.e. via env vars)
	*/

	if _, err := os.Stat(catalogsFile); os.IsNotExist(err) {
		fmt.Println("Catalog Info file not found")
		os.Exit(0)
	}

	plugin := types.Plugin{}

	dronePlugin.Param("vargs", &plugin)
	dronePlugin.MustParse()

	if err := exec(&plugin); err != nil {
		fmt.Println("There was an error :(")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("plugin completed, exiting")
	os.Exit(0)
}

func exec(p *types.Plugin) error {
	var catalogs []types.CatalogInfo
	rawData, err := ioutil.ReadFile(catalogsFile)
	if err != nil {
		fmt.Println("Error Reading File")
		return err
	}
	yaml.Unmarshal(rawData, &catalogs)

	for _, catalog := range catalogs {
		request, err := cowpokeRequest(catalog.Version, catalog.Branch, p.CatalogRepo, p.RancherCatalogName, p.GitHubToken, p.CowpokeURL, p.BearerToken)
		if err != nil {
			fmt.Println("Error building request")
			return err
		}
		client := http.Client{
			Timeout: time.Second * 60,
		}
		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("error sending request: %s\n", request.URL)
			return err
		}
		contents, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			fmt.Println("Error reading response")
			return err
		}
		fmt.Println("response status code:", response.StatusCode)
		fmt.Println("content:", string(contents))
	}
	return nil
}

//calls cowpoke after catalog is built
func cowpokeRequest(catalogNo int, branchName string, CatalogRepo string, rancherCatalogName string, token string, CowpokeURL string, BearerToken string) (*http.Request, error) {
	var jsonStr = []byte(fmt.Sprintf(`{"catalog":"%s","rancherCatalogName":"%s","githubToken":"%s","catalogVersion":"%s","branch":"%s"}`, CatalogRepo, rancherCatalogName, token, strconv.Itoa(catalogNo), branchName))
	request, err := http.NewRequest("PATCH", CowpokeURL+"/api/stack", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("bearer", BearerToken)
	if err != nil {
		fmt.Println("Error making request object to cowpoke")
		return nil, err
	}
	request.Close = true
	return request, nil
}
