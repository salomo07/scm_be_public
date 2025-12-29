package main

import (
	"log"
	"os"
	"scm/config"
	"scm/routers"
	"scm/utils"
	"time"

	_ "scm/docs"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

var app_port = os.Getenv("APP_PORT")

func main() {
	log.Println(config.DecryptAES(config.DecryptAES("nfSajo0XBmYbURO/BbdqEiNs1qGgPay04OeT6p+SWVghcj1bmDd4hxyVS6Q0YgJeWQo3NVs4hNaQOzyHfDjiTV/nwoZqhLzi986yvAn39U9yrhAYk/pfGg==")))
	utils.GenerateSuperAdminToken()

	router := fasthttprouter.New()

	router.GET("/", func(ctx *fasthttp.RequestCtx) {
		utils.ShowResponseDefault(ctx, 200, "ok", "Welcome to SCM API")
	})

	router.GET("/company/", func(ctx *fasthttp.RequestCtx) {
		utils.ShowResponseDefault(ctx, 200, "ok", "Ini router /company/")
	})

	routers.AdminRouters(router)
	routers.FCMRouters(router)
	routers.WebsocketRouters(router)
	routers.BranchRouters(router)
	routers.IDM_Routers(router)

	router.HandleMethodNotAllowed = false
	router.MethodNotAllowed = func(ctx *fasthttp.RequestCtx) {
		log.Printf("❌ METHOD NOT ALLOWED: %s %s", ctx.Method(), ctx.Path())
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.Response.Header.Set("Content-Type", "application/json")
		utils.ShowResponseDefault(ctx, fasthttp.StatusMethodNotAllowed, "error", "Method not allowed")
	}

	router.NotFound = func(ctx *fasthttp.RequestCtx) {
		log.Printf("❌ NOT FOUND: %s %s", ctx.Method(), ctx.Path())
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.Set("Content-Type", "application/json")
		utils.ShowResponseDefault(ctx, fasthttp.StatusNotFound, "error", "The endpoint is not found")
	}

	if app_port == "" {
		app_port = "7777"
	}

	// ✅ Chain middlewares TANPA compression dulu
	handler := withLogging(withCORS(router.Handler))

	server := &fasthttp.Server{
		Handler:            handler,
		MaxRequestBodySize: 10 * 1024 * 1024, // 10MB
		ReadTimeout:        30 * time.Second, // Naikkan timeout
		WriteTimeout:       30 * time.Second,
		IdleTimeout:        120 * time.Second,
		Name:               "SCM-API-Server",
	}

	log.Printf("🚀 Server starting on :%s", app_port)
	if err := server.ListenAndServe("0.0.0.0:" + app_port); err != nil {
		log.Fatalf("❌ Error starting server: %s", err)
	}
}

// Middleware untuk logging
func withLogging(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()

		// Log incoming request dengan detail header
		acceptEncoding := string(ctx.Request.Header.Peek("Accept-Encoding"))
		log.Printf("📨 %s %s | Accept-Encoding: %s | Body: %d bytes",
			ctx.Method(),
			ctx.Path(),
			acceptEncoding,
			len(ctx.Request.Body()),
		)

		// Process request
		next(ctx)

		// Log response
		duration := time.Since(start)
		log.Printf("📤 %s %s | Status: %d | Duration: %v | Response: %d bytes",
			ctx.Method(),
			ctx.Path(),
			ctx.Response.StatusCode(),
			duration,
			len(ctx.Response.Body()),
		)
	}
}

// Middleware CORS dengan whitelist origin
func withCORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	allowedOrigins := map[string]bool{
		"http://localhost:9000":     true,
		"http://localhost:3000":     true,
		"http://192.168.99.47:9000": true,
		"http://127.0.0.1:9000":     true,
		"http://127.0.0.1:3000":     true,
	}

	return func(ctx *fasthttp.RequestCtx) {
		origin := string(ctx.Request.Header.Peek("Origin"))

		// Set CORS headers
		if origin != "" {
			if allowedOrigins[origin] {
				ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
				ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
			} else {
				// Development: allow all
				ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
			}
		} else {
			ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		}

		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Accept-Language")
		ctx.Response.Header.Set("Access-Control-Max-Age", "86400")
		ctx.Response.Header.Set("Vary", "Origin")

		// Handle preflight OPTIONS request
		if string(ctx.Method()) == "OPTIONS" {
			log.Printf("⚡ OPTIONS preflight: %s", ctx.Path())
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}

		// Process actual request
		next(ctx)
	}
}
