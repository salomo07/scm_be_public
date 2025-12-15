package admin

import (
	"scm/config"
	"scm/consts"
	"scm/controllers"
	"scm/controllers/master_controller"
	"scm/models"
	"scm/utils"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func BranchRouters(router *fasthttprouter.Router) {
	print("\n -branch\n")
	print(consts.URL_Branch_Upsert + "\n")
	print(consts.URL_Branch_Find + "\n")
	print(consts.URL_Branch_Delete + "\n")

	router.POST(consts.URL_Branch_Upsert, func(ctx *fasthttp.RequestCtx) {
		user, err, _, isCompanyAdmin := controllers.CheckSession(ctx)
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
	router.POST(consts.URL_Branch_Find, func(ctx *fasthttp.RequestCtx) {
		user, err, _, _ := controllers.CheckSession(ctx)
		var dataBranch models.BranchRequest
		utils.JsonToStruct(string(ctx.Request.Body()), &dataBranch)
		ctx.Response.Header.Set("Content-Type", "application/json")
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
			return
		} else if config.DecodingBase64(user.IdCompany) == dataBranch.IdCompany {
			master_controller.FindBranches(ctx, user)
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
