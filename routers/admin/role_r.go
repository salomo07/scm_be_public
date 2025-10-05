package admin

import (
	"scm/consts"
	"scm/controllers"
	"scm/controllers/master_controller"
	"scm/utils"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func RolesRouters(router *fasthttprouter.Router) {
	print("\n -roles\n")
	print(" -" + consts.URL_Role_Find + "\n")
	print(" -" + consts.URL_Role_Upsert + "\n")
	print(" -" + consts.URL_Role_Delete + "\n")

	router.POST(consts.URL_Role_Find, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		user, err, isSuperAdmin, isCompanyAdmin := controllers.CheckSession(ctx)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
			return
		} else if isSuperAdmin || isCompanyAdmin {
			master_controller.FindRoles(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
			return
		}
	})
	router.POST(consts.URL_Role_Upsert, func(ctx *fasthttp.RequestCtx) {
		user, err, _, isCompanyAdmin := controllers.CheckSession(ctx)
		ctx.Response.Header.Set("Content-Type", "application/json")
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
			return
		} else if isCompanyAdmin {
			master_controller.UpsertRole(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
		}
	})
	router.POST(consts.URL_Role_Delete, func(ctx *fasthttp.RequestCtx) {
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
