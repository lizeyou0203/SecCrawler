package bot

import (
    . "SecCrawler/config"
    "SecCrawler/register"
    "SecCrawler/utils"
    "fmt"
    "net/http"
    "encoding/json"
    "bytes"
    "errors"
)

type MessagePusher struct{}

func (bot MessagePusher) Config() register.BotConfig {
    return register.BotConfig{
        Name: "MessagePusher",
    }
}

type request struct {
    Title        string `json:"title"`
    Description string `json:"description"`
    Content     string `json:"content"`
    URL         string `json:"url"`
    Channel     string `json:"channel"`
    Token        string `json:"token"`
}

type response struct {
    Success bool    `json:"success"`
    Message string `json:"message"`
}

// Send 推送消息给企业微信机器人。
func (bot MessagePusher) Send(crawlerResult [][]string, description string) error {
    var msg string
    var description_a string
    
    for _, i := range crawlerResult {
        text := fmt.Sprintf(" **%s**:--->[%s](%s)  \n", i[1], i[0], i[0])
        msg += text
    }
    title := fmt.Sprintf("## %s ## %s", description, utils.CurrentTime())
    
    for _, i := range crawlerResult {
        text := fmt.Sprintf("|---%s    \n", i[1])
        description_a += text
    }

    fmt.Printf(msg)

    reqs := request{
        Title:        title,
        Description: description_a,
        Content:     msg,
        Token:        Cfg.Bot.MessagePusher.Token,
    }

    data, err := json.Marshal(reqs)
    if err != nil {
        return err
    }

    resp, err := http.Post(fmt.Sprintf("%s/push/%s", Cfg.Bot.MessagePusher.ServerAddress, Cfg.Bot.MessagePusher.Username), "application/json", bytes.NewBuffer(data))
    if err != nil {
        return err
    }
    var res response
    err = json.NewDecoder(resp.Body).Decode(&res)
    if err != nil {
        return err
    }
    if !res.Success {
        return errors.New(res.Message)
    }
    return nil
}
