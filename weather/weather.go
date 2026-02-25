package weather

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
)

type WeatherResponse struct {
    Name string `json:"name"`
    Main struct {
        Temp     float64 `json:"temp"`
        Humidity int     `json:"humidity"`
    } `json:"main"`
    Weather []struct {
        Description string `json:"description"`
    } `json:"weather"`
}

func GetWeather(city string) (string, error) {
    apiKey := os.Getenv("OPENWEATHER_API_KEY")
    if apiKey == "" {
        return "", fmt.Errorf("OPENWEATHER_API_KEY 未设置")
    }

    url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=zh_cn", city, apiKey)

    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return "", fmt.Errorf("查询失败，状态码: %d", resp.StatusCode)
    }

    body, _ := io.ReadAll(resp.Body)
    var weather WeatherResponse
    if err := json.Unmarshal(body, &weather); err != nil {
        return "", err
    }

    return fmt.Sprintf("🌍 **%s**\n🌡️ 温度: %.1f°C\n💧 湿度: %d%%\n☁️ %s",
        weather.Name, weather.Main.Temp, weather.Main.Humidity, weather.Weather[0].Description), nil
}