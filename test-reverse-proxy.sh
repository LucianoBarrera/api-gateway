#!/bin/bash

# Test script for reverse proxy functionality
# Make sure to run: docker-compose up -d first

echo "Testing API Gateway Reverse Proxy with Mock Servers"
echo "=================================================="

# Wait for services to be ready
echo "Waiting for services to be ready..."
sleep 10

# Test health check
echo -e "\n1. Testing API Gateway health check:"
curl -s http://localhost:8080/liveness | jq .

# Test middleware requirements
echo -e "\n2. Testing middleware requirements:"
echo "Testing missing X-Request-ID header:"
curl -s -H "x-api-key: example-api-key-local-env" http://localhost:8080/api/users/users | jq .

echo -e "\nTesting missing x-api-key header:"
curl -s -H "X-Request-ID: test-request-id" http://localhost:8080/api/users/users | jq .

echo -e "\nTesting invalid x-api-key:"
curl -s -H "X-Request-ID: test-request-id" -H "x-api-key: invalid-key" http://localhost:8080/api/users/users | jq .

# Test users service through API gateway with proper headers
echo -e "\n3. Testing users service through API gateway (with proper headers):"
echo "GET /api/users/users"
curl -s -H "X-Request-ID: test-request-id-1" -H "x-api-key: example-api-key-local-env" http://localhost:8080/api/users/users | jq .

echo -e "\nGET /api/users/users/123"
curl -s -H "X-Request-ID: test-request-id-2" -H "x-api-key: example-api-key-local-env" http://localhost:8080/api/users/users/123 | jq .

echo -e "\nPOST /api/users/users"
curl -s -X POST http://localhost:8080/api/users/users \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: test-request-id-3" \
  -H "x-api-key: example-api-key-local-env" \
  -d '{"name": "Test User", "email": "test@example.com"}' | jq .

# Test auth service through API gateway with proper headers
echo -e "\n4. Testing auth service through API gateway (with proper headers):"
echo "GET /api/auth/auth/status"
curl -s -H "X-Request-ID: test-request-id-4" -H "x-api-key: example-api-key-local-env" http://localhost:8080/api/auth/auth/status | jq .

echo -e "\nPOST /api/auth/auth/login"
curl -s -X POST http://localhost:8080/api/auth/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: test-request-id-5" \
  -H "x-api-key: example-api-key-local-env" \
  -d '{"email": "user@example.com", "password": "password123"}' | jq .

# Test non-existent service with proper headers
echo -e "\n5. Testing non-existent service (with proper headers):"
curl -s -H "X-Request-ID: test-request-id-6" -H "x-api-key: example-api-key-local-env" http://localhost:8080/api/nonexistent/test | jq .

# Test direct access to mock servers (for comparison)
echo -e "\n6. Testing direct access to mock servers:"
echo "Direct access to users service:"
curl -s http://localhost:8081/users | jq .

echo -e "\nDirect access to auth service:"
curl -s http://localhost:8082/auth/status | jq .

echo -e "\nTest completed!" 