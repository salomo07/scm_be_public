package master_controller

import (
	"fmt"
	"scm/consts"
	"scm/models"
	"scm/services"
	"scm/utils"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

func GenerateQR(ctx *fasthttp.RequestCtx, user models.User) {
	type ReqGenerateQR struct {
		Prefix   string `json:"prefix"`
		JumlahHU int    `json:"jumlahhu" validate:"required"`
		TypeHUID string `json:"idtypehu" validate:"required"`
	}

	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}

	var req ReqGenerateQR
	utils.JsonToStruct(string(ctx.Request.Body()), &req)

	if msg := utils.ValidateRequiredFields(req); msg != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, msg, "")
		return
	}

	// Ambil HU terakhir berdasarkan type
	var lastHU models.HandlingUnits
	query := fmt.Sprintf(`{"is_active":true,"handling_unit_type_id":"%s"}`, req.TypeHUID)
	res, errStr, _ := services.FindOne(
		user.IdCompany,
		consts.Coll_Handling_Units,
		query,
		"{}",
		`{"seq_num": -1}`,
	)

	if errStr != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", errStr)
		return
	}

	utils.JsonToStruct(res, &lastHU)
	startSeq := lastHU.Seq_Num
	currentPrefix := req.Prefix
	if currentPrefix == "" {
		if lastHU.HUCode != "" && lastHU.Seq_Num != 0 {
			currentPrefix = strings.TrimSuffix(lastHU.HUCode, fmt.Sprintf("%d", lastHU.Seq_Num))
		}
		if currentPrefix == "" {
			currentPrefix = "HU-"
		}
	}

	// Buat array map hasil generate tanpa _id
	var huList []map[string]interface{}
	now := int(time.Now().Unix())

	for i := 1; i <= req.JumlahHU; i++ {
		hu := models.HandlingUnits{
			Seq_Num:            startSeq + i,
			CreateTime:         now,
			HandlingUnitTypeId: req.TypeHUID,
			HUCode:             fmt.Sprintf("%s%d", currentPrefix, startSeq+i),
			IsActive:           true,
		}

		// Hapus _id
		huMap := utils.RemoveField(hu, "_id")
		huList = append(huList, huMap)
	}

	// Insert ke Mongo
	insertRes, insertErr, _ := services.InsertMany(
		user.IdCompany,
		consts.Coll_Handling_Units,
		utils.StructToJson(huList),
	)

	if insertErr != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusOK, "success", insertErr)
		return
	}

	utils.ShowResponseJson(ctx, fasthttp.StatusOK, "success", insertRes)
}
func InsertHUType(ctx *fasthttp.RequestCtx, user models.User) {
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}
	var hutype models.HandlingUnitType
	hutype.IsActive = true
	utils.JsonToStruct(string(ctx.Request.Body()), &hutype)
	if utils.ValidateRequiredFields(hutype) != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, utils.ValidateRequiredFields(hutype), "")
	} else {
		res, err, _ := services.FindOne(user.IdCompany, consts.Coll_Handling_Unit_Types, `{"name":"`+hutype.Name+`"}`, "", "")
		if err == "" {
			if res != "" {
				utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.HUTypeMustUnique)
			} else {
				tempData := utils.RemoveField(hutype, "_id")
				_, errStr, _ := services.InsertOne(user.IdCompany, consts.Coll_Handling_Unit_Types, utils.StructToJson(tempData))
				if errStr == "" {
					utils.ShowResponseDefault(ctx, fasthttp.StatusOK, "success", consts.HUTypeSuccessInsert)
				} else {
					utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", errStr)
				}
			}
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
		}
	}
}
func InsertProductMenu(ctx *fasthttp.RequestCtx, user models.User) {
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}
	var product_menu models.PruductMenu
	product_menu.IsActive = true
	utils.JsonToStruct(string(ctx.Request.Body()), &product_menu)
	if utils.ValidateRequiredFields(product_menu) != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, utils.ValidateRequiredFields(product_menu), "")
	} else {
		res, err, _ := services.FindOne(user.IdCompany, consts.Coll_Product_Menu, `{"name":"`+product_menu.Name+`"}`, "", "")
		if err == "" {
			if res != "" {
				utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.ProductMenuMustUnique)
			} else {
				tempData := utils.RemoveField(product_menu, "_id")
				_, errStr, _ := services.InsertOne(user.IdCompany, consts.Coll_Product_Menu, utils.StructToJson(tempData))
				if errStr == "" {
					utils.ShowResponseDefault(ctx, fasthttp.StatusOK, "success", consts.ProductMenuSuccessInsert)
				} else {
					utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", errStr)
				}
			}
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
		}
	}
}

func GetHUAll(ctx *fasthttp.RequestCtx, user models.User) {
	var paginationReq models.PaginationRequest
	utils.JsonToStruct(string(ctx.Request.Body()), &paginationReq)

	// validasi
	if utils.ValidateRequiredFields(paginationReq) != "" {
		utils.ShowResponseDefault(
			ctx,
			fasthttp.StatusBadRequest,
			utils.ValidateRequiredFields(paginationReq),
			"",
		)
		return
	}

	// tentukan urutan sort
	sortOrder := 1
	if !paginationReq.IsASC {
		sortOrder = -1
	}
	sortJSON := fmt.Sprintf(`{"seq_num": %d}`, sortOrder)
	filterJSON := `{"is_active": true}`

	// panggil service
	res, errStr, code := services.FindMany(
		user.IdCompany,
		consts.Coll_Handling_Units,
		filterJSON,
		sortJSON,
		paginationReq.Skip,
		paginationReq.Limit,
	)

	// cek error
	if errStr != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", errStr)
		return
	}

	// return langsung hasil slice (bukan string JSON lagi)
	utils.ShowResponseJson(ctx, code, "success", res)
}

func GetHUType(ctx *fasthttp.RequestCtx, user models.User) {
	res, errStr, code := services.FindMany(
		user.IdCompany,
		consts.Coll_Handling_Unit_Types,
		`{"is_active": true}`,
		``,
		0,
		0,
	)

	// jika ada error
	if errStr != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", errStr)
		return
	}

	// langsung return slice hasil query
	utils.ShowResponseJson(ctx, code, "success", res)
}
