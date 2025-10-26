package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"scm/config"
	"scm/consts"
	"scm/models"
	"scm/utils"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/ssh"
)

var (
	sshClient       atomic.Value
	poolLocks       sync.Map
	pool            sync.Map
	clientInstances sync.Map
)

const pemPath = "./private_key_mongo.pem"
const sshHost = "103.127.132.161:22"
const sshUser = "rejoice"

type Operation func(int, int) int

var MongoClient *mongo.Client

func GetURI(creddb models.CredDB) (uri string) {
	val := os.Getenv("USING_MONGO_TUNNEL")
	isUSingTunnel, _ := strconv.ParseBool(val)
	if isUSingTunnel {
		userpass := creddb.User + ":" + url.QueryEscape(creddb.Pass) + "@" + os.Getenv("MONGODB_HOST_ONLINE")
		uri = "mongodb://" + userpass
	} else {
		uri = "mongodb://" + os.Getenv("MONGODB_HOST_OFFLINE")
	}
	println("\n" + uri + "\n")
	return uri
}
func GetMongoClient() *mongo.Client {
	return MongoClient
}
func SetMongoClient(client *mongo.Client) {
	MongoClient = client
}

func connectDB(uri string) (*mongo.Client, string) {
	val := os.Getenv("USING_MONGO_TUNNEL")
	isUSingTunnel, err := strconv.ParseBool(val)
	if !isUSingTunnel {
		return connectDBWithAuth(uri)
	}

	if err := ensureSSHTunnel(); err != nil {
		return nil, fmt.Sprintf("Gagal memastikan SSH Tunnel: %v", err)
	}

	// Cek apakah sudah ada client
	if v, ok := clientInstances.Load(uri); ok {
		client := v.(*mongo.Client)
		if err := checkConnection(client); err == nil {
			return client, ""
		}
		log.Println("❗ Koneksi MongoDB tidak valid, membuat ulang...")
	}

	// Buat koneksi baru
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Sprintf("Gagal koneksi ke database: %v", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Sprintf("Gagal ping database: %v", err)
	}

	log.Println("✅ Koneksi ke MongoDB berhasil!")
	clientInstances.Store(uri, client)
	return client, ""
}
func connectDBWithAuth(uri string) (*mongo.Client, string) {
	// Cek apakah sudah ada koneksi yang aktif
	if v, ok := clientInstances.Load(uri); ok {
		client := v.(*mongo.Client)
		if err := checkConnection(client); err == nil {
			return client, ""
		}
		log.Println("❗ Koneksi MongoDB tidak valid, membuat ulang...")
	}

	// Buat koneksi baru ke MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Sprintf("Gagal koneksi ke database: %v", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Sprintf("Gagal ping database: %v", err)
	}

	log.Println("✅ Koneksi ke MongoDB dengan URI autentikasi berhasil!")
	clientInstances.Store(uri, client)
	return client, ""
}

func ensureSSHTunnel() error {
	// Cek apakah sudah ada client aktif
	if v := sshClient.Load(); v != nil {
		client := v.(*ssh.Client)
		_, _, err := client.SendRequest("keepalive@openssh.com", true, nil)
		if err == nil {
			return nil // masih aktif
		}
		log.Println("🔁 SSH Tunnel terputus, mencoba reconnect...")
		client.Close()
		sshClient.Store(nil)
	}

	// Reconnect
	client, err := startSSHTunnel(pemPath, sshUser, sshHost)
	if err != nil {
		return fmt.Errorf("gagal membuat SSH Tunnel: %w", err)
	}

	sshClient.Store(client)
	log.Println("✅ SSH Tunnel berhasil dibuat kembali.")
	return nil
}

