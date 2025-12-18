package controllers

import (
	"scm/consts"
	"scm/models"
	"scm/services"
	"scm/utils"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

func FindManyAccessMenu(ctx *fasthttp.RequestCtx, user models.User) {
	if len(ctx.PostBody()) == 0 {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}

	var accessModel models.AccessMenu
	utils.JsonToStruct(string(ctx.PostBody()), &accessModel)

	res, errStr, code := services.FindMany(
		user.IdCompany,
		consts.Coll_AccessMenu,
		string(ctx.PostBody()),
		`{"idmenu": 1}`,
		0,
		0,
	)

	if errStr != "" {
		utils.ShowResponseDefault(ctx, code, "error", errStr)
		return
	}

	// langsung return hasil query (slice of bson.M atau struct tergantung implementasi services)
	utils.ShowResponseJson(ctx, code, "success", res)
}

func AddMenu(ctx *fasthttp.RequestCtx, user models.User) {
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}
	var menuModel models.Menu
	utils.JsonToStruct(string(ctx.PostBody()), &menuModel)
	menuModel.Id = "m_" + strconv.FormatInt(time.Now().UnixNano()/1000, 10)
	credRoot := utils.GetMongoDBRoot()
	err := utils.ValidateRequiredFields(menuModel)
	if err != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", err)
	} else {

		//check urlmenu sdh isExist g
		res, err, _ := services.FindOneRootDBUsingURI(services.GetURI(credRoot), consts.DB_CORE_NAME, consts.Coll_Menu, `{"url":"`+menuModel.Url+`"}`, "", "")
		if err == "" && res == nil {
			menuModel.Id = "m_" + strconv.FormatInt(time.Now().UnixNano()/1000, 10)
			//JIka url tidak duplicate, maka insert
			_, err, _ := services.InsertOne(user.IdCompany, consts.Coll_Menu, utils.StructToJson(menuModel))
			if err == "" {
				utils.ShowResponseDefault(ctx, fasthttp.StatusOK, "success", consts.MenuAdd)
			} else {
				utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", consts.FailAddMenu)
			}
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.URLAlreadyExist)
		}
	}
}

func AddMenuBulk(ctx *fasthttp.RequestCtx, user models.User) {
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}
	var menusModel []models.Menu
	errBody := utils.JsonToStruct(string(ctx.PostBody()), &menusModel)
	if errBody != nil {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", errBody.Error())
		return
	}
	for i, _ := range menusModel {
		menusModel[i].Id = "m_" + strconv.FormatInt(time.Now().UnixNano()/1000, 10)
	}
	err := utils.ValidateRequiredFields(menusModel)
	if err != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", err)
		return
	} else {
		var arrayDuplicate []models.Menu
		var arrayNotValid []models.Menu
		for _, value := range menusModel {
			err := utils.ValidateRequiredFields(value)
			if err != "" {
				arrayNotValid = append(arrayNotValid, value)
			} else {
				credRoot := utils.GetMongoDBRoot()
				// Check urlmenu sdh isExist g, pengecekan 1 persatu menu. Ini berat, namun harusnya tidak masalah, karena endpoint bulk ini jarang dipakai, hanya superadmin saja.

				res, err, _ := services.FindOneRootDBUsingURI(services.GetURI(credRoot), consts.DB_CORE_NAME, consts.Coll_Menu, `{"url":"`+value.Url+`"}`, "", "")
				if err == "" && res == nil {
				} else {
					arrayDuplicate = append(arrayDuplicate, value)
				}
			}
		}
		// credRoot:=GetMongoDBRoot()
		// services.FindMany(services.GetURI(credRoot),credRoot.DBName,consts.Coll_Menu,`{"ur}`)
		if len(arrayDuplicate) > 0 {
			utils.ShowResponseJson(ctx, fasthttp.StatusBadRequest, "There are duplicate data", arrayDuplicate)
		} else if len(arrayNotValid) > 0 {
			utils.ShowResponseJson(ctx, fasthttp.StatusBadRequest, "There are invalid data", arrayNotValid)
		} else {
			credRoot := utils.GetMongoDBRoot()
			services.InsertManyUsingURI(services.GetURI(credRoot), credRoot.DBName, consts.Coll_Menu, utils.StructToJson(menusModel))
			utils.ShowResponseDefault(ctx, fasthttp.StatusOK, "success", strconv.Itoa(len(menusModel))+" datas was inserted")
		}
	}
}

func ShowAllMenuBulk(ctx *fasthttp.RequestCtx, user models.User) {
	res, errStr, code := services.FindMany(
		user.IdCompany,
		consts.Coll_Menu,
		string(ctx.PostBody()),
		"",
		0,
		0,
	)

	if errStr != "" {
		utils.ShowResponseDefault(ctx, code, "error", errStr)
		return
	}

	utils.ShowResponseJson(ctx, code, "success", res)
}

// Ini hanya boleh diakses oleh SuperAdmin
func UpdateDocument(ctx *fasthttp.RequestCtx, creddb models.CredDB) {
	var inputUpdate models.UpdateDocumentBySuperAdmin
	utils.JsonToStruct(string(ctx.PostBody()), &inputUpdate)
	err := utils.ValidateRequiredFields(inputUpdate)
	if err != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", err)
	} else {
		credRoot := utils.GetMongoDBRoot()
		services.UpdateOneUsingURI(services.GetURI(credRoot), inputUpdate.DBName, inputUpdate.CollName, inputUpdate.Query, inputUpdate.UpdateField, false)
	}
}

// Ini hanya boleh diakses oleh SuperAdmin
func InsertDocument(ctx *fasthttp.RequestCtx, user models.User) {
	credRoot := utils.GetMongoDBRoot()
	var inputInsert models.InsertDocumentBySuperAdmin
	utils.JsonToStruct(string(ctx.PostBody()), &inputInsert)
	err := utils.ValidateRequiredFields(inputInsert)
	if err != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", err)
	} else {
		services.InsertOneUsingURI(services.GetURI(credRoot), inputInsert.DBName, inputInsert.CollName, inputInsert.Document)
	}
}

// Ini hanya boleh diakses oleh SuperAdmin
func FindDocument(ctx *fasthttp.RequestCtx, user models.User) {
	var findDocument models.FindDocumentBySuperAdmin
	credRoot := utils.GetMongoDBRoot()
	utils.JsonToStruct(string(ctx.PostBody()), &findDocument)
	err := utils.ValidateRequiredFields(findDocument)
	if err != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", err)
	} else {
		services.FindOneRootDBUsingURI(services.GetURI(credRoot), findDocument.DBName, findDocument.CollName, findDocument.Query, findDocument.Projection, "")
	}
}
