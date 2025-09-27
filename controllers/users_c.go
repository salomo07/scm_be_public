package controllers

import (
	"fmt"
	"scm/config"
	"scm/consts"
	"scm/models"
	"scm/services"
	"scm/utils"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Encrypted fields ["nik","name","nickname","username","contact.email","contact.mobile","contact.phone","contact.whatsApp"]
//Semua field tersebut diencrypt ketika diterima BE, namun sebelum masuk ke DB di decrypt LG.
//Field "password" dikenakan encrypt 2 kali sebelum dimasukkan ke DB

func AddUser(ctx *fasthttp.RequestCtx, user models.User) {
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}
	var userModel models.User
	utils.JsonToStruct(string(ctx.PostBody()), &userModel)
	userModel.Id = "u_" + strconv.FormatInt(time.Now().UnixNano()/1000, 10)
	err := utils.ValidateRequiredFields(userModel)
	if err == "" {
		userModel.Nik = config.DecryptAES(userModel.Nik)
		userModel.Name = config.DecryptAES(userModel.Name)
		userModel.Nickname = config.DecryptAES(userModel.Nickname)
		userModel.Username = config.DecryptAES(userModel.Username)
		userModel.Contact.Email = config.DecryptAES(userModel.Contact.Email)
		userModel.Contact.Mobile = config.DecryptAES(userModel.Contact.Mobile)
		userModel.IdCompany = user.IdCompany

		pipeline := `[{"$facet":{"userCheck":[{"$match":{"$or":[{"username":"` + userModel.Username + `"},{"contact.email":"` + userModel.Contact.Email + `"},{"contact.mobile":"` + userModel.Contact.Mobile + `"}]}},{"$count":"userCount"}],"roleCheck":[{"$lookup":{"from":"companies","pipeline":[{"$unwind":"$roles"},{"$match":{"roles._id":"` + userModel.IdRole + `"}},{"$count":"roleCount"}],"as":"roleData"}},{"$unwind":"$roleData"}]}},{"$project":{"userExist":{"$cond":{"if":{"$gt":[{"$arrayElemAt":["$userCheck.userCount",0]},0]},"then":1,"else":0}},"roleExist":{"$cond":{"if":{"$gt":[{"$arrayElemAt":["$roleCheck.roleData.roleCount",0]},0]},"then":1,"else":0}}}}]`

		type CheckUserExistResponse struct {
			RoleExist int `json:"roleExist"`
			UserExist int `json:"userExist"`
		}
		var userExist CheckUserExistResponse
		credRoot := utils.GetMongoDBRoot()
		resAgg, err, code := services.AggregationOneUsingURI(services.GetURI(credRoot), credRoot.DBName, consts.Coll_Users, pipeline)
		if err != "" {
			utils.ShowResponseDefault(ctx, code, "error", err)
			return
		} else {
			utils.JsonToStruct(resAgg, &userExist)
			if userExist.UserExist > 0 {
				utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.UserAlreadyExist)
			} else {
				if userExist.RoleExist > 0 {
					// Insert User ke scm_core
					userModel.Password = config.EncryptAES(userModel.Password)
					resInsert, errStr, _ := services.InsertOneUsingURI(services.GetURI(credRoot), credRoot.DBName, consts.Coll_Users, utils.StructToJson(userModel))
					if errStr != "" {
						utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", consts.FailAddUser)
					} else {
						var insertedID string
						if oid, ok := resInsert.InsertedID.(primitive.ObjectID); ok {
							insertedID = oid.Hex()
						} else {
							insertedID = fmt.Sprintf("%v", resInsert.InsertedID)
						}

						userModel.Id = insertedID

						go services.InsertOne(user.IdCompany, consts.Coll_Users, utils.StructToJson(userModel))
						utils.ShowResponseDefault(ctx, fasthttp.StatusOK, "success", consts.UserAdd+" ("+userModel.Id+")")
					}
				} else {
					utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.RoleNotFound)
					return
				}
			}
		}
	} else {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", err)
	}
}
func GetUserOne(ctx *fasthttp.RequestCtx, user models.User) {
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}
	var userModel models.User
	res, err, _ := services.FindOne(user.IdCompany, consts.Coll_Users, string(ctx.PostBody()), `{"password":0}`, "")
	if err == "" {
		if res != "" {
			utils.JsonToStruct(res, &userModel)
			// ["nik","name","nickname","username","contact.email","contact.mobile","contact.phone","contact.whatsApp"]
			userModel.Nik = config.EncryptAES(userModel.Nik)
			userModel.Name = config.EncryptAES(userModel.Name)
			userModel.Nickname = config.EncryptAES(userModel.Nickname)
			userModel.Username = config.EncryptAES(userModel.Username)
			userModel.Contact.Email = config.EncryptAES(userModel.Contact.Email)
			userModel.Contact.Mobile = config.EncryptAES(userModel.Contact.Mobile)
			userModel.Contact.Phone = config.EncryptAES(userModel.Contact.Phone)
			userModel.Contact.WhatsApp = config.EncryptAES(userModel.Contact.WhatsApp)
			userModel.Password = "**hiddenpassword**"
			utils.ShowResponseJson(ctx, fasthttp.StatusOK, "success", userModel)
		} else {
			var result map[string]interface{}
			utils.JsonToStruct("{}", &result)
			utils.ShowResponseJson(ctx, fasthttp.StatusOK, "success", result)
		}
	} else {
		utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
	}
}
func GetUserMany(ctx *fasthttp.RequestCtx, user models.User) {
	if len(ctx.PostBody()) == 0 {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}

	// query Mongo
	res, errStr, code := services.FindMany(
		user.IdCompany,
		consts.Coll_Users,
		string(ctx.PostBody()),
		"{}",
		0,
		0,
	)

	if errStr != "" {
		utils.ShowResponseDefault(ctx, code, "error", errStr)
		return
	}

	// kalau kosong
	if len(res) == 0 {
		utils.ShowResponseJson(ctx, fasthttp.StatusOK, "success", []models.User{})
		return
	}

	// konversi hasil query ke struct
	var users []models.User
	for _, doc := range res {
		var u models.User
		bsonBytes, _ := bson.Marshal(doc)
		_ = bson.Unmarshal(bsonBytes, &u)

		// enkripsi field sensitif
		u.Nik = config.EncryptAES(u.Nik)
		u.Name = config.EncryptAES(u.Name)
		u.Nickname = config.EncryptAES(u.Nickname)
		u.Username = config.EncryptAES(u.Username)
		u.Contact.Email = config.EncryptAES(u.Contact.Email)
		u.Contact.Mobile = config.EncryptAES(u.Contact.Mobile)
		u.Contact.Phone = config.EncryptAES(u.Contact.Phone)
		u.Contact.WhatsApp = config.EncryptAES(u.Contact.WhatsApp)

		// sembunyikan password
		u.Password = "**hiddenpassword**"

		users = append(users, u)
	}

	utils.ShowResponseJson(ctx, fasthttp.StatusOK, "success", users)
}

