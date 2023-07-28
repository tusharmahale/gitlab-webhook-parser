package slack

import (
	"fmt"
	"log"
	"os"
	"time"

	// "github.com/joho/godotenv"
	"github.com/go-playground/webhooks/v6/gitlab"
	"github.com/slack-go/slack"
)

func SendSlackNotification(mrpl gitlab.MergeRequestEventPayload) {

	token := os.Getenv("SLACK_AUTH_TOKEN")
	if token == "" {
		log.Fatalf("Please provide valid SLACK_AUTH_TOKEN")
	}
	channelID := os.Getenv("SLACK_CHANNEL_ID")
	if channelID == "" {
		log.Fatalf("Please provide valid SLACK_CHANNEL_ID")
	}

	// Create a new client to slack by giving token
	// Set debug to true while developing
	client := slack.New(token, slack.OptionDebug(true))
	// Create the Slack attachment that we will send to the channel

	titleText := fmt.Sprintf("MR - %d - State - %s [No new MR will be merged till Build completes]", mrpl.ObjectAttributes.ID, mrpl.ObjectAttributes.State)
	frmtedText := fmt.Sprintf("Merge request\t%s \nMerge State\t %s \nSource branch\t %s\nTarget branch\t %s \n", mrpl.ObjectAttributes.URL, mrpl.ObjectAttributes.State, mrpl.ObjectAttributes.SourceBranch, mrpl.ObjectAttributes.TargetBranch)
	attachment := slack.Attachment{
		Title:   titleText,
		Pretext: titleText,
		Text:    frmtedText,
		// Color Styles the Text, making it possible to have like Warnings etc.
		Color: "#36a64f",
		// Fields are Optional extra data!
		Fields: []slack.AttachmentField{
			{
				Title: "Date",
				Value: time.Now().String(),
			},
		},
	}
	// PostMessage will send the message away.
	// First parameter is just the channelID, makes no sense to accept it
	_, timestamp, err := client.PostMessage(
		channelID,
		// uncomment the item below to add a extra Header to the message, try it out :)
		//slack.MsgOptionText("New message from bot", false),
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Message sent at %s", timestamp)
}
