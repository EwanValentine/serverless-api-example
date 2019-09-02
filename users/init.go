package users

import (
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Init sets up an instance of this domains
// usecase, pre-configured with the dependencies.
func Init() (*Usecase, error) {
	region := os.Getenv("AWS_REGION")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, err
	}

	tableName := os.Getenv("TABLE_NAME")
	repository := NewDynamoDBRepository(dynamodb.New(sess), tableName)
	usecase := &Usecase{Repository: repository}
	return usecase, nil
}