func startSSHTunnel(pemPath, sshUser, sshHost string) (*ssh.Client, error) {
	pemBytes, err := os.ReadFile(pemPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PEM file: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	config := &ssh.ClientConfig{
		User:            sshUser,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := ssh.Dial("tcp", sshHost, config)
	if err != nil {
		return nil, fmt.Errorf("failed to establish SSH connection: %w", err)
	}

	// forwarding dalam goroutine
	go startTunnelForwarding(client)
	return client, nil
}

func startTunnelForwarding(client *ssh.Client) {
	localPort := os.Getenv("MONGO_LOCAL_PORT")
	if localPort == "" {
		localPort = "27018"
	}
	remotePort := os.Getenv("MONGO_REMOTE_PORT")
	if remotePort == "" {
		remotePort = "27017"
	}

	listener, err := net.Listen("tcp", "localhost:"+localPort)
	if err != nil {
		log.Printf("⚠️ Gagal listen port lokal %s: %v", localPort, err)
		return
	}
	log.Printf("🚀 SSH Tunnel forwarding aktif di localhost:%s → %s", localPort, remotePort)
	defer listener.Close()

	for {
		localConn, err := listener.Accept()
		if err != nil {
			log.Println("⚠️ SSH listener accept error:", err)
			continue
		}

		go func() {
			remoteConn, err := client.Dial("tcp", "127.0.0.1:"+remotePort)
			if err != nil {
				log.Println("❌ Remote dial SSH gagal:", err)
				localConn.Close()
				return
			}
			go io.Copy(remoteConn, localConn)
			go io.Copy(localConn, remoteConn)
		}()
	}
}

// checkConnection checks if the MongoDB client is still connected.
func checkConnection(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return client.Ping(ctx, readpref.Primary())
}

func FindMany(idCompany, collectionName, query, sort string, skip, limit int64) ([]bson.M, string, int) {
	db, err := GetMongoPool(idCompany, "")
	if err != nil {
		return nil, consts.ErrDatabase + " : " + err.Error(), fasthttp.StatusInternalServerError
	}
	collection := db.Collection(collectionName)

	// Parse filter
	var filter bson.M
	if query != "" {
		if err := bson.UnmarshalExtJSON([]byte(query), true, &filter); err != nil {
			return nil, consts.ErrMarshalling + " query: " + err.Error(), 400
		}
	} else {
		filter = bson.M{}
	}

	// Parse sort
	var sortOption bson.M
	if sort != "" {
		if err := bson.UnmarshalExtJSON([]byte(sort), true, &sortOption); err != nil {
			return nil, consts.ErrMarshalling + " sort: " + err.Error(), 400
		}
	}

	// Find options
	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(limit)
	}
	if skip > 0 {
		findOptions.SetSkip(skip)
	}
	if len(sortOption) > 0 {
		findOptions.SetSort(sortOption)
	}

	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []bson.M{}, "", 200
		}
		return nil, consts.ErrFindingDoc + " : " + err.Error(), 500
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	for cursor.Next(context.Background()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, consts.ErrDecodingDoc + " : " + err.Error(), 500
		}
		if oid, ok := result["_id"].(primitive.ObjectID); ok {
			result["_id"] = oid.Hex()
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, consts.ErrCursor + " : " + err.Error(), 500
	}

	return results, "", 200
}

func FindOne(idCompany, collectionName, query, projection, sort string) (string, string, int) {
	db, errConnect := GetMongoPool(idCompany, "")
	if errConnect != nil {
		return "", consts.ErrDatabase + " : " + errConnect.Error(), fasthttp.StatusInternalServerError
	}
	collection := db.Collection(collectionName)

	var filter bson.M
	if err := bson.UnmarshalExtJSON([]byte(query), true, &filter); err != nil {
		return "", consts.ErrMarshalling + " query", 400
	}

	var proj bson.M
	if projection != "" {
		if err := bson.UnmarshalExtJSON([]byte(projection), true, &proj); err != nil {
			return "", consts.ErrMarshalling + " projection", 400
		}
	}

	var sortOpt bson.M
	if sort != "" {
		if err := bson.UnmarshalExtJSON([]byte(sort), true, &sortOpt); err != nil {
			return "", consts.ErrMarshalling + " sort", 400
		}
	}

	opts := options.FindOne()
	if len(proj) > 0 {
		opts.SetProjection(proj)
	}
	if len(sortOpt) > 0 {
		opts.SetSort(sortOpt)
	}

	var result bson.M
	err := collection.FindOne(context.Background(), filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", "", 200
		}
		return "", consts.ErrFindingDoc, 500
	}

	if oid, ok := result["_id"].(primitive.ObjectID); ok {
		result["_id"] = oid.Hex()
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return "", consts.ErrMarshalling, 500
	}

	return string(resultBytes), "", 200
}

