package consts

var ServerMongoOnline = true
var UsingRedisOnline = true

var APP_ID = "scm_app"
var ISSUER_ID = "Salomo07"

var DB_CORE_NAME = "scm_core"
var Coll_Initiate = "initiate"
var Coll_Users = "users"
var Coll_Companies = "companies"
var Coll_Role = "role"
var Coll_Menu = "menu"
var Coll_Handling_Units = "handling_units"
var Coll_Handling_Unit_Types = "handling_unit_types"
var Coll_Product_Menu = "product_menu"
var Coll_Products = "products"
var Coll_Category = "category"
var Coll_Branch = "branch"
var Coll_AccessMenu = "accessmenu"
var InitiateData = `[{"collection":"users","json":'{"name":"Company Administrator","code":"adm","desc":"Fullaccess for all company data"}'}]`
var Default_Role_Company_Json = `{"_id":"r_owner","name":"Company Administrator","code":"adm","desc":"Fullaccess for all company's data"}`
var MONGODB_CRED_ADMIN = ""

var CDB_HOST = "http://localhost"
var KEY_ADMIN = ""
var REDIS_CRED_ADMIN = ""
