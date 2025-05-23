// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.15.8
// source: sro/character/character.proto

package pb

import (
	pb "github.com/ShatteredRealms/go-common-service/pkg/pb"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetCharacterRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Mask          *fieldmaskpb.FieldMask `protobuf:"bytes,1,opt,name=mask,proto3" json:"mask,omitempty"`
	Id            string                 `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCharacterRequest) Reset() {
	*x = GetCharacterRequest{}
	mi := &file_sro_character_character_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCharacterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCharacterRequest) ProtoMessage() {}

func (x *GetCharacterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sro_character_character_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCharacterRequest.ProtoReflect.Descriptor instead.
func (*GetCharacterRequest) Descriptor() ([]byte, []int) {
	return file_sro_character_character_proto_rawDescGZIP(), []int{0}
}

func (x *GetCharacterRequest) GetMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.Mask
	}
	return nil
}

func (x *GetCharacterRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetCharactersRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Mask          *fieldmaskpb.FieldMask `protobuf:"bytes,1,opt,name=mask,proto3" json:"mask,omitempty"`
	Filters       *pb.QueryFilters       `protobuf:"bytes,2,opt,name=filters,proto3" json:"filters,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCharactersRequest) Reset() {
	*x = GetCharactersRequest{}
	mi := &file_sro_character_character_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCharactersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCharactersRequest) ProtoMessage() {}

func (x *GetCharactersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sro_character_character_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCharactersRequest.ProtoReflect.Descriptor instead.
func (*GetCharactersRequest) Descriptor() ([]byte, []int) {
	return file_sro_character_character_proto_rawDescGZIP(), []int{1}
}

func (x *GetCharactersRequest) GetMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.Mask
	}
	return nil
}

func (x *GetCharactersRequest) GetFilters() *pb.QueryFilters {
	if x != nil {
		return x.Filters
	}
	return nil
}

type GetUserCharactersRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Mask          *fieldmaskpb.FieldMask `protobuf:"bytes,1,opt,name=mask,proto3" json:"mask,omitempty"`
	Filters       *pb.QueryFilters       `protobuf:"bytes,2,opt,name=filters,proto3" json:"filters,omitempty"`
	OwnerId       string                 `protobuf:"bytes,3,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserCharactersRequest) Reset() {
	*x = GetUserCharactersRequest{}
	mi := &file_sro_character_character_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserCharactersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserCharactersRequest) ProtoMessage() {}

func (x *GetUserCharactersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sro_character_character_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserCharactersRequest.ProtoReflect.Descriptor instead.
func (*GetUserCharactersRequest) Descriptor() ([]byte, []int) {
	return file_sro_character_character_proto_rawDescGZIP(), []int{2}
}

func (x *GetUserCharactersRequest) GetMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.Mask
	}
	return nil
}

func (x *GetUserCharactersRequest) GetFilters() *pb.QueryFilters {
	if x != nil {
		return x.Filters
	}
	return nil
}

func (x *GetUserCharactersRequest) GetOwnerId() string {
	if x != nil {
		return x.OwnerId
	}
	return ""
}

type EditCharacterRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Mask          *fieldmaskpb.FieldMask `protobuf:"bytes,1,opt,name=mask,proto3" json:"mask,omitempty"`
	Character     *Character             `protobuf:"bytes,2,opt,name=character,proto3" json:"character,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EditCharacterRequest) Reset() {
	*x = EditCharacterRequest{}
	mi := &file_sro_character_character_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EditCharacterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EditCharacterRequest) ProtoMessage() {}

func (x *EditCharacterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sro_character_character_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EditCharacterRequest.ProtoReflect.Descriptor instead.
func (*EditCharacterRequest) Descriptor() ([]byte, []int) {
	return file_sro_character_character_proto_rawDescGZIP(), []int{3}
}

func (x *EditCharacterRequest) GetMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.Mask
	}
	return nil
}

func (x *EditCharacterRequest) GetCharacter() *Character {
	if x != nil {
		return x.Character
	}
	return nil
}

type CreateCharacterRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OwnerId       string                 `protobuf:"bytes,1,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Gender        string                 `protobuf:"bytes,3,opt,name=gender,proto3" json:"gender,omitempty"`
	Realm         string                 `protobuf:"bytes,4,opt,name=realm,proto3" json:"realm,omitempty"`
	DimensionId   string                 `protobuf:"bytes,5,opt,name=dimension_id,json=dimensionId,proto3" json:"dimension_id,omitempty"`
	Profession    string                 `protobuf:"bytes,6,opt,name=profession,proto3" json:"profession,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateCharacterRequest) Reset() {
	*x = CreateCharacterRequest{}
	mi := &file_sro_character_character_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateCharacterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCharacterRequest) ProtoMessage() {}

func (x *CreateCharacterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sro_character_character_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCharacterRequest.ProtoReflect.Descriptor instead.
func (*CreateCharacterRequest) Descriptor() ([]byte, []int) {
	return file_sro_character_character_proto_rawDescGZIP(), []int{4}
}

func (x *CreateCharacterRequest) GetOwnerId() string {
	if x != nil {
		return x.OwnerId
	}
	return ""
}

func (x *CreateCharacterRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateCharacterRequest) GetGender() string {
	if x != nil {
		return x.Gender
	}
	return ""
}

func (x *CreateCharacterRequest) GetRealm() string {
	if x != nil {
		return x.Realm
	}
	return ""
}

func (x *CreateCharacterRequest) GetDimensionId() string {
	if x != nil {
		return x.DimensionId
	}
	return ""
}

func (x *CreateCharacterRequest) GetProfession() string {
	if x != nil {
		return x.Profession
	}
	return ""
}

type Character struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	OwnerId       string                 `protobuf:"bytes,2,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Gender        string                 `protobuf:"bytes,4,opt,name=gender,proto3" json:"gender,omitempty"`
	Realm         string                 `protobuf:"bytes,5,opt,name=realm,proto3" json:"realm,omitempty"`
	PlayTime      int32                  `protobuf:"varint,6,opt,name=play_time,json=playTime,proto3" json:"play_time,omitempty"`
	Location      *pb.Location           `protobuf:"bytes,7,opt,name=location,proto3" json:"location,omitempty"`
	DimensionId   string                 `protobuf:"bytes,8,opt,name=dimension_id,json=dimensionId,proto3" json:"dimension_id,omitempty"`
	CreatedAt     int64                  `protobuf:"varint,9,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     int64                  `protobuf:"varint,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt     int64                  `protobuf:"varint,11,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	Profession    string                 `protobuf:"bytes,12,opt,name=profession,proto3" json:"profession,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Character) Reset() {
	*x = Character{}
	mi := &file_sro_character_character_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Character) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Character) ProtoMessage() {}

func (x *Character) ProtoReflect() protoreflect.Message {
	mi := &file_sro_character_character_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Character.ProtoReflect.Descriptor instead.
func (*Character) Descriptor() ([]byte, []int) {
	return file_sro_character_character_proto_rawDescGZIP(), []int{5}
}

func (x *Character) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Character) GetOwnerId() string {
	if x != nil {
		return x.OwnerId
	}
	return ""
}

func (x *Character) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Character) GetGender() string {
	if x != nil {
		return x.Gender
	}
	return ""
}

func (x *Character) GetRealm() string {
	if x != nil {
		return x.Realm
	}
	return ""
}

func (x *Character) GetPlayTime() int32 {
	if x != nil {
		return x.PlayTime
	}
	return 0
}

func (x *Character) GetLocation() *pb.Location {
	if x != nil {
		return x.Location
	}
	return nil
}

func (x *Character) GetDimensionId() string {
	if x != nil {
		return x.DimensionId
	}
	return ""
}

func (x *Character) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Character) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *Character) GetDeletedAt() int64 {
	if x != nil {
		return x.DeletedAt
	}
	return 0
}

func (x *Character) GetProfession() string {
	if x != nil {
		return x.Profession
	}
	return ""
}

