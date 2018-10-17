package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
)

var debugExists = false
var debugVerboseExists = false

var logToSlack = false
var slackAPI *slack.Client
var slackParams slack.PostMessageParameters

func init() {
	_, debugExists = os.LookupEnv("DEBUG")
	_, debugVerboseExists = os.LookupEnv("DEBUGVERBOSE")

	apiToken, apiTokenExists := os.LookupEnv("SLACKAPI_TOKEN")
	if !apiTokenExists {
		log.Printf("logger [" + LogErr + "] SLACKAPI_TOKEN env variable not set\n")
	}
	if debugVerboseExists {
		log.Printf("logger ["+LogDebugVerbose+"] slack api token: %s", apiToken)
	}

	channelName, channelNameExists := os.LookupEnv("SLACK_CHANNEL")
	if !channelNameExists {
		log.Printf("logger [" + LogErr + "] SLACK_CHANNEL env variable not set")
	}
	if debugVerboseExists {
		log.Printf("logger ["+LogDebug+"] slack channel: %s", channelName)
	}

	if apiTokenExists && channelNameExists {
		logToSlack = true

		slackAPI = slack.New(apiToken)
	}
}

// logger levels
const (
	LogPanic        = "PANIC"
	LogFatal        = "FATAL"
	LogErr          = "ERR"
	LogWarn         = "WARN"
	LogInfo         = "INFO"
	LogDebug        = "DEBUG"
	LogDebugVerbose = "DEBUGVERBOSE"
)

// Log formatted messages
// env variables:
//   DEBUG: include DEBUG level messages
//   DEBUGVERBOSE: include DEBUGVERBOSE level messages
func Log(module string, level string, format string, v ...interface{}) {
	switch level {
	case LogPanic:
		log.Panicf(module+" ["+LogPanic+"] "+format, v...)
		return
	case LogFatal:
		log.Fatalf(module+" ["+LogFatal+"] "+format, v...)
		return
	case LogDebug:
		if debugExists || debugVerboseExists {
			log.Printf(module+" ["+level+"] "+format, v...)
		}
		return
	case LogDebugVerbose:
		if debugVerboseExists {
			log.Printf(module+" ["+level+"] "+format, v...)
		}
		return
	default:
		log.Printf(module+" ["+level+"] "+format, v...)
	}
}

// LogToSlack fomatted messages to SLACK channel
// env variables:
//   SLACKAPI_TOKEN: SLACK API application token with permissions to write to specified channel
//   SLACK_CHANNEL: SLACK channel name
func LogToSlack(module string, level string, format string, v ...interface{}) {
	if logToSlack {
		message := fmt.Sprintf("```"+module+" ["+level+"] "+format+"```", v...)
		channelID, timestamp, err := slackAPI.PostMessage("logs", slack.MsgOptionText(message, false), slack.MsgOptionAsUser(true), slack.MsgOptionUser(module))
		if err != nil {
			Log("logger", LogErr, "%s logger error: %s", module, err)
		}
		Log("logger", LogDebugVerbose, "%s message successfully send to channel %s at %s", module, channelID, timestamp)
	}
	Log(module, level, format, v...)
}
