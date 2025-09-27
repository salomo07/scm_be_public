package consts

var ErrNotFoundDoc = "No document found"
var ErrFindingDoc = "Error finding documents"
var ErrDecodingDoc = "Error decoding document"
var ErrCursor = "Cursor error"
var ErrMarshalling = "Error marshalling"
var EmptyBody = "Request body cant be empty"
var ErrDatabase = "Error database"

var UserAdd = "User added successfully"
var RoleAdd = "Role added successfully"
var MenuAdd = "Menu added successfully"
var AccessMenuAdd = "AccessMenu added successfully"
var AccessMenuExist = "AccessMenu already exist"
var AccessDeleteMandatory = "'Id' field is required and cannot be empty"

var RoleNotFound = "Role not found"
var FailAddUser = "Failed to add User"
var FailGetUser = "Failed to get User"
var FailAddRole = "Failed to add Role"
var FailAddMenu = "Failed to add Menu"
var FailAddAccessMenu = "Failed to add AccessMenu"
var FailGetSession = "Fail to get session (Redis)"

var CompanyAlreadyExist = `Company has been registered` // Alias terdaftar
var URLAlreadyExist = `URL has already been used`
var UserAlreadyExist = `Username, Mobile Phone or Email already registered`
var CodeCompanyExist = "Code Company already taken"

var PayloadEmpty = "Payload cant be empty"
var AdminKeyTidakDikenali = "AdminKey tidak dikenali"
var SessionNotFound = "Session is not found / expired, please re-login"
var UserNotFound = "User not found"
var PasswordIncorrect = "Password is incorrect"
var TokenExpired = "Token is expired"
var Unauthorized = "You have not access to this endpoint."
var HUTypeMustUnique = "Handling Unit Type name must be unique"
var HUTypeSuccessInsert = "Type of Handling Unit is saved"
var ProductMenuSuccessInsert = "Product is saved"
var ProductMenuMustUnique = "Product name must be unique"