type Characters struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Total         int64                  `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Characters    []*Character           `protobuf:"bytes,2,rep,name=characters,proto3" json:"characters,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Characters) Reset() {
	*x = Characters{}
	mi := &file_sro_character_character_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Characters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Characters) ProtoMessage() {}

func (x *Characters) ProtoReflect() protoreflect.Message {
	mi := &file_sro_character_character_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Characters.ProtoReflect.Descriptor instead.
func (*Characters) Descriptor() ([]byte, []int) {
	return file_sro_character_character_proto_rawDescGZIP(), []int{6}
}

func (x *Characters) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *Characters) GetCharacters() []*Character {
	if x != nil {
		return x.Characters
	}
	return nil
}

type AddPlayTimeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Time          int32                  `protobuf:"varint,2,opt,name=time,proto3" json:"time,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddPlayTimeRequest) Reset() {
	*x = AddPlayTimeRequest{}
	mi := &file_sro_character_character_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddPlayTimeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddPlayTimeRequest) ProtoMessage() {}

func (x *AddPlayTimeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sro_character_character_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddPlayTimeRequest.ProtoReflect.Descriptor instead.
func (*AddPlayTimeRequest) Descriptor() ([]byte, []int) {
	return file_sro_character_character_proto_rawDescGZIP(), []int{7}
}

func (x *AddPlayTimeRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AddPlayTimeRequest) GetTime() int32 {
	if x != nil {
		return x.Time
	}
	return 0
}

var File_sro_character_character_proto protoreflect.FileDescriptor

const file_sro_character_character_proto_rawDesc = "" +
	"\n" +
	"\x1dsro/character/character.proto\x12\rsro.character\x1a\x1cgoogle/api/annotations.proto\x1a\x1bgoogle/protobuf/empty.proto\x1a google/protobuf/field_mask.proto\x1a\x11sro/globals.proto\x1a\x10sro/filter.proto\"U\n" +
	"\x13GetCharacterRequest\x12.\n" +
	"\x04mask\x18\x01 \x01(\v2\x1a.google.protobuf.FieldMaskR\x04mask\x12\x0e\n" +
	"\x02id\x18\x02 \x01(\tR\x02id\"s\n" +
	"\x14GetCharactersRequest\x12.\n" +
	"\x04mask\x18\x01 \x01(\v2\x1a.google.protobuf.FieldMaskR\x04mask\x12+\n" +
	"\afilters\x18\x02 \x01(\v2\x11.sro.QueryFiltersR\afilters\"\x92\x01\n" +
	"\x18GetUserCharactersRequest\x12.\n" +
	"\x04mask\x18\x01 \x01(\v2\x1a.google.protobuf.FieldMaskR\x04mask\x12+\n" +
	"\afilters\x18\x02 \x01(\v2\x11.sro.QueryFiltersR\afilters\x12\x19\n" +
	"\bowner_id\x18\x03 \x01(\tR\aownerId\"~\n" +
	"\x14EditCharacterRequest\x12.\n" +
	"\x04mask\x18\x01 \x01(\v2\x1a.google.protobuf.FieldMaskR\x04mask\x126\n" +
	"\tcharacter\x18\x02 \x01(\v2\x18.sro.character.CharacterR\tcharacter\"\xb8\x01\n" +
	"\x16CreateCharacterRequest\x12\x19\n" +
	"\bowner_id\x18\x01 \x01(\tR\aownerId\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x16\n" +
	"\x06gender\x18\x03 \x01(\tR\x06gender\x12\x14\n" +
	"\x05realm\x18\x04 \x01(\tR\x05realm\x12!\n" +
	"\fdimension_id\x18\x05 \x01(\tR\vdimensionId\x12\x1e\n" +
	"\n" +
	"profession\x18\x06 \x01(\tR\n" +
	"profession\"\xe0\x02\n" +
	"\tCharacter\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x19\n" +
	"\bowner_id\x18\x02 \x01(\tR\aownerId\x12\x12\n" +
	"\x04name\x18\x03 \x01(\tR\x04name\x12\x16\n" +
	"\x06gender\x18\x04 \x01(\tR\x06gender\x12\x14\n" +
	"\x05realm\x18\x05 \x01(\tR\x05realm\x12\x1b\n" +
	"\tplay_time\x18\x06 \x01(\x05R\bplayTime\x12)\n" +
	"\blocation\x18\a \x01(\v2\r.sro.LocationR\blocation\x12!\n" +
	"\fdimension_id\x18\b \x01(\tR\vdimensionId\x12\x1d\n" +
	"\n" +
	"created_at\x18\t \x01(\x03R\tcreatedAt\x12\x1d\n" +
	"\n" +
	"updated_at\x18\n" +
	" \x01(\x03R\tupdatedAt\x12\x1d\n" +
	"\n" +
	"deleted_at\x18\v \x01(\x03R\tdeletedAt\x12\x1e\n" +
	"\n" +
	"profession\x18\f \x01(\tR\n" +
	"profession\"\\\n" +
	"\n" +
	"Characters\x12\x14\n" +
	"\x05total\x18\x01 \x01(\x03R\x05total\x128\n" +
	"\n" +
	"characters\x18\x02 \x03(\v2\x18.sro.character.CharacterR\n" +
	"characters\"8\n" +
	"\x12AddPlayTimeRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04time\x18\x02 \x01(\x05R\x04time2\xb8\x06\n" +
	"\x10CharacterService\x12l\n" +
	"\fGetCharacter\x12\".sro.character.GetCharacterRequest\x1a\x18.sro.character.Character\"\x1e\x82\xd3\xe4\x93\x02\x18\x12\x16/v1/characters/id/{id}\x12g\n" +
	"\rGetCharacters\x12#.sro.character.GetCharactersRequest\x1a\x19.sro.character.Characters\"\x16\x82\xd3\xe4\x93\x02\x10\x12\x0e/v1/characters\x12\x83\x01\n" +
	"\x14GetCharactersForUser\x12'.sro.character.GetUserCharactersRequest\x1a\x19.sro.character.Characters\"'\x82\xd3\xe4\x93\x02!\x12\x1f/v1/characters/owner/{owner_id}\x12m\n" +
	"\x0fCreateCharacter\x12%.sro.character.CreateCharacterRequest\x1a\x18.sro.character.Character\"\x19\x82\xd3\xe4\x93\x02\x13:\x01*\"\x0e/v1/characters\x12X\n" +
	"\x0fDeleteCharacter\x12\r.sro.TargetId\x1a\x16.google.protobuf.Empty\"\x1e\x82\xd3\xe4\x93\x02\x18*\x16/v1/characters/id/{id}\x12{\n" +
	"\rEditCharacter\x12#.sro.character.EditCharacterRequest\x1a\x18.sro.character.Character\"+\x82\xd3\xe4\x93\x02%:\x01*2 /v1/characters/id/{character.id}\x12\x80\x01\n" +
	"\x14AddCharacterPlayTime\x12!.sro.character.AddPlayTimeRequest\x1a\x16.google.protobuf.Empty\"-\x82\xd3\xe4\x93\x02':\x04time2\x1f/v1/characters/id/{id}/playtimeB8Z6github.com/ShatteredRealms/character-service/pkg/pb;pbb\x06proto3"

var (
	file_sro_character_character_proto_rawDescOnce sync.Once
	file_sro_character_character_proto_rawDescData []byte
)

func file_sro_character_character_proto_rawDescGZIP() []byte {
	file_sro_character_character_proto_rawDescOnce.Do(func() {
		file_sro_character_character_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_sro_character_character_proto_rawDesc), len(file_sro_character_character_proto_rawDesc)))
	})
	return file_sro_character_character_proto_rawDescData
}

var file_sro_character_character_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_sro_character_character_proto_goTypes = []any{
	(*GetCharacterRequest)(nil),      // 0: sro.character.GetCharacterRequest
	(*GetCharactersRequest)(nil),     // 1: sro.character.GetCharactersRequest
	(*GetUserCharactersRequest)(nil), // 2: sro.character.GetUserCharactersRequest
	(*EditCharacterRequest)(nil),     // 3: sro.character.EditCharacterRequest
	(*CreateCharacterRequest)(nil),   // 4: sro.character.CreateCharacterRequest
	(*Character)(nil),                // 5: sro.character.Character
	(*Characters)(nil),               // 6: sro.character.Characters
	(*AddPlayTimeRequest)(nil),       // 7: sro.character.AddPlayTimeRequest
	(*fieldmaskpb.FieldMask)(nil),    // 8: google.protobuf.FieldMask
	(*pb.QueryFilters)(nil),          // 9: sro.QueryFilters
	(*pb.Location)(nil),              // 10: sro.Location
	(*pb.TargetId)(nil),              // 11: sro.TargetId
	(*emptypb.Empty)(nil),            // 12: google.protobuf.Empty
}
var file_sro_character_character_proto_depIdxs = []int32{
	8,  // 0: sro.character.GetCharacterRequest.mask:type_name -> google.protobuf.FieldMask
	8,  // 1: sro.character.GetCharactersRequest.mask:type_name -> google.protobuf.FieldMask
	9,  // 2: sro.character.GetCharactersRequest.filters:type_name -> sro.QueryFilters
	8,  // 3: sro.character.GetUserCharactersRequest.mask:type_name -> google.protobuf.FieldMask
	9,  // 4: sro.character.GetUserCharactersRequest.filters:type_name -> sro.QueryFilters
	8,  // 5: sro.character.EditCharacterRequest.mask:type_name -> google.protobuf.FieldMask
	5,  // 6: sro.character.EditCharacterRequest.character:type_name -> sro.character.Character
	10, // 7: sro.character.Character.location:type_name -> sro.Location
	5,  // 8: sro.character.Characters.characters:type_name -> sro.character.Character
	0,  // 9: sro.character.CharacterService.GetCharacter:input_type -> sro.character.GetCharacterRequest
	1,  // 10: sro.character.CharacterService.GetCharacters:input_type -> sro.character.GetCharactersRequest
	2,  // 11: sro.character.CharacterService.GetCharactersForUser:input_type -> sro.character.GetUserCharactersRequest
	4,  // 12: sro.character.CharacterService.CreateCharacter:input_type -> sro.character.CreateCharacterRequest
	11, // 13: sro.character.CharacterService.DeleteCharacter:input_type -> sro.TargetId
	3,  // 14: sro.character.CharacterService.EditCharacter:input_type -> sro.character.EditCharacterRequest
	7,  // 15: sro.character.CharacterService.AddCharacterPlayTime:input_type -> sro.character.AddPlayTimeRequest
	5,  // 16: sro.character.CharacterService.GetCharacter:output_type -> sro.character.Character
	6,  // 17: sro.character.CharacterService.GetCharacters:output_type -> sro.character.Characters
	6,  // 18: sro.character.CharacterService.GetCharactersForUser:output_type -> sro.character.Characters
	5,  // 19: sro.character.CharacterService.CreateCharacter:output_type -> sro.character.Character
	12, // 20: sro.character.CharacterService.DeleteCharacter:output_type -> google.protobuf.Empty
	5,  // 21: sro.character.CharacterService.EditCharacter:output_type -> sro.character.Character
	12, // 22: sro.character.CharacterService.AddCharacterPlayTime:output_type -> google.protobuf.Empty
	16, // [16:23] is the sub-list for method output_type
	9,  // [9:16] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_sro_character_character_proto_init() }
func file_sro_character_character_proto_init() {
	if File_sro_character_character_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_sro_character_character_proto_rawDesc), len(file_sro_character_character_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sro_character_character_proto_goTypes,
		DependencyIndexes: file_sro_character_character_proto_depIdxs,
		MessageInfos:      file_sro_character_character_proto_msgTypes,
	}.Build()
	File_sro_character_character_proto = out.File
	file_sro_character_character_proto_goTypes = nil
	file_sro_character_character_proto_depIdxs = nil
}
