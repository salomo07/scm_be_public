package routers

import (
	"fmt"
	"scm/config"
	"scm/consts"
	"scm/controllers"
	"scm/utils"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func IDM_Routers(router *fasthttprouter.Router) {
	print("\n -IDM\n")
	print(consts.URL_Auth_Login + "\n")
	print(consts.URL_Auth_Logout + "\n")
	print(consts.URL_Auth_Enc + "\n")
	router.POST(consts.URL_Auth_Login, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		controllers.Login(ctx)
	})

	router.POST(consts.URL_Auth_Logout, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		controllers.Logout(ctx)
	})
	router.POST(consts.URL_Auth_Enc, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		type Payload struct {
			Payload string `json:"payload"`
		}
		var pay Payload
		utils.JsonToStruct(string(ctx.Request.Body()), &pay)
		utils.ShowResponseDefault(ctx, fasthttp.StatusOK, "encrypt", config.EncryptAES(pay.Payload))
	})
	router.POST(consts.URL_Auth_Dec, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		type Payload struct {
			Payload string `json:"payload"`
		}
		var pay Payload
		utils.JsonToStruct(string(ctx.Request.Body()), &pay)
		utils.ShowResponseDefault(ctx, fasthttp.StatusOK, "decrypt", config.DecryptAES(pay.Payload))
	})
	router.POST("/test", func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		ctx.Response.Header.Set("Content-Type", "application/json")
		fmt.Fprintf(ctx, `[{"id":1,"nik":"10091062","jobcode":"staff","employeename":"Siganteng"},{"id":2,"nik":"10091063","jobcode":"am","employeename":"Sicakep"}]`)
	})

	// router.POST(consts.URL_Auth_RefreshToken, func(ctx *fasthttp.RequestCtx) {
	// 	refreshToken := string(ctx.Request.Header.Cookie("refresh_token"))

	// 	if refreshToken == "" {
	// 		utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "error", "No refresh token")
	// 		return
	// 	}

	// 	// Verify refresh token
	// 	claims, err := jwt.VerifyRefreshToken(refreshToken)
	// 	if err != nil {
	// 		utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "error", "Invalid refresh token")
	// 		return
	// 	}

	// 	// Generate new access token
	// 	newAccessToken, _ := jwt.GenerateAccessToken(claims.Data)

	// 	utils.ShowResponseDefault(ctx, fasthttp.StatusOK, "success", map[string]string{
	// 		"access_token": newAccessToken,
	// 	})
	// })
}
