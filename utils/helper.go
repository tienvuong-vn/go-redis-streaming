package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"go-streaming/model"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func LoadConfiguration(file string) (model.Config, error) {
	var config model.Config
	configFile, err := os.Open(file)
	if err != nil {
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	defer configFile.Close()
	return config, err
}

func SendData(channelManager *model.ChannelManager, msg *redis.Message) {
	channelManager.Submit(msg.Channel, msg.Payload)
}

func SendPing(channelManager *model.ChannelManager, sseInstanceId string, rdb *redis.Client) {
	channelManager.Submit("PING", time.Now().Format("01-02-2006 15:04:05"))
	go func() {
		str := fmt.Sprintf("%s, SSE-Total:%d, SSE-Live:%d, SSE-Closed:%d, WS-Total:%d, WS-Live:%d, WS-Closed:%d, Messages:%d", sseInstanceId, channelManager.SseTotal, channelManager.SseLive, channelManager.SseClosed, channelManager.WsTotal, channelManager.WsLive, channelManager.WsClosed, channelManager.TotalMessage)
		rdb.Publish(context.Background(), "sse:admin", str).Err()
	}()
}
