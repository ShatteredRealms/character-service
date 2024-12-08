// Code generated by MockGen. DO NOT EDIT.
// Source: /home/wil/dev/sro/character-service/pkg/repository/character_r.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=/home/wil/dev/sro/character-service/pkg/repository/character_r.go -destination=/home/wil/dev/sro/character-service/pkg/mocks/character_r_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	character "github.com/ShatteredRealms/character-service/pkg/model/character"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockCharacterRepository is a mock of CharacterRepository interface.
type MockCharacterRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCharacterRepositoryMockRecorder
	isgomock struct{}
}

// MockCharacterRepositoryMockRecorder is the mock recorder for MockCharacterRepository.
type MockCharacterRepositoryMockRecorder struct {
	mock *MockCharacterRepository
}

// NewMockCharacterRepository creates a new mock instance.
func NewMockCharacterRepository(ctrl *gomock.Controller) *MockCharacterRepository {
	mock := &MockCharacterRepository{ctrl: ctrl}
	mock.recorder = &MockCharacterRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCharacterRepository) EXPECT() *MockCharacterRepositoryMockRecorder {
	return m.recorder
}

// CreateCharacter mocks base method.
func (m *MockCharacterRepository) CreateCharacter(ctx context.Context, newCharacter *character.Character) (*character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCharacter", ctx, newCharacter)
	ret0, _ := ret[0].(*character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCharacter indicates an expected call of CreateCharacter.
func (mr *MockCharacterRepositoryMockRecorder) CreateCharacter(ctx, newCharacter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCharacter", reflect.TypeOf((*MockCharacterRepository)(nil).CreateCharacter), ctx, newCharacter)
}

// DeleteCharacter mocks base method.
func (m *MockCharacterRepository) DeleteCharacter(ctx context.Context, characterId *uuid.UUID) (*character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCharacter", ctx, characterId)
	ret0, _ := ret[0].(*character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteCharacter indicates an expected call of DeleteCharacter.
func (mr *MockCharacterRepositoryMockRecorder) DeleteCharacter(ctx, characterId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCharacter", reflect.TypeOf((*MockCharacterRepository)(nil).DeleteCharacter), ctx, characterId)
}

// DeleteCharactersByOwner mocks base method.
func (m *MockCharacterRepository) DeleteCharactersByOwner(ctx context.Context, ownerId string) (*character.Characters, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCharactersByOwner", ctx, ownerId)
	ret0, _ := ret[0].(*character.Characters)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteCharactersByOwner indicates an expected call of DeleteCharactersByOwner.
func (mr *MockCharacterRepositoryMockRecorder) DeleteCharactersByOwner(ctx, ownerId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCharactersByOwner", reflect.TypeOf((*MockCharacterRepository)(nil).DeleteCharactersByOwner), ctx, ownerId)
}

// GetCharacterById mocks base method.
func (m *MockCharacterRepository) GetCharacterById(ctx context.Context, characterId *uuid.UUID) (*character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCharacterById", ctx, characterId)
	ret0, _ := ret[0].(*character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCharacterById indicates an expected call of GetCharacterById.
func (mr *MockCharacterRepositoryMockRecorder) GetCharacterById(ctx, characterId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharacterById", reflect.TypeOf((*MockCharacterRepository)(nil).GetCharacterById), ctx, characterId)
}

// GetCharacters mocks base method.
func (m *MockCharacterRepository) GetCharacters(ctx context.Context) (*character.Characters, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCharacters", ctx)
	ret0, _ := ret[0].(*character.Characters)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCharacters indicates an expected call of GetCharacters.
func (mr *MockCharacterRepositoryMockRecorder) GetCharacters(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharacters", reflect.TypeOf((*MockCharacterRepository)(nil).GetCharacters), ctx)
}

// GetCharactersByOwner mocks base method.
func (m *MockCharacterRepository) GetCharactersByOwner(ctx context.Context, ownerId string) (*character.Characters, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCharactersByOwner", ctx, ownerId)
	ret0, _ := ret[0].(*character.Characters)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCharactersByOwner indicates an expected call of GetCharactersByOwner.
func (mr *MockCharacterRepositoryMockRecorder) GetCharactersByOwner(ctx, ownerId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharactersByOwner", reflect.TypeOf((*MockCharacterRepository)(nil).GetCharactersByOwner), ctx, ownerId)
}

// UpdateCharacter mocks base method.
func (m *MockCharacterRepository) UpdateCharacter(ctx context.Context, updatedCharacter *character.Character) (*character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCharacter", ctx, updatedCharacter)
	ret0, _ := ret[0].(*character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCharacter indicates an expected call of UpdateCharacter.
func (mr *MockCharacterRepositoryMockRecorder) UpdateCharacter(ctx, updatedCharacter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCharacter", reflect.TypeOf((*MockCharacterRepository)(nil).UpdateCharacter), ctx, updatedCharacter)
}
