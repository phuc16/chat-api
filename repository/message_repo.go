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

func (r *Repo) messageColl() *mongo.Collection {
	return r.db.Database(config.Cfg.DB.DBName).Collection("messages")
}

func (r *Repo) CreateMessageIndexes(ctx context.Context) ([]string, error) {
	indexes, err := r.messageColl().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{
			{"id", 1},
		}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"sender", 1},
		}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"message", 1},
		}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"chat_id", 1},
		}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"created_at", 1},
		}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		return nil, err
	}
	return indexes, nil
}

func (r *Repo) SaveMessage(ctx context.Context, message *entity.Message) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	opts := options.Update().SetUpsert(true)
	update := bson.D{
		{"$set", message},
	}
	_, err = r.messageColl().UpdateOne(ctx, bson.M{"id": message.ID}, update, opts)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error collection") {
			return errors.MessageExists()
		}
		return err
	}
	return nil
}

func (r *Repo) GetMessageListByChatId(ctx context.Context, chatId string, params *QueryParams) (res []*entity.Message, total int64, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	coll := r.messageColl()

	pipeLine := mongo.Pipeline{}
	if params.Search != "" {
		pipeLine = append(pipeLine, partialMatchingSearchPipeline([]string{"name", "message", "email"}, params.Search)...)
	}
	for k, v := range params.Filter {
		pipeLine = append(pipeLine, matchFieldPipeline(k, v))
	}
	pipeLine = append(pipeLine, matchFieldPipeline("chat_id", chatId))

	cursor, err := coll.Aggregate(ctx, append(pipeLine, bson.D{{"$count", "total_count"}}))
	if err != nil {
		return nil, 0, err
	}
	result := bson.M{}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, 0, err
		}
	}
	totalCount, _ := result["total_count"].(int32)
	total = int64(totalCount)

	pipeLine = append(pipeLine, params.SkipLimitSortPipeline()...)

	cursor, err = coll.Aggregate(ctx, pipeLine, collationAggregateOption)
	if err != nil {
		return res, 0, err
	}
	if err = cursor.All(ctx, &res); err != nil {
		return res, 0, err
	}
	return res, total, nil
}
