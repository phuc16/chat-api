package service

import (
	"app/dto"
	"app/entity"
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
)

type UserService struct {
	UserRepo       IUserRepo
	TokenRepo      ITokenRepo
	UpdateAsyncSvc IUpdateAsyncSvc
}

func NewUserService(userRepo IUserRepo, tokenRepo ITokenRepo, updateAsyncSvc IUpdateAsyncSvc) *UserService {
	return &UserService{UserRepo: userRepo, TokenRepo: tokenRepo, UpdateAsyncSvc: updateAsyncSvc}
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserReq) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	user := &entity.User{
		ID: req.ID,
	}
	err = s.UserRepo.SaveUser(ctx, user)
	return
}

func (s *UserService) GetUser(ctx context.Context, id string) (res *entity.User, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	return s.UserRepo.FindUserByID(ctx, id)
}

func (s *UserService) UpdateAvatarAsync(ctx context.Context, req *dto.UpdateAvatarAsyncReq) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	s.UpdateAsyncSvc.UpdateAvatarAsync(ctx, req.OldAvatar, req.NewAvatar)
	return
}

// func (s *UserService) SuggestFriend(ctx context.Context, e *entity.User) (res []*entity.User, err error) {
// 	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
// 	defer span.End()

// 	targetUser, err := s.UserRepo.GetUserById(ctx, e.ID)
// 	if err != nil {
// 		return
// 	}

// 	allUsers, err := s.UserRepo.GetAllUsers(ctx)
// 	if err != nil {
// 		return
// 	}

// 	// Find nearest neighbors of the target user
// 	nearestNeighbors := s.FindNearestNeighbors(targetUser, allUsers)

// 	// Suggest friends of nearest neighbors
// 	for _, neighbor := range nearestNeighbors {
// 		for _, friend := range neighbor.Friends {
// 			if !s.Contains(targetUser.FriendIds, friend.ID) && friend.ID != targetUser.ID {
// 				res = append(res, friend)
// 			}
// 		}
// 	}
// 	return
// }

// // Function to find nearest neighbors based on cosine similarity
// func (s *UserService) FindNearestNeighbors(targetUser *entity.User, allUsers []*entity.User) []*entity.User {
// 	// Calculate similarity with all users
// 	var similarities []float64
// 	for _, user := range allUsers {
// 		similarity := s.CosineSimilarity(targetUser.FriendIds, user.FriendIds)
// 		similarities = append(similarities, similarity)
// 	}

// 	// Find nearest neighbors
// 	var nearestNeighbors []*entity.User
// 	for i, similarity := range similarities {
// 		if similarity > 0.8 && allUsers[i].ID != targetUser.ID { // Adjust threshold as needed
// 			nearestNeighbors = append(nearestNeighbors, allUsers[i])
// 		}
// 	}

// 	return nearestNeighbors
// }

// // Function to calculate cosine similarity between two vectors
// func (s *UserService) CosineSimilarity(a, b []string) float64 {
// 	// Convert friend lists to binary vectors
// 	vectorA := make([]int, len(a))
// 	vectorB := make([]int, len(b))
// 	for i, friendID := range a {
// 		if s.Contains(b, friendID) {
// 			vectorA[i] = 1
// 		}
// 	}
// 	for i, friendID := range b {
// 		if s.Contains(a, friendID) {
// 			vectorB[i] = 1
// 		}
// 	}

// 	// Calculate dot product
// 	dotProduct := 0.0
// 	for i := 0; i < len(vectorA); i++ {
// 		dotProduct += float64(vectorA[i] * vectorB[i])
// 	}

// 	// Calculate magnitudes
// 	magnitudeA := 0.0
// 	magnitudeB := 0.0
// 	for _, val := range vectorA {
// 		magnitudeA += float64(val * val)
// 	}
// 	for _, val := range vectorB {
// 		magnitudeB += float64(val * val)
// 	}
// 	magnitudeA = math.Sqrt(magnitudeA)
// 	magnitudeB = math.Sqrt(magnitudeB)

// 	// Calculate cosine similarity
// 	similarity := dotProduct / (magnitudeA * magnitudeB)
// 	return similarity
// }

// // Function to check if a user ID exists in a list of friend IDs
// func (s *UserService) Contains(friends []string, userID string) bool {
// 	for _, friendID := range friends {
// 		if friendID == userID {
// 			return true
// 		}
// 	}
// 	return false
// }
