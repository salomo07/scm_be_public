package master_controller

import (
	"fmt"
	"reflect"
	"scm/consts"
	"scm/controllers"
	"scm/models"
	"scm/services"
	"scm/utils"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var EntityRegistry = map[string]struct {
	New        func() interface{}
	Collection string
}{
	"product": {
		New:        func() interface{} { return &models.Product{} },
		Collection: consts.Coll_Products,
	},
	"category": {
		New:        func() interface{} { return &models.Category{} },
		Collection: consts.Coll_Category,
	},
}

func FindProducts(ctx *fasthttp.RequestCtx, user models.User) { //querynya menggunakan "_id"
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}
	queryFindProduct := string(ctx.Request.Body())

	sort, skip, limit := controllers.GetSortSkipLimit(ctx)
	resFindMany, err, code := services.FindMany(user.IdCompany, consts.Coll_Products, queryFindProduct, sort, skip, limit)
	if err != "" {
		utils.ShowResponseDefault(ctx, code, "error", err)
	} else {
		if len(resFindMany) == 0 {
			utils.ShowResponseJson(ctx, code, "success", []interface{}{})
			return
		}
		utils.ShowResponseJson(ctx, code, "success", resFindMany)
	}
}

func FindCategory(ctx *fasthttp.RequestCtx, user models.User) { //querynya menggunakan "_id"
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}
	queryFindCategory := string(ctx.Request.Body())

	sort, skip, limit := controllers.GetSortSkipLimit(ctx)
	resFindMany, err, code := services.FindMany(user.IdCompany, consts.Coll_Category, queryFindCategory, sort, skip, limit)
	if err != "" {
		utils.ShowResponseDefault(ctx, code, "error", err)
	} else {
		if len(resFindMany) == 0 {
			utils.ShowResponseJson(ctx, code, "success", []interface{}{})
			return
		}
		utils.ShowResponseJson(ctx, code, "success", resFindMany)
	}
}
func UpsertCategoryProduct(ctx *fasthttp.RequestCtx, user models.User) {
	UpsertEntity(ctx, user, "category")
}
func UpsertProduct(ctx *fasthttp.RequestCtx, user models.User) {
	UpsertEntity(ctx, user, "product")
}

// func UpsertCategoryProduct(ctx *fasthttp.RequestCtx, user models.User) {
// 	if len(ctx.Request.Body()) == 0 {
// 		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
// 		return
// 	}

// 	var dataCategory models.CategoryRequest
// 	if err := utils.JsonToStruct(string(ctx.Request.Body()), &dataCategory); err != nil {
// 		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", "Invalid JSON body")
// 		return
// 	}

// 	// Validasi field wajib
// 	if msg := utils.ValidateRequiredFields(dataCategory); msg != "" {
// 		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, msg, "")
// 		return
// 	}

// 	// Buat filter
// 	var filter string
// 	if dataCategory.Id != "" {
// 		// User kirim _id string → convert ke filter MongoDB
// 		if !primitive.IsValidObjectID(dataCategory.Id) {
// 			utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", "Invalid _id format")
// 			return
// 		}
// 		// manual bikin JSON supaya tetap ObjectID
// 		filter = fmt.Sprintf(`{"_id":{"$oid":"%s"}}`, dataCategory.Id)
// 	} else {
// 		// Tidak ada _id → biarkan kosong supaya jadi insert
// 		filter = `{"_id":{"$exists":false}}`
// 	}

// 	// Data untuk update/insert (hapus field _id agar tidak overwrite)
// 	removedId := utils.RemoveField(dataCategory, "_id")
// 	data := utils.StructToJson(removedId)

// 	// Debug log
// 	fmt.Println("FILTER JSON:", filter)
// 	fmt.Println("DATA JSON:", data)

// 	// Jalankan upsert
// 	resInsert, err, code := services.UpdateOne(
// 		user.IdCompany,
// 		consts.Coll_Category,
// 		filter,
// 		data,
// 		true, // upsert
// 	)

//		if resInsert == nil {
//			utils.ShowResponseDefault(ctx, code, "error", err)
//		} else {
//			// Response jangan diubah
//			utils.ShowResponseJson(ctx, code, "success", resInsert)
//		}
//	}
// func UpsertProduct(ctx *fasthttp.RequestCtx, user models.User) {
// 	if len(ctx.Request.Body()) == 0 {
// 		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
// 		return
// 	}

