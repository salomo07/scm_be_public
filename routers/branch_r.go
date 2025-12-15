package routers

import (
	"log"
	"scm/consts"
	"scm/controllers"
	"scm/controllers/master_controller"
	"scm/models"
	"scm/utils"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func BranchRouters(router *fasthttprouter.Router) {
	log.Println("\n========== REGISTERING BRANCH ROUTES ==========")
	log.Println("📍 POST", consts.URL_Branch_Find)
	log.Println("📍 POST", consts.URL_Branch_Upsert)
	log.Println("📍 POST", consts.URL_Branch_Delete)
	log.Println("===============================================\n")

	router.POST(consts.URL_Branch_Find, func(ctx *fasthttp.RequestCtx) {
		log.Println("✅ ===== BRANCH FIND HANDLER HIT =====")
		log.Println("Method:", string(ctx.Method()))
		log.Println("Path:", string(ctx.Path()))
		log.Println("Body:", string(ctx.Request.Body()))

		ctx.Response.Header.Set("Content-Type", "application/json")

		// Uncomment untuk testing
		user, err, _, _ := controllers.CheckSession(ctx)
		var dataBranch models.BranchRequest
		utils.JsonToStruct(string(ctx.Request.Body()), &dataBranch)

		if err != "" {
			log.Println("❌ CheckSession error:", err)
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
			return
		} else {
			log.Println("✅ CheckSession success, calling FindBranches")
			master_controller.FindBranches(ctx, user)
		}
	})
	// router.POST(consts.URL_Branch_Find, func(ctx *fasthttp.RequestCtx) {
	// 	// Tambahkan log ini di paling awal
	// 	log.Println("========================================")
	// 	log.Println("ROUTE HIT: URL_Branch_Find")
	// 	log.Println("Path:", string(ctx.Path()))
	// 	log.Println("Method:", string(ctx.Method()))
	// 	log.Println("========================================")

	// 	ctx.Response.Header.Set("Content-Type", "application/json")
	// 	log.Println("ccccccccccc")

	// 	// Uncomment untuk testing
	// 	user, err, _, _ := controllers.CheckSession(ctx)
	// 	var dataBranch models.BranchRequest
	// 	utils.JsonToStruct(string(ctx.Request.Body()), &dataBranch)

	// 	if err != "" {
	// 		utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
	// 		return
	// 	} else {
	// 		master_controller.FindBranches(ctx, user)
	// 	}
	// })
	router.POST(consts.URL_Branch_Upsert, func(ctx *fasthttp.RequestCtx) {
		log.Println("Route hit:", consts.URL_Branch_Upsert)
		user, err, _, isCompanyAdmin := controllers.CheckSession(ctx)
		log.Println("CheckSession returned:", err, isCompanyAdmin)
		ctx.Response.Header.Set("Content-Type", "application/json")
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
			return
		} else if isCompanyAdmin {
			master_controller.UpsertBranch(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
		}
	})

	router.POST(consts.URL_Branch_Delete, func(ctx *fasthttp.RequestCtx) {
		creddb, err, _, isCompanyAdmin := controllers.CheckSession(ctx)
		ctx.Response.Header.Set("Content-Type", "application/json")
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
			return
		} else if isCompanyAdmin {
			master_controller.UpsertBranch(ctx, creddb)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
		}
	})
}
