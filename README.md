# Jenkins-Slack Bot 

## Table of contents
- [Requirements](#requirements)
- [Setup](#setup)
- [Run the app](#running-the-app)

### Requirements
- Slack app with socket mode enabled
- Jenkins slack plugin
- Go >=1.20

### Setup
#### Create Slack App 

You have options regarding slack app. You can create a new app dedicated for jenkins-bot-go and slack jenkins plugin. Or use existing app. Install it to your slack workspace.

Get the App Token in `Settings -> Basic Information` section, scroll down and you will see App Token. Create a new and add these scopes:
- `connections:write`
- `authorization:read`

![Screenshot](/screenshots/1-slack-app-token.png?raw=true "Slack App Token")

Next, you need to obtain the Bot Token. Go to `Features -> OAuth & Permissions`. 

Go to Scopes section and add these OAuth Scope:
- `chat:write`
- `im:history`

Scroll up and copy the bot token.
![Screenshot](/screenshots/2-slack-bot-token.png?raw=true "Slack bot token")

Go to 

#### Install and configure Slack jenkins plugin 
#### Setup jenkins-bot-go
- Clone the repository
- run 
```shell
go mod download
```
- Copy `.env.example` and modify them according to your configuration

### Running the app
