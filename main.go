package main

import (
	"context"

	"github.com/elaletovic/slacksnitch/models"
	"github.com/sethvargo/go-githubactions"
	"github.com/shurcooL/githubv4"
	"github.com/slack-go/slack"
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
	slackClient := slack.New(slackAccessToken)

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

	var query models.VulnerabilityQuery

	err = githubClient.Query(context.Background(), &query, values)
	if err != nil {
		githubactions.Fatalf("failed to call a vulnerability query, error :%v", err)
	}

	githubactions.Infof("vulnerability query: %v", query)

	if query.Repository.VulnerabilityAlerts != nil {
		for _, edge := range query.Repository.VulnerabilityAlerts.Edges {
			githubactions.Infof("node ID %v, severity %v", edge.Node.ID, edge.Node.SecurityAdvisory.Severity)
		}
	}

	if query.Repository.VulnerabilityAlerts != nil && len(query.Repository.VulnerabilityAlerts.Edges) > 0 {

		blocks := query.GetMessage()
		a, b, c, err := slackClient.SendMessage("security-alerts", slack.MsgOptionBlocks(blocks.BlockSet...))
		githubactions.Infof("a: %s, b: %s, c: %s, err: %v", a, b, c, err)
	}
}
