//tests will be reworked after this proof of concept plugin
package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/franela/goblin"
)

func TestHookImage(t *testing.T) {

	g := goblin.Goblin(t)

	g.Describe("Make a request object for a request to cowpoke", func() {
		g.It("should return the correct request", func() {
			catalogNo := 1
			branchName := "test"
			CatalogRepo := "repo"
			rancherCatalogName := "catalog"
			token := "secret"
			CowpokeURL := "cowpoke.mydomain.io"
			BearerToken := "token"
			var args map[string]interface{}
			req, err := cowpokeRequest(catalogNo, branchName, CatalogRepo, rancherCatalogName, token, CowpokeURL, BearerToken)
			body, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(body, &args)
			g.Assert(err).Equal(nil)
			g.Assert(req.Header.Get("bearer")).Equal(BearerToken)
			g.Assert(req.Header.Get("Content-Type")).Equal("application/json")
			g.Assert(args["catalog"].(string)).Equal(CatalogRepo)
			g.Assert(args["rancherCatalogName"].(string)).Equal(rancherCatalogName)
			g.Assert(args["githubToken"].(string)).Equal(token)
			g.Assert(args["catalogVersion"].(string)).Equal("1")
			g.Assert(args["branch"].(string)).Equal(branchName)
		})
	})
}
