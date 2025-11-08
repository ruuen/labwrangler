package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// command definition
var ChannelClearCommand = discordgo.ApplicationCommand{
	Name: "clear-channel-all",
	Description: "Clear all messages from current channel",
}

// handler func
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

	// channelbulkmessagedelete has a restriction of 2 weeks old; will 400 if anything older
	// bulk delete under 2 weeks old, then individual delete the remainder
	var allMessages []*discordgo.Message
	// TODO: move this to own fn called in recursion, pass in ptr to bulk and individual slices, keep allmessages inside fn
	// TODO: bulkdeletes and individualdeletes slices should only store the messageid; storing the whole message struct is moving a lot of extra data around
	for true {
		beforeId := ""
		if len(allMessages) > 0 {
			beforeId = allMessages[len(allMessages)-1].ID
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

	var (
		individualDeletes []*discordgo.Message
		bulkDeletes []*discordgo.Message
		bufferMins int8 = 5
	)
	bulkCutoffTime := time.Now().UTC().AddDate(0, 0, -14).Add(-time.Duration(bufferMins) * time.Minute)

 	for _, v := range allMessages {
		timestamp, err := discordgo.SnowflakeTimestamp(v.ID)
		if err != nil {
			log.Println("Failed to parse snowflake to timestamp")
		}

		if timestamp.Before(bulkCutoffTime) {
			individualDeletes = append(individualDeletes, v)
			continue
		}

		bulkDeletes = append(bulkDeletes, v)
	}

	log.Printf("Bulk deletes: %v, Single deletes: %v\n", len(bulkDeletes), len(individualDeletes))

	// perform bulk deletes in batches of 100
	for true {
		count := len(bulkDeletes)

		if count < 100 {
			batch := bulkDeletes[:count]
			log.Printf("Batch count: %v", len(batch)) // bulk delete call goes here when you're ready to light this candle
			break
		}
		
		batch := bulkDeletes[:100]
		bulkDeletes = bulkDeletes[100:]
		
		log.Printf("Batch count: %v", len(batch)) // bulk delete call goes here when you're ready to light this candle
	}

	// perform individual deletes
	for _, v := range individualDeletes {
		log.Printf("Deleting message id %v", v.ID) // individual delete call goes here when you need to waste this sucker
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

// any other helper funcs used inside this command
