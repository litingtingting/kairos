package commands

import (
    "strings"
    "kairos/reminders"
    "kairos/weather"
    "kairos/dice"

    "github.com/bwmarrin/discordgo"
    "log"
    "kairos/ai"      // 引入刚才写的 ai 包
    "context"
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
    
    case "ai", "ask", "chat": // 新增 AI 命令
        if len(args) == 0 {
            s.ChannelMessageSend(m.ChannelID, "你想让我用 AI 帮你做什么？在后面加上你的问题，例如 `!ai 写一首关于大海的诗`")
            return
        }
        // 将用户的所有输入合并成一个提示词
        prompt := strings.Join(args, " ")
        
        // 在 goroutine 中处理，避免阻塞消息接收
        go handleAIRequest(s, m, prompt)

    case "help", "帮助":
        helpMsg := `**🤖 Kairos 机器人命令列表**
!天气 [城市] - 查询天气
!提醒 [时间] [消息] - 设置提醒（支持：每3分钟、10:30、2025-03-01 15:04）
!我的提醒 - 查看当前提醒
!取消提醒 [ID] - 取消提醒
!骰子 [表达式] - 掷骰子，如 !骰子 2d6+3
!ai [对话] - 跟AI聊天
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

// handleAIRequest 是一个独立的函数来处理 AI 请求
func handleAIRequest(s *discordgo.Session, m *discordgo.MessageCreate, prompt string) {
    // 先发送一个“正在思考”的提示，因为 AI 响应可能需要几秒钟
    thinkingMsg, _ := s.ChannelMessageSend(m.ChannelID, "🤔 让我想想...")

    // 创建一个上下文
    ctx := context.Background()
    
    // 初始化 AI 客户端
    aiClient, err := ai.NewClient(ctx)
    if err != nil {
        log.Printf("AI 客户端初始化失败: %v", err)
        s.ChannelMessageEdit(m.ChannelID, thinkingMsg.ID, "抱歉，AI 大脑暂时无法连接，请检查服务器配置（GEMINI_API_KEY）。")
        return
    }
    //defer aiClient.Close() // 记得关闭

    // 调用 AI 获取回答
    answer, err := aiClient.Ask(prompt)
    if err != nil {
        log.Printf("AI 请求失败: %v", err)
        s.ChannelMessageEdit(m.ChannelID, thinkingMsg.ID, "抱歉，AI 思考时出了点小差错，请稍后再试。")
        return
    }

    // 编辑之前的“思考中”消息，替换为 AI 的回答
    // 注意：Discord 消息有长度限制（2000 字符），如果答案太长可能需要分段发送
    s.ChannelMessageEdit(m.ChannelID, thinkingMsg.ID, "🤖 **AI 回答**:\n"+answer)
}