func FindManyRootDBUsingURI(uri string, dbname string, collectionName string, query, projection, sort string, limit, skip int64) ([]bson.M, string, int) {
	client, errConnect := connectDB(uri)
	if errConnect != "" {
		return nil, errConnect, 500
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(dbname).Collection(collectionName)

	// Parse filter
	var filter bson.M
	if query != "" {
		if err := bson.UnmarshalExtJSON([]byte(query), true, &filter); err != nil {
			return nil, consts.ErrMarshalling + " query", 400
		}
	} else {
		filter = bson.M{}
	}

	// Parse projection
	var proj bson.M
	if projection != "" {
		if err := bson.UnmarshalExtJSON([]byte(projection), true, &proj); err != nil {
			return nil, consts.ErrMarshalling + " projection", 400
		}
	}

	// Parse sort
	var sortOpt bson.M
	if sort != "" {
		if err := bson.UnmarshalExtJSON([]byte(sort), true, &sortOpt); err != nil {
			return nil, consts.ErrMarshalling + " sort", 400
		}
	}

	// Build find options
	opts := options.Find()
	if len(proj) > 0 {
		opts.SetProjection(proj)
	}
	if len(sortOpt) > 0 {
		opts.SetSort(sortOpt)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}
	if skip > 0 {
		opts.SetSkip(skip)
	}

	// Execute query
	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, consts.ErrFindingDoc, 500
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	for cursor.Next(context.Background()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, consts.ErrDecodingDoc, 500
		}

		if oid, ok := doc["_id"].(primitive.ObjectID); ok {
			doc["_id"] = oid.Hex()
		}

		results = append(results, doc)
	}

	if err := cursor.Err(); err != nil {
		return nil, consts.ErrCursor + " : " + err.Error(), 500
	}

	return results, "", 200
}

func FindOneRootDBUsingURI(uri string, dbname string, collectionName string, query, projection, sort string) (bson.M, string, int) {

	// Connect DB via URI
	client, errConnect := connectDB(uri)
	if errConnect != "" {
		return nil, errConnect, 500
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(dbname).Collection(collectionName)

	// Parse filter
	var filter bson.M
	if err := bson.UnmarshalExtJSON([]byte(query), true, &filter); err != nil {
		return nil, consts.ErrMarshalling + " query", 400
	}

	// Parse projection
	var proj bson.M
	if projection != "" {
		if err := bson.UnmarshalExtJSON([]byte(projection), true, &proj); err != nil {
			return nil, consts.ErrMarshalling + " projection", 400
		}
	}

	// Parse sort
	var sortOpt bson.M
	if sort != "" {
		if err := bson.UnmarshalExtJSON([]byte(sort), true, &sortOpt); err != nil {
			return nil, consts.ErrMarshalling + " sort", 400
		}
	}

	// Build options
	opts := options.FindOne()
	if len(proj) > 0 {
		opts.SetProjection(proj)
	}
	if len(sortOpt) > 0 {
		opts.SetSort(sortOpt)
	}

	// Execute query
	var result bson.M
	err := collection.FindOne(context.Background(), filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, "", 200
		}
		return nil, consts.ErrFindingDoc + " : " + err.Error(), 500
	}

	// Convert _id ke string
	if oid, ok := result["_id"].(primitive.ObjectID); ok {
		result["_id"] = oid.Hex()
	}

	return result, "", 200
}

