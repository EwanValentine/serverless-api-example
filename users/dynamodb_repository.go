package users

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

// DynamoDBRepository -
type DynamoDBRepository struct {
	session   *dynamodb.DynamoDB
	tableName string
}

// NewDynamoDBRepository -
func NewDynamoDBRepository(ddb *dynamodb.DynamoDB, tableName string) *DynamoDBRepository {
	return &DynamoDBRepository{ddb, tableName}
}

// Get a user
func (r *DynamoDBRepository) Get(ctx context.Context, id string) (*User, error) {
	user := &User{}
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	result, err := r.session.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if err := dynamodbattribute.UnmarshalMap(result.Item, &user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetAll users
func (r *DynamoDBRepository) GetAll(ctx context.Context) ([]*User, error) {
	users := make([]*User, 0)
	result, err := r.session.ScanWithContext(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return nil, err
	}

	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &users); err != nil {
		return nil, err
	}

	return users, nil
}

type updateUser struct {
	Name  string `json:":n"`
	Age   uint32 `json:":a"`
	Email string `json:":e"`
}

type userKey struct {
	ID string `json:":id"`
}

// Update a user
func (r *DynamoDBRepository) Update(ctx context.Context, id string, user *UpdateUser) error {
	log.Println("id", id)
	update, err := dynamodbattribute.MarshalMap(&updateUser{
		Name:  user.Name,
		Age:   user.Age,
		Email: user.Email,
	})
	if err != nil {
		return nil
	}

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		ExpressionAttributeValues: update,
		TableName:                 aws.String(r.tableName),
		UpdateExpression:          aws.String("set #uname = :n, age = :a, email = :e"),
		ReturnValues:              aws.String("UPDATED_NEW"),
		ExpressionAttributeNames: map[string]*string{
			"#uname": aws.String("name"),
		},
	}
	_, err = r.session.UpdateItemWithContext(ctx, input)
	return err
}

// Create a user
func (r *DynamoDBRepository) Create(ctx context.Context, user *User) error {
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(r.tableName),
	}
	_, err = r.session.PutItemWithContext(ctx, input)
	return err
}

// Delete a user
func (r *DynamoDBRepository) Delete(ctx context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	_, err := r.session.DeleteItemWithContext(ctx, input)
	return err
}
