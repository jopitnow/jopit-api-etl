package authorization

import (
	"github.com/gin-gonic/gin"
)

type FirebaseMock struct {
	HandleGetUserId func(ctx *gin.Context) (string, error)
}

func NewAuthorizationMock() FirebaseMock {
	return FirebaseMock{}
}

func (mock FirebaseMock) GetUserId(ctx *gin.Context) (string, error) {
	if mock.HandleGetUserId != nil {
		return mock.HandleGetUserId(ctx)
	}
	return "", nil
}
