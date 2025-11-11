package commands

import (
	"math/rand/v2"
	"fmt"
	"log"
	"slices"
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

	log.Printf("%v: Started deleting messages from channel: %v", ChannelClearCommand.Name, event.ChannelID)
	handlerStartTime := time.Now()
	const batchLimit int = 100 // maximum discord page fetch/batch delete size

	// Return initial deferred response to avoid timeout
	err := session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Failed to send defer message: %v", err)
		return
	}

	// Fetch all channel messages in batches
	allMessages := make([]*discordgo.Message, 0, batchLimit)
	for true {
		length := len(allMessages)
		beforeId := ""
		if length > 0 {
			beforeId = allMessages[length-1].ID
		}

		messages, err := session.ChannelMessages(event.ChannelID, batchLimit, beforeId, "", "")
		if err != nil {
			log.Printf("error fetching messages: %v", err)
			break
		}

		if len(messages) == 0 {
			break
		}

		allMessages = append(allMessages, messages...)
	}

	// We need to avoid deleting the app interaction message so we can reply to it later
	// It will be one of the first ones, as disc messages are returned in descending order
	allMessages = slices.DeleteFunc(allMessages, func(m *discordgo.Message) bool {
		if m.Interaction != nil {
			if m.Interaction.ID == event.Interaction.ID {
				return true
			}
		}
		return false
	})

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

	deletedCount := 0
	// perform bulk deletes in batches
	for true {
		count := len(bulkDeleteIds)

		if count < batchLimit {
			batch := bulkDeleteIds[:count]
			err = session.ChannelMessagesBulkDelete(event.ChannelID, batch)
			if err != nil {
				log.Printf("Failed to delete batch of messages: %v", err)
				break
			}
			deletedCount += count
			break
		}
		batch := bulkDeleteIds[:batchLimit]
		bulkDeleteIds = bulkDeleteIds[batchLimit:]
		err = session.ChannelMessagesBulkDelete(event.ChannelID, batch)
		if err != nil {
			log.Printf("Failed to delete batch of messages: %v", err)
			continue
		}
		deletedCount += batchLimit

	}

	// perform individual deletes
	deletedCount += massDeleteSingleMessages(session, event.ChannelID, &individualDeleteIds)

	// update reply to user and tell them how much hellfire they just rained
	replyContent := fmt.Sprintf("Deleted %v messages.", deletedCount)
	_, err = session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
		Content: &replyContent,
	})
	if err != nil {
		log.Printf("Failed to edit initial interaction msg: %v", err)
	}

	handlerDuration := time.Since(handlerStartTime).String()
	log.Printf("%v: Deleted %v messages in %v", ChannelClearCommand.Name, deletedCount, handlerDuration)
}

func massDeleteSingleMessages(session *discordgo.Session, channelId string, messageIds *[]string) int {
	var deletedCount int
	var minDeleteDelay int = 600
	var maxDeleteDelay int = 1600

	for _, v := range *messageIds {
		// buffer for preventing ratelimiting
		seed := rand.IntN(maxDeleteDelay-minDeleteDelay) + minDeleteDelay
		deleteDelay := time.Duration(seed) * time.Millisecond 
		log.Printf("Waiting %v until next delete...", deleteDelay.String()) // TODO: make these logs debug-level
		time.Sleep(deleteDelay)

		log.Printf("Sending delete request for id: %v", v) // TODO: make these logs debug-level
		err := session.ChannelMessageDelete(channelId, v)
		if err != nil {
			log.Printf("Failed to delete message %v: %v", v, err)
			continue
		}
		log.Printf("Send complete: %v", v) // TODO: make these logs debug-level
		deletedCount++
	}

	return deletedCount
}
