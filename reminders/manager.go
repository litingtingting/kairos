package reminders

import (
    "fmt"
    "strings"
    "sync"
    "time"

    "github.com/bwmarrin/discordgo"
    "github.com/robfig/cron/v3"
)

type Reminder struct {
    ID          string
    UserID      string
    ChannelID   string
    Message     string
    Schedule    string
    IsRecurring bool
}

type Manager struct {
    session *discordgo.Session
    cron    *cron.Cron
    entries map[string]cron.EntryID
    mu      sync.RWMutex
}

var instance *Manager
var once sync.Once

func StartManager(s *discordgo.Session) {
    once.Do(func() {
        instance = &Manager{
            session: s,
            cron:    cron.New(cron.WithSeconds()),
            entries: make(map[string]cron.EntryID),
        }
        instance.cron.Start()
        log.Println("提醒管理器已启动")
    })
}

func GetManager() *Manager {
    return instance
}

// 添加提醒
func (m *Manager) AddReminder(userID, channelID, message, timeSpec string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    // 解析时间表达式
    schedule, err := parseTimeSpec(timeSpec)
    if err != nil {
        return err
    }

    id := fmt.Sprintf("%d", time.Now().UnixNano())

    entryID, err := m.cron.AddFunc(schedule, func() {
        m.sendReminder(userID, channelID, message)
    })
    if err != nil {
        return err
    }

    m.entries[id] = entryID

    // 发送确认消息
    m.session.ChannelMessageSend(channelID,
        fmt.Sprintf("✅ 提醒已设置！\nID: `%s`\n时间: %s\n消息: %s", id, timeSpec, message))

    return nil
}

func (m *Manager) sendReminder(userID, channelID, message string) {
    m.session.ChannelMessageSend(channelID,
        fmt.Sprintf("<@%s> ⏰ 提醒你：%s", userID, message))
}

// 解析用户输入的时间格式
func parseTimeSpec(input string) (string, error) {
    input = strings.TrimSpace(input)

    // 每X分钟
    if strings.HasPrefix(input, "每") && strings.Contains(input, "分钟") {
        var minutes int
        fmt.Sscanf(input, "每%d分钟", &minutes)
        if minutes < 1 {
            minutes = 1
        }
        return fmt.Sprintf("0 */%d * * * *", minutes), nil
    }

    // 每X秒 (最小30秒)
    if strings.HasPrefix(input, "每") && strings.Contains(input, "秒") {
        var seconds int
        fmt.Sscanf(input, "每%d秒", &seconds)
        if seconds < 30 {
            seconds = 30
        }
        return fmt.Sprintf("@every %ds", seconds), nil
    }

    // 每天固定时间，如 "10:30"
    if strings.Contains(input, ":") && len(input) <= 5 {
        return fmt.Sprintf("0 %s * * *", input), nil
    }

    // 一次性提醒：日期时间 "2025-03-01 15:04"
    if t, err := time.Parse("2006-01-02 15:04", input); err == nil {
        return fmt.Sprintf("%d %d %d %d %d *",
            t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month())), nil
    }

    return "", fmt.Errorf("无法识别的时间格式")
}

// 列出用户的所有提醒
func (m *Manager) ListReminders(userID string) []string {
    m.mu.RLock()
    defer m.mu.RUnlock()

    var list []string
    for id, entryID := range m.entries {
        entry := m.cron.Entry(entryID)
        list = append(list, fmt.Sprintf("ID: `%s` | 下次执行: %s", id, entry.Next.Format("2006-01-02 15:04:05")))
    }
    return list
}

// 删除提醒
func (m *Manager) RemoveReminder(id string) bool {
    m.mu.Lock()
    defer m.mu.Unlock()

    if entryID, ok := m.entries[id]; ok {
        m.cron.Remove(entryID)
        delete(m.entries, id)
        return true
    }
    return false
}