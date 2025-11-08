package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

var ChannelClearCommand = discordgo.ApplicationCommand{
	Name: "clear-channel-all",
	Description: "Clear all messages from current channel",
}

func ChannelClearHandler(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if event.ApplicationCommandData().Name != ChannelClearCommand.Name {
		return
	}

	// Return initial deferred response to avoid timeout
	err := session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Failed to send defer message: %v", err)
		return
	}

	// Fetch all channel messages in batches
	allMessages := make([]*discordgo.Message, 0, 100)
	for true {
		length := len(allMessages)
		beforeId := ""
		if length > 0 {
			beforeId = allMessages[length-1].ID
		}

		messages, err := session.ChannelMessages(event.ChannelID, 100, beforeId, "", "")
		if err != nil {
			log.Printf("error fetching messages: %v", err)
			break
		}

		if len(messages) == 0 {
			break
		}

		allMessages = append(allMessages, messages...)
	}

	// channelbulkmessagedelete has a restriction of 2 weeks old; will 400 if anything older
	// therefore we split messages <2 weeks and >2 weeks to be processed separately
	var (
		individualDeleteIds []string
		bulkDeleteIds []string
		bufferMins int8 = 5
	)
	bulkCutoffTime := time.Now().UTC().AddDate(0, 0, -14).Add(-time.Duration(bufferMins) * time.Minute) // 2 weeks and 5 mins in past

	for i := 0; i < len(allMessages); i++ {
		v := allMessages[i]
		timestamp, err := discordgo.SnowflakeTimestamp(v.ID)
		if err != nil {
			log.Println("Failed to parse snowflake to timestamp")
		}

		if timestamp.Before(bulkCutoffTime) {
			individualDeleteIds = append(individualDeleteIds, v.ID)
			continue
		}

		bulkDeleteIds = append(bulkDeleteIds, v.ID)
	}

	// perform bulk deletes in batches of 100
	for true {
		count := len(bulkDeleteIds)

		if count < 100 {
			batch := bulkDeleteIds[:count]
			log.Printf("Batch count: %v", len(batch)) // bulk delete call goes here when you're ready to light this candle
			break
		}
		
		batch := bulkDeleteIds[:100]
		bulkDeleteIds = bulkDeleteIds[100:]
		
		log.Printf("Batch count: %v", len(batch)) // bulk delete call goes here when you're ready to light this candle
	}

	// perform individual deletes
	for _, v := range individualDeleteIds {
		log.Printf("Deleting message id %v", v) // individual delete call goes here when you need to waste this sucker
	}

	// update reply to user and tell them how much hellfire they just rained
	replyContent := fmt.Sprintf("Deleted %v messages.", len(allMessages))
	_, err = session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
		Content: &replyContent,
	})
	if err != nil {
		log.Printf("Failed to edit initial interaction msg: %v", err)
	}
}

