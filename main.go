package main

import (
	"context"

	"github.com/elaletovic/slacksnitch/queries"
	"github.com/sethvargo/go-githubactions"
	"github.com/shurcooL/githubv4"
	_ "github.com/slack-go/slack"
	"golang.org/x/oauth2"
)

func main() {
	// get and check inputs
	githubAcessToken := githubactions.GetInput("github_access_token")
	if githubAcessToken == "" {
		githubactions.Fatalf("missing input 'github_access_token'")
	}

	slackAccessToken := githubactions.GetInput("slack_access_token")
	if slackAccessToken == "" {
		githubactions.Fatalf("missing input 'slack_access_token'")
	}

	slackChannel := githubactions.GetInput("slack_channel")
	if slackChannel == "" {
		githubactions.Fatalf("missing input 'slack_channel'")
	}

	githubactions.Infof("github access token: %v", githubAcessToken)
	githubactions.Infof("slack access token: %v", slackAccessToken)
	githubactions.Infof("slack channel: %v", slackChannel)

	// slack client
	//slackClient := slack.New(slackAccessToken)

	//get GH client
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAcessToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	githubClient := githubv4.NewClient(httpClient)

	gCtx, err := githubactions.Context()
	if err != nil {
		githubactions.Fatalf("could not get context for this action, error: %v", err)
	}

	owner, repository := gCtx.Repo()
	values := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repository),
		"count": githubv4.Int(30),
	}

	err = githubClient.Query(context.Background(), &queries.VulnerabilityQuery, values)
	if err != nil {
		githubactions.Fatalf("failed to call a vulnerability query, error :%v", err)
	}

	githubactions.Infof("vulnerability query: %v", queries.VulnerabilityQuery)

	/*
	   payload, err := json.Marshal(&queries.VulnerabilityQuery)

	   	if err != nil {
	   		githubactions.Fatalf("failed to marshal the vulnerability query, error :%v", err)
	   	} else {

	   		a, b, c, err := slackClient.SendMessage("security-alerts", slack.MsgOptionText(string(payload), true))
	   		githubactions.Infof("a: %s, b: %s, c: %s, err: %v", a, b, c, err)
	   	}
	*/
}