func UpdateUser(ctx *fasthttp.RequestCtx, user models.User) {
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}
	type UserMandatory struct { //Field mandatory
		Id        string `json:"_id" validate:"required"`
		IdCompany string `json:"idcompany" validate:"required"`
	}
	// var userModel models.User
	var userModelMandatory UserMandatory
	// utils.JsonToStruct(string(ctx.PostBody()), &userModel)
	utils.JsonToStruct(string(ctx.PostBody()), &userModelMandatory)
	err := utils.ValidateRequiredFields(userModelMandatory)
	if err == "" {
		res, err, _ := services.UpdateOne(user.IdCompany, consts.Coll_Users, utils.StructToJson(userModelMandatory), string(ctx.PostBody()), false)
		if err != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", err)
			return
		} else {
			utils.ShowResponseJson(ctx, fasthttp.StatusOK, "success", res)
			return
		}
	} else {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", err)
		return
	}
}
func AddRole(ctx *fasthttp.RequestCtx, user models.User) {
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}
	var roleModel models.Role
	utils.JsonToStruct(string(ctx.PostBody()), &roleModel)
	roleModel.Id = strconv.FormatInt(time.Now().UnixNano()/1000, 10)
	err := utils.ValidateRequiredFields(roleModel)
	if err == "" {
		credRoot := utils.GetMongoDBRoot()
		query := `{"$and":[{"$or":[{"roles.name":"` + roleModel.Name + `"},{"roles.code":"` + roleModel.Code + `"}]},{"_id":"` + user.IdCompany + `"}]}`
		res, errFind, _ := services.FindOneRootDBUsingURI(services.GetURI(credRoot), credRoot.DBName, consts.Coll_Companies, query, "", "")
		if res == "" && errFind == "" {
			//Jika tidak ditemukan, ambil dulu Roles yang terdapat di "companies"
			resBody, errStr, statuscode := services.FindOneRootDBUsingURI(services.GetURI(credRoot), consts.DB_CORE_NAME, consts.Coll_Companies, `{"_id":"`+user.IdCompany+`"}`, "", "")
			if resBody != "" {
				// Insert data Role di DB "scm_core->companies"
				var companyObj models.Company
				utils.JsonToStruct(resBody, &companyObj)

				_, errUp, _ := services.UpdateOneUsingURI(services.GetURI(credRoot), credRoot.DBName, consts.Coll_Companies, `{"_id":"`+user.IdCompany+`"}`, `{"$addToSet":{"roles":`+utils.StructToJson(roleModel)+`}}`, false)
				if errUp != "" {
					utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", consts.FailAddRole)
				} else {
					go services.InsertOne(user.IdCompany, consts.Coll_Role, utils.StructToJson(roleModel))
					utils.ShowResponseJson(ctx, fasthttp.StatusOK, "success", utils.StructToJson(roleModel))
					return
				}
			} else {
				utils.ShowResponseDefault(ctx, statuscode, "error", errStr)
			}
		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, consts.CodeCompanyExist, "")
		}
	} else {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, err, "")
	}
}
