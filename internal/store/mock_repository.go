package store

type MockLogRepository struct{}

func NewMockLogRepository() *MockLogRepository {
	return &MockLogRepository{}
}
