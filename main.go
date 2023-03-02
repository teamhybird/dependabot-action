package main

import (
	"context"
	"strconv"

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

	numberOfItems := 20
	count := githubactions.GetInput("number_of_records")
	if count != "" {
		val, err := strconv.Atoi(count)
		if err != nil || val <= 0 {
			githubactions.Fatalf("'number_of_records' must be an integer bigger than zero")
		} else {
			numberOfItems = val
		}
	}

	githubactions.Infof("github access token: %v", githubAcessToken)
	githubactions.Infof("slack access token: %v", slackAccessToken)
	githubactions.Infof("slack channel: %v", slackChannel)
	githubactions.Infof("number of items: %v", numberOfItems)

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
		"count": githubv4.Int(numberOfItems),
	}

	var query models.VulnerabilityQuery

	err = githubClient.Query(context.Background(), &query, values)
	if err != nil {
		githubactions.Fatalf("failed to call a vulnerability query, error :%v", err)
	}

	githubactions.Infof("vulnerability query: %v", query)

	if query.Repository.VulnerabilityAlerts != nil && len(query.Repository.VulnerabilityAlerts.Edges) > 0 {

		blocks := query.GetMessage()
		_, _, _, err := slackClient.SendMessage(slackChannel, slack.MsgOptionBlocks(blocks...))
		if err != nil {
			githubactions.Fatalf("failed to send Slack message, error: %v", err)
		}
	}
}
