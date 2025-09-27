package controllers

import (
	"strconv"

	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"
)

func GetSortSkipLimit(ctx *fasthttp.RequestCtx) (string, int64, int64) {
	// default sort: _id desc
	sortParam := string(ctx.QueryArgs().Peek("sort"))
	var sort string
	if sortParam == "asc" {
		sort = `{ "_id": 1 }`
	} else if sortParam == "desc" {
		sort = `{ "_id": -1 }`
	} else {
		sort = `{ "_id": -1 }` // default
	}

	// skip
	skipStr := string(ctx.QueryArgs().Peek("skip"))
	var skip int64
	if skipStr != "" {
		if v, err := strconv.ParseInt(skipStr, 10, 64); err == nil {
			skip = v
		}
	}

	// limit
	limitStr := string(ctx.QueryArgs().Peek("limit"))
	var limit int64
	if limitStr != "" {
		if v, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			limit = v
		}
	}

	return sort, skip, limit
}
func HashPassword(password string) (string, error) {
	// Generate a salt with a cost factor of 10
	salt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		print("\n" + err.Error() + "\n")
		return "", err
	}

	// Combine the salt and hashed password
	hashedPassword := string(salt)

	return hashedPassword, nil
}
