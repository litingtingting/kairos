package dice

import (
    "fmt"
    "math/rand"
    "strconv"
    "strings"
    "time"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

func RollDice(input string) string {
    input = strings.ReplaceAll(input, " ", "")
    input = strings.ToLower(input)

    // 处理 + 和 -
    parts := strings.FieldsFunc(input, func(r rune) bool {
        return r == '+' || r == '-'
    })

    var total int
    var detail []string

    for _, part := range parts {
        if strings.Contains(part, "d") {
            // 骰子部分
            diceParts := strings.Split(part, "d")
            count, _ := strconv.Atoi(diceParts[0])
            sides, _ := strconv.Atoi(diceParts[1])

            var sum int
            var rolls []string
            for i := 0; i < count; i++ {
                roll := rand.Intn(sides) + 1
                sum += roll
                rolls = append(rolls, strconv.Itoa(roll))
            }

            // 判断是加法还是减法
            if strings.Contains(input, "+"+part) {
                total += sum
            } else if strings.Contains(input, "-"+part) {
                total -= sum
            } else {
                total += sum
            }
            detail = append(detail, fmt.Sprintf("%s:[%s]", part, strings.Join(rolls, ",")))
        } else {
            // 常数部分
            num, _ := strconv.Atoi(part)
            if strings.Contains(input, "+"+part) || (!strings.HasPrefix(input, "-") && len(detail) == 0) {
                total += num
                detail = append(detail, fmt.Sprintf("+%d", num))
            } else {
                total -= num
                detail = append(detail, fmt.Sprintf("-%d", num))
            }
        }
    }

    return fmt.Sprintf("🎲 **%s**\n= **%d**", strings.Join(detail, " "), total)
}