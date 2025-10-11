package master_controller

import (
	"fmt"
	"scm/consts"
	"scm/controllers"
	"scm/models"
	"scm/services"
	"scm/utils"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindUsers(ctx *fasthttp.RequestCtx) { //querynya menggunakan "_id"
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}
	queryFindRole := string(ctx.Request.Body())
	credRootDB := utils.GetMongoDBRoot()
	sort, skip, limit := controllers.GetSortSkipLimit(ctx)
	resFindMany, err, code := services.FindManyRootDBUsingURI(services.GetURI(credRootDB), consts.DB_CORE_NAME, consts.Coll_Users, queryFindRole, "", sort, skip, limit)
	if err != "" {
		utils.ShowResponseDefault(ctx, code, "error", err)
	} else {
		if len(resFindMany) == 0 {
			utils.ShowResponseJson(ctx, code, "success", []interface{}{})
			return
		} else {
			utils.ShowResponseJson(ctx, code, "success", resFindMany)
		}
	}
}
func UpsertUser(ctx *fasthttp.RequestCtx) {
	if len(ctx.Request.Body()) == 0 {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}

	var dataRole models.RoleRequest
	if err := utils.JsonToStruct(string(ctx.Request.Body()), &dataRole); err != nil {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", "Invalid JSON body")
		return
	}

	// Validasi field wajib
	if msg := utils.ValidateRequiredFields(dataRole); msg != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, msg, "")
		return
	}

	// Buat filter
	var filter string
	if dataRole.Id != "" {
		// User kirim _id string → convert ke filter MongoDB
		if !primitive.IsValidObjectID(dataRole.Id) {
			utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", "Invalid _id format")
			return
		}
		// manual bikin JSON supaya tetap ObjectID
		filter = fmt.Sprintf(`{"_id":{"$oid":"%s"}}`, dataRole.Id)
	} else {
		// Tidak ada _id → biarkan kosong supaya jadi insert
		filter = `{"_id":{"$exists":false}}`
	}

	// Data untuk update/insert (hapus field _id agar tidak overwrite)
	removedId := utils.RemoveField(dataRole, "_id")
	data := utils.StructToJson(removedId)

	// Debug log
	fmt.Println("FILTER JSON:", filter)
	fmt.Println("DATA JSON:", data)

	credRoot := utils.GetMongoDBRoot()
	resUpsert, err, code := services.UpdateOneUsingURI(
		services.GetURI(credRoot), consts.DB_CORE_NAME,
		consts.Coll_Users,
		filter,
		data,
		true,
	)

	if resUpsert == nil {
		utils.ShowResponseDefault(ctx, code, "error", err)
	} else {
		// Response jangan diubah
		utils.ShowResponseJson(ctx, code, "success", resUpsert)
	}
}

func FindRoles(ctx *fasthttp.RequestCtx) { //querynya menggunakan "_id"
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}
	queryFindRole := string(ctx.Request.Body())
	credRootDB := utils.GetMongoDBRoot()
	sort, skip, limit := controllers.GetSortSkipLimit(ctx)
	resFindMany, err, code := services.FindManyRootDBUsingURI(services.GetURI(credRootDB), consts.DB_CORE_NAME, consts.Coll_Role, queryFindRole, "", sort, skip, limit)
	if err != "" {
		utils.ShowResponseDefault(ctx, code, "error", err)
	} else {
		if len(resFindMany) == 0 {
			utils.ShowResponseJson(ctx, code, "success", []interface{}{})
			return
		} else {
			utils.ShowResponseJson(ctx, code, "success", resFindMany)
		}
	}
}

func UpsertRole(ctx *fasthttp.RequestCtx) {
	if len(ctx.Request.Body()) == 0 {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}

	var dataRole models.RoleRequest
	if err := utils.JsonToStruct(string(ctx.Request.Body()), &dataRole); err != nil {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", "Invalid JSON body")
		return
	}

	// Validasi field wajib
	if msg := utils.ValidateRequiredFields(dataRole); msg != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, msg, "")
		return
	}

	// Buat filter
	var filter string
	if dataRole.Id != "" {
		// User kirim _id string → convert ke filter MongoDB
		if !primitive.IsValidObjectID(dataRole.Id) {
			utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", "Invalid _id format")
			return
		}
		// manual bikin JSON supaya tetap ObjectID
		filter = fmt.Sprintf(`{"_id":{"$oid":"%s"}}`, dataRole.Id)
	} else {
		// Tidak ada _id → biarkan kosong supaya jadi insert
		filter = `{"_id":{"$exists":false}}`
	}

	// Data untuk update/insert (hapus field _id agar tidak overwrite)
	removedId := utils.RemoveField(dataRole, "_id")
	data := utils.StructToJson(removedId)

	// Debug log
	fmt.Println("FILTER JSON:", filter)
	fmt.Println("DATA JSON:", data)

	credRoot := utils.GetMongoDBRoot()
	resUpsert, err, code := services.UpdateOneUsingURI(
		services.GetURI(credRoot), consts.DB_CORE_NAME,
		consts.Coll_Role,
		filter,
		data,
		true,
	)

	if resUpsert == nil {
		utils.ShowResponseDefault(ctx, code, "error", err)
	} else {
		// Response jangan diubah
		utils.ShowResponseJson(ctx, code, "success", resUpsert)
	}
}

func DeleteRoles(ctx *fasthttp.RequestCtx, idcompany string) { //querynya menggunakan "_id"
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}
	var dataProduct models.Branch
	utils.JsonToStruct(string(ctx.Request.Body()), &dataProduct)

	if msg := utils.ValidateRequiredFields(dataProduct); msg != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, msg, "")
		return
	}
	queryBranch := `{"idbranch":"` + dataProduct.Id + `"}`
	resUpsert, err, code := services.DeleteOne(idcompany, consts.Coll_Products, queryBranch)
	if err != "" {
		utils.ShowResponseDefault(ctx, code, "error", err)
	} else {
		utils.ShowResponseJson(ctx, code, "success", resUpsert)

	}
}
