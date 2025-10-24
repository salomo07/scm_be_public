package controllers

import (
	"encoding/json"
	"os"
	"scm/config"
	"scm/consts"
	"scm/models"
	"scm/services"
	"scm/utils"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

// START_SAAS_FEATURE
func RegisterCompany(ctx *fasthttp.RequestCtx, user models.User) {
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.Unauthorized)
		return
	}
	var companyModel models.Company
	var userOwner models.User
	var userOwnerExist models.User
	companyModel.Owner.Id = strconv.FormatInt(time.Now().UnixNano()/1000, 10)
	utils.JsonToStruct(string(ctx.PostBody()), &companyModel)
	credRoot := utils.GetMongoDBRoot()
	errCompanyMandatory := utils.ValidateRequiredFields(companyModel) // Cek mandatory di data Company
	if errCompanyMandatory != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", errCompanyMandatory)
		return
	} else if companyModel.Owner.Username != "" {
		userOwner = companyModel.Owner
		errUserMandatory := utils.ValidateRequiredFields(userOwner) // Cek mandatory di data User
		if errUserMandatory != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", errUserMandatory)
			return
		}

		// credRoot := creddbRoot
		existUserResult, errExistUserResult, _ := services.FindOneRootDBUsingURI(services.GetURI(credRoot), credRoot.DBName, consts.Coll_Users, `{"username":"`+config.DecryptAES(userOwner.Username)+`"}`, "", "")
		if errExistUserResult != "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "warning", consts.FailGetUser)
			return
		} else {
			if existUserResult != nil {
				jsonBytes, _ := json.Marshal(existUserResult)
				utils.JsonToStruct(string(jsonBytes), &userOwnerExist)
			}
			if userOwnerExist.Username != "" {
				utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.UserAlreadyExist)
				return
			}
		}
	}

	query := `{"alias":"` + companyModel.Alias + `","appid":"` + companyModel.AppId + `"}`
	existCompany, errFind, statuscode := services.FindOneRootDBUsingURI(services.GetURI(credRoot), consts.DB_CORE_NAME, consts.Coll_Companies, query, "", "")
	jsonBytes, _ := json.Marshal(existCompany)
	utils.JsonToStruct(string(jsonBytes), &companyModel)
	if errFind != "" {
		utils.ShowResponseDefault(ctx, statuscode, "error", errFind)
		return
	} else if companyModel.IdCompany != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.CompanyAlreadyExist)
		return
	} else {
		// Insert document company
		id := strconv.FormatInt(time.Now().UnixNano()/1000, 10)
		companyModel.IdCompany = "c_" + id
		if companyModel.LevelMembership == "" {
			companyModel.LevelMembership = "default"
		}

		_, errInsert, statuscode := services.InsertOneUsingURI(services.GetURI(credRoot), consts.DB_CORE_NAME, consts.Coll_Companies, utils.StructToJson(companyModel))
		if errInsert == "" {
			print("\nCompany registered\n")
			createCompanyDB(ctx, credRoot, companyModel.IdCompany, id, userOwner)
		} else {
			utils.ShowResponseDefault(ctx, statuscode, "error", errInsert)
		}
	}
}

// END_SAAS_FEATURE

