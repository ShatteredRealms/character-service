// Code generated by MockGen. DO NOT EDIT.
// Source: /home/wil/sro/git/character-service/pkg/pb/character_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -source=/home/wil/sro/git/character-service/pkg/pb/character_grpc.pb.go -destination=/home/wil/sro/git/character-service/pkg/pb/mocks/character_grpc.pb.go
//

// Package mock_pb is a generated GoMock package.
package mock_pb

import (
	context "context"
	reflect "reflect"

	pb "github.com/ShatteredRealms/character-service/pkg/pb"
	pb0 "github.com/ShatteredRealms/go-common-service/pkg/pb"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockCharacterServiceClient is a mock of CharacterServiceClient interface.
type MockCharacterServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockCharacterServiceClientMockRecorder
	isgomock struct{}
}

// MockCharacterServiceClientMockRecorder is the mock recorder for MockCharacterServiceClient.
type MockCharacterServiceClientMockRecorder struct {
	mock *MockCharacterServiceClient
}

// NewMockCharacterServiceClient creates a new mock instance.
func NewMockCharacterServiceClient(ctrl *gomock.Controller) *MockCharacterServiceClient {
	mock := &MockCharacterServiceClient{ctrl: ctrl}
	mock.recorder = &MockCharacterServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCharacterServiceClient) EXPECT() *MockCharacterServiceClientMockRecorder {
	return m.recorder
}

