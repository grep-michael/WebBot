package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/grep-michael/WebBot/Bots/GovDeals"
	_ "github.com/grep-michael/WebBot/Bots/MTGSecretLair"
	_ "github.com/grep-michael/WebBot/Caches/MapCache"
	_ "github.com/grep-michael/WebBot/NotificationDestination/DiscordDestinations"

	dynamicconfiguration "github.com/grep-michael/WebBot/DynamicConfiguration"
	"github.com/grep-michael/WebBot/globals"
)

func loadBots() []globals.Bot {

	data, _ := os.ReadFile("config.json")
	var botList dynamicconfiguration.BotList
	json.Unmarshal(data, &botList)

	bots := make([]globals.Bot, 0)
	for _, botCfg := range botList.Bots {
		bot, err := dynamicconfiguration.CreateBot(botCfg)
		if err != nil {
			log.Printf("[ERROR] Failed to make bot %s:\n\t%v\n", botCfg.InstanceName, err)
			continue
		}
		bots = append(bots, bot)
	}
	return bots
}

func main() {
	log.SetOutput(os.Stdout)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	bots := loadBots()

	if len(bots) <= 0 {
		log.Printf("No bots loaded, exiting ...")
		return
	}
	errCh := make(chan error, len(bots))
	finChan := make(chan string, len(bots))

	for _, bot := range bots {
		go func() {
			err := bot.Run(ctx)
			if err != nil {
				errCh <- fmt.Errorf("bot %s: %+v", bot.Name(), err)
			} else {
				finChan <- fmt.Sprintf("bot %s finished", bot.Name())
			}
		}()
	}

	select {
	case <-ctx.Done():
		log.Println("Context Canceled, Exiting...")
	case err := <-errCh:
		log.Printf("[ERROR]: %v\n", err)
	case msg := <-finChan:
		log.Println(msg)
	}
}
