package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/romanpitatelev/denet/internal/controller/rest"
	taskshandler "github.com/romanpitatelev/denet/internal/controller/rest/tasks-handler"
	usershandler "github.com/romanpitatelev/denet/internal/controller/rest/users-handler"
	"github.com/romanpitatelev/denet/internal/entity"
	"github.com/romanpitatelev/denet/internal/repository/store"
	tasksrepo "github.com/romanpitatelev/denet/internal/repository/tasks-repo"
	usersrepo "github.com/romanpitatelev/denet/internal/repository/users-repo"
	tasksservice "github.com/romanpitatelev/denet/internal/usecase/task-service"
	usersservice "github.com/romanpitatelev/denet/internal/usecase/users-service"
	"github.com/rs/zerolog/log"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/suite"
)

const (
	pgDSN     = "postgresql://postgres:my_pass@localhost:5432/denet_db"
	port      = 5003
	usersPath = "/api/v1/users"
)

type IntegrationTestSuite struct {
	suite.Suite
	cancelFunc   context.CancelFunc
	db           *store.DataStore
	usersrepo    *usersrepo.Repo
	tasksrepo    *tasksrepo.Repo
	usersservice *usersservice.Service
	tasksservice *tasksservice.Service
	usershandler *usershandler.Handler
	taskshandler *taskshandler.Handler
	server       *rest.Server
}

func (s *IntegrationTestSuite) SetupSuite() {
	log.Info().Msg("starting SetupSuite ...")

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel

	var err error

	s.db, err = store.New(ctx, store.Config{Dsn: pgDSN})
	s.Require().NoError(err)

	log.Info().Msg("starting new db ...")

	err = s.db.Migrate(migrate.Up)
	s.Require().NoError(err)

	log.Info().Msg("migrations are ready")

	s.usersrepo = usersrepo.New(s.db)
	s.tasksrepo = tasksrepo.New(s.db)

	s.usersservice = usersservice.New(s.usersrepo)
	s.tasksservice = tasksservice.New(s.tasksrepo)

	s.usershandler = usershandler.New(s.usersservice)
	s.taskshandler = taskshandler.New(s.tasksservice)

	s.server = rest.New(
		rest.Config{BindAddress: fmt.Sprintf(":%d", port)},
		s.usershandler,
		s.taskshandler,
		rest.GetPublicKey(),
	)

	//nolint:testifylint
	go func() {
		err = s.server.Run(ctx)
		s.Require().NoError(err)
	}()

	time.Sleep(50 * time.Millisecond)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.cancelFunc()
}

func (s *IntegrationTestSuite) TearDownTest() {
	err := s.db.Truncate(context.Background(),
		"reference",
		"tasks",
		"users",
	)
	s.Require().NoError(err)
}

func TestIntegrationSetupSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) sendRequest(method, path string, status int, entity, result any, user entity.User) {
	body, err := json.Marshal(entity)
	s.Require().NoError(err)

	requestURL := fmt.Sprintf("http://localhost:%d%s", port, path)
	s.T().Logf("Sending request to %s", requestURL)

	request, err := http.NewRequestWithContext(context.Background(), method,
		fmt.Sprintf("http://localhost:%d%s", port, path), bytes.NewReader(body))
	s.Require().NoError(err, "fail to create request")

	token := s.getToken(user)

	request.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}

	response, err := client.Do(request)

	s.Require().NoError(err, "fail to execute request")

	s.Require().NotNil(response, "response object is nil")

	defer func() {
		err = response.Body.Close()
		s.Require().NoError(err)
	}()

	s.T().Logf("Response Status Code: %d", response.StatusCode)

	if status != response.StatusCode {
		responseBody, err := io.ReadAll(response.Body)
		s.Require().NoError(err)

		s.T().Logf("Response Body: %s", string(responseBody))

		s.Require().Equal(status, response.StatusCode, "unexpected status code")

		return
	}

	if result == nil {
		return
	}

	err = json.NewDecoder(response.Body).Decode(result)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) getToken(user entity.User) string {
	now := time.Now()
	updatedAt := now
	claims := entity.Claims{
		UserID:    user.ID,
		Email:     &user.Email,
		Role:      &user.Role,
		Points:    user.Points,
		CreatedAt: now,
		UpdatedAt: &updatedAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	privateKey, err := readPrivateKey()
	s.Require().NoError(err)

	token, err := generateToken(&claims, privateKey)
	s.Require().NoError(err)

	return token
}
