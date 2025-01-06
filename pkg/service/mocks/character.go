// Code generated by MockGen. DO NOT EDIT.
// Source: /home/wil/dev/sro/character-service/pkg/service/character.go
//
// Generated by this command:
//
//	mockgen -source=/home/wil/dev/sro/character-service/pkg/service/character.go -destination=/home/wil/dev/sro/character-service/pkg/service/mocks/character.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	character "github.com/ShatteredRealms/character-service/pkg/model/character"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockCharacterService is a mock of CharacterService interface.
type MockCharacterService struct {
	ctrl     *gomock.Controller
	recorder *MockCharacterServiceMockRecorder
	isgomock struct{}
}

// MockCharacterServiceMockRecorder is the mock recorder for MockCharacterService.
type MockCharacterServiceMockRecorder struct {
	mock *MockCharacterService
}

// NewMockCharacterService creates a new mock instance.
func NewMockCharacterService(ctrl *gomock.Controller) *MockCharacterService {
	mock := &MockCharacterService{ctrl: ctrl}
	mock.recorder = &MockCharacterServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCharacterService) EXPECT() *MockCharacterServiceMockRecorder {
	return m.recorder
}

// AddCharacterPlaytime mocks base method.
func (m *MockCharacterService) AddCharacterPlaytime(ctx context.Context, char *character.Character, seconds int32) (*character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCharacterPlaytime", ctx, char, seconds)
	ret0, _ := ret[0].(*character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCharacterPlaytime indicates an expected call of AddCharacterPlaytime.
func (mr *MockCharacterServiceMockRecorder) AddCharacterPlaytime(ctx, char, seconds any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCharacterPlaytime", reflect.TypeOf((*MockCharacterService)(nil).AddCharacterPlaytime), ctx, char, seconds)
}

// CreateCharacter mocks base method.
func (m *MockCharacterService) CreateCharacter(ctx context.Context, ownerId, name, gender, realm string, dimensionId *uuid.UUID) (*character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCharacter", ctx, ownerId, name, gender, realm, dimensionId)
	ret0, _ := ret[0].(*character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCharacter indicates an expected call of CreateCharacter.
func (mr *MockCharacterServiceMockRecorder) CreateCharacter(ctx, ownerId, name, gender, realm, dimensionId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCharacter", reflect.TypeOf((*MockCharacterService)(nil).CreateCharacter), ctx, ownerId, name, gender, realm, dimensionId)
}

// DeleteCharacter mocks base method.
func (m *MockCharacterService) DeleteCharacter(ctx context.Context, characterId *uuid.UUID) (*character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCharacter", ctx, characterId)
	ret0, _ := ret[0].(*character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteCharacter indicates an expected call of DeleteCharacter.
func (mr *MockCharacterServiceMockRecorder) DeleteCharacter(ctx, characterId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCharacter", reflect.TypeOf((*MockCharacterService)(nil).DeleteCharacter), ctx, characterId)
}

// EditCharacter mocks base method.
func (m *MockCharacterService) EditCharacter(ctx context.Context, newCharacter *character.Character) (*character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditCharacter", ctx, newCharacter)
	ret0, _ := ret[0].(*character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditCharacter indicates an expected call of EditCharacter.
func (mr *MockCharacterServiceMockRecorder) EditCharacter(ctx, newCharacter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditCharacter", reflect.TypeOf((*MockCharacterService)(nil).EditCharacter), ctx, newCharacter)
}

// GetCharacterById mocks base method.
func (m *MockCharacterService) GetCharacterById(ctx context.Context, characterId *uuid.UUID) (*character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCharacterById", ctx, characterId)
	ret0, _ := ret[0].(*character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCharacterById indicates an expected call of GetCharacterById.
func (mr *MockCharacterServiceMockRecorder) GetCharacterById(ctx, characterId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharacterById", reflect.TypeOf((*MockCharacterService)(nil).GetCharacterById), ctx, characterId)
}

// GetCharacters mocks base method.
func (m *MockCharacterService) GetCharacters(ctx context.Context) (character.Characters, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCharacters", ctx)
	ret0, _ := ret[0].(character.Characters)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCharacters indicates an expected call of GetCharacters.
func (mr *MockCharacterServiceMockRecorder) GetCharacters(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharacters", reflect.TypeOf((*MockCharacterService)(nil).GetCharacters), ctx)
}

// GetCharactersByOwner mocks base method.
func (m *MockCharacterService) GetCharactersByOwner(ctx context.Context, ownerId string) (character.Characters, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCharactersByOwner", ctx, ownerId)
	ret0, _ := ret[0].(character.Characters)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCharactersByOwner indicates an expected call of GetCharactersByOwner.
func (mr *MockCharacterServiceMockRecorder) GetCharactersByOwner(ctx, ownerId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharactersByOwner", reflect.TypeOf((*MockCharacterService)(nil).GetCharactersByOwner), ctx, ownerId)
}

// GetDeletedCharacters mocks base method.
func (m *MockCharacterService) GetDeletedCharacters(ctx context.Context) (character.Characters, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeletedCharacters", ctx)
	ret0, _ := ret[0].(character.Characters)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeletedCharacters indicates an expected call of GetDeletedCharacters.
func (mr *MockCharacterServiceMockRecorder) GetDeletedCharacters(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeletedCharacters", reflect.TypeOf((*MockCharacterService)(nil).GetDeletedCharacters), ctx)
}
