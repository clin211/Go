package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/option"
)

var users []map[string]string = []map[string]string{
	{"767425412@qq.com": "e6GTz7hZTnSmOTrE1MC3fv:APA91bFj5PydDZkS8Nyck-sRaH9e2rl2HGU7BSlpBKQ6yN3cykMjRI1_i6Evg6SMen2Huj2n7xOpmjLp3kcSejtDZdJB-UMDN2MQTly6V6jByFWdHkIKTL0"},
	{"178588646@qq.com": "e6GTz7hZTnSmOTrE1MC3fv:APA91bFj5PydDZkS8Nyck-sRaH9e2rl2HGU7BSlpBKQ6yN3cykMjRI1_i6Evg6SMen2Huj2n7xOpmjLp3kcSejtDZdJB-UMDN2MQTly6V6jByFWdHkIKTL1"},
	{"1785886467@qq.com": "e6GTz7hZTnSmOTrE1MC3fv:APA91bFj5PydDZkS8Nyck-sRaH9e2rl2HGU7BSlpBKQ6yN3cykMjRI1_i6Evg6SMen2Huj2n7xOpmjLp3kcSejtDZdJB-UMDN2MQTly6V6jByFWdHkIKTL2"},
	{"178588648@qq.com": "e6GTz7hZTnSmOTrE1MC3fv:APA91bFj5PydDZkS8Nyck-sRaH9e2rl2HGU7BSlpBKQ6yN3cykMjRI1_i6Evg6SMen2Huj2n7xOpmjLp3kcSejtDZdJB-UMDN2MQTly6V6jByFWdHkIKTL3"},
	{"178588649@qq.com": "e6GTz7hZTnSmOTrE1MC3fv:APA91bFj5PydDZkS8Nyck-sRaH9e2rl2HGU7BSlpBKQ6yN3cykMjRI1_i6Evg6SMen2Huj2n7xOpmjLp3kcSejtDZdJB-UMDN2MQTly6V6jByFWdHkIKTL4"},
}

var emails []string = []string{
	"767425412@qq.com",
	"178588646@qq.com",
	"1785886467@qq.com",
	"178588648@qq.com",
	"178588649@qq.com",
}

var tokens []string = []string{
	"e6GTz7hZTnSmOTrE1MC3fv:APA91bFj5PydDZkS8Nyck-sRaH9e2rl2HGU7BSlpBKQ6yN3cykMjRI1_i6Evg6SMen2Huj2n7xOpmjLp3kcSejtDZdJB-UMDN2MQTly6V6jByFWdHkIKTL0",
}

type MessageData struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	ID     string `json:"id"`
	WebUrl string `json:"web_url"`
}

func (m *MessageData) ToMap() map[string]string {
	return map[string]string{
		"type":    m.Type,
		"title":   m.Title,
		"id":      m.ID,
		"web_url": m.WebUrl,
	}
}
func main() {
	ctx := context.Background()
	app, err := instance()
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// 写一个 API 接口
	http.HandleFunc("GET /send", func(w http.ResponseWriter, r *http.Request) {
		sendMessages(ctx, client)
	})
	http.HandleFunc("GET /single", func(w http.ResponseWriter, r *http.Request) {
		singleMessage(ctx, client)
	})
	http.HandleFunc("GET /batch", func(w http.ResponseWriter, r *http.Request) {
		sendMessagesAdvanced(ctx, client)
	})

	http.HandleFunc("GET /concurrently", func(w http.ResponseWriter, r *http.Request) {
		sendMessagesConcurrently(ctx, client)
	})

	http.HandleFunc("GET /topic", func(w http.ResponseWriter, r *http.Request) {
		sendMessages11(ctx, client)
	})
	http.ListenAndServe(":8080", nil)
}

func singleMessage(ctx context.Context, client *messaging.Client) {
	data := &MessageData{
		Type:   "direct-top-up",
		Title:  "Harry Potter: Magic Awakened Jewels",
		ID:     "16384",
		WebUrl: "https://www.seagm.com/harry-potter-magic-awakened-jewels?ps=Home-Special-Deals&item_id=16384",
	}
	message := &messaging.Message{
		Token: "e6GTz7hZTnSmOTrE1MC3fv:APA91bFj5PydDZkS8Nyck-sRaH9e2rl2HGU7BSlpBKQ6yN3cykMjRI1_i6Evg6SMen2Huj2n7xOpmjLp3kcSejtDZdJB-UMDN2MQTly6V6jByFWdHkIKTL0",
		Notification: &messaging.Notification{
			Title: "Hello",
			Body:  "World",
		},
		Data: data.ToMap(),
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				Title: "Hello",
				Body:  "World",
			},
			Data: data.ToMap(),
		},
	}

	// res, err := client.Send(ctx, message)
	br, err := client.Send(ctx, message)
	if err != nil {
		panic("firebase messaging error: " + err.Error())
	}
	fmt.Printf("Successfully sent message: %s\n", br)
}

