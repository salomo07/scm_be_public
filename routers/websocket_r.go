package routers

import (
	"log"
	"scm/controllers"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func FCMRouters(router *fasthttprouter.Router) {
	router.POST("/fcm/send", func(ctx *fasthttp.RequestCtx) { //Connect to WS
		controllers.SendPushNotification(ctx)
	})
	print("--FCMRouters--\n")
}
func WebsocketRouters(router *fasthttprouter.Router) {
	router.GET("/ws", func(ctx *fasthttp.RequestCtx) { //Connect to WS
		controllers.WebSocketHandler(ctx, "")
	})
	router.POST("/send", func(ctx *fasthttp.RequestCtx) {
		err := controllers.SendMessage(ctx)
		if err != nil {
			log.Println("Error sending message:", err)
		}
	})
	router.POST("/creategroup", func(ctx *fasthttp.RequestCtx) {
		err := controllers.SendMessage(ctx)
		if err != nil {
			log.Println("Error sending message:", err)
		}
	})
	router.POST("/create-channel", func(ctx *fasthttp.RequestCtx) {
		controllers.CreateChannel(ctx)
	})
	print("--WebsocketRouters--\n")
}
