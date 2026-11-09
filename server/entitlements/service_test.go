package entitlements_test

import (
	"context"
	"testing"

	"github.com/abhinash-kml/nova/server/entitlements"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Create(t *testing.T) {
	t.Parallel()

	// 1. Define the test case structural schema
	type testCase struct {
		name          string
		request       entitlements.CreateRequest
		setupMock     func(m *MockRepository)
		expectedError bool
		expectedRes   entitlements.Entitlement
	}

	// 2. Define the matrix of scenarios
	tests := []testCase{
		{
			name:    "Success - Valid Entitlement",
			request: entitlements.CreateRequest{Source: "e", IsActive: true},
			setupMock: func(m *MockRepository) {
				m.EXPECT().
					Create(mock.Anything, entitlements.CreateRequest{Source: "e", IsActive: true}).
					Return(entitlements.Entitlement{Source: "e", IsActive: true}, nil).
					Once()
			},
			expectedError: false,
			expectedRes:   entitlements.Entitlement{Source: "e", IsActive: true},
		},
		{
			name:          "Failure - Inactive Entitlement",
			request:       entitlements.CreateRequest{Source: "e", IsActive: false},
			setupMock:     func(m *MockRepository) { /* No repo call expected */ },
			expectedError: true,
			expectedRes:   entitlements.Entitlement{},
		},
		{
			name:          "Failure - Empty source",
			request:       entitlements.CreateRequest{},
			setupMock:     func(m *MockRepository) { /* No repo call expected */ },
			expectedError: true,
			expectedRes:   entitlements.Entitlement{},
		},
	}

	// 3. Iterate through the table matrix seamlessly
	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange - Freshly isolated for each table loop iteration
			mockRepo := NewMockRepository(t)
			tc.setupMock(mockRepo)

			svc := entitlements.NewLocalEntitlementService(mockRepo, nil)

			// Test
			res, err := svc.Create(context.Background(), tc.request)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Zero(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			}
		})
	}
}

func Test_Delete(t *testing.T) {
	t.Parallel()

	// 1. Define the test case structural schema
	type testCase struct {
		name          string
		request       entitlements.DeleteRequest
		setupMock     func(m *MockRepository)
		expectedError bool
		expectedRes   entitlements.DeleteResponse
	}

	id := uuid.New()

	tests := []testCase{
		{
			name:    "Success - Valid request",
			request: entitlements.DeleteRequest{entitlements.EntitlementID{Id: id.String()}},
			setupMock: func(m *MockRepository) {
				m.EXPECT().Delete(mock.Anything, entitlements.DeleteRequest{entitlements.EntitlementID{Id: id.String()}}).
					Return(entitlements.DeleteResponse{Id: id}, nil).
					Once()
			},
			expectedError: false,
			expectedRes:   entitlements.DeleteResponse{Id: id},
		},
		{
			name:          "Failure - Invalid requirest",
			request:       entitlements.DeleteRequest{},
			setupMock:     func(m *MockRepository) { /* No db call needed */ },
			expectedError: true,
			expectedRes:   entitlements.DeleteResponse{},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			mockRepo := NewMockRepository(t)
			tc.setupMock(mockRepo)

			svc := entitlements.NewLocalEntitlementService(mockRepo, nil)

			// Test
			res, err := svc.Delete(context.Background(), tc.request)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Zero(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			}
		})
	}
}

func Test_CheckIfActive(t *testing.T) {
	t.Parallel()

	// 1. Define the test case structural schema
	type testCase struct {
		name          string
		request       entitlements.CheckIfActiveRequest
		setupMock     func(m *MockRepository)
		expectedError bool
		expectedRes   bool
	}

	id := uuid.New()

	tests := []testCase{
		{
			name:    "Success - Valid request",
			request: entitlements.CheckIfActiveRequest{entitlements.EntitlementID{Id: id.String()}},
			setupMock: func(m *MockRepository) {
				m.EXPECT().CheckIfActive(mock.Anything, entitlements.CheckIfActiveRequest{entitlements.EntitlementID{Id: id.String()}}).
					Return(true, nil).
					Once()
			},
			expectedError: false,
			expectedRes:   true,
		},
		{
			name:          "Failure - Invalid requirest",
			request:       entitlements.CheckIfActiveRequest{},
			setupMock:     func(m *MockRepository) { /* No db call needed */ },
			expectedError: true,
			expectedRes:   false,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			mockRepo := NewMockRepository(t)
			tc.setupMock(mockRepo)

			svc := entitlements.NewLocalEntitlementService(mockRepo, nil)

			// Test
			res, err := svc.CheckIfActive(context.Background(), tc.request)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Zero(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			}
		})
	}
}

func Test_GetAllOfUser(t *testing.T) {
	t.Parallel()

	// 1. Define the test case structural schema
	type testCase struct {
		name          string
		request       entitlements.GetAllOfUserRequest
		setupMock     func(m *MockRepository)
		expectedError bool
		expectedRes   []entitlements.Entitlement
	}

	id := uuid.New()

	tests := []testCase{
		{
			name:    "Success - Valid request",
			request: entitlements.GetAllOfUserRequest{UserID: id.String()},
			setupMock: func(m *MockRepository) {
				m.EXPECT().GetAllOfUser(mock.Anything, entitlements.GetAllOfUserRequest{UserID: id.String()}).
					Return(nil, nil).
					Once()
			},
			expectedError: false,
			expectedRes:   nil,
		},
		{
			name:          "Failure - Invalid request",
			request:       entitlements.GetAllOfUserRequest{},
			setupMock:     func(m *MockRepository) { /* No db call needed */ },
			expectedError: true,
			expectedRes:   nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			mockRepo := NewMockRepository(t)
			tc.setupMock(mockRepo)

			svc := entitlements.NewLocalEntitlementService(mockRepo, nil)

			// Test
			res, err := svc.GetAllOfUser(context.Background(), tc.request)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Zero(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			}
		})
	}
}