func instance() (*firebase.App, error) {
	opt := option.WithCredentialsFile("./configs/lin-rn-firebase-admin-sdk.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}
	return app, nil
}

// 发送 100000 消息
func sendMessages(ctx context.Context, client *messaging.Client) {
	// 生成 100000 条 token
	tokens := make([]string, 100000)
	for i := range tokens {
		tokens[i] = fmt.Sprintf("token-%d", i)
	}
	tokens[0] = "e6GTz7hZTnSmOTrE1MC3fv:APA91bFj5PydDZkS8Nyck-sRaH9e2rl2HGU7BSlpBKQ6yN3cykMjRI1_i6Evg6SMen2Huj2n7xOpmjLp3kcSejtDZdJB-UMDN2MQTly6V6jByFWdHkIKTL0"
	// 生成 100000 条消息
	messages := make([]*messaging.MulticastMessage, 100000)
	for i := range messages {
		messages[i] = &messaging.MulticastMessage{
			Tokens: tokens[i : i+1],
			Android: &messaging.AndroidConfig{
				Priority: "high",
				Notification: &messaging.AndroidNotification{
					Title:    "6480 + 2000 Jewels",
					Body:     "Harry Potter: Magic Awakened Jewels",
					ImageURL: "https://seagm-media.seagmcdn.com/icon_400/1803.jpg?x-oss-process=image/resize,w_108",
				},
				Data: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
		}
	}

	// 将 token 和 消息按 500 条进行分组发送
	batchSize := 500
	batches := make([][]*messaging.MulticastMessage, 0, len(messages)/batchSize+1)
	tokenBatches := make([][]string, 0, len(tokens)/batchSize+1)
	for i := 0; i < len(messages); i += batchSize {
		batch := messages[i:min(i+batchSize, len(messages))]
		batches = append(batches, batch)
		tokenBatches = append(tokenBatches, tokens[i:min(i+batchSize, len(tokens))])
	}

	// 发送消息
	startTime := time.Now()
	for i := range batches {
		fmt.Printf("Sending batch %d...\n", i)
		res, err := client.SendEachForMulticast(ctx, &messaging.MulticastMessage{
			Tokens: tokenBatches[i],
			Android: &messaging.AndroidConfig{
				Priority: "high",
				Notification: &messaging.AndroidNotification{
					Title:    "6480 + 2000 Jewels",
					Body:     "Harry Potter: Magic Awakened Jewels",
					ImageURL: "https://seagm-media.seagmcdn.com/icon_400/1803.jpg?x-oss-process=image/resize,w_108",
				},
				Data: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
		})
		if err != nil {
			panic("firebase messaging error: " + err.Error())
		}
		fmt.Println("res: ", res.SuccessCount, res.FailureCount)
		if res.FailureCount > 0 {
			var failedTokens []string
			for idx, resp := range res.Responses {
				if !resp.Success {
					// The order of responses corresponds to the order of the registration tokens.
					failedTokens = append(failedTokens, tokenBatches[i][idx])
				}
			}

			fmt.Printf("List of tokens that caused failures: %v\n", failedTokens)
		}
	}
	elapsedTime := time.Since(startTime)
	fmt.Println("Total time: ", elapsedTime)
}

func sendMessagesBatch(ctx context.Context, client *messaging.Client, tokenBatch []string) {
	res, err := client.SendEachForMulticast(ctx, &messaging.MulticastMessage{
		Tokens: tokenBatch,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				Title:    "6480 + 2000 Jewels",
				Body:     "Harry Potter: Magic Awakened Jewels",
				ImageURL: "https://seagm-media.seagmcdn.com/icon_400/1803.jpg?x-oss-process=image/resize,w_108",
			},
			Data: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	})
	if err != nil {
		panic("firebase messaging error: " + err.Error())
	}
	fmt.Println("res: ", res.SuccessCount, res.FailureCount)
	if res.FailureCount > 0 {
		var failedTokens []string
		for idx, resp := range res.Responses {
			if !resp.Success {
				failedTokens = append(failedTokens, tokenBatch[idx])
			}
		}
		fmt.Printf("List of tokens that caused failures: %v\n", failedTokens)
	}
}

func sendMessagesAdvanced(ctx context.Context, client *messaging.Client) error {
	// 生成 100000 条 token 和消息
	tokens := make([]string, 100000)
	messages := make([]*messaging.MulticastMessage, 100000)
	for i := range tokens {
		tokens[i] = fmt.Sprintf("token-%d", i)
		messages[i] = &messaging.MulticastMessage{
			Tokens: tokens[i : i+1],
			Android: &messaging.AndroidConfig{
				Priority: "high",
				Notification: &messaging.AndroidNotification{
					Title:    "6480 + 2000 Jewels",
					Body:     "Harry Potter: Magic Awakened Jewels",
					ImageURL: "https://seagm-media.seagmcdn.com/icon_400/1803.jpg?x-oss-process=image/resize,w_108",
				},
				Data: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
		}
	}
	tokens[0] = "e6GTz7hZTnSmOTrE1MC3fv:APA91bFj5PydDZkS8Nyck-sRaH9e2rl2HGU7BSlpBKQ6yN3cykMjRI1_i6Evg6SMen2Huj2n7xOpmjLp3kcSejtDZdJB-UMDN2MQTly6V6jByFWdHkIKTL0"

	// 将 token 和 消息按 500 条进行分组
	batchSize := 500
	batches := make([][]*messaging.MulticastMessage, 0, len(messages)/batchSize+1)
	tokenBatches := make([][]string, 0, len(tokens)/batchSize+1)
	for i := 0; i < len(messages); i += batchSize {
		batch := messages[i:min(i+batchSize, len(messages))]
		batches = append(batches, batch)
		tokenBatches = append(tokenBatches, tokens[i:min(i+batchSize, len(tokens))])
	}

	// 并行发送消息
	startTime := time.Now()
	var wg sync.WaitGroup
	msgCh := make(chan [][]*messaging.MulticastMessage, len(batches))
	tokenCh := make(chan [][]string, len(tokenBatches))

	for i := range batches {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			msgCh <- [][]*messaging.MulticastMessage{batches[i]}
			tokenCh <- [][]string{tokenBatches[i]}
		}(i)
	}

	go func() {
		wg.Wait()
		close(msgCh)
		close(tokenCh)
	}()

	// for msgs := range msgCh {

	// }
	for tokens := range tokenCh {
		sendMessagesBatch(ctx, client, tokens[0])
	}
	elapsedTime := time.Since(startTime)
	fmt.Println("Total time: ", elapsedTime)
	return nil
}

func sendMessagesConcurrently(ctx context.Context, client *messaging.Client) {
	batchSize := 500
	// 生成 100000 条 token 和消息
	tokens := make([]string, 100000)
	messages := make([]*messaging.MulticastMessage, 100000)
	for i := range tokens {
		tokens[i] = fmt.Sprintf("token-%d", i)
		messages[i] = &messaging.MulticastMessage{
			Tokens: tokens[i : i+1],
			Android: &messaging.AndroidConfig{
				Priority: "high",
				Notification: &messaging.AndroidNotification{
					Title:    "6480 + 2000 Jewels",
					Body:     "Harry Potter: Magic Awakened Jewels",
					ImageURL: "https://seagm-media.seagmcdn.com/icon_400/1803.jpg?x-oss-process=image/resize,w_108",
				},
				Data: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
		}
	}
	tokens[0] = "e6GTz7hZTnSmOTrE1MC3fv:APA91bFj5PydDZkS8Nyck-sRaH9e2rl2HGU7BSlpBKQ6yN3cykMjRI1_i6Evg6SMen2Huj2n7xOpmjLp3kcSejtDZdJB-UMDN2MQTly6V6jByFWdHkIKTL0"

	var wg sync.WaitGroup

	var sc, fc int
	startTime := time.Now()
	for i := 0; i < len(messages); i += batchSize {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			// batch := messages[offset:min(offset+batchSize, len(messages))]
			tokenBatch := tokens[offset:min(offset+batchSize, len(tokens))]
			res, err := client.SendEachForMulticast(ctx, &messaging.MulticastMessage{
				Tokens: tokenBatch,
				Android: &messaging.AndroidConfig{
					Priority: "high",
					Notification: &messaging.AndroidNotification{
						Title:    "6480 + 2000 Jewels",
						Body:     "Harry Potter: Magic Awakened Jewels",
						ImageURL: "https://seagm-media.seagmcdn.com/icon_400/1803.jpg?x-oss-process=image/resize,w_108",
					},
					Data: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
			})
			if err != nil {
				log.Printf("firebase messaging error: %v", err)
				return
			}
			log.Printf("Batch %d: SuccessCount %d, FailureCount %d", offset/batchSize, res.SuccessCount, res.FailureCount)
			sc += res.SuccessCount
			fc += res.FailureCount
			if res.FailureCount > 0 {
				// 处理失败的tokens
			}
		}(i)
	}

	wg.Wait()
	elapsedTime := time.Since(startTime)
	fmt.Printf("Total time: %d; success count: %d; failure count: %d", elapsedTime, sc, fc)
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sendMessages11(ctx context.Context, client *messaging.Client) {
	// 生成 100000 条 token
	tokens := make([]string, 100000)
	for i := range tokens {
		tokens[i] = fmt.Sprintf("token-%d", i)
	}
	tokens[0] = "e6GTz7hZTnSmOTrE1MC3fv:APA91bFj5PydDZkS8Nyck-sRaH9e2rl2HGU7BSlpBKQ6yN3cykMjRI1_i6Evg6SMen2Huj2n7xOpmjLp3kcSejtDZdJB-UMDN2MQTly6V6jByFWdHkIKTL0"

	// 将 token 按 500 条进行分组
	batchSize := 500
	tokenBatches := make([][]string, 0, len(tokens)/batchSize+1)
	for i := 0; i < len(tokens); i += batchSize {
		end := min(i+batchSize, len(tokens))
		tokenBatches = append(tokenBatches, tokens[i:end])
	}

	// 创建消息配置
	androidConfig := &messaging.AndroidConfig{
		Priority: "high",
		Notification: &messaging.AndroidNotification{
			Title:    "6480 + 2000 Jewels",
			Body:     "Harry Potter: Magic Awakened Jewels",
			ImageURL: "https://seagm-media.seagmcdn.com/icon_400/1803.jpg?x-oss-process=image/resize,w_108",
		},
		Data: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	// 并行发送消息
	startTime := time.Now()
	var g errgroup.Group
	var mu sync.Mutex
	totalSuccess := 0
	totalFailure := 0

	for i, tokens := range tokenBatches {
		tokensCopy := tokens // 为每个 go routine 创建独立的 token 切片
		i := i               // 保护索引变量
		g.Go(func() error {
			fmt.Printf("Sending batch %d...\n", i)
			res, err := client.SendEachForMulticast(ctx, &messaging.MulticastMessage{
				Tokens:  tokensCopy,
				Android: androidConfig,
			})
			if err != nil {
				return fmt.Errorf("firebase messaging error: %w", err)
			}

			fmt.Println("res: ", res.SuccessCount, res.FailureCount)
			mu.Lock()
			totalSuccess += res.SuccessCount
			totalFailure += res.FailureCount
			mu.Unlock()

			if res.FailureCount > 0 {
				var failedTokens []string
				for idx, resp := range res.Responses {
					if !resp.Success {
						failedTokens = append(failedTokens, tokensCopy[idx])
					}
				}
				fmt.Printf("List of tokens that caused failures: %v\n", failedTokens)
			}
			return nil
		})
	}

	// 等待所有消息发送完成
	if err := g.Wait(); err != nil {
		log.Fatalf("Error in sending messages: %v", err)
	}

	// 计算成功率
	totalMessages := totalSuccess + totalFailure
	if totalMessages == 0 {
		fmt.Println("No messages to send")
		return
	}
	successPercentage := (float64(totalSuccess) / float64(totalMessages)) * 100
	elapsedTime := time.Since(startTime)

	fmt.Printf("Total messages: %d\n", totalMessages)
	fmt.Printf("Successful messages: %d\n", totalSuccess)
	fmt.Printf("Failed messages: %d\n", totalFailure)
	fmt.Printf("Success rate: %.2f%%\n", successPercentage)
	fmt.Println("Total time: ", elapsedTime)
}