func TryLoginToDB(usernameDecrypted string, ctx *fasthttp.RequestCtx, loginReq models.LoginRequest) {
	// pipeline := `[{"$match":{"$or":[{"username":"` + usernameDecrypted + `"},{"contact.mobile":"` + config.DecryptAES(loginReq.Mobile) + `"}]}},{"$lookup":{"from":"` + consts.Coll_Companies + `","localField":"idcompany","foreignField":"_id","as":"company"}},{"$unwind":{"path":"$company","preserveNullAndEmptyArrays":true}}]`
	pipeline := `[
		{
			"$match": {
				"$or": [
					{ "username": "` + usernameDecrypted + `" },
					{ "contact.mobile": "` + config.DecryptAES(loginReq.Mobile) + `" }
				]
			}
		},
		{
			"$lookup": {
				"from": "` + consts.Coll_Companies + `",
				"localField": "idcompany",
				"foreignField": "_id",
				"as": "company"
			}
		},
		{
			"$unwind": {
				"path": "$company",
				"preserveNullAndEmptyArrays": true
			}
		},
		{
			"$lookup": {
				"from": "` + consts.Coll_Role + `",
				"let": { "idroleStr": "$idrole" },
				"pipeline": [
					{
						"$match": {
							"$expr": { "$eq": ["$_id", { "$toObjectId": "$$idroleStr" }] }
						}
					}
				],
				"as": "role"
			}
		},
		{
			"$unwind": {
				"path": "$role",
				"preserveNullAndEmptyArrays": true
			}
		},
		{
			"$unset": ["company.roles"]
		}
	]`

	// pipeline := `[
	// 	{
	// 		"$match": {
	// 		"$or": [
	// 			{ "username": "` + usernameDecrypted + `" },
	// 			{ "contact.mobile": "` + config.DecryptAES(loginReq.Mobile) + `" }
	// 		]
	// 		}
	// 	},
	// 	{
	// 		"$lookup": {
	// 		"from": "` + consts.Coll_Companies + `",
	// 		"localField": "idcompany",
	// 		"foreignField": "_id",
	// 		"as": "company"
	// 		}
	// 	},
	// 	{
	// 		"$unwind": {
	// 		"path": "$company",
	// 		"preserveNullAndEmptyArrays": true
	// 		}
	// 	},
	// 	{
	// 		 "$lookup": {
	// 			"from":"` + consts.Coll_Role + `",
	// 			"let": { "idroleStr": "$idrole" },
	// 			"pipeline": [
	// 				{
	// 				"$match": {
	// 					"$expr": {
	// 					"$eq": [
	// 						"$_id",
	// 						{ "$toObjectId": "$$idroleStr" }
	// 					]
	// 					}
	// 				}
	// 				}
	// 			],
	// 			"as": "role"
	// 		}
	// 	},
	// 	{
	// 		"$unwind": {
	// 		"path": "$role",
	// 		"preserveNullAndEmptyArrays": true
	// 		}
	// 	}
	// ]`
	print(pipeline)
	res, err, code := AggregationOneUsingURI(GetURI(utils.GetMongoDBRoot()), consts.DB_CORE_NAME, consts.Coll_Users, pipeline)
	if err != "" {
		if err == consts.ErrNotFoundDoc {
			utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, "warning", consts.UserNotFound)
			return
		} else {
			utils.ShowResponseDefault(ctx, code, "error", err)
			return
		}
	} else {
		if res == "" {
			utils.ShowResponseDefault(ctx, fasthttp.StatusNotFound, "warning", consts.UserNotFound)
			return
		} else {
			var dataLogin models.LoginFromDB
			utils.JsonToStruct(res, &dataLogin)
			if config.DecryptAES(loginReq.Password) == config.DecryptAES(config.DecryptAES(dataLogin.Password)) {
				hours := 4
				if loginReq.Duration > 0 && loginReq.Duration < (24*7) {
					hours = loginReq.Duration
				} else if loginReq.RememberPassword {
					hours = 24 * 7
				}
				expTime := time.Now().Local().Add(time.Duration(hours) * time.Hour).Unix()
				expTime1Day := time.Now().Local().Add(time.Duration(24) * time.Hour).Unix()
				securedUserData := models.LoginResponseJWT{
					AppId:     os.Getenv("APP_ID"),
					Id:        config.EncodingBase64(dataLogin.Id),
					NIK:       config.EncodingBase64(dataLogin.NIK),
					Name:      config.EncodingBase64(dataLogin.Name),
					Username:  config.EncodingBase64(dataLogin.Username),
					IdCompany: config.EncodingBase64(dataLogin.IdCompany),
					IdBranch:  config.EncodingBase64(dataLogin.IdBranch),
					IdRole:    config.EncodingBase64(dataLogin.IdRole),
				}
				jwt := utils.GenerateJWT(securedUserData, expTime)
				jwt1Day := utils.GenerateJWT(securedUserData, expTime1Day)
				// go func() {
				// 	go services.SaveValueRedis("cred_"+userData.IdCompany, `{"cred":"`+UserCompany.Company.Cred+`","nonce":"`+UserCompany.Company.Nonce+`"}`, strconv.FormatInt(expTime, 10))
				// }()
				SaveValueRedis("cred_"+dataLogin.IdCompany, `{"cred":"`+dataLogin.Company.Cred+`","nonce":"`+dataLogin.Company.Nonce+`"}`, strconv.FormatInt(expTime1Day, 10))

				GetMongoPool(config.EncodingBase64(dataLogin.IdCompany), GetURI(models.CredDB{DBName: dataLogin.Company.IdCompany, User: dataLogin.IdCompany, Pass: config.EncryptAES(dataLogin.Company.IdCompany), Nonce: dataLogin.Company.Nonce}))
				// accessmenu,err := GetAccessMenuForLoginResponse(ctx, userData.IdCompany, userData.IdRole)

				// go services.SaveValueRedis(userData.Username, jwt, strconv.FormatInt(expTime, 10))

				go SaveValueRedis(dataLogin.Username+"_refreshtoken", jwt1Day, strconv.FormatInt(expTime1Day, 10))
				print(res)
				utils.ShowResponseJson(ctx, fasthttp.StatusOK, "success", models.LoginResponse{
					Username:     securedUserData.Username,
					IdUser:       config.EncryptAES(dataLogin.Id),
					IdCompany:    config.EncryptAES(dataLogin.IdCompany),
					Fullname:     securedUserData.Name,
					RoleName:     securedUserData.Role.Name,
					Access:       dataLogin.Role.AccessMenu,
					Token:        jwt,
					RefreshToken: jwt1Day,
					Expired:      time.Unix(expTime1Day, 0).String(),
				})
			} else {
				utils.ShowResponseDefault(ctx, fasthttp.StatusUnauthorized, consts.PasswordIncorrect, "")
				return
			}
		}
	}
}

