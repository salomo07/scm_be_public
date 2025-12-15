package routers

import (
	"scm/consts"
	"scm/controllers"
	"scm/routers/admin"
	"scm/utils"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func AdminRouters(router *fasthttprouter.Router) {

	// master_routers(router)
	accessmenu_routers(router)
	// role_routers(router)
	menu_routers(router)

	print("--ADMIN Router--\n")
	admin.UserRouters(router)
	admin.CompanyRouters(router)
	admin.RolesRouters(router)
	admin.ProductRouters(router)
	injectToDB(router)
}

// func master_routers(router *fasthttprouter.Router) {
// 	// START_SAAS_FEATURE
// 	router.POST(consts.URL_Master_AddHUType, func(ctx *fasthttp.RequestCtx) {
// 		user, err, _, isOwner := controllers.CheckSession(ctx)
// 		ctx.Response.Header.Set("Content-Type", "application/json")
// 		if err != "" {
// 			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
// 			return
// 		} else if isOwner {
// 			master_controller.InsertHUType(ctx, user)
// 		} else {
// 			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
// 		}
// 	})
// 	router.POST(consts.URL_Master_AddProduct_Menu, func(ctx *fasthttp.RequestCtx) {
// 		user, err, _, isOwner := controllers.CheckSession(ctx)
// 		ctx.Response.Header.Set("Content-Type", "application/json")
// 		if err != "" {
// 			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
// 			return
// 		} else if isOwner {
// 			master_controller.InsertProductMenu(ctx, user)
// 		} else {
// 			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
// 		}
// 	})
// 	router.POST(consts.URL_Master_AddHU, func(ctx *fasthttp.RequestCtx) {
// 		user, err, _, _ := controllers.CheckSession(ctx)
// 		ctx.Response.Header.Set("Content-Type", "application/json")
// 		if err != "" {
// 			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
// 			return
// 		} else {
// 			master_controller.GenerateQR(ctx, user)
// 		}
// 	})
// 	router.POST(consts.URL_Master_GetHUType, func(ctx *fasthttp.RequestCtx) {
// 		user, err, _, isOwner := controllers.CheckSession(ctx)
// 		ctx.Response.Header.Set("Content-Type", "application/json")
// 		if err != "" {
// 			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
// 			return
// 		} else if isOwner {
// 			master_controller.GetHUType(ctx, user)
// 		} else {
// 			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
// 		}
// 	})
// 	router.POST(consts.URL_Master_GetHU, func(ctx *fasthttp.RequestCtx) {
// 		user, err, _, _ := controllers.CheckSession(ctx)
// 		ctx.Response.Header.Set("Content-Type", "application/json")
// 		if err != "" {
// 			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, consts.AdminKeyTidakDikenali+" / terjadi error")
// 			return
// 		} else {
// 			master_controller.GetHUAll(ctx, user)
// 		}
// 	})
// 	print(" -master\n")
// }

func accessmenu_routers(router *fasthttprouter.Router) {
	router.POST(consts.URL_AccessMenu_Create, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		_, err, isSuperAdmin, isCompanyAdmin := controllers.CheckSession(ctx)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
			return
		} else if isSuperAdmin || isCompanyAdmin {
			// controllers.AddAccessMenu(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
			return
		}
	})
	router.POST(consts.URL_AccessMenu_CreateMany, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		_, err, isSuperAdmin, isCompanyAdmin := controllers.CheckSession(ctx)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
			return
		} else if isSuperAdmin || isCompanyAdmin {
			// controllers.AddAccessMenuBulk(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
			return
		}
	})
	router.POST(consts.URL_AccessMenu_FindMany, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		user, err, isSuperAdmin, isCompanyAdmin := controllers.CheckSession(ctx)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
			return
		} else if isSuperAdmin || isCompanyAdmin {
			controllers.FindManyAccessMenu(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
			return
		}
	})
	router.POST(consts.URL_AccessMenu_Update, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		_, err, isSuperAdmin, isCompanyAdmin := controllers.CheckSession(ctx)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
			return
		} else if isSuperAdmin || isCompanyAdmin {
			// controllers.UpdateAccessMenu(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
			return
		}
	})
	router.POST(consts.URL_AccessMenu_Delete, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		_, err, isSuperAdmin, isCompanyAdmin := controllers.CheckSession(ctx)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
			return
		} else if isSuperAdmin || isCompanyAdmin {
			// controllers.DeleteAccessMenu(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
			return
		}
	})
	print(" -accessmenu\n")
}
func menu_routers(router *fasthttprouter.Router) {
	router.POST(consts.URL_Menu_Create, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		user, err, isSuperAdmin, _ := controllers.CheckSession(ctx)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, err, "")
			return
		} else if isSuperAdmin {
			controllers.AddMenu(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
		}
	})
	router.POST(consts.URL_Menu_Create_Bulk, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		user, err, isSuperAdmin, _ := controllers.CheckSession(ctx)
		if isSuperAdmin && err == "" {
			controllers.AddMenuBulk(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
		}
	})
	router.POST(consts.URL_Menu_Show_All, func(ctx *fasthttp.RequestCtx) {
		user, err, isSuperAdmin, _ := controllers.CheckSession(ctx)
		ctx.Response.Header.Set("Content-Type", "application/json")
		//Endpoint ini hanya bisa diakses oleh SuperAdmin (bukan company)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, err, consts.AdminKeyTidakDikenali)
			return
		}
		if isSuperAdmin {
			controllers.ShowAllMenuBulk(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
		}
	})
	print(" -menu\n")
}

