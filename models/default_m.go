package models

type PaginationRequest struct {
	Skip  int64 `json:"skip" validate:"min=0"`
	Limit int64 `json:"limit" validate:"required"`
	IsASC bool  `json:"isasc"`
}

type Session struct {
	//contoh format IdAppCompanyUser scm*c_1324353452345*u_34345345
	AdminKey  string `json:"adminkey"` //Hanya untuk SuperAdmin
	AppId     string `json:"appid" validate:"required"`
	IdCompany string `json:"idcompany" validate:"required"`
	IdUser    string `json:"iduser" validate:"required"`
	IdRole    string `json:"idrole" validate:"required"`
	UserAgent string `json:"useragent"`
	IpClient  string `json:"ipclient"`
}

type RequestInsertBulk []RequestInsertBulkModel

type RequestInsertBulkModel struct {
	Ok  bool   `json:"ok"`
	ID  string `json:"id"`
	Rev string `json:"rev"`
}

type CredDB struct {
	DBName string `json:"dbname"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Nonce  string `json:"nonce"`
}
type CredJson struct {
	Cred  string `json:"cred"`
	Nonce string `json:"nonce"`
}

type UpdateDocumentBySuperAdmin struct {
	DBName      string `json:"dbname" validate:"required"`
	CollName    string `json:"collname" validate:"required"`
	Query       string `json:"query" validate:"required"`
	UpdateField string `json:"updatefield" validate:"required"`
}
type InsertDocumentBySuperAdmin struct {
	DBName   string `json:"dbname" validate:"required"`
	CollName string `json:"collname" validate:"required"`
	Document string `json:"document" validate:"required"`
}
type FindDocumentBySuperAdmin struct {
	DBName     string `json:"dbname" validate:"required"`
	CollName   string `json:"collname" validate:"required"`
	Query      string `json:"query" validate:"required"`
	Projection string `json:"projection" validate:"required"`
}

// func ValidateStruct(myStruct any, ctx *fasthttp.RequestCtx) (err string) {
// 	validate := validator.New()
// 	errMsg := validate.Struct(myStruct)
// 	if errMsg != nil {
// 		// Handle kesalahan validasi
// 		validationErrors := errMsg.(validator.ValidationErrors)
// 		for i, e := range validationErrors {
// 			err = err + e.Namespace() + " is " + e.Tag()
// 			if i < len(validationErrors)-1 {
// 				err = err + ", "
// 			}
// 			if i == len(validationErrors)-1 {
// 				err = err + "."
// 			}
// 		}
// 		ShowResponseDefault(ctx, fasthttp.StatusBadRequest, err, "")
// 		return err
// 	}
// 	return err
// }

// func validateStruct(data interface{}, index int) string {
// 	var missingFields []string
// 	err := validate.Struct(data)

// 	if err != nil {
// 		if _, ok := err.(*validator.InvalidValidationError); ok {
// 			return "Invalid input"
// 		}

// 		for _, err := range err.(validator.ValidationErrors) {
// 			missingFields = append(missingFields, fmt.Sprintf("'%s'", err.StructField()))
// 		}
// 		missingFieldsStr := strings.Join(missingFields, ", ")

// 		if len(missingFields) > 1 {
// 			if index >= 0 {
// 				return fmt.Sprintf("Index %d: %s fields are required and cannot be empty", index, missingFieldsStr)
// 			}
// 			return fmt.Sprintf("%s fields are required and cannot be empty", missingFieldsStr)
// 		} else {
// 			if index >= 0 {
// 				return fmt.Sprintf("Index %d: %s field is required and cannot be empty", index, missingFieldsStr)
// 			}
// 			return fmt.Sprintf("%s field is required and cannot be empty", missingFieldsStr)
// 		}
// 	}

// 	return ""
// }
