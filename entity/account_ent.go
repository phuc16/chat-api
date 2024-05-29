package entity

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	accountKey     = "account"
	accountNameKey = "phoneNumber"
)
const (
	ROLE_USER  string = "account"
	ROLE_ADMIN string = "admin"
)

type Account struct {
	ID          string    `bson:"id" json:"id"`
	PhoneNumber string    `bson:"phone_number" json:"phoneNumber"`
	Pw          string    `bson:"pw" json:"pw"`
	Type        string    `bson:"type" json:"type"`
	Profile     Profile   `bson:"profile" json:"profile"`
	Role        string    `bson:"role" json:"role"`
	Setting     Setting   `bson:"setting" json:"setting"`
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
}

func (e *Account) SetToContext(ctx *gin.Context) {
	ctx.Set(accountNameKey, e.PhoneNumber)
	newCtx := context.WithValue(ctx.Request.Context(), accountKey, *e)
	ctx.Request = ctx.Request.WithContext(newCtx)
}

func GetAccountFromContext(ctx context.Context) Account {
	account := Account{}
	value := ctx.Value(accountKey)
	if value != nil {
		account = value.(Account)
	}
	return account
}
