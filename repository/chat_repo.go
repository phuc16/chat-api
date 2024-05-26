package repository

import (
	"app/config"
	"app/entity"
	"app/errors"
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repo) chatColl() *mongo.Collection {
	return r.db.Database(config.Cfg.DB.DBName).Collection("chats")
}

func (r *Repo) CreateChatIndexes(ctx context.Context) ([]string, error) {
	indexes, err := r.chatColl().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{
			{"id", 1},
		}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"chat_name", 1},
		}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"group_admin", 1},
		}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"latest_message", 1},
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

func (r *Repo) SaveChat(ctx context.Context, chat *entity.Chat) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	opts := options.Update().SetUpsert(true)
	update := bson.D{
		{"$set", chat},
	}
	_, err = r.chatColl().UpdateOne(ctx, bson.M{"id": chat.ID}, update, opts)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error collection") {
			return errors.ChatExists()
		}
		return err
	}
	return nil
}

func (r *Repo) GetChatById(ctx context.Context, id string) (res *entity.Chat, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	var d []*entity.Chat

	pipeLine := mongo.Pipeline{}
	pipeLine = append(pipeLine, matchFieldPipeline("id", id))
	pipeLine = append(pipeLine, limitPipeline(1))

	cursor, err := r.chatColl().Aggregate(ctx, pipeLine, collationAggregateOption)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &d); err != nil {
		return nil, err
	}
	if len(d) <= 0 {
		return nil, errors.ChatNotFound()
	}
	return d[0], nil
}

func (r *Repo) GetChatList(ctx context.Context, params *QueryParams) (res []*entity.Chat, total int64, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	coll := r.chatColl()

	pipeLine := mongo.Pipeline{}
	if params.Search != "" {
		pipeLine = append(pipeLine, partialMatchingSearchPipeline([]string{"name", "chat_name", "email"}, params.Search)...)
	}
	for k, v := range params.Filter {
		pipeLine = append(pipeLine, matchFieldPipeline(k, v))
	}

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

func (r *Repo) UpdateChat(ctx context.Context, chat *entity.Chat) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	filter := bson.D{{"id", chat.ID}}
	update := bson.M{"$set": chat}
	_, err = r.chatColl().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) AddToGroup(ctx context.Context, chat *entity.Chat) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	filter := bson.M{"_id": chat.ID}
	update := bson.M{
		"$addToSet": bson.M{"users": chat.Users[0]},
		"$set":      bson.M{"updated_at": time.Now()},
	}
	_, err = r.chatColl().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil

}

func (r *Repo) RemoveFromGroup(ctx context.Context, chat *entity.Chat) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	filter := bson.M{"_id": chat.ID}
	update := bson.M{
		"$pull": bson.M{"users": chat.Users[0]},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err = r.chatColl().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil

}
