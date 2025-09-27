package admin

import (
	"scm/consts"
	"scm/controllers"
	"scm/controllers/master_controller"
	"scm/utils"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func CompanyRouters(router *fasthttprouter.Router) {
	// START_SAAS_FEATURE
	router.POST(consts.URL_Company_Create, func(ctx *fasthttp.RequestCtx) {
		// master_controller.RegisterCompany(ctx, controllers.GetMongoDBRoot())
		user, err, isSuperAdmin, _ := controllers.CheckSession(ctx)
		ctx.Response.Header.Set("Content-Type", "application/json")
		//Endpoint ini hanya bisa diakses oleh SuperAdmin (bukan company)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
			return
		} else if isSuperAdmin {
			controllers.RegisterCompany(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
		}
	})
	// END_SAAS_FEATURE

	router.POST(consts.URL_Company_Initiate, func(ctx *fasthttp.RequestCtx) {
		// controllers.CopyInitiateData(adminCred, ctx)
		ctx.Response.Header.Set("Content-Type", "application/json")
	})
	print("\n -company\n")
}
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
		user, err, _, isCompanyAdmin := controllers.CheckSession(ctx)
		ctx.Response.Header.Set("Content-Type", "application/json")
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
			return
		} else if isCompanyAdmin {
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
