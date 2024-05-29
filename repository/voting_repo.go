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

func (r *Repo) votingColl() *mongo.Collection {
	return r.db.Database(config.Cfg.DB.DBName).Collection("votings")
}

func (r *Repo) SaveVoting(ctx context.Context, voting *entity.Voting) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	opts := options.Update().SetUpsert(true)
	update := bson.D{
		{"$set", voting},
	}
	_, err = r.votingColl().UpdateOne(ctx, bson.M{"id": voting.ID}, update, opts)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error collection") {
			return errors.VotingExists()
		}
		return err
	}
	return nil
}

func (r *Repo) FindVotingByID(ctx context.Context, id string) (*entity.Voting, error) {
	var d entity.Voting
	filter := bson.D{
		{"id", id},
	}
	if err := r.votingColl().FindOne(ctx, filter).Decode(&d); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.VotingNotFound()
		}
		return nil, err
	}
	return &d, nil
}

func (r *Repo) AppendVoter(ctx context.Context, id string, name string, info entity.PersonInfo) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": id, "choices.name": name}
	update := bson.M{"$push": bson.M{"choices.$.voters": info}}
	return r.votingColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) RemoveVoter(ctx context.Context, id string, name string, userID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": id, "choices.name": name}
	update := bson.M{"$pull": bson.M{"choices.$.voters": bson.M{"user_id": userID}}}
	return r.votingColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) LockVoting(ctx context.Context, id string, isLock bool, dateLock time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"lock": isLock, "date_lock": dateLock}}
	return r.votingColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) DeleteVotingByID(ctx context.Context, id string) error {
	filter := bson.D{{"id", id}}
	_, err := r.votingColl().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
