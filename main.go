package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ghAuth struct {
	AuthToken string
	Owner     string
	Repo      string
}

func NewGithubAuth(authToken, owner, repo string) ghAuth {
	return ghAuth{
		AuthToken: authToken,
		Owner:     owner,
		Repo:      repo,
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("error loading .env file: %w", err))
	}

	authToken := os.Getenv("GITHUB_ACCESS_TOKEN") // Your Github PAT token
	OWNER := os.Getenv("GITHUB_USERNAME")         // cosmos
	REPO := os.Getenv("GITHUB_REPOSITORY")        // cosmos-sdk
	gh := NewGithubAuth(authToken, OWNER, REPO)

	login, _ := gh.loginExample()
	fmt.Println("[Login]:", login.Data.Viewer.User)

	d, _ := gh.queryDiscussionByID(19391)
	fmt.Println("[Discussion]:", d.Data.Repository.Discussion.Comments)
}

// loginExample returns the login of the authenticated user via the AuthToken
// and returns the username.
func (ga ghAuth) loginExample() (Login, error) {
	return MakeQuery(ga.AuthToken, `query { viewer { login }}`, Login{})
}

// queryDiscussionByID returns the discussion and comments for a given discussion by its ID.
// This ID is unique, Reminder: issues, prs, and discussions all share the same auto-incrementing ID space.
func (ga ghAuth) queryDiscussionByID(id uint) (DiscussionComments, error) {
	req := fmt.Sprintf(`query { repository(owner: "%s", name: "%s") { discussion(number: %d) { body comments(first: 100){edges{node{author{login} body}}  }} } }`, ga.Owner, ga.Repo, id)
	return MakeQuery(ga.AuthToken, req, DiscussionComments{})
}

// MakeQuery makes a query request and handles the response.
func MakeQuery[T any](authToken, query string, t T) (T, error) {
	var ret T
	b, err := makeReq(authToken, query)
	if err != nil {
		return ret, err
	}
	if err := json.Unmarshal(b, &ret); err != nil {
		return ret, err
	}
	return ret, nil
}

// makeReq, you could probably use the https://github.com/graphql-go/graphql package for this
func makeReq(authToken, rawQuery string) ([]byte, error) {
	client := &http.Client{}

	rawQuery = strings.ReplaceAll(rawQuery, "\"", "\\\"")
	query := fmt.Sprintf("{\"query\": \"%s\"}", rawQuery)
	fmt.Println("[GraphQL Query]:", query)

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", strings.NewReader(query))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "bearer "+authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
