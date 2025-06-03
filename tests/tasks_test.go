package tests

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/romanpitatelev/denet/internal/entity"
)

func (s *IntegrationTestSuite) TestTasks() {
	user := entity.User{
		ID:    entity.UserID(uuid.New()),
		Name:  "task function name",
		Email: "task@mail.ru",
		Role:  "Task student",
	}

	err := s.db.UpsertUser(context.Background(), user)
	s.Require().NoError(err)

	s.Run("complete task successfully", func() {
		taskPath := usersPath + "/" + uuid.UUID(user.ID).String() + "/task/complete"

		task := entity.Task{
			Type: "telegram",
		}

		var response entity.TaskResponse

		s.sendRequest(http.MethodPost, taskPath, http.StatusOK, task, &response, user)

		s.Require().Equal(task.Type, response.Type)
		s.Require().Equal(3, response.Points)
		s.Require().Equal(3, response.TotatlPoints)
	})

	s.Run("invalid task type", func() {
		taskPath := usersPath + "/" + uuid.UUID(user.ID).String() + "/task/complete"

		task := entity.Task{
			Type: "invalid type",
		}

		s.sendRequest(http.MethodPost, taskPath, http.StatusBadRequest, task, nil, user)
	})

	s.Run("invalid task points", func() {
		taskPath := usersPath + "/" + uuid.UUID(user.ID).String() + "/task/complete"

		task := entity.Task{
			Type:   "twitter",
			Points: 10,
		}

		s.sendRequest(http.MethodPost, taskPath, http.StatusBadRequest, task, nil, user)
	})

	s.Run("task for non-existent user", func() {
		taskPath := usersPath + "/" + uuid.New().String() + "/task/complete"

		task := entity.Task{
			Type: "telegram",
		}

		s.sendRequest(http.MethodPost, taskPath, http.StatusNotFound, task, nil, entity.User{})
	})
}

func (s *IntegrationTestSuite) TestReferral() {
	referrer := entity.User{
		ID:    entity.UserID(uuid.New()),
		Name:  "referral function name",
		Email: "referral@mail.ru",
		Role:  "Referral student",
	}

	err := s.db.UpsertUser(context.Background(), referrer)
	s.Require().NoError(err)

	referee := entity.User{
		ID:    entity.UserID(uuid.New()),
		Name:  "referee name",
		Email: "referee@mail.ru",
		Role:  "Referee student",
	}

	err = s.db.UpsertUser(context.Background(), referee)
	s.Require().NoError(err)

	s.Run("referral task successful", func() {
		referral := entity.Reference{
			UserReferenceID: referee.ID,
		}

		var response entity.ReferenceResponse

		referencePath := usersPath + "/" + uuid.UUID(referrer.ID).String() + "/referrer"

		s.sendRequest(http.MethodPost, referencePath, http.StatusOK, referral, &response, referrer)

		s.Require().Equal(referrer.ID, response.UserID)
		s.Require().Equal(referee.ID, response.UserReferenceID)
		s.Require().Equal(1, response.Points)
		s.Require().Equal(1, response.TotatlPoints)
	})

	s.Run("self referral fails", func() {
		referral := entity.Reference{
			UserReferenceID: referrer.ID,
		}

		referencePath := usersPath + "/" + uuid.UUID(referrer.ID).String() + "/referrer"
		s.sendRequest(http.MethodPost, referencePath, http.StatusBadRequest, referral, nil, referrer)
	})

	s.Run("set force points to default value", func() {
		referral := entity.Reference{
			UserReferenceID: referee.ID,
			Points:          1000,
		}

		var response entity.ReferenceResponse

		referencePath := usersPath + "/" + uuid.UUID(referrer.ID).String() + "/referrer"

		s.sendRequest(http.MethodPost, referencePath, http.StatusOK, referral, &response, referrer)

		s.Require().Equal(referrer.ID, response.UserID)
		s.Require().Equal(referee.ID, response.UserReferenceID)
		s.Require().Equal(1, response.Points)
		s.Require().Equal(2, response.TotatlPoints)
	})
}
