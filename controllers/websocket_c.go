package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"scm/consts"
	"scm/utils"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/fasthttp/websocket"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"google.golang.org/api/option"
)

// Map untuk menyimpan koneksi WebSocket pengguna
var clients = make(map[string]*websocket.Conn)
var mu sync.Mutex

var upgrader = websocket.FastHTTPUpgrader{
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true // Izinkan semua origin
	},
}

type Chat struct {
	Id         string `json:"_id"`
	SenderID   string `json:"sid"`
	ReceiverID string `json:"rid"`
	Message    string `json:"msg"`
}

func generateUserID() string {
	return uuid.New().String()
}

type RequestPayload struct {
	Token string `json:"token"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func SendPushNotification(ctx *fasthttp.RequestCtx) error {
	if len(ctx.Request.Body()) == 0 {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.PayloadEmpty)
		return fmt.Errorf("%s", consts.PayloadEmpty)
	}

	var reqPayload RequestPayload
	err := json.Unmarshal(ctx.PostBody(), &reqPayload)
	if err != nil {
		return fmt.Errorf("invalid json payload: %w", err)
	}

	// Inisialisasi Firebase
	opt := option.WithCredentialsFile("./serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing firebase app: %w", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return fmt.Errorf("error getting messaging client: %w", err)
	}

	// Gunakan hanya bagian `Data`
	// message := &messaging.Message{
	// 	Token: reqPayload.Token,
	// 	Data: map[string]string{
	// 		"title":        reqPayload.Title,
	// 		"body":         reqPayload.Body,
	// 		"click_action": "https://yourweb.com/pesan", // Bisa diarahkan ke halaman tujuan
	// 	},
	// }

	message := &messaging.Message{
		Token: reqPayload.Token,
		Data: map[string]string{
			"title":        reqPayload.Title,
			"body":         reqPayload.Body,
			"click_action": "http://localhost:9000/", // sesuaikan
		},
	}

	// Kirim pesan
	response, err := client.Send(context.Background(), message)
	if err != nil {
		utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err.Error())
		log.Println("FCM error:", err)
		return fmt.Errorf("error sending FCM message: %w", err)
	}

	utils.ShowResponseDefault(ctx, fasthttp.StatusOK, "success", response)
	return nil
}

func WebSocketHandler(ctx *fasthttp.RequestCtx, clientID string) {
	var userID string

	if clientID == "" {
		userID = generateUserID()
		ctx.Response.Header.Set("userID", userID)
		log.Printf("Generated new userID: %s", userID)
	} else {
		userID = clientID
		log.Printf("Using provided clientID: %s", userID)
	}

	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		log.Println("WebSocket connection upgraded successfully for userID:", userID) // Log ini untuk upgrade sukses
		defer conn.Close()

		// Simpan koneksi WebSocket pengguna berdasarkan userID
		mu.Lock()
		clients[userID] = conn
		mu.Unlock()

		// Terus-menerus menerima pesan dari pengguna
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}

			log.Printf("Received message from %s: %s", userID, message)

			// Kirim balasan atau tangani pesan yang diterima
			err = conn.WriteMessage(websocket.TextMessage, []byte("Message received: "+string(message)))
			if err != nil {
				log.Println("Write error:", err)
				break
			}
		}

		// Hapus pengguna dari map saat koneksi putus
		mu.Lock()
		delete(clients, userID)
		mu.Unlock()
		log.Printf("Connection closed for userID: %s", userID)
	})

	if err != nil {
		log.Printf("Failed to upgrade WebSocket for userID %s: %+v", userID, err)
	}
}

// Fungsi untuk mengirim pesan dari satu pengguna ke pengguna lain
func SendMessage(ctx *fasthttp.RequestCtx) error {

	var chat Chat
	utils.JsonToStruct(string(ctx.PostBody()), &chat)
	mu.Lock()
	recipientConn, ok := clients[chat.ReceiverID]
	mu.Unlock()

	if !ok {
		log.Printf("User %s is not connected", chat.ReceiverID)
		return nil // Atau kembalikan error jika ingin menandai bahwa user tidak online
	}

	// Kirim pesan ke penerima
	err := recipientConn.WriteMessage(websocket.TextMessage, []byte(chat.SenderID+": "+chat.Message))
	if err != nil {
		log.Println("write error:", err)
		return err
	}

	log.Printf("Message sent from %s to %s: %s", chat.SenderID, chat.ReceiverID, chat.Message)
	return nil
}

func CreateChannel(ctx *fasthttp.RequestCtx) {
	// Ambil userA dan userB dari request body
	userA := string(ctx.PostArgs().Peek("userA"))
	userB := string(ctx.PostArgs().Peek("userB"))

	// Periksa apakah user A dan B terhubung
	mu.Lock()
	_, userAExists := clients[userA]
	_, userBExists := clients[userB]
	mu.Unlock()

	if !userAExists || !userBExists {
		ctx.Error("One or both users are not connected", fasthttp.StatusBadRequest)
		return
	}

	// Buat ID channel dan simpan
	// channelID := generateChannelID(userA, userB)

	// newChannel := &Channel{
	// 	ID:        channelID,
	// 	Users:     []string{userA, userB},
	// 	Messages:  []Message{},
	// 	CreatedAt: time.Now(),
	// }

	// Simpan channel
	// mu.Lock()
	// channels[channelID] = newChannel
	// mu.Unlock()

	// // Kirim response sukses
	// ctx.SetStatusCode(fasthttp.StatusOK)
	// ctx.SetBodyString(fmt.Sprintf("Channel created with ID: %s", channelID))
}
