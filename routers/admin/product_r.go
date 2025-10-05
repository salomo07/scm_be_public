package admin

import (
	"scm/consts"
	"scm/controllers"
	"scm/controllers/master_controller"
	"scm/utils"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func ProductRouters(router *fasthttprouter.Router) {
	print("\n -product\n")
	print(" -" + consts.URL_Product_Upsert + "\n")
	print(" -" + consts.URL_Product_Delete + "\n")
	print(" -" + consts.URL_Category_Product_Upsert + "\n")

	router.POST(consts.URL_Category_Product_Upsert, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		user, err, _, isCompanyAdmin := controllers.CheckSession(ctx)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
			return
		} else if isCompanyAdmin {
			master_controller.UpsertCategoryProduct(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
			return
		}
	})
	router.POST(consts.URL_Category_Product_Find, func(ctx *fasthttp.RequestCtx) {
		user, err, _, isCompanyAdmin := controllers.CheckSession(ctx)
		ctx.Response.Header.Set("Content-Type", "application/json")
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
			return
		} else if isCompanyAdmin {
			master_controller.FindCategory(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
		}
	})
	router.POST(consts.URL_Product_Upsert, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		user, err, _, isCompanyAdmin := controllers.CheckSession(ctx)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
			return
		} else if isCompanyAdmin {
			master_controller.UpsertProduct(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
			return
		}
	})
	router.POST(consts.URL_Product_Find, func(ctx *fasthttp.RequestCtx) {
		user, err, _, isCompanyAdmin := controllers.CheckSession(ctx)
		ctx.Response.Header.Set("Content-Type", "application/json")
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
			return
		} else if isCompanyAdmin {
			master_controller.FindProducts(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
		}
	})
	router.POST(consts.URL_Product_Delete, func(ctx *fasthttp.RequestCtx) {
		user, err, _, isCompanyAdmin := controllers.CheckSession(ctx)
		ctx.Response.Header.Set("Content-Type", "application/json")
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
			return
		} else if isCompanyAdmin {
			master_controller.DeleteProduct(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
		}
	})
}
