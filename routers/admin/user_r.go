package admin

import (
	"scm/consts"
	"scm/controllers"
	"scm/controllers/master_controller"
	"scm/utils"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func UserRouters(router *fasthttprouter.Router) {
	print("\n -user\n")
	print(" -" + consts.URL_User_Find + "\n")
	print(" -" + consts.URL_User_Delete + "\n")
	print(" -" + consts.URL_User_Upsert + "\n")

	router.POST(consts.URL_User_Find, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		_, err, isSuperAdmin, isCompanyAdmin := controllers.CheckSession(ctx)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
			return
		} else if isSuperAdmin || isCompanyAdmin {
			master_controller.FindUsers(ctx)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
			return
		}
	})
	router.POST(consts.URL_User_Upsert, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		_, err, isSuperAdmin, isCompanyAdmin := controllers.CheckSession(ctx)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
			return
		} else if isSuperAdmin || isCompanyAdmin {
			master_controller.FindUsers(ctx)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
			return
		}
	})

	// router.POST(consts.URL_User_Find, func(ctx *fasthttp.RequestCtx) {
	// 	ctx.Response.Header.Set("Content-Type", "application/json")
	// 	user, err, isSuperAdmin, isCompanyAdmin := controllers.CheckSession(ctx)
	// 	if err != "" {
	// 		utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
	// 		return
	// 	} else if isSuperAdmin || isCompanyAdmin {
	// 		controllers.GetUserOne(ctx, user)
	// 	} else {
	// 		utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
	// 		return
	// 	}
	// })
	// router.POST(consts.URL_User_FindMany, func(ctx *fasthttp.RequestCtx) {
	// 	ctx.Response.Header.Set("Content-Type", "application/json")
	// 	user, err, isSuperAdmin, isCompanyAdmin := controllers.CheckSession(ctx)
	// 	if err != "" {
	// 		utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
	// 		return
	// 	} else if isSuperAdmin || isCompanyAdmin {
	// 		controllers.GetUserMany(ctx, user)
	// 	} else {
	// 		utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
	// 		return
	// 	}
	// })
	// router.POST(consts.URL_User_Update, func(ctx *fasthttp.RequestCtx) {
	// 	ctx.Response.Header.Set("Content-Type", "application/json")
	// 	user, err, isSuperAdmin, isCompanyAdmin := controllers.CheckSession(ctx)
	// 	if err != "" {
	// 		utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
	// 		return
	// 	} else if isSuperAdmin || isCompanyAdmin {
	// 		controllers.UpdateUser(ctx, user)
	// 	} else {
	// 		utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
	// 		return
	// 	}

	// })
	// router.POST(consts.URL_User_Create, func(ctx *fasthttp.RequestCtx) {
	// 	ctx.Response.Header.Set("Content-Type", "application/json")
	// 	user, err, isSuperAdmin, isCompanyAdmin := controllers.CheckSession(ctx)
	// 	if err != "" {
	// 		utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
	// 		return
	// 	} else if isSuperAdmin || isCompanyAdmin {
	// 		controllers.AddUser(ctx, user)
	// 	} else {
	// 		utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", "You have not access to this endpoint (Company is unregistered).")
	// 	}
	// })
}
