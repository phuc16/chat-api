package repository

import (
	"app/config"
	"app/entity"
	"app/errors"
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repo) accountColl() *mongo.Collection {
	return r.db.Database(config.Cfg.DB.DBName).Collection("accounts")
}

func (r *Repo) CreateAccountIndexes(ctx context.Context) ([]string, error) {
	indexes, err := r.accountColl().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{
			{"phone_number", 1},
		}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		return nil, err
	}
	return indexes, nil
}

func (r *Repo) SaveAccount(ctx context.Context, account *entity.Account) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	opts := options.Update().SetUpsert(true)
	update := bson.D{
		{"$set", account},
	}
	_, err = r.accountColl().UpdateOne(ctx, bson.M{"id": account.ID}, update, opts)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error collection") {
			return errors.AccountExists()
		}
		return err
	}
	return nil
}

func (r *Repo) FindAccountByID(ctx context.Context, id string) (*entity.Account, error) {
	var d entity.Account
	filter := bson.D{
		{"id", id},
	}
	if err := r.accountColl().FindOne(ctx, filter).Decode(&d); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.AccountNotFound()
		}
		return nil, err
	}
	return &d, nil
}

func (r *Repo) GetAllAccounts(ctx context.Context) ([]*entity.Account, error) {
	var res []*entity.Account
	pipeline := mongo.Pipeline{
		{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "profile.user_id"},
			{"foreignField", "id"},
			{"as", "user"},
		}}},
	}

	pipeline = append(pipeline, bson.D{
		{"$unwind", bson.D{
			{"path", "$user"},
			{"preserveNullAndEmptyArrays", false},
		}},
	})
	cursor, err := r.accountColl().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repo) SearchByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.Account, error) {
	filter := bson.M{"phone_number": phoneNumber}
	var account entity.Account
	err := r.accountColl().FindOne(ctx, filter).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.AccountNotFound()
		}
		return nil, err
	}
	return &account, nil
}

func (r *Repo) ChangePassword(ctx context.Context, phoneNumber, password string) (*mongo.UpdateResult, error) {
	filter := bson.M{"phone_number": phoneNumber}
	update := bson.M{"$set": bson.M{"pw": password}}
	return r.accountColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) ChangeAvatar(ctx context.Context, phoneNumber string, profile entity.Profile) (*mongo.UpdateResult, error) {
	filter := bson.M{"phone_number": phoneNumber}
	update := bson.M{"$set": bson.M{"profile": profile}}
	return r.accountColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) SearchByUserID(ctx context.Context, userID string) (*entity.Account, error) {
	filter := bson.M{"profile.user_id": userID}
	var account entity.Account
	err := r.accountColl().FindOne(ctx, filter).Decode(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *Repo) DeleteAccountByID(ctx context.Context, id string) error {
	filter := bson.D{{"id", id}}
	_, err := r.accountColl().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) UpdateAccount(ctx context.Context, account *entity.Account) error {
	filter := bson.D{{"id", account.ID}}
	update := bson.D{{"$set", account}}
	_, err := r.accountColl().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
