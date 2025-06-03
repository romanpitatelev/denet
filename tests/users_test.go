package tests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/romanpitatelev/denet/internal/entity"
	"github.com/romanpitatelev/denet/internal/utils"
)

const defaultTopLimit = 20

func (s *IntegrationTestSuite) TestGetUpdateDeleteUser() {
	user := entity.User{
		ID:     entity.UserID(uuid.New()),
		Name:   "Some name",
		Email:  "email@mail.ru",
		Role:   "student",
		Points: 10,
	}

	err := s.db.UpsertUser(context.Background(), user)
	s.Require().NoError(err)

	var newUser entity.User

	s.Run("get user successfully", func() {
		userIDPath := usersPath + "/" + uuid.UUID(user.ID).String() + "/status"

		s.sendRequest(http.MethodGet, userIDPath, http.StatusOK, nil, &newUser, user)
		s.Require().Equal(user.ID, newUser.ID)
		s.Require().Equal(user.Name, newUser.Name)
		s.Require().Equal(user.Email, newUser.Email)
		s.Require().Equal(user.Role, newUser.Role)
		s.Require().Equal(user.Points, newUser.Points)
	})

	s.Run("user not found", func() {
		userIDPath := usersPath + "/" + uuid.New().String() + "/status"

		s.sendRequest(http.MethodGet, userIDPath, http.StatusNotFound, nil, nil, user)
	})

	s.Run("update does not exist", func() {
		userIDPath := usersPath + "/" + uuid.New().String()

		s.sendRequest(http.MethodPatch, userIDPath, http.StatusNotFound, nil, nil, user)
	})

	user2 := entity.User{
		ID:     entity.UserID(uuid.New()),
		Name:   "Roman",
		Email:  "roman@mail.ru",
		Role:   "nobody",
		Points: 20,
	}

	err = s.db.UpsertUser(context.Background(), user2)
	s.Require().NoError(err)

	updateUser := entity.UserUpdate{
		Name:  utils.Pointer("Biba"),
		Email: &user.Email,
		Role:  utils.Pointer("test role"),
	}

	s.Run("update conflict contact", func() {
		s.sendRequest(http.MethodPatch, usersPath+"/"+uuid.UUID(user2.ID).String(), http.StatusConflict, updateUser, nil, user2)
	})

	var updatedUser entity.User

	updateUser.Email = nil

	s.Run("name updated successfully", func() {
		userIDPath := usersPath + "/" + uuid.UUID(user2.ID).String()

		s.sendRequest(http.MethodPatch, userIDPath, http.StatusOK, updateUser, &updatedUser, user2)
		s.Require().Equal(user2.ID, updatedUser.ID)
		s.Require().Equal(*updateUser.Name, updatedUser.Name)
		s.Require().Equal(*updateUser.Role, updatedUser.Role)
	})

	s.Run("delete non-existent user", func() {
		userIDPath := usersPath + "/" + uuid.New().String()
		s.sendRequest(http.MethodDelete, userIDPath, http.StatusNotFound, nil, nil, entity.User{})
	})

	s.Run("delete existing user", func() {
		userIDPath := usersPath + "/" + uuid.UUID(user.ID).String()
		s.sendRequest(http.MethodDelete, userIDPath, http.StatusNoContent, nil, nil, user)
	})
}

func (s *IntegrationTestSuite) TestGetUsers() {
	users := []entity.User{
		{
			ID:     entity.UserID(uuid.New()),
			Name:   "a",
			Email:  "a@mail.ru",
			Role:   "manager1",
			Points: 100,
		},
		{
			ID:     entity.UserID(uuid.New()),
			Name:   "b",
			Email:  "b@mail.ru",
			Role:   "manager2",
			Points: 200,
		},
		{
			ID:     entity.UserID(uuid.New()),
			Name:   "c",
			Email:  "c@mail.ru",
			Role:   "manager3",
			Points: 150,
		},
		{
			ID:     entity.UserID(uuid.New()),
			Name:   "d",
			Email:  "d@mail.ru",
			Role:   "manager4",
			Points: 400,
		},
		{
			ID:     entity.UserID(uuid.New()),
			Name:   "e",
			Email:  "e@mail.ru",
			Role:   "manager5",
			Points: 1,
		},
	}

	for _, user := range users {
		err := s.db.UpsertUser(context.Background(), user)
		s.Require().NoError(err)
	}

	s.Run("get users with default sorting - by name", func() {
		var result []entity.User

		s.sendRequest(http.MethodGet, usersPath, http.StatusOK, nil, &result, users[0])

		s.Require().Len(result, 5)
		s.Require().Equal("a", result[0].Name)
		s.Require().Equal("b", result[1].Name)
		s.Require().Equal("c", result[2].Name)
	})

	s.Run("sort by points with limit 2", func() {
		var result []entity.User

		path := usersPath + "?sorting=points&limit=2"

		s.sendRequest(http.MethodGet, path, http.StatusOK, nil, &result, users[0])

		s.Require().Len(result, 2)
		s.Require().Equal(result[0].ID, users[4].ID)
		s.Require().Equal(result[1].ID, users[0].ID)
	})

	s.Run("sort by points with limit 2 and offset 2", func() {
		var result []entity.User

		path := usersPath + "?sorting=points&limit=2&offset=2"

		s.sendRequest(http.MethodGet, path, http.StatusOK, nil, &result, users[0])

		s.Require().Len(result, 2)
		s.Require().Equal(result[0].ID, users[2].ID)
		s.Require().Equal(result[1].ID, users[1].ID)
	})

	s.Run("sort by points with limit 2 and offset 2, descending true", func() {
		var result []entity.User

		path := usersPath + "?sorting=points&limit=2&offset=2&descending=true"

		s.sendRequest(http.MethodGet, path, http.StatusOK, nil, &result, users[0])

		s.Require().Len(result, 2)
		s.Require().Equal(result[0].ID, users[2].ID)
		s.Require().Equal(result[1].ID, users[0].ID)
	})
}

func (s *IntegrationTestSuite) TestGetTopUsers() {
	users := make([]entity.User, 30)
	for i := range 30 {
		users[i] = entity.User{
			ID:     entity.UserID(uuid.New()),
			Name:   fmt.Sprintf("User %d", i+1),
			Email:  fmt.Sprintf("user%d@mail.ya", i+1),
			Role:   "student",
			Points: i + 1,
		}

		err := s.db.UpsertUser(context.Background(), users[i])
		s.Require().NoError(err)
	}

	s.Run("get top users with default limit", func() {
		var result []entity.User

		path := usersPath + "/leaderboard"
		s.sendRequest(http.MethodGet, path, http.StatusOK, nil, &result, users[0])

		s.Require().Len(result, defaultTopLimit)
		s.Require().Equal(30, result[0].Points)
	})
}