// START_SAAS_FEATURE
func createCompanyDB(ctx *fasthttp.RequestCtx, creddbRoot models.CredDB, dbName string, user string, userOwner models.User) {
	_, err, statuscode := services.CreateDB(services.GetURI(creddbRoot), dbName)
	if err != "" {
		utils.ShowResponseDefault(ctx, statuscode, "error", err)
	} else {
		if statuscode == 200 {
			print("\nCompany DB is created\n")
			//Tambahkan user untuk akses DB dari company masing2
			var creddbCompany models.CredDB
			creddbCompany.User = user
			creddbCompany.Pass = config.EncryptAES(user)
			creddbCompany.DBName = dbName
			_, err := services.AddUserDB(services.GetURI(creddbRoot), creddbCompany)
			if err != nil {
				utils.ShowResponseDefault(ctx, statuscode, "error", err.Error())
				return
			} else {
				type RegisterCompanyResult struct {
					IdCompany string `json:"idcompany"`
					Token     string `json:"token"`
					Message   string `json:"message"`
				}
				print("\nUser for DB Company was created\n(mongodb://" + user + ":" + creddbCompany.Pass + "@" + os.Getenv("MONGODB_HOST_ONLINE") + ")  \n")
				expTime := time.Now().Local().Add(time.Hour*24*30).UnixNano() / 1000
				jwt := utils.GenerateJWT(models.Session{AppId: consts.APP_ID, IdCompany: dbName, IdUser: "admin_" + dbName, IdRole: "r_owner"}, expTime)
				utils.ShowResponseJson(ctx, statuscode, "success", RegisterCompanyResult{IdCompany: dbName, Token: jwt, Message: "Company was saved"})

				// Jalankan secara asinkron
				go func() {
					// Insert data user Administrator di "scm_core"->"users"
					// Insert data user Administrator di "c_xxxxxxx"->"users"
					// Memanggil fungsi UpdateDocument company, menambahkan cred & nonce
					var defaultRole models.Role
					utils.JsonToStruct(consts.Default_Role_Company_Json, &defaultRole)
					// userOwner.Id = "admin_" + dbName
					userOwner.AppId = consts.APP_ID
					// userOwner.Name = config.DecryptAES(userOwner.Name)
					// userOwner.Username = config.DecryptAES(userOwner.Username)
					// userOwner.Username = userOwner.Username
					userOwner.IdCompany = dbName
					userOwner.IdRole = defaultRole.Id
					if userOwner.Nickname == "" {
						userOwner.Nickname = "adm"
					}
					// userOwner.Nickname = config.DecryptAES(userOwner.Nickname)
					userOwner.Password = config.EncryptAES(config.EncryptAES(user))
					print("\n\nPassword:", userOwner.Password)
					_, errInsert, _ := services.InsertOneUsingURI(services.GetURI(creddbRoot), consts.DB_CORE_NAME, consts.Coll_Users, utils.StructToJson(userOwner))
					if errInsert == "" {
						updateField := `{"owner": ` + utils.StructToJson(userOwner) + `}`
						services.UpdateOneUsingURI(services.GetURI(creddbRoot), consts.DB_CORE_NAME, consts.Coll_Companies, `{"_id":"`+dbName+`"}`, updateField, false)
						AfterRegistration(creddbRoot, creddbCompany, dbName)
					}
					// go services.InsertOne(creddbCompany.DBName, consts.Coll_Users, utils.StructToJson(userOwner))
				}()
			}
		}
	}
}

// END_SAAS_FEATURE

// START_SAAS_FEATURE
func AfterRegistration(creddbRoot models.CredDB, creddbCompany models.CredDB, dbName string) {
	//Tambahkan role baru
	//Menambahkan creddb ke collection companies
	nonce, enc, err := config.EncryptChacha20poly1305(utils.StructToJson(creddbCompany))

	if err != "" {
		return
	}

	var defaultRole models.Role
	utils.JsonToStruct(consts.Default_Role_Company_Json, &defaultRole)
	creddbCompany.Nonce = nonce

	//Masukkan role baru ke dalam "roles" di coll "users" serta masukkan cred & nonce
	_, err1, _ := services.UpdateOneUsingURI(services.GetURI(creddbRoot), consts.DB_CORE_NAME, consts.Coll_Companies, `{"_id":"`+dbName+`"}`, `{"nonce":"`+nonce+`","cred":"`+enc+`","roles":[`+utils.StructToJson(defaultRole)+`]}`, false)

	//update "idrole" pada coll "users"
	_, err2, _ := services.UpdateOneUsingURI(services.GetURI(creddbRoot), consts.DB_CORE_NAME, consts.Coll_Users, `{"_id":"admin_c_`+creddbCompany.User+`"}`, `{"idrole":"`+defaultRole.Id+`"}`, false)
	print("\nerr1", err1, "\nerr2", err2)

	// InitiateDefaultCompanyData(creddbCompany)
	// InitiateDefaultCoreData(creddbRoot)
}

// END_SAAS_FEATURE

//	func InitiateDefaultCompanyData(creddbCompany models.CredDB) {
//		rolesDefaultJson := services.ReadFileToString("./data/roles.json")
//		services.InsertMany(creddbCompany.DBName, consts.Coll_Role, rolesDefaultJson)
//	}
func InitiateDefaultCoreData(creddbRoot models.CredDB) {
	menusDefaultJson := services.ReadFileToString("./data/menus.json")
	services.InsertManyUsingURI(services.GetURI(creddbRoot), consts.DB_CORE_NAME, consts.Coll_Menu, menusDefaultJson)
}

func GetCredDBCompany(idcompany string) (admReturn models.CredDB, err string) {
	val, errStr := services.GetValueRedis("cred_" + idcompany)
	if errStr != "" || val == "" {
		print(errStr)
		return admReturn, errStr
	} else {
		var credJson models.CredJson
		utils.JsonToStruct(val, &credJson)
		var cred models.CredDB
		res, _ := config.DecryptChacha20poly1305(credJson.Cred, credJson.Nonce)
		utils.JsonToStruct(res, &cred)
		return cred, ""
	}
}