// 	var dataProduct models.Product
// 	if err := utils.JsonToStruct(string(ctx.Request.Body()), &dataProduct); err != nil {
// 		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", "Invalid JSON body")
// 		return
// 	}

// 	// Validasi field wajib
// 	if msg := utils.ValidateRequiredFields(dataProduct); msg != "" {
// 		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, msg, "")
// 		return
// 	}

// 	// Buat filter
// 	var filter string
// 	if dataProduct.Id != "" {
// 		// User kirim _id string → convert ke filter MongoDB
// 		if !primitive.IsValidObjectID(dataProduct.Id) {
// 			utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", "Invalid _id format")
// 			return
// 		}
// 		// manual bikin JSON supaya tetap ObjectID
// 		filter = fmt.Sprintf(`{"_id":{"$oid":"%s"}}`, dataProduct.Id)
// 	} else {
// 		// Tidak ada _id → biarkan kosong supaya jadi insert
// 		filter = `{"_id":{"$exists":false}}`
// 	}

// 	// Data untuk update/insert (hapus field _id agar tidak overwrite)
// 	removedId := utils.RemoveField(dataProduct, "_id")
// 	data := utils.StructToJson(removedId)

// 	// Debug log
// 	fmt.Println("FILTER JSON:", filter)
// 	fmt.Println("DATA JSON:", data)

// 	// Jalankan upsert
// 	resInsert, err, code := services.UpdateOne(
// 		user.IdCompany,
// 		consts.Coll_Products,
// 		filter,
// 		data,
// 		true, // upsert
// 	)

// 	if resInsert == nil {
// 		utils.ShowResponseDefault(ctx, code, "error", err)
// 	} else {
// 		// Response jangan diubah
// 		utils.ShowResponseJson(ctx, code, "success", resInsert)
// 	}
// }

func UpsertEntity(ctx *fasthttp.RequestCtx, user models.User, entityName string) {
	if len(ctx.Request.Body()) == 0 {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}

	// Cari entity di registry
	reg, ok := EntityRegistry[entityName]
	if !ok {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", "Unknown entity: "+entityName)
		return
	}

	// Buat struct target kosong
	entity := reg.New() // biasanya pointer (*Category, *Product, dll)

	// Decode JSON ke struct target
	if err := utils.JsonToStruct(string(ctx.Request.Body()), entity); err != nil {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", "Invalid JSON body")
		return
	}

	// Dereference entity sebelum divalidasi
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if msg := utils.ValidateRequiredFields(val.Interface()); msg != "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, msg, "")
		return
	}

	// Ambil ID (pakai reflection supaya generic)
	id := utils.GetIdField(entity) // bikin helper sendiri
	var filter string
	if id != "" {
		if !primitive.IsValidObjectID(id) {
			utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "error", "Invalid _id format")
			return
		}
		filter = fmt.Sprintf(`{"_id":{"$oid":"%s"}}`, id)
	} else {
		filter = `{"_id":{"$exists":false}}`
	}

	// Hapus _id sebelum insert/update
	removedId := utils.RemoveField(entity, "_id")
	data := utils.StructToJson(removedId)

	fmt.Println("FILTER JSON:", filter)
	fmt.Println("DATA JSON:", data)

	resInsert, err, code := services.UpdateOne(
		user.IdCompany,
		reg.Collection, // ← pakai collection dari registry
		filter,
		data,
		true,
	)

	if resInsert == nil {
		utils.ShowResponseDefault(ctx, code, "error", err)
	} else {
		utils.ShowResponseJson(ctx, code, "success", resInsert)
	}
}

func DeleteProduct(ctx *fasthttp.RequestCtx, user models.User) { //querynya menggunakan "_id"
	if string(ctx.Request.Body()) == "" {
		utils.ShowResponseDefault(ctx, fasthttp.StatusBadRequest, "warning", consts.EmptyBody)
		return
	}

	queryProduct := string(ctx.Request.Body())
	resUpsert, err, code := services.DeleteOne(user.IdCompany, consts.Coll_Products, queryProduct)
	if err != "" {
		utils.ShowResponseDefault(ctx, code, "error", err)
	} else {
		utils.ShowResponseJson(ctx, code, "success", resUpsert)

	}
}
