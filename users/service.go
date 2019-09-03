package users

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"
	"go.uber.org/zap"
	"os"
)

// UseService is the top level signature of this service
type UserService interface {
	Get(ctx context.Context, id string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, id string, user *UpdateUser) error
	Create(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

// Init sets up an instance of this domains
// usecase, pre-configured with the dependencies.
func Init(integration bool) (UserService, error) {
	region := os.Getenv("AWS_REGION")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, err
	}

	ddb := dynamodb.New(sess)
	if integration == false {
		xray.Configure(xray.Config{LogLevel: "trace"})
		xray.AWS(ddb.Client)
	}

	logger, _ := zap.NewProduction()

	tableName := os.Getenv("TABLE_NAME")
	repository := NewDynamoDBRepository(ddb, tableName)
	usecase := &LoggerAdapter{
		Logger:  logger,
		Usecase: &Usecase{Repository: repository},
	}
	return usecase, nil
}
