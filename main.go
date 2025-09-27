package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"scm/config"
	"scm/routers"
	"scm/utils"

	_ "scm/docs" // Pastikan path ini sesuai dengan lokasi file docs.go yang dihasilkan oleh `swag init`
	"scm/services"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var app_port = os.Getenv("APP_PORT")

// @title SCM  API (Go Swagger)
// @version 1.0.0
// @description All API for SCM.
// @contact.name salomo07
// @contact.email sitompulsalomo@gmail.com
// @basePath /api/v1
// @host localhost:8081
// @schemes http https
func main() {
	print(config.EncryptAES("solideogloria"))
	// print(config.DecryptChacha20poly1305("TsI8DiAkFoiGXpgZbYJM3w7wKvBrI51OT0zWq/GEvfv9iiHce6HCGW3jXaxzfdlPFcbMZLhSsmKFyMQZQ3Wg0z4dmaLsCEZ8JioO", "lt2OQvpizbly5OhXfhYHIm1HJX5Fqtxr"),"")
	// print(config.EncryptAES("solideogloria"))
	// print(config.DecryptAES("1J60PTcNaoT+7PbhglQA6uxjSUvGTuS4zZ3epNAgGKX2M/Zl9BCxqAk="))
	// print(config.EncryptAES("1504070208920001"),
	// 	"\n", config.EncryptAES("Raja Salomo Sitompul"),
	// 	"\n", config.EncryptAES("Salomo07"),
	// 	"\n", config.EncryptAES(config.EncryptAES("xxx")),
	// 	"\n", config.EncryptAES("085186803737"),
	// 	"\n", config.EncryptAES("sitompulsalomo@gmail.com"),
	// 	"\n", config.EncryptAES("solideogloria"),
	// )
	utils.GenerateSuperAdminToken()
	// print(config.DecryptAES(config.DecryptAES("/cNs20MYUnvmb+JB9ocC620B0fDNo/opw+9vLe+/Qy/GJoAX3RCa3RJpZfdKFMlGQg6KCYd8PiaXWb60YGFVcRzs5YYMmBM4tBwF7yoHQNhBhbdPooAy0Q==")), "\n\n")
	// print("\nusername", config.DecryptAES("s5Y/hWnF9MOFG06Pz6I3YsYjsICkA67Qon80mdNrs1HEShSXus2dLj9ae1k="), "\npass", config.DecryptAES("WQqKh9P34fX2WDozYYwQZzLYUu3pBVWlRkxIb+SBT390WDsZoZsG3vqrDzH9S4N3vEZGxChLK78ZwPYe1xnddNoDSSWz/G67XzVcZ21QSgx52BOXIuyuyQ=="), "\n\n")
	// print("\n", config.EncryptAES("Salomo07"), "\n")
	// services.GenerateRSAPrivatePublicKey()

	router := fasthttprouter.New()

	router.GET("/", func(ctx *fasthttp.RequestCtx) {
		utils.ShowResponseDefault(ctx, 200, "ok", "Welcome to SCM API")
	})
	router.GET("/company/", func(ctx *fasthttp.RequestCtx) {
		utils.ShowResponseDefault(ctx, 200, "ok", "Ini router /company/")
	})

	routers.FCMRouters(router)
	routers.WebsocketRouters(router)
	routers.AdminRouters(router)
	routers.IDM_Routers(router)

	router.HandleMethodNotAllowed = false
	router.MethodNotAllowed = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.Response.Header.Set("Content-Type", "application/json")
	}
	router.NotFound = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		utils.ShowResponseDefault(ctx, fasthttp.StatusNotFound, "warning", "The endpoint is not found")
		ctx.Response.Header.Set("Content-Type", "application/json")
	}
	if app_port == "" {
		app_port = "8080"
	}
	server := &fasthttp.Server{
		Handler: withCORS(router.Handler),
	}

	// server := &fasthttp.Server{Handler: router.Handler}
	go func() {
		if err := server.ListenAndServe("0.0.0.0:" + app_port); err != nil {
			log.Printf("Error starting server: %s\n", err)
		}
	}()
	log.Println("Server listen on :" + app_port)
	select {}
}

// ✅ Middleware CORS
func withCORS(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	allowedOrigins := map[string]bool{
		"http://localhost:9000":   true,
		"http://localhost:3000":   true,
		"http://10.10.20.14:9000": true,
	}
	return func(ctx *fasthttp.RequestCtx) {
		origin := string(ctx.Request.Header.Peek("Origin"))
		// Jika origin diizinkan, set header CORS
		if allowedOrigins[origin] {
			ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
			ctx.Response.Header.Set("Vary", "Origin") // agar caching aware origin
			ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight (OPTIONS)
		if string(ctx.Method()) == fasthttp.MethodOptions {
			ctx.SetStatusCode(fasthttp.StatusOK)
			return
		}

		handler(ctx)
	}
}
func initConnectionMongoDB() {
	uri := "mongodb://superuser:inipassmongo@localhost:27017/admin"
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		// return nil, fmt.Sprintf("Failed to connect to database: %v", err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		// return nil, fmt.Sprintf("Failed to ping database: %v", err)
	}
	services.SetMongoClient(client)
}
func handleGithubWebhook(ctx *fasthttp.RequestCtx) {
	// Eksekusi perintah untuk memperbarui proyek
	cmd := exec.Command("sh", "-c", "cd /path/to/your/project && git pull origin main_dev")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Output: %s\n", stdoutStderr)

	// Kirim respons ke GitHub
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte("Webhook received!"))
}