func InsertOneUsingURI(uri string, dbName string, collectionName string, document string) (*mongo.InsertOneResult, string, int) {
	// Connect langsung via URI
	client, errConnect := connectDB(uri)
	if errConnect != "" {
		return nil, errConnect, 500
	}

	collection := client.Database(dbName).Collection(collectionName)

	// Parse dokumen JSON
	var doc bson.M
	if err := bson.UnmarshalExtJSON([]byte(document), true, &doc); err != nil {
		return nil, "Error unmarshalling document : " + err.Error(), 400
	}

	// Insert ke DB
	result, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		return nil, "Error inserting document : " + err.Error(), 500
	}

	return result, "", 200
}

func InsertOne(idCompany, collectionName, document string) (*mongo.InsertOneResult, string, int) {
	db, err := GetMongoPool(idCompany, "")
	if err != nil {
		return nil, consts.ErrDatabase + " : " + err.Error(), fasthttp.StatusInternalServerError
	}
	collection := db.Collection(collectionName)

	var doc bson.M
	if err := bson.UnmarshalExtJSON([]byte(document), true, &doc); err != nil {
		return nil, "Error unmarshalling document : " + err.Error(), 400
	}

	result, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		return nil, "Error inserting document : " + err.Error(), 500
	}

	if err != nil {
		return nil, "Error marshalling result : " + err.Error(), 500
	}

	return result, "", 200
}

func getCredFromDB(idCompanyDecoded string) string {
	res, errGetCred, _ := FindOneRootDBUsingURI(GetURI(utils.GetMongoDBRoot()), consts.DB_CORE_NAME, consts.Coll_Companies, consts.QueryFindCompany(idCompanyDecoded), `{}`, `{}`)
	if errGetCred != "" {
		return errGetCred
	}
	var company models.Company
	jsonBytes, _ := json.Marshal(res)
	utils.JsonToStruct(string(jsonBytes), &company)
	GetMongoPool(config.EncodingBase64(idCompanyDecoded), GetURI(models.CredDB{DBName: idCompanyDecoded, User: idCompanyDecoded, Pass: config.EncryptAES(idCompanyDecoded), Nonce: company.Nonce}))
	return ""
}

