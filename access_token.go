package gh

import "os"

func AccessToken() string {
	accessToken := os.Getenv("GITHUB_TOKEN")
	if accessToken == "" {
		panic("`GITHUB_TOKEN' environment variable is not set")
	}
	return accessToken
}
