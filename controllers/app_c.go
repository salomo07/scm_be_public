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

// func AddAccessMenuBulk(ctx *fasthttp.RequestCtx, user models.User) {
// 	if string(ctx.Request.Body()) == "" {
// 		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
// 		return
// 	}
// 	type AccessMenuModified struct {
// 		Id            string                 `json:"_id" validate:"required"`
// 		IdRole        string                 `json:"idrole" validate:"required"`
// 		Idmenu        string                 `json:"idmenu" validate:"required"`
// 		AccessSubmenu []models.AccessSubmenu `json:"accesssubmenu" validate:"dive"`
// 		Create        *bool                  `json:"create" validate:"required"`
// 		Read          *bool                  `json:"read" validate:"required"`
// 		Update        *bool                  `json:"update" validate:"required"`
// 		Delete        *bool                  `json:"delete" validate:"required"`
// 		Error         string                 `json:"error" validate:"required"`
// 	}
// 	var accessModelBulk []models.AccessMenu
// 	var accessModelInvalid []AccessMenuModified = nil
// 	i := 0
// 	utils.JsonToStruct(string(ctx.PostBody()), &accessModelBulk)
// 	for _, value := range accessModelBulk {
// 		var accessmenumodified AccessMenuModified
// 		value.Id = "a_" + strconv.FormatInt(time.Now().UnixNano()/1000, 10)
// 		utils.JsonToStruct(utils.StructToJson(value), &accessmenumodified)
// 		err := utils.ValidateRequiredFieldsOld(value)
// 		if err != "" {
// 			accessmenumodified.Error = err
// 			accessModelInvalid = append(accessModelInvalid, accessmenumodified)
// 		} else {
// 			res, err, _ := services.FindOne(user.IdCompany, consts.Coll_AccessMenu, `{"idrole":"`+value.IdRole+`","idmenu":"`+value.Idmenu+`"}`, "", ``)
// 			if err != "" {
// 				accessmenumodified.Error = err
// 				accessModelInvalid = append(accessModelInvalid, accessmenumodified)
// 			} else {
// 				if res == "" {
// 					_, err, _ := services.InsertOne(user.IdCompany, consts.Coll_AccessMenu, utils.StructToJson(value))
// 					if err == "" {
// 						i = i + 1
// 					} else {
// 						accessmenumodified.Error = err
// 						accessModelInvalid = append(accessModelInvalid, accessmenumodified)
// 					}
// 				} else {
// 					accessmenumodified.Error = consts.AccessMenuExist
// 					accessModelInvalid = append(accessModelInvalid, accessmenumodified)
// 				}
// 			}
// 		}
// 	}
// 	utils.ShowResponseJson(ctx, fasthttp.StatusOK, "success", accessModelInvalid)
// }
// func AddAccessMenu(ctx *fasthttp.RequestCtx, user models.User) {
// 	if string(ctx.Request.Body()) == "" {
// 		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
// 		return
// 	}
// 	var accessModel models.AccessMenu
// 	utils.JsonToStruct(string(ctx.PostBody()), &accessModel)
// 	// accessModel.Id = "a_" + strconv.FormatInt(time.Now().UnixNano()/1000, 10)
// 	err := utils.ValidateRequiredFieldsOld(accessModel)
// 	if err != "" {
// 		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", err)
// 	} else {
// 		res, err, _ := services.FindOne(user.IdCompany, consts.Coll_AccessMenu, `{"idmenu":"`+accessModel.Idmenu+`"}`, "", "")
// 		if err != "" {
// 			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
// 			return
// 		} else {
// 			if res == "" {
// 				_, err, _ := services.InsertOne(user.IdCompany, consts.Coll_AccessMenu, utils.StructToJson(accessModel))
// 				if err == "" {
// 					utils.ShowResponseDefault(ctx, fasthttp.StatusOK, "success", consts.AccessMenuAdd)
// 				} else {
// 					utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", consts.FailAddAccessMenu)
// 				}
// 			} else {
// 				utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.AccessMenuExist)
// 			}
// 		}
// 	}
// }
// func UpdateAccessMenu(ctx *fasthttp.RequestCtx, user models.User) {
// 	if string(ctx.Request.Body()) == "" {
// 		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
// 		return
// 	}
// 	var accessModel models.AccessMenu
// 	utils.JsonToStruct(string(ctx.PostBody()), &accessModel)
// 	// accessModel.Id = "a_" + strconv.FormatInt(time.Now().UnixNano()/1000, 10)
// 	err := utils.ValidateRequiredFields(accessModel)
// 	if err != "" {
// 		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", err)
// 	} else {
// 		tempData := utils.RemoveField(accessModel, "_id")
// 		queryUpdateAccessMenu := `{"idmenu":"` + accessModel.Idmenu + `"}`
// 		res, err, _ := services.UpdateOne(user.IdCompany, consts.Coll_AccessMenu, queryUpdateAccessMenu, utils.StructToJson(tempData), false)
// 		if err != "" {
// 			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
// 			return
// 		} else {
// 			utils.ShowResponseJson(ctx, fasthttp.StatusOK, "success", res)
// 			return
// 		}
// 	}
// }

// func DeleteAccessMenu(ctx *fasthttp.RequestCtx, user models.User) {
// 	if string(ctx.Request.Body()) == "" {
// 		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
// 		return
// 	}
// 	var accessModel models.AccessMenu
// 	utils.JsonToStruct(string(ctx.PostBody()), &accessModel)

// 	if accessModel.Id != "" {
// 		res, err, _ := services.DeleteOne(user.IdCompany, consts.Coll_AccessMenu, `{"_id":"`+accessModel.Id+`"}`)
// 		if err != "" {
// 			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
// 			return
// 		} else {
// 			utils.ShowResponseJson(ctx, fasthttp.StatusOK, "success", res)
// 			return
// 		}
// 	} else {
// 		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", consts.AccessDeleteMandatory)
// 	}
// }

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
		res, err, _ := services.FindOneRootDBUsingURI(services.GetURI(credRoot), consts.DB_CORE_NAME, consts.Coll_Menu, `{"url":"`+menuModel.Url+`","appid":"`+menuModel.AppId+`"}`, "", "")
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

				res, err, _ := services.FindOneRootDBUsingURI(services.GetURI(credRoot), consts.DB_CORE_NAME, consts.Coll_Menu, `{"url":"`+value.Url+`","appid":"`+value.AppId+`"}`, "", "")
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