func GetMongoPool(idCompany, uri string) (*mongo.Database, error) {
	decoded := config.DecodingBase64(idCompany)

	if v, ok := pool.Load(decoded); ok {
		client := v.(*mongo.Client)
		if err := client.Ping(context.Background(), nil); err == nil {
			return client.Database(decoded), nil
		}
		pool.Delete(decoded)
	}

	if uri == "" {
		errReconnect := getCredFromDB(decoded)
		if errReconnect != "" {
			return nil, fmt.Errorf("Failed to access data (%s). Please re-login", decoded)
		}
	}

	// ambil lock khusus untuk company ini
	muIface, _ := poolLocks.LoadOrStore(decoded, &sync.Mutex{})
	mu := muIface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	// cek ulang setelah dapat lock (double-check)
	if v, ok := pool.Load(decoded); ok {
		return v.(*mongo.Client).Database(decoded), nil
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect mongo: %w", err)
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	pool.Store(decoded, client)
	return client.Database(decoded), nil
}

func UpdateOne(idCompany string, collectionName string, query string, update string, upsert bool) (*mongo.UpdateResult, string, int) {
	db, err := GetMongoPool(idCompany, "")
	if err != nil {
		return nil, consts.ErrDatabase + " : " + err.Error(), fasthttp.StatusInternalServerError
	}
	collection := db.Collection(collectionName)

	// Unmarshal query
	var filter bson.M
	if err := bson.UnmarshalExtJSON([]byte(query), true, &filter); err != nil {
		return nil, "Error unmarshalling query : " + err.Error(), 400
	}

	// Unmarshal update
	if !json.Valid([]byte(update)) {
		return nil, "Invalid JSON format for update", 400
	}
	var updateFields map[string]interface{}
	if err := json.Unmarshal([]byte(update), &updateFields); err != nil {
		return nil, "Error unmarshalling update: " + err.Error(), 400
	}

	// Optional decrypt for user collection
	if collectionName == consts.Coll_Users {
		updateFields = decryptUserColl(updateFields)
	}

	updateDoc := bson.M{"$set": updateFields}
	opts := options.Update().SetUpsert(upsert)

	result, err := collection.UpdateOne(context.Background(), filter, updateDoc, opts)
	if err != nil {
		return nil, "Error updating document : " + err.Error(), 500
	}

	return result, "", 200
}

func UpdateOneUsingURI(uri string, dbname string, collectionName string, query string, update string, upsert bool) (*mongo.UpdateResult, string, int) {
	// Connect DB langsung via URI
	client, errConnect := connectDB(uri)
	if errConnect != "" {
		return nil, errConnect, 500
	}
	collection := client.Database(dbname).Collection(collectionName)

	// Unmarshal query
	var filter bson.M
	if err := bson.UnmarshalExtJSON([]byte(query), true, &filter); err != nil {
		return nil, "Error unmarshalling query : " + err.Error(), 400
	}

	// Unmarshal update
	if !json.Valid([]byte(update)) {
		return nil, "Invalid JSON format for update", 400
	}
	var updateFields map[string]interface{}
	if err := json.Unmarshal([]byte(update), &updateFields); err != nil {
		return nil, "Error unmarshalling update: " + err.Error(), 400
	}

	// Optional decrypt untuk collection user
	// if collectionName == consts.Coll_Users {
	// 	updateFields = decryptUserColl(updateFields)
	// }

	updateDoc := bson.M{"$set": updateFields}
	opts := options.Update().SetUpsert(upsert)

	// Execute update
	result, err := collection.UpdateOne(context.Background(), filter, updateDoc, opts)
	if err != nil {
		return nil, "Error updating document : " + err.Error(), 500
	}

	return result, "", 200
}

func decryptUserColl(updateFields map[string]interface{}) map[string]interface{} {
	if name, ok := updateFields["name"].(string); ok && name != "" {
		updateFields["name"] = config.DecryptAES(name)
	}
	if nickname, ok := updateFields["nickname"].(string); ok && nickname != "" {
		updateFields["nickname"] = config.DecryptAES(nickname)
	}
	if username, ok := updateFields["username"].(string); ok && username != "" {
		updateFields["username"] = config.DecryptAES(username)
	}
	if contact, ok := updateFields["contact"].(map[string]interface{}); ok {
		if email, ok := contact["email"].(string); ok && email != "" {
			contact["email"] = config.DecryptAES(email)
		}
		if mobile, ok := contact["mobile"].(string); ok && mobile != "" {
			contact["mobile"] = config.DecryptAES(mobile)
		}
		if whatsapp, ok := contact["whatsapp"].(string); ok && whatsapp != "" {
			contact["whatsapp"] = config.DecryptAES(whatsapp)
		}
		if phone, ok := contact["phone"].(string); ok && phone != "" {
			contact["phone"] = config.DecryptAES(phone)
		}
		updateFields["contact"] = contact
	}
	return updateFields
}

func UpdateMany(idCompany, collectionName, query, update string) (*mongo.UpdateResult, string, int) {
	// Ambil koneksi dari pool
	db, err := GetMongoPool(idCompany, "")
	if err != nil {
		return nil, consts.ErrDatabase + " : " + err.Error(), fasthttp.StatusInternalServerError
	}
	collection := db.Collection(collectionName)

	// Unmarshal query filter
	var filter bson.M
	if err := bson.UnmarshalExtJSON([]byte(query), true, &filter); err != nil {
		return nil, "Error unmarshalling query: " + err.Error(), 400
	}

	// Unmarshal update document
	var updateDoc bson.M
	if err := bson.UnmarshalExtJSON([]byte(update), true, &updateDoc); err != nil {
		return nil, "Error unmarshalling update document: " + err.Error(), 400
	}

	// Eksekusi UpdateMany
	result, err := collection.UpdateMany(context.Background(), filter, updateDoc)
	if err != nil {
		return nil, "Error updating documents: " + err.Error(), 500
	}

	if err != nil {
		return nil, "Error marshalling result: " + err.Error(), 500
	}

	return result, "", 200
}

func DeleteOne(idCompany, collectionName, query string) (*mongo.DeleteResult, string, int) {
	fmt.Print(query)
	db, err := GetMongoPool(idCompany, "")
	if err != nil {
		return nil, consts.ErrDatabase + " : " + err.Error(), fasthttp.StatusInternalServerError
	}
	collection := db.Collection(collectionName)

	// Parse query filter
	var filter bson.M
	if err := bson.UnmarshalExtJSON([]byte(query), true, &filter); err != nil {
		return nil, "Error unmarshalling query : " + err.Error(), 400
	}

	// Eksekusi delete
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, "Error deleting document : " + err.Error(), 500
	}

	return result, "", 200
}

