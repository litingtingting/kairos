package commands

import (
    "strings"
    "kairos/reminders"
    "kairos/weather"
    "kairos/dice"

    "github.com/bwmarrin/discordgo"
)

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID {
        return
    }

    // 处理命令（以 ! 开头）
    if !strings.HasPrefix(m.Content, "!") {
        // 简单的闲聊回应
        casualChat(s, m)
        return
    }

    parts := strings.Fields(m.Content[1:])
    if len(parts) == 0 {
        return
    }

    cmd := strings.ToLower(parts[0])
    args := parts[1:]

    switch cmd {
    case "天气", "weather":
        if len(args) == 0 {
            s.ChannelMessageSend(m.ChannelID, "请指定城市，例如 `!天气 北京`")
            return
        }
        go func() {
            result, err := weather.GetWeather(args[0])
            if err != nil {
                s.ChannelMessageSend(m.ChannelID, "天气查询失败: "+err.Error())
                return
            }
            s.ChannelMessageSend(m.ChannelID, result)
        }()

    case "提醒", "remind":
        if len(args) < 2 {
            s.ChannelMessageSend(m.ChannelID, "用法：`!提醒 [时间] [消息]`\n例如：`!提醒 每3分钟 喝水`\n`!提醒 10:30 开会`")
            return
        }
        timeSpec := args[0]
        message := strings.Join(args[1:], " ")
        err := reminders.GetManager().AddReminder(m.Author.ID, m.ChannelID, message, timeSpec)
        if err != nil {
            s.ChannelMessageSend(m.ChannelID, "设置提醒失败: "+err.Error())
        }

    case "我的提醒", "list":
        list := reminders.GetManager().ListReminders(m.Author.ID)
        if len(list) == 0 {
            s.ChannelMessageSend(m.ChannelID, "你还没有任何提醒")
        } else {
            s.ChannelMessageSend(m.ChannelID, "📋 **你的提醒列表:**\n"+strings.Join(list, "\n"))
        }

    case "取消提醒", "取消":
        if len(args) == 0 {
            s.ChannelMessageSend(m.ChannelID, "请指定提醒ID，例如 `!取消提醒 123456`")
            return
        }
        if reminders.GetManager().RemoveReminder(args[0]) {
            s.ChannelMessageSend(m.ChannelID, "✅ 提醒已取消")
        } else {
            s.ChannelMessageSend(m.ChannelID, "❌ 未找到该提醒ID")
        }

    case "骰子", "roll":
        var diceInput string
        if len(args) == 0 {
            diceInput = "1d6"
        } else {
            diceInput = args[0]
        }
        result := dice.RollDice(diceInput)
        s.ChannelMessageSend(m.ChannelID, result)

    case "help", "帮助":
        helpMsg := `**🤖 Kairos 机器人命令列表**
!天气 [城市] - 查询天气
!提醒 [时间] [消息] - 设置提醒（支持：每3分钟、10:30、2025-03-01 15:04）
!我的提醒 - 查看当前提醒
!取消提醒 [ID] - 取消提醒
!骰子 [表达式] - 掷骰子，如 !骰子 2d6+3
!ping - 测试机器人是否在线`
        s.ChannelMessageSend(m.ChannelID, helpMsg)

    case "ping":
        s.ChannelMessageSend(m.ChannelID, "Pong! 🏓")
    }
}

func casualChat(s *discordgo.Session, m *discordgo.MessageCreate) {
    content := strings.ToLower(m.Content)
    switch {
    case strings.Contains(content, "你好") || strings.Contains(content, "hello"):
        s.ChannelMessageSend(m.ChannelID, "你好呀！👋 需要帮忙吗？输入 `!help` 查看命令")
    case strings.Contains(content, "在吗"):
        s.ChannelMessageSend(m.ChannelID, "在的在的，随时待命！")
    case strings.Contains(content, "谢谢"):
        s.ChannelMessageSend(m.ChannelID, "不客气～有什么需要随时叫我")
    }
}