//	func role_routers(router *fasthttprouter.Router) {
//		router.POST(consts.URL_Role_Create, func(ctx *fasthttp.RequestCtx) {
//			ctx.Response.Header.Set("Content-Type", "application/json")
//			user, err, isSuperAdmin, isCompanyAdmin := controllers.CheckSession(ctx)
//			if err != "" {
//				utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
//				return
//			} else if isSuperAdmin || isCompanyAdmin {
//				controllers.AddRole(ctx, user)
//			} else {
//				utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
//				return
//			}
//		})
//		print(" -role\n")
//	}
func injectToDB(router *fasthttprouter.Router) {
	router.POST(consts.URL_Insert_DB, func(ctx *fasthttp.RequestCtx) {
		user, err, isSuperAdmin, _ := controllers.CheckSession(ctx)
		ctx.Response.Header.Set("Content-Type", "application/json")
		//Endpoint ini hanya bisa diakses oleh SuperAdmin (bukan company)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, err, consts.AdminKeyTidakDikenali)
			return
		}
		if isSuperAdmin {
			controllers.InsertDocument(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
		}
	})
	router.POST(consts.URL_Find_DB, func(ctx *fasthttp.RequestCtx) {
		user, err, isSuperAdmin, _ := controllers.CheckSession(ctx)
		ctx.Response.Header.Set("Content-Type", "application/json")
		//Endpoint ini hanya bisa diakses oleh SuperAdmin (bukan company)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, err, consts.AdminKeyTidakDikenali)
			return
		}
		if isSuperAdmin {
			controllers.FindDocument(ctx, user)
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
		}
	})
	// router.POST(consts.URL_Update_DB, func(ctx *fasthttp.RequestCtx) {
	// 	creddb, err, isSuperAdmin, _ := controllers.CheckSession(ctx)
	// 	ctx.Response.Header.Set("Content-Type", "application/json")
	// 	//Endpoint ini hanya bisa diakses oleh SuperAdmin (bukan company)
	// 	if err != "" {
	// 		utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, err, consts.AdminKeyTidakDikenali)
	// 		return
	// 	}
	// 	if isSuperAdmin {
	// 		controllers.UpdateDocument(ctx, creddb)
	// 	} else {
	// 		utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.Unauthorized)
	// 	}
	// })
}