func DeleteMany(idCompany, collectionName, query string) (*mongo.DeleteResult, string, int) {
	db, err := GetMongoPool(idCompany, "")
	if err != nil {
		return nil, consts.ErrDatabase + " : " + err.Error(), fasthttp.StatusInternalServerError
	}
	collection := db.Collection(collectionName)

	var filter bson.M
	if err := bson.UnmarshalExtJSON([]byte(query), true, &filter); err != nil {
		return nil, "Error unmarshalling query : " + err.Error(), 400
	}

	result, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		return nil, "Error deleting documents : " + err.Error(), 500
	}

	return result, "", 200
}

func InsertMany(idCompany, collectionName, jsonString string) (*mongo.InsertManyResult, string, int) {
	db, err := GetMongoPool(idCompany, "")
	if err != nil {
		return nil, consts.ErrDatabase + " : " + err.Error(), fasthttp.StatusInternalServerError
	}
	collection := db.Collection(collectionName)

	var documents []bson.M
	if err := json.Unmarshal([]byte(jsonString), &documents); err != nil {
		return nil, "Error unmarshalling JSON: " + err.Error(), 400
	}

	docs := make([]interface{}, len(documents))
	for i, doc := range documents {
		docs[i] = doc
	}

	result, err := collection.InsertMany(context.Background(), docs)
	if err != nil {
		return nil, "Error inserting documents: " + err.Error(), 500
	}

	return result, "", 200
}

func InsertManyUsingURI(uri string, dbname string, collectionName string, jsonString string) (*mongo.InsertManyResult, string, int) {
	// Connect DB langsung via URI
	client, errConnect := connectDB(uri)
	if errConnect != "" {
		return nil, errConnect, 500
	}
	collection := client.Database(dbname).Collection(collectionName)

	// Parse JSON menjadi slice of bson.M
	var documents []bson.M
	if err := json.Unmarshal([]byte(jsonString), &documents); err != nil {
		return nil, "Error unmarshalling JSON: " + err.Error(), 400
	}

	// Konversi ke slice of interface{}
	docs := make([]interface{}, len(documents))
	for i, doc := range documents {
		docs[i] = doc
	}

	// InsertMany
	result, err := collection.InsertMany(context.Background(), docs)
	if err != nil {
		return nil, "Error inserting documents: " + err.Error(), 500
	}

	return result, "", 200
}

func AggregationOneUsingURI(uri, dbName, collectionName, jsonPipeline string) (string, string, int) {
	// 1. Connect pakai URI
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return "", "Error creating Mongo client: " + err.Error(), 500
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return "", "Error connecting to Mongo: " + err.Error(), 500
	}
	defer client.Disconnect(ctx)

	// 2. Ambil collection
	collection := client.Database(dbName).Collection(collectionName)

	// 3. Parse pipeline dari JSON
	var pipeline []interface{}
	if err := json.Unmarshal([]byte(jsonPipeline), &pipeline); err != nil {
		return "", "Error parsing pipeline: " + err.Error(), 400
	}

	// 4. Jalankan aggregation
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return "", "Error aggregation: " + err.Error(), 500
	}
	defer cursor.Close(ctx)

	// 5. Ambil hasil pertama
	var result bson.M
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return "", "Error decoding result: " + err.Error(), 500
		}
	} else {
		return "", "Document not found", 404
	}

	// 6. Convert ke JSON string
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", "Error marshal JSON: " + err.Error(), 500
	}

	return string(resultJSON), "", 200
}

