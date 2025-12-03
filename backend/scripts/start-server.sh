#!/bin/bash

echo "🧹 Cleaning up old processes..."
lsof -ti :8080 | xargs kill -9 2>/dev/null

echo "🚀 Starting server..."
go run cmd/api/main.go
