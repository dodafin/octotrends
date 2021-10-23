package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const JSONFILENAME = "data/repo-info.json"

var GHP = os.Getenv("GHP")

func GetGHRepoInfo(data DataTable) map[string]github.Repository {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GHP},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	jsonMap := make(map[string]github.Repository)
	log.Println("Loaded JSON Map", len(jsonMap))
	for repoName := range data {
		if _, ok := jsonMap[repoName]; ok {
			// We have this info, skip it
			log.Printf("Skipping %s", repoName)
			continue
		}
		log.Println("Getting", repoName)
		spluts := strings.Split(repoName, "/")
		owner, name := spluts[0], spluts[1]
		var repo *github.Repository
		var err error
		for {
			repo, _, err = client.Repositories.Get(ctx, owner, name)
			if _, ok := err.(*github.RateLimitError); ok {
				log.Println("Hit rate limit, sleeping 1 minute")
				time.Sleep(time.Minute)
			} else {
				break
			}
		}
		if repo == nil {
			log.Println("Repo is nil:", repoName)
			jsonMap[repoName] = github.Repository{}
			continue
		}

		jsonMap[repoName] = *repo
	}
	return jsonMap
}
