package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "kairos/commands"
    "kairos/reminders"

    "github.com/bwmarrin/discordgo"
    "github.com/joho/godotenv"
)

func main() {
    // 加载 .env 文件
    err := godotenv.Load()
    if err != nil {
        log.Println("警告: 未找到 .env 文件，将使用系统环境变量")
    }

    token := os.Getenv("DISCORD_BOT_TOKEN")
    if token == "" {
        log.Fatal("环境变量 DISCORD_BOT_TOKEN 未设置")
    }

    // 创建 Discord 会话
    dg, err := discordgo.New("Bot " + token)
    if err != nil {
        log.Fatal("创建 Discord 会话失败:", err)
    }

    // 注册事件处理
    dg.AddHandler(commands.MessageHandler)

    // 设置需要的 intents
    dg.Identify.Intents = discordgo.IntentsGuildMessages |
        discordgo.IntentsMessageContent |
        discordgo.IntentsDirectMessages

    // 打开连接
    err = dg.Open()
    if err != nil {
        log.Fatal("无法打开连接:", err)
    }
    defer dg.Close()

    // 启动定时提醒管理器
    reminders.StartManager(dg)

    log.Println("Kairos 机器人已启动！按 Ctrl+C 停止")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
    <-sc
}