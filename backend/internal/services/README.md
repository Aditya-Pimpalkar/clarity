# Services Package

Business logic layer for the LLM Observability Platform.

## Services

### TraceService
Handles trace creation, retrieval, and processing.

**Key Features:**
- Request validation
- Automatic cost calculation
- Status determination
- Async metric recording
- ID and timestamp generation

**Usage:**
```go
traceService := services.NewTraceService(repo)
resp, err := traceService.CreateTrace(ctx, traceRequest)
```

### AnalyticsService
Provides analytics, insights, and aggregations.

**Key Features:**
- Dashboard summaries
- Cost analysis with projections
- Performance metrics
- Model comparisons
- Automated insights generation

**Usage:**
```go
analyticsService := services.NewAnalyticsService(repo)
summary, err := analyticsService.GetDashboardSummary(ctx, orgID, projectID, "24h")
```

### UserService
Manages users, organizations, and projects.

**Key Features:**
- User creation and retrieval
- Organization management
- Project management
- Access validation
- Role-based permissions

**Usage:**
```go
userService := services.NewUserService(repo)
user, err := userService.CreateUser(ctx, email, name, orgID, role)
```

## Architecture
API Handlers
↓
Services  ← Business Logic, Validation, Calculations
↓
Repository ← Data Access
↓
ClickHouse

## Testing

Run service tests:
```bash
go test ./internal/services -v
```

Run with coverage:
```bash
go test ./internal/services -cover
```

## Design Principles

1. **Single Responsibility**: Each service handles one domain
2. **Dependency Injection**: Services receive repository via constructor
3. **Context Propagation**: All methods accept context for cancellation
4. **Error Wrapping**: Errors include context about what failed
5. **Validation**: Input validation happens at service layer
6. **Testability**: Mock repository for unit testing
