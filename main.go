package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v29/github"
)

func main() {
	owner, repo, err := getRepo()
	if err != nil {
		log.Fatal(err)
	}

	appID, err := strconv.ParseInt(os.Getenv("GH_APP_ID"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	key, err := getPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	tr := http.DefaultTransport
	itr, err := ghinstallation.NewAppsTransport(tr, appID, key)
	if err != nil {
		log.Fatal(err)
	}

	client := github.NewClient(&http.Client{Transport: itr})
	install, _, err := client.Apps.FindRepositoryInstallation(context.Background(), owner, repo)
	if err != nil {
		log.Fatal(err)
	}

	opts := &github.InstallationTokenOptions{}
	token, _, err := client.Apps.CreateInstallationToken(context.Background(), *install.ID, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(*token.Token)
}

func getRepo() (owner string, repo string, err error) {
	repoSlug, found := os.LookupEnv("TRAVIS_REPO_SLUG")
	if !found {
		repoSlug, found = os.LookupEnv("GITHUB_REPOSITORY")
	}
	if !found {
		err = fmt.Errorf("TRAVIS_REPO_SLUG or GITHUB_REPOSITORY not set")
		return
	}

	slugParts := strings.Split(repoSlug, "/")
	if len(slugParts) != 2 {
		err = fmt.Errorf("Could not get owner and repo")
	}

	owner = slugParts[0]
	repo = slugParts[1]
	return
}

func getPrivateKey() ([]byte, error) {
	encodedKey := os.Getenv("GH_PRIVATE_KEY")
	return base64.StdEncoding.DecodeString(encodedKey)
}