func AggregationOne(idCompany, collectionName, jsonPipeline string) (string, string, int) {
	db, err := GetMongoPool(idCompany, "")
	if err != nil {
		return "", consts.ErrDatabase + " : " + err.Error(), fasthttp.StatusInternalServerError
	}
	collection := db.Collection(collectionName)

	var pipeline []interface{}
	if err := utils.JsonToStruct(jsonPipeline, &pipeline); err != nil {
		return "", "Error parsing pipeline: " + err.Error(), 400
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return "", err.Error(), 500
	}
	defer cursor.Close(context.Background())

	var result bson.M
	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&result); err != nil {
			return "", err.Error(), 500
		}
	} else {
		return "", consts.ErrNotFoundDoc, 404
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", err.Error(), 500
	}

	return string(resultJSON), "", 200
}

func AggregationMany(idCompany, local, foreign, localField, foreignField, as, query string) (string, string, int) {
	db, err := GetMongoPool(idCompany, "")
	if err != nil {
		return "", consts.ErrDatabase + " : " + err.Error(), fasthttp.StatusInternalServerError
	}
	localsCollection := db.Collection(local)

	var matchStage bson.D
	if err := bson.UnmarshalExtJSON([]byte(`{"$match":`+query+`}`), true, &matchStage); err != nil {
		return "", "Error building match stage: " + err.Error(), 400
	}

	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: foreign},
			{Key: "localField", Value: localField},
			{Key: "foreignField", Value: foreignField},
			{Key: "as", Value: as},
		}},
	}

	pipeline := mongo.Pipeline{matchStage, lookupStage}

	cursor, err := localsCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return "", err.Error(), 500
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		return "", err.Error(), 500
	}

	resultJSON, err := json.Marshal(results)
	if err != nil {
		return "", err.Error(), 500
	}

	return string(resultJSON), "", 200
}

// START_SAAS_FEATURE
func CreateDB(uri, dbname string) (resBody string, errStr string, statuscode int) {
	client, errConnect := connectDB(uri)
	if errConnect != "" {
		log.Printf("Error connecting to database: %v\n", errConnect)
		return "", errConnect, 500
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Printf("%v\n", err)
		}
	}()

	collection := client.Database(dbname).Collection(consts.Coll_Role)
	document := consts.Default_Role_Company_Json
	var doc bson.M
	err := bson.UnmarshalExtJSON([]byte(document), true, &doc)
	if err != nil {
		log.Printf("Error unmarshalling document: %v\n", err)
		return "", "Error unmarshalling document: " + err.Error(), 400
	}

	result, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		log.Printf("Error inserting document: %v\n", err)
		return "", "Error inserting document: " + err.Error(), 500
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error marshalling result: %v\n", err)
		return "", "Internal Server Error", 500
	}

	return string(resultBytes), "", 200
}

// END_SAAS_FEATURE

// START_SAAS_FEATURE
func AddUserDB(uri string, creddbCompany models.CredDB) (models.CredDB, error) {
	client, errConnect := connectDB(uri)
	if errConnect != "" {
		return models.CredDB{}, fmt.Errorf("Failed to connect to MongoDB : " + errConnect)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			fmt.Printf("Failed to disconnect from MongoDB: %v\n", err)
		}
	}()

	// Access admin database to create user
	adminDB := client.Database("admin")

	// Command to create user
	cmd := bson.D{
		{Key: "createUser", Value: creddbCompany.User},
		{Key: "pwd", Value: creddbCompany.Pass},
		{Key: "roles", Value: bson.A{
			bson.D{
				{Key: "role", Value: "readWrite"},
				{Key: "db", Value: creddbCompany.DBName},
			},
		}},
	}

	// Run the createUser command
	var result bson.M
	err := adminDB.RunCommand(context.Background(), cmd).Decode(&result)
	if err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
		return models.CredDB{}, err
	}
	return creddbCompany, nil
}

// END_SAAS_FEATURE
func insertManyDocuments(client *mongo.Client, dbName string, collectionName string, documents []interface{}) {
	collection := client.Database(dbName).Collection(collectionName)

	result, err := collection.InsertMany(context.Background(), documents)
	if err != nil {
		log.Fatal(err)
	}

	print("Inserted documents: %v\n", result.InsertedIDs)
}
