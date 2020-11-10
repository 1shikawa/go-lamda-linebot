package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
	"linebot/gurunavi"
	"log"
	"net/http"
	"os"
)

type Webhook struct {
	Destination string           `json:"destination"`
	Events      []*linebot.Event `json:"events"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Print(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusInternalServerError)),
		}, nil
	}

	log.Print(request.Headers)
	log.Print(request.Body)

	var webhook Webhook

	// jsonパースしてwebhook構造体に格納
	if err := json.Unmarshal([]byte(request.Body), &webhook); err != nil {
		log.Print(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusBadRequest)),
		}, nil
	}
	
	for _, event := range webhook.Events {
		switch event.Type {
		case linebot.EventTypeMessage:
			// event.Messageの型によって場合分け
			switch m := event.Message.(type) {
			// テキストメッセージの場合
			case *linebot.TextMessage:
				switch request.Path {
				case "/parrot":
					// オウム返しとしてテキスト返信
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(m.Text)).Do(); err != nil {
						log.Print(err)
						return events.APIGatewayProxyResponse{
							StatusCode: http.StatusInternalServerError,
							Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusBadRequest)),
						}, nil
					}
				case "/restaurants":
					// レストラン情報を検索して結果json取得
					g, err := gurunavi.SearchRestaurants(m.Text)
					if err != nil {
						log.Print(err)
						return events.APIGatewayProxyResponse{
							StatusCode: http.StatusInternalServerError,
							Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusInternalServerError)),
						}, nil
					}

					// var t string
					var sm linebot.SendingMessage

					switch {
					case g.Error != nil:
						t := g.Error[0].Message
						sm = linebot.NewTextMessage(t)
					default:
						// t = TextRestaurants(g)
						f := FlexRestaurants(g) //バブルコンテナ生成→カルーセルコンテナ生成情報
						sm = linebot.NewFlexMessage("飲食店検索結果", f)
					}

					// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(t)).Do(); err != nil {
					// 	log.Fatal(err)
					// }

					// カルーセルコンテナとしてLine返信
					if _, err = bot.ReplyMessage(event.ReplyToken, sm).Do(); err != nil {
						log.Print(err)
						return events.APIGatewayProxyResponse{
							StatusCode: http.StatusInternalServerError,
							Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusInternalServerError)),
						}, err
					}
				}
			}
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
