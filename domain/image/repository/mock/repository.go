// Code generated by MockGen. DO NOT EDIT.
// Source: domain/image/repository/repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	utility "github.com/Mitra-Apps/be-utility-service/domain/proto/utility"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockImageRepository is a mock of ImageRepository interface.
type MockImageRepository struct {
	ctrl     *gomock.Controller
	recorder *MockImageRepositoryMockRecorder
}

// MockImageRepositoryMockRecorder is the mock recorder for MockImageRepository.
type MockImageRepositoryMockRecorder struct {
	mock *MockImageRepository
}

// NewMockImageRepository creates a new mock instance.
func NewMockImageRepository(ctrl *gomock.Controller) *MockImageRepository {
	mock := &MockImageRepository{ctrl: ctrl}
	mock.recorder = &MockImageRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageRepository) EXPECT() *MockImageRepositoryMockRecorder {
	return m.recorder
}

// GetImagesByIds mocks base method.
func (m *MockImageRepository) GetImagesByIds(ctx context.Context, ids []string) ([]*utility.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImagesByIds", ctx, ids)
	ret0, _ := ret[0].([]*utility.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImagesByIds indicates an expected call of GetImagesByIds.
func (mr *MockImageRepositoryMockRecorder) GetImagesByIds(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImagesByIds", reflect.TypeOf((*MockImageRepository)(nil).GetImagesByIds), ctx, ids)
}

// RemoveImage mocks base method.
func (m *MockImageRepository) RemoveImage(ctx context.Context, ids []string, groupName, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveImage", ctx, ids, groupName, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveImage indicates an expected call of RemoveImage.
func (mr *MockImageRepositoryMockRecorder) RemoveImage(ctx, ids, groupName, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveImage", reflect.TypeOf((*MockImageRepository)(nil).RemoveImage), ctx, ids, groupName, userID)
}

// UploadImage mocks base method.
func (m *MockImageRepository) UploadImage(ctx context.Context, imageBase64Str, groupName, userID string) (*uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadImage", ctx, imageBase64Str, groupName, userID)
	ret0, _ := ret[0].(*uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadImage indicates an expected call of UploadImage.
func (mr *MockImageRepositoryMockRecorder) UploadImage(ctx, imageBase64Str, groupName, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadImage", reflect.TypeOf((*MockImageRepository)(nil).UploadImage), ctx, imageBase64Str, groupName, userID)
}
