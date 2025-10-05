package consts

var URL_Req_OTP = "/api/v1/auth/reqotp"
var URL_Validate_OTP = "/api/v1/auth/validotp"

var URL_Auth_Login = "/api/v1/auth/login"
var URL_Auth_Logout = "/api/v1/auth/logout"
var URL_Auth_Enc = "/api/v1/auth/encryptaes"
var URL_Auth_Dec = "/api/v1/auth/decryptaes"

// SuperAdmin - Company URL
var company = "/api/v1/admin/company/"
var URL_Company_Create = company + "create/"
var URL_Company_Role_Create = company + "role/addrole"
var URL_Company_Initiate = company + "copyinitiatedata"

var URL_Branch_Create = company + "branch/create"
var URL_Branch_Find = company + "branch/find"
var URL_Branch_Upsert = company + "branch/upsert"
var URL_Branch_Delete = company + "branch/delete"

var master = "/api/v1/master/"
var URL_Role_Find = master + "roles/find"
var URL_Role_Upsert = master + "roles/upsert"
var URL_Role_Create = master + "roles/create"
var URL_Role_Delete = master + "roles/delete"

// SuperAdmin - Menu URL
var URL_Menu_Show_All = "/api/v1/admin/menu/showall"
var URL_Menu_Create = "/api/v1/admin/menu/add"
var URL_Menu_Create_Bulk = "/api/v1/admin/menu/addBulk"

// SuperAdmin, CompanyAdmin - Role URL
// var URL_Role_Create = "/api/v1/admin/role/create"
// var URL_Role_CreateBulk = "/api/v1/admin/role/create" // Soon

// SuperAdmin, CompanyAdmin - Access URL
var URL_AccessMenu_Find = "/api/v1/admin/access/find"
var URL_AccessMenu_FindMany = "/api/v1/admin/access/findMany"
var URL_AccessMenu_Create = "/api/v1/admin/access/create"
var URL_AccessMenu_CreateMany = "/api/v1/admin/access/createMany"
var URL_AccessMenu_Update = "/api/v1/admin/access/update"
var URL_AccessMenu_Delete = "/api/v1/admin/access/delete"

// SuperAdmin, CompanyAdmin - User URL
var URL_User_Find = "/api/v1/admin/user/find"
var URL_User_FindMany = "/api/v1/admin/user/findMany"
var URL_User_Update = "/api/v1/admin/user/update"
var URL_User_Create = "/api/v1/admin/user/create"

// SuperAdmin - Hanya boleh dipakai URGENT (backdoor)
var URL_Find_DB = "/api/v1/db/find"
var URL_Insert_DB = "/api/v1/db/insert"
var URL_Update_DB = "/api/v1/db/update"

var URL_Master = "/api/v1/master/"
var URL_Master_AddHUType = URL_Master + "addHUType"
var URL_Master_AddProduct_Menu = URL_Master + "insertproductmenu"
var URL_Master_AddHU = URL_Master + "addHU"
var URL_Master_GetHUType = URL_Master + "getHUType"
var URL_Master_GetHU = URL_Master + "getHU"

var URL_Product_Find = URL_Master + "product/find"
var URL_Product_Upsert = URL_Master + "product/upsert"
var URL_Product_Delete = URL_Master + "product/delete"
var URL_Category_Product_Upsert = URL_Master + "product/cat_upsert"
var URL_Category_Product_Find = URL_Master + "product/cat_find"
