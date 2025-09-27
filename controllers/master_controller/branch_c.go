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

func FindBranches(ctx *fasthttp.RequestCtx, user models.User) { //querynya menggunakan "_id"
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}
	queryFindProduct := string(ctx.Request.Body())

	sort, skip, limit := controllers.GetSortSkipLimit(ctx)
	resFindMany, err, code := services.FindMany(user.IdCompany, consts.Coll_Branch, queryFindProduct, sort, skip, limit)
	if err != "" {
		utils.ShowResponseDefault(ctx, code, "error", err)
	} else {
		if len(resFindMany) == 0 {
			utils.ShowResponseJson(ctx, code, "success", []interface{}{})
			return
		}
		utils.ShowResponseJson(ctx, code, "success", resFindMany)
	}
}

func UpsertBranch(ctx *fasthttp.RequestCtx, user models.User) {
	if len(ctx.Request.Body()) == 0 {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}

	var dataBranch models.BranchRequest
	if err := utils.JsonToStruct(string(ctx.Request.Body()), &dataBranch); err != nil {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", "Invalid JSON body")
		return
	}

	// Validasi field wajib
	if msg := utils.ValidateRequiredFields(dataBranch); msg != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, msg, "")
		return
	}

	// Buat filter
	var filter string
	if dataBranch.Id != "" {
		// User kirim _id string → convert ke filter MongoDB
		if !primitive.IsValidObjectID(dataBranch.Id) {
			utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", "Invalid _id format")
			return
		}
		// manual bikin JSON supaya tetap ObjectID
		filter = fmt.Sprintf(`{"_id":{"$oid":"%s"}}`, dataBranch.Id)
	} else {
		// Tidak ada _id → biarkan kosong supaya jadi insert
		filter = `{"_id":{"$exists":false}}`
	}

	// Data untuk update/insert (hapus field _id agar tidak overwrite)
	removedId := utils.RemoveField(dataBranch, "_id")
	data := utils.StructToJson(removedId)

	// Debug log
	fmt.Println("FILTER JSON:", filter)
	fmt.Println("DATA JSON:", data)

	// Jalankan upsert
	resInsert, err, code := services.UpdateOne(
		user.IdCompany,
		consts.Coll_Branch,
		filter,
		data,
		true, // upsert
	)

	if resInsert == nil {
		utils.ShowResponseDefault(ctx, code, "error", err)
	} else {
		// Response jangan diubah
		utils.ShowResponseJson(ctx, code, "success", resInsert)
	}
}

func DeleteBranches(ctx *fasthttp.RequestCtx, idcompany string) { //querynya menggunakan "_id"
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
