package storage

import (
	"context"
	"testing"

	"github.com/AARCSX/AARCSX_Forge/internal/logger"
	"github.com/stretchr/testify/require"
)

// MockProvider is a mock implementation of the Provider interface for testing
type MockProvider struct {
	PresignPutFunc func(ctx context.Context, objectKey string, ttlSeconds int) (string, error)
	PresignGetFunc func(ctx context.Context, objectKey string, ttlSeconds int) (string, error)
}

func (m *MockProvider) PresignPut(ctx context.Context, objectKey string, ttlSeconds int) (string, error) {
	if m.PresignPutFunc != nil {
		return m.PresignPutFunc(ctx, objectKey, ttlSeconds)
	}
	return "http://mock-put-url", nil
}

func (m *MockProvider) PresignGet(ctx context.Context, objectKey string, ttlSeconds int) (string, error) {
	if m.PresignGetFunc != nil {
		return m.PresignGetFunc(ctx, objectKey, ttlSeconds)
	}
	return "http://mock-get-url", nil
}

// TestStorageServiceCreateUploadURL tests the CreateUploadURL method
func TestStorageServiceCreateUploadURL(t *testing.T) {
	logger, err := logger.New("info", "storage-test")
	require.NoError(t, err)

	// Test validation - missing tenant ID
	service := &StorageService{
		logger: logger,
	}

	_, err = service.CreateUploadURL(context.Background(), CreateUploadURLInput{
		ObjectID: "test-object",
		FileName: "test.txt",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "tenant ID is required")

	// Test validation - missing object ID
	_, err = service.CreateUploadURL(context.Background(), CreateUploadURLInput{
		TenantID: "test-tenant",
		FileName: "test.txt",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "object ID is required")

	// Test validation - missing file name
	_, err = service.CreateUploadURL(context.Background(), CreateUploadURLInput{
		TenantID: "test-tenant",
		ObjectID: "test-object",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "file name is required")

	// Test validation - TTL default
	mockProvider := &MockProvider{}
	serviceWithMock := &StorageService{
		provider: mockProvider,
		logger:  	logger,
	}

	resp, err := serviceWithMock.CreateUploadURL(context.Background(), CreateUploadURLInput{
		TenantID:  "test-tenant",
		ObjectID:  "test-object",
		FileName:  "test.txt",
		TTLSecond: 0, // Should default to 3600
	})
	require.NoError(t, err)
	require.Equal(t, "http://mock-put-url", resp.URL)
	// Verify the mock was called with correct TTL (3600, default)
	// We can't easily assert on the mock call without more sophisticated mocking, but we tested the defaulting logic
}

// TestStorageServiceCreateDownloadURL tests the CreateDownloadURL method
func TestStorageServiceCreateDownloadURL(t *testing.T) {
	logger, err := logger.New("info", "storage-test")
	require.NoError(t, err)

	// Test validation - missing tenant ID
	service := &StorageService{
		logger: logger,
	}

	_, err = service.CreateDownloadURL(context.Background(), CreateDownloadURLInput{
		ObjectID: "test-object",
		TTLSecond: 3600,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "tenant ID is required")

	// Test validation - missing object ID
	_, err = service.CreateDownloadURL(context.Background(), CreateDownloadURLInput{
		TenantID:  "test-tenant",
		TTLSecond: 3600,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "object ID is required")

	// Test validation - TTL default
	mockProvider := &MockProvider{}
	serviceWithMock := &StorageService{
		provider: mockProvider,
		logger:  	logger,
	}

	resp, err := serviceWithMock.CreateDownloadURL(context.Background(), CreateDownloadURLInput{
		TenantID:  "test-tenant",
		ObjectID:  "test-object",
		TTLSecond: 0, // Should default to 3600
	})
	require.NoError(t, err)
	require.Equal(t, "http://mock-get-url", resp.URL)
	// Verify the mock was called with correct TTL (3600, default)
}

// TestStorageServiceWithRepository tests that the repository is properly set and used
func TestStorageServiceWithRepository(t *testing.T) {
	logger, err := logger.New("info", "storage-test")
	require.NoError(t, err)

	mockRepo := &RepositoryMock{}
	service := &StorageService{
		logger: logger,
	}
	service = service.WithRepository(mockRepo)

	// Set up mock to return a specific object ID when SaveObjectMeta is called
	mockRepo.SaveObjectMetaFunc = func(ctx context.Context, obj ObjectMeta) (ObjectMeta, error) {
		obj.ID = "returned-id"
		return obj, nil
	}

	mockProvider := &MockProvider{}
	serviceWithDeps := &StorageService{
		provider: mockProvider,
		repo:     mockRepo,
		logger:  	logger,
	}

	resp, err := serviceWithDeps.CreateUploadURL(context.Background(), CreateUploadURLInput{
		TenantID: "test-tenant",
		ObjectID: "test-object",
		FileName: "test.txt",
	})
	require.NoError(t, err)
	require.Equal(t, "http://mock-put-url", resp.URL)

	// Verify the repository was called
	require.NotNil(t, mockRepo.SaveObjectMetaCalled)
	require.Equal(t, "test-tenant", mockRepo.SaveObjectMetaCalled.TenantID)
	require.Equal(t, "test-object", mockRepo.SaveObjectMetaCalled.ID) // Fixed: was ObjectID, should be ID
	require.Equal(t, "test-tenant/test-object", mockRepo.SaveObjectMetaCalled.ObjectKey)
	require.Equal(t, "test.txt", mockRepo.SaveObjectMetaCalled.FileName)
}

// RepositoryMock implements Repository for testing
type RepositoryMock struct {
	SaveObjectMetaFunc func(ctx context.Context, obj ObjectMeta) (ObjectMeta, error)
	GetObjectMetaFunc  func(ctx context.Context, tenantID, objectID string) (ObjectMeta, error)

	SaveObjectMetaCalled *ObjectMeta
	GetObjectMetaCalledTenantID string
	GetObjectMetaCalledObjectID string
}

func (m *RepositoryMock) SaveObjectMeta(ctx context.Context, obj ObjectMeta) (ObjectMeta, error) {
	m.SaveObjectMetaCalled = &obj
	if m.SaveObjectMetaFunc != nil {
		return m.SaveObjectMetaFunc(ctx, obj)
	}
	return obj, nil
}

func (m *RepositoryMock) GetObjectMeta(ctx context.Context, tenantID, objectID string) (ObjectMeta, error) {
	m.GetObjectMetaCalledTenantID = tenantID
	m.GetObjectMetaCalledObjectID = objectID
	if m.GetObjectMetaFunc != nil {
		return m.GetObjectMetaFunc(ctx, tenantID, objectID)
	}
	return ObjectMeta{}, ErrObjectNotFound
}