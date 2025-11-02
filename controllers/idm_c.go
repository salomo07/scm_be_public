package controllers

import (
	"encoding/json"
	"fmt"
	"os"
	"scm/config"
	"scm/consts"
	"scm/models"
	"scm/services"
	"scm/utils"
	"sort"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
)

type Claims struct {
	Data json.RawMessage `json:"data"`
	jwt.RegisteredClaims
}

// LoginRequest defines the structure of the login request body.
// @Description Login
// @Param body body models.LoginRequest true "Login request"
// @Success 200 {object} models.LoginResponse "Successful login"
// @Failure 400 {object} models.DefaultResponse "Bad Request"
// @Router /Auth/Login [post]
func Login(ctx *fasthttp.RequestCtx) {
	var loginReq models.LoginRequest
	utils.JsonToStruct(string(ctx.Request.Body()), &loginReq)
	usernameDecrypted := config.DecryptAES(loginReq.Username)
	if utils.ValidateRequiredFields(loginReq) != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, utils.ValidateRequiredFields(loginReq), "")
	} else {
		if loginReq.IdCompany == "" {
			services.TryLoginToDB(usernameDecrypted, ctx, loginReq)

		} else {
			utils.ShowResponseDefault(ctx, fasthttp.StatusOK, "error", "Ini mau ngapain, ada iddcompany, terus ngapain?")
		}
	}
}

func Logout(ctx *fasthttp.RequestCtx) {
	type LogoutRequest struct {
		Username string `json:"username" validate:"required"`
	}
	var logoutRequest LogoutRequest
	utils.JsonToStruct(string(ctx.Request.Body()), &logoutRequest)
	if utils.ValidateRequiredFields(logoutRequest) != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, utils.ValidateRequiredFields(logoutRequest), "")
		return
	}
	_, err := services.RemoveValueRedis(config.DecryptAES(logoutRequest.Username)) // Remove session on Redis
	if err == "" {
		utils.ShowResponseJson(ctx, fasthttp.StatusOK, "success", "Logout berhasil")
		return
	} else {
		utils.ShowResponseJson(ctx, fasthttp.StatusInternalServerError, "error", err)
		return
	}
}

func GetAccessMenuForLoginResponse(ctx *fasthttp.RequestCtx, idcompany string, idrole string) ([]models.AccessMenu, string) {
	// ambil credential DB perusahaan
	credDB, errGetCred := GetCredDBCompany(idcompany)
	if errGetCred != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusInternalServerError, "error", errGetCred)
		return nil, errGetCred
	}

	// query ke Mongo
	res, errStr, code := services.FindMany(
		credDB.DBName,
		consts.Coll_AccessMenu,
		fmt.Sprintf(`{"idrole":"%s"}`, idrole),
		`{"idmenu":1}`,
		0,
		0,
	)

	if errStr != "" {
		utils.ShowResponseDefault(ctx, code, "error", errStr)
		return nil, errStr
	}

	// konversi hasil query ke struct
	var accessMenu []models.AccessMenu
	for _, doc := range res {
		var m models.AccessMenu
		bsonBytes, _ := bson.Marshal(doc)
		_ = bson.Unmarshal(bsonBytes, &m)
		accessMenu = append(accessMenu, m)
	}

	// urutkan berdasarkan idmenu
	sort.Slice(accessMenu, func(i, j int) bool {
		return accessMenu[i].Idmenu < accessMenu[j].Idmenu
	})

	return accessMenu, ""
}

func CheckAdminKey(key string) string {
	val, err := services.GetValueRedis(key)
	if err != "" {
		print("Error : " + err + "\n\n")
		return ""
	} else {
		return val
	}
}

func extractToken(tokenString string) (data string, expiredTime int64, err error) {
	// Memparsing token tanpa memverifikasi tanda tangan
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", 0, fmt.Errorf("error parsing token: %v", err)
	}

	// Mengambil klaim dari token
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Mengambil data dari klaim "data"
		if dataStr, ok := claims["data"].(string); ok {
			data = dataStr
		} else {
			return "", 0, fmt.Errorf("token does not contain 'data' claim or it is not a string")
		}

		// Mengambil waktu kedaluwarsa dari klaim "exp"
		if exp, ok := claims["exp"].(float64); ok {
			expiredTime = int64(exp)
		} else {
			return "", 0, fmt.Errorf("token does not contain 'exp' claim or it is not a float64")
		}

		return data, expiredTime, nil
	}

	return "", 0, fmt.Errorf("invalid token claims")
}

func CheckSession(ctx *fasthttp.RequestCtx) (
	user models.User,
	errStr string,
	isRootAdmin bool,
	isAdminCompany bool,
) {
	var session models.Session
	authHeader := ctx.Request.Header.Peek("Authorization")
	tokenString, err := extractBearerToken(authHeader)
	if err != nil {
		return user, err.Error(), false, false
	}

	pubKeyData, err := os.ReadFile("public.key")
	if err != nil {
		return user, "Error reading public key", false, false
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
	if err != nil {
		return user, "Error parsing public key", false, false
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil || !token.Valid {
		return user, "Invalid token : " + err.Error(), false, false
	}

	// Validasi issuer
	if claims.Issuer != consts.ISSUER_ID {
		return user, "Invalid issuer", false, false
	}

	// Validasi audience -> jangan pakai .Contains(), lakukan pengecekan manual
	audOK := false
	for _, a := range claims.Audience {
		if a == consts.APP_ID {
			audOK = true
			break
		}
	}
	if !audOK {
		return user, "Invalid audience", false, false
	}
	session.IpClient = ctx.RemoteIP().String()
	session.UserAgent = string(ctx.UserAgent())
	session.IdRole = user.IdRole

	// Decode ke user
	_ = json.Unmarshal(claims.Data, &user)
	var loginResponseJWT models.LoginResponseJWT
	_ = json.Unmarshal(claims.Data, &loginResponseJWT)
	// Decode ke session (punya AdminKey)
	_ = json.Unmarshal(claims.Data, &session)

	// SuperAdmin check
	if session.AdminKey != "" {
		if config.DecryptAES(session.AdminKey) != os.Getenv("KEY_SUPERADMIN") {
			return user, consts.AdminKeyTidakDikenali, false, false
		}
		return user, "", true, false
	}

	// Role-based flag
	if loginResponseJWT.RoleType == "ADM" || strings.Contains(loginResponseJWT.RoleType, "ADM") {
		return user, "", false, true
	}
	return user, "", false, false
}

/*NOTE :
1) Perlu buat semacam function untuk TRANSACTION. Jadi ketika ada error dalam transaksi, lakukan ROLLBACK sehingga seluruh step dianggap tidak terjadi.
2) Perlu dibuat pengecekan Redis, jika session sudah ada beri flag alreadyLogin pada response yang menandakan Device lain sedang login dengan user yang sama, nanti di FE akan ada dialog konfirmasi apakah melanjutkan login (token yang lalu sudah tidak berlaku), atau batal (gagal login).
*/

func extractBearerToken(authHeader []byte) (string, error) {
	if !strings.HasPrefix(string(authHeader), "Bearer ") {
		return "", fmt.Errorf("Invalid Bearer token format")
	}

	token := strings.TrimPrefix(string(authHeader), "Bearer ")
	return token, nil
}
