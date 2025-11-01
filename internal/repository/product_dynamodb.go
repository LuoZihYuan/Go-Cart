package repository

import (
	"context"
	"strconv"

	"github.com/LuoZihYuan/Go-Cart/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ProductDynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewProductDynamoDBRepository(client *dynamodb.Client) *ProductDynamoDBRepository {
	return &ProductDynamoDBRepository{
		client:    client,
		tableName: "Products",
	}
}

// GetByID retrieves a product by its ID
func (r *ProductDynamoDBRepository) GetByID(productID int) (*models.Product, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"product_id": &types.AttributeValueMemberN{Value: strconv.Itoa(productID)},
		},
	}

	result, err := r.client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, ErrProductNotFound
	}

	var product models.Product
	err = attributevalue.UnmarshalMap(result.Item, &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// Upsert creates or updates a product's details
func (r *ProductDynamoDBRepository) Upsert(product *models.Product) error {
	item, err := attributevalue.MarshalMap(product)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.client.PutItem(context.TODO(), input)
	return err
}

// Exists checks if a product exists
func (r *ProductDynamoDBRepository) Exists(productID int) (bool, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"product_id": &types.AttributeValueMemberN{Value: strconv.Itoa(productID)},
		},
		ProjectionExpression: aws.String("product_id"),
	}

	result, err := r.client.GetItem(context.TODO(), input)
	if err != nil {
		return false, err
	}

	return result.Item != nil, nil
}
