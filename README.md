# Jenkins-Slack Bot

## Table of contents

- [Overview](#overview)
- [Requirements](#requirements)
- [Setup](#setup)
- [Run the app](#running-the-app)
- [Pipeline Example](#pipeline-example)
- [Contribution](#contribution)
- [TODO](#todo)

### Overview

[![Coverage Status](https://coveralls.io/repos/github/mgufrone/jenkins-bot-go/badge.svg)](https://coveralls.io/github/mgufrone/jenkins-bot-go)

Here's the diagram of the flow on what this repo is about
![flow](https://www.planttext.com/api/plantuml/svg/DOun4i8m30HxlK8ym4DFmIIKj7c14ZGXn9Q2fOJlOpEcrPxshZEdx7kASF8d9qRwMD1CCZF0dMLTn31SyQP-mO7bWRHjMG-AcFczipaKL1D3f6bjcSHcD3EwejKp_226lwYVV551FbZQVo6jhIdwLdEcKRDNKzc7Bnq1kiBBuYy0)

here is also the working demonstration on how this repository should work
[![Demonstration](https://img.youtube.com/vi/kbix7WRzgLI/0.jpg)](https://youtu.be/kbix7WRzgLI)

### Requirements

- Slack app with socket mode enabled
- Go >=1.20

### Setup

#### Create Slack App

You have options regarding Slack app. You can create a new app dedicated for jenkins-bot-go and slack jenkins plugin. Or
use existing app. Install it to your slack workspace.

Get the App Token in `Settings -> Basic Information` section, scroll down and you will see App Token. Create a new and
add these scopes:

- `connections:write`
- `authorization:read`

![Screenshot](./screenshots/1-slack-app-token.png?raw=true "Slack App Token")

Next, you need to obtain the Bot Token. Go to `Features -> OAuth & Permissions`.

Go to Scopes section and add these OAuth Scope:

- `chat:write`
- `im:history`

Scroll up and copy the bot token.
![Screenshot](./screenshots/2-slack-bot-token.png?raw=true "Slack bot token")

#### Setup jenkins-bot-go

- Clone the repository
- run this to download the dependencies

```shell
go mod download
```

- Copy `.env.example` and modify them according to your configuration
    - for `JENKINS_URL`, you need to put your jenkins URL. for example: `http://jenkins.localhost`
    - for `JENKINS_USERNAME`, you can use the current username (`admin` for example) you're logged in to your Jenkins.
      Make sure you have access to the necessary pipeline jobs that is about to be integrated with the
    - for `JENKINS_USER_API_TOKEN`, you can obtain it in your profile by clicking your name next to `Logout` button then
      click `Configure`. Go to `API Token` section and create one for the jenkins bot go
    - `SLACK_BOT_TOKEN` will be your slack bot token
    - `SLACK_APP_TOKEN` will be your slack bot token
    - `SLACK_DEFAULT_CHANNEL` set default channel to send the approval to 
    - `APP_DEBUG` you can choose between `true` or `false`
    - `APP_MODE` you can choose between `aio` or `standlone`. See more in [running the app](#running-the-app)

### Running the app

If you don't set the `APP_MODE`, the default value would be `aio`. That means both HTTP API and Slack bot will run at the same command

```shell
go run .
```

If you choose run them separately, set the `APP_MODE` to `standalone`. Then you can run the slack bot with
```shell
go run . artisan slack:socket
```

### Pipeline Example

Create a new pipeline and copy this to see if it's working fine

```
stage("Setup") {
  node('master') {
    def host = "http://slack-bot.dev:3000/"
    def data = [
        "build_number": env.BUILD_NUMBER,
        "build_name": env.JOB_NAME,
        "message": "approval for ${env.JOB_NAME}"
    ]
    def json = JsonOutput.toJson(data)
    sh "curl -X POST -H 'Content-Type: application/json' ${host}/approval/jenkins --data '${json}'" 
  } 
  input("Approval required before proceeding deployment")
}
```

You can also build the binary and put it in your server or even build it in container and deploy it in your cluster.

### Contribution

PR or Issue submission are very much welcome. I will do my best to engage and resolve the said items.

### Support

Issues will be reviewed and resolved at best effort.
However, if you need further assistance, do not hesitate to contact me at [my email](mailto:mgufronefendi@gmail.com)

### TODO

- [x] Initial Release
- [ ] Expand to other communication channel
    - [x] Slack
    - [ ] Telegram
    - [ ] Mattermost
    - [ ] ...
- [ ] Integrate to other CI/CD tools if possible
  - [x] Jenkins
  - [ ] Gitlab
- [x] Allow Jenkins to directly send API request to jenkins-bot-go that will send approval message to Slack
- [ ] Security Layer