// AddCharacterPlayTime mocks base method.
func (m *MockCharacterServiceClient) AddCharacterPlayTime(ctx context.Context, in *pb.AddPlayTimeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddCharacterPlayTime", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCharacterPlayTime indicates an expected call of AddCharacterPlayTime.
func (mr *MockCharacterServiceClientMockRecorder) AddCharacterPlayTime(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCharacterPlayTime", reflect.TypeOf((*MockCharacterServiceClient)(nil).AddCharacterPlayTime), varargs...)
}

// CreateCharacter mocks base method.
func (m *MockCharacterServiceClient) CreateCharacter(ctx context.Context, in *pb.CreateCharacterRequest, opts ...grpc.CallOption) (*pb.Character, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateCharacter", varargs...)
	ret0, _ := ret[0].(*pb.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCharacter indicates an expected call of CreateCharacter.
func (mr *MockCharacterServiceClientMockRecorder) CreateCharacter(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCharacter", reflect.TypeOf((*MockCharacterServiceClient)(nil).CreateCharacter), varargs...)
}

// DeleteCharacter mocks base method.
func (m *MockCharacterServiceClient) DeleteCharacter(ctx context.Context, in *pb0.TargetId, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteCharacter", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteCharacter indicates an expected call of DeleteCharacter.
func (mr *MockCharacterServiceClientMockRecorder) DeleteCharacter(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCharacter", reflect.TypeOf((*MockCharacterServiceClient)(nil).DeleteCharacter), varargs...)
}

// EditCharacter mocks base method.
func (m *MockCharacterServiceClient) EditCharacter(ctx context.Context, in *pb.EditCharacterRequest, opts ...grpc.CallOption) (*pb.Character, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EditCharacter", varargs...)
	ret0, _ := ret[0].(*pb.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditCharacter indicates an expected call of EditCharacter.
func (mr *MockCharacterServiceClientMockRecorder) EditCharacter(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditCharacter", reflect.TypeOf((*MockCharacterServiceClient)(nil).EditCharacter), varargs...)
}

// GetCharacter mocks base method.
func (m *MockCharacterServiceClient) GetCharacter(ctx context.Context, in *pb.GetCharacterRequest, opts ...grpc.CallOption) (*pb.Character, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetCharacter", varargs...)
	ret0, _ := ret[0].(*pb.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCharacter indicates an expected call of GetCharacter.
func (mr *MockCharacterServiceClientMockRecorder) GetCharacter(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharacter", reflect.TypeOf((*MockCharacterServiceClient)(nil).GetCharacter), varargs...)
}

// GetCharacters mocks base method.
func (m *MockCharacterServiceClient) GetCharacters(ctx context.Context, in *pb.GetCharactersRequest, opts ...grpc.CallOption) (*pb.Characters, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetCharacters", varargs...)
	ret0, _ := ret[0].(*pb.Characters)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCharacters indicates an expected call of GetCharacters.
func (mr *MockCharacterServiceClientMockRecorder) GetCharacters(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharacters", reflect.TypeOf((*MockCharacterServiceClient)(nil).GetCharacters), varargs...)
}

// GetCharactersForUser mocks base method.
func (m *MockCharacterServiceClient) GetCharactersForUser(ctx context.Context, in *pb.GetUserCharactersRequest, opts ...grpc.CallOption) (*pb.Characters, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetCharactersForUser", varargs...)
	ret0, _ := ret[0].(*pb.Characters)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCharactersForUser indicates an expected call of GetCharactersForUser.
func (mr *MockCharacterServiceClientMockRecorder) GetCharactersForUser(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharactersForUser", reflect.TypeOf((*MockCharacterServiceClient)(nil).GetCharactersForUser), varargs...)
}

// MockCharacterServiceServer is a mock of CharacterServiceServer interface.
type MockCharacterServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockCharacterServiceServerMockRecorder
	isgomock struct{}
}

// MockCharacterServiceServerMockRecorder is the mock recorder for MockCharacterServiceServer.
type MockCharacterServiceServerMockRecorder struct {
	mock *MockCharacterServiceServer
}

// NewMockCharacterServiceServer creates a new mock instance.
func NewMockCharacterServiceServer(ctrl *gomock.Controller) *MockCharacterServiceServer {
	mock := &MockCharacterServiceServer{ctrl: ctrl}
	mock.recorder = &MockCharacterServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCharacterServiceServer) EXPECT() *MockCharacterServiceServerMockRecorder {
	return m.recorder
}

// AddCharacterPlayTime mocks base method.
func (m *MockCharacterServiceServer) AddCharacterPlayTime(arg0 context.Context, arg1 *pb.AddPlayTimeRequest) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCharacterPlayTime", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCharacterPlayTime indicates an expected call of AddCharacterPlayTime.
func (mr *MockCharacterServiceServerMockRecorder) AddCharacterPlayTime(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCharacterPlayTime", reflect.TypeOf((*MockCharacterServiceServer)(nil).AddCharacterPlayTime), arg0, arg1)
}

// CreateCharacter mocks base method.
func (m *MockCharacterServiceServer) CreateCharacter(arg0 context.Context, arg1 *pb.CreateCharacterRequest) (*pb.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCharacter", arg0, arg1)
	ret0, _ := ret[0].(*pb.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCharacter indicates an expected call of CreateCharacter.
func (mr *MockCharacterServiceServerMockRecorder) CreateCharacter(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCharacter", reflect.TypeOf((*MockCharacterServiceServer)(nil).CreateCharacter), arg0, arg1)
}

// DeleteCharacter mocks base method.
func (m *MockCharacterServiceServer) DeleteCharacter(arg0 context.Context, arg1 *pb0.TargetId) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCharacter", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteCharacter indicates an expected call of DeleteCharacter.
func (mr *MockCharacterServiceServerMockRecorder) DeleteCharacter(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCharacter", reflect.TypeOf((*MockCharacterServiceServer)(nil).DeleteCharacter), arg0, arg1)
}

// EditCharacter mocks base method.
func (m *MockCharacterServiceServer) EditCharacter(arg0 context.Context, arg1 *pb.EditCharacterRequest) (*pb.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditCharacter", arg0, arg1)
	ret0, _ := ret[0].(*pb.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditCharacter indicates an expected call of EditCharacter.
func (mr *MockCharacterServiceServerMockRecorder) EditCharacter(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditCharacter", reflect.TypeOf((*MockCharacterServiceServer)(nil).EditCharacter), arg0, arg1)
}

// GetCharacter mocks base method.
func (m *MockCharacterServiceServer) GetCharacter(arg0 context.Context, arg1 *pb.GetCharacterRequest) (*pb.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCharacter", arg0, arg1)
	ret0, _ := ret[0].(*pb.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCharacter indicates an expected call of GetCharacter.
func (mr *MockCharacterServiceServerMockRecorder) GetCharacter(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharacter", reflect.TypeOf((*MockCharacterServiceServer)(nil).GetCharacter), arg0, arg1)
}

// GetCharacters mocks base method.
func (m *MockCharacterServiceServer) GetCharacters(arg0 context.Context, arg1 *pb.GetCharactersRequest) (*pb.Characters, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCharacters", arg0, arg1)
	ret0, _ := ret[0].(*pb.Characters)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCharacters indicates an expected call of GetCharacters.
func (mr *MockCharacterServiceServerMockRecorder) GetCharacters(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharacters", reflect.TypeOf((*MockCharacterServiceServer)(nil).GetCharacters), arg0, arg1)
}

// GetCharactersForUser mocks base method.
func (m *MockCharacterServiceServer) GetCharactersForUser(arg0 context.Context, arg1 *pb.GetUserCharactersRequest) (*pb.Characters, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCharactersForUser", arg0, arg1)
	ret0, _ := ret[0].(*pb.Characters)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCharactersForUser indicates an expected call of GetCharactersForUser.
func (mr *MockCharacterServiceServerMockRecorder) GetCharactersForUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharactersForUser", reflect.TypeOf((*MockCharacterServiceServer)(nil).GetCharactersForUser), arg0, arg1)
}

// mustEmbedUnimplementedCharacterServiceServer mocks base method.
func (m *MockCharacterServiceServer) mustEmbedUnimplementedCharacterServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedCharacterServiceServer")
}

// mustEmbedUnimplementedCharacterServiceServer indicates an expected call of mustEmbedUnimplementedCharacterServiceServer.
func (mr *MockCharacterServiceServerMockRecorder) mustEmbedUnimplementedCharacterServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedCharacterServiceServer", reflect.TypeOf((*MockCharacterServiceServer)(nil).mustEmbedUnimplementedCharacterServiceServer))
}

// MockUnsafeCharacterServiceServer is a mock of UnsafeCharacterServiceServer interface.
type MockUnsafeCharacterServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeCharacterServiceServerMockRecorder
	isgomock struct{}
}

// MockUnsafeCharacterServiceServerMockRecorder is the mock recorder for MockUnsafeCharacterServiceServer.
type MockUnsafeCharacterServiceServerMockRecorder struct {
	mock *MockUnsafeCharacterServiceServer
}

// NewMockUnsafeCharacterServiceServer creates a new mock instance.
func NewMockUnsafeCharacterServiceServer(ctrl *gomock.Controller) *MockUnsafeCharacterServiceServer {
	mock := &MockUnsafeCharacterServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeCharacterServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeCharacterServiceServer) EXPECT() *MockUnsafeCharacterServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedCharacterServiceServer mocks base method.
func (m *MockUnsafeCharacterServiceServer) mustEmbedUnimplementedCharacterServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedCharacterServiceServer")
}

// mustEmbedUnimplementedCharacterServiceServer indicates an expected call of mustEmbedUnimplementedCharacterServiceServer.
func (mr *MockUnsafeCharacterServiceServerMockRecorder) mustEmbedUnimplementedCharacterServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedCharacterServiceServer", reflect.TypeOf((*MockUnsafeCharacterServiceServer)(nil).mustEmbedUnimplementedCharacterServiceServer))
}
