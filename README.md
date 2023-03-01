# slacksnitch
Github action that sends Dependabot vulnerability alerts to Slack.

## Description

This action retrieves Github Dependabot security vulnerability alerts from your repository and send them to Slack. It uses a Github token to retrieve teh security data via a GraphQL query, formats that data and sends it to the designated Slack channel, using a Slack API token.

## Usage
```
name: 'Check for Vulnerabilities'

on:
  schedule:
    - cron: '59 23 * * 0'

jobs:
  main:
    runs-on: ubuntu-latest
    steps:
      - uses: elaletovic/slacksnitch@main
        with:
          github_access_token: ${{ secrets.PERSONAL_ACCESS_TOKEN }}
          slack_access_token: ${{ secrets.SLACK_ACCESS_TOKEN }}
          slack_channel: security-alerts
          number_of_records: 10 
```
