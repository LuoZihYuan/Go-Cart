package repository

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/LuoZihYuan/Go-Cart/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type CartDynamoDBRepository struct {
	client     *dynamodb.Client
	tableName  string
	nextCartID int64
}

func NewCartDynamoDBRepository(client *dynamodb.Client) *CartDynamoDBRepository {
	// Initialize with timestamp-based ID
	return &CartDynamoDBRepository{
		client:     client,
		tableName:  "Carts",
		nextCartID: time.Now().Unix(),
	}
}

// Create creates a new cart
func (r *CartDynamoDBRepository) Create(customerID int) (*models.Cart, error) {
	cartID := int(atomic.AddInt64(&r.nextCartID, 1))

	cart := models.Cart{
		CartID:     cartID,
		CustomerID: customerID,
		Items:      []models.CartItem{},
	}

	item, err := attributevalue.MarshalMap(cart)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.client.PutItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

// GetByID retrieves a cart by its ID
func (r *CartDynamoDBRepository) GetByID(cartID int) (*models.Cart, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"cart_id": &types.AttributeValueMemberN{Value: strconv.Itoa(cartID)},
		},
	}

	result, err := r.client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, ErrCartNotFound
	}

	var cart models.Cart
	err = attributevalue.UnmarshalMap(result.Item, &cart)
	if err != nil {
		return nil, err
	}

	// Ensure items is never nil
	if cart.Items == nil {
		cart.Items = []models.CartItem{}
	}

	return &cart, nil
}

// AddItem adds an item to a cart
func (r *CartDynamoDBRepository) AddItem(cartID int, item models.CartItem) error {
	// First, get the current cart to check if product exists
	cart, err := r.GetByID(cartID)
	if err != nil {
		return err
	}

	// Check if product already exists in cart items
	found := false
	for i, existingItem := range cart.Items {
		if existingItem.ProductID == item.ProductID {
			cart.Items[i].Quantity += item.Quantity
			found = true
			break
		}
	}

	if !found {
		cart.Items = append(cart.Items, item)
	}

	// Marshal updated items
	itemsAttr, err := attributevalue.Marshal(cart.Items)
	if err != nil {
		return err
	}

	// Update the cart with new items list
	update := expression.Set(
		expression.Name("items"),
		expression.Value(itemsAttr),
	)

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"cart_id": &types.AttributeValueMemberN{Value: strconv.Itoa(cartID)},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	_, err = r.client.UpdateItem(context.TODO(), input)
	return err
}

// Delete removes a cart (used after checkout)
func (r *CartDynamoDBRepository) Delete(cartID int) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"cart_id": &types.AttributeValueMemberN{Value: strconv.Itoa(cartID)},
		},
		ReturnValues: types.ReturnValueAllOld,
	}

	result, err := r.client.DeleteItem(context.TODO(), input)
	if err != nil {
		return err
	}

	// Check if item existed before deletion
	if len(result.Attributes) == 0 {
		return ErrCartNotFound
	}

	return nil
}
