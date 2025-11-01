#!/bin/sh
set -e

echo "Waiting for DynamoDB Local to be ready..."
sleep 5

ENDPOINT="http://dynamo.gocart-dev:8000"

echo "Creating Products table..."
aws dynamodb create-table \
  --endpoint-url $ENDPOINT \
  --region us-east-1 \
  --table-name Products \
  --attribute-definitions AttributeName=product_id,AttributeType=N \
  --key-schema AttributeName=product_id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  2>/dev/null && echo "✓ Products table created" || echo "✓ Products table already exists"

echo "Creating Carts table..."
aws dynamodb create-table \
  --endpoint-url $ENDPOINT \
  --region us-east-1 \
  --table-name Carts \
  --attribute-definitions \
    AttributeName=cart_id,AttributeType=N \
  --key-schema AttributeName=cart_id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  2>/dev/null && echo "✓ Carts table created" || echo "✓ Carts table already exists"

echo ""
echo "DynamoDB Local tables initialized successfully!"
echo "Tables: Products, Carts"