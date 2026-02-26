#!/bin/sh
# mkdir kairos && cd kairos
go mod init kairos
go get github.com/bwmarrin/discordgo
go get github.com/joho/godotenv
go get github.com/robfig/cron/
#go get github.com/google/generative-ai-go
go get cloud.google.com/go/ai
go get google.golang.org/genai
