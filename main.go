package main

import (
  "fmt"
  "log"
  "os"
  "os/signal"
  "syscall"

  "github.com/bwmarrin/discordgo"
  "github.com/joho/godotenv"
)

func main() {
  err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))
  if err != nil {
    log.Fatal(err)
  }

  Token := os.Getenv("DISCORD_TOKEN")
  log.Println("Token: ", Token)
  if Token == "" {
    return
  }

  dg, err := discordgo.New("Bot " + Token)
  if err != nil {
    log.Println("error creating Discord session,", err)
    return
  }

  dg.AddHandler(ready)
  dg.AddHandler(messageCreate)
  dg.AddHandler(createChannel)

  dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages)

  err = dg.Open()
  if err != nil {
    log.Println("error opening connection,", err)
    return
  }

  log.Println("Bot is now running. Press CTRL-C to exit.")
  sc := make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
  <-sc

  err = dg.Close()
  if err != nil {
    log.Println("error closing connection,", err)
    return
  }
}

func ready(s *discordgo.Session, event *discordgo.MessageCreate) {
  err := s.UpdateGameStatus(0, "demoapp!")
  if err != nil {
    log.Println("error updating game status,", err)
    return
  }
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
  log.Println("messageCreate started")

  if m.Author.ID == s.State.User.ID {
    return
  }

  if m.Content == "ServerName" {
    g, err := s.Guild(m.GuildID)
    if err != nil {
      log.Fatal(err)
    }
    log.Println(g.Name)
    _, err = s.ChannelMessageSend(m.ChannelID, g.Name)
    if err != nil {
      log.Println("error sending message")
    }
  }

  if m.Content == "!Hello" {
    _, err := s.ChannelMessageSend(m.ChannelID, "Hello")
    if err != nil {
      log.Println("error sending message")
    }
  }
}

func createChannel(s *discordgo.Session, m *discordgo.MessageCreate) {
  log.Println("createChannel started")

  if m.Content == "!CC" {
    g, err := s.Guild(m.GuildID)
    if err != nil {
      log.Println("error getting guild id")
    }

    _, err = s.GuildChannelCreate(g.ID, "channel-created-by-bot", discordgo.ChannelTypeGuildText)
    if err != nil {
      log.Println("error creating channel")
    }
  }
}
