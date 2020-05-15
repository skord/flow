// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: go/protocol/flow.proto

package protocol

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	protocol1 "go.gazette.dev/core/broker/protocol"
	protocol "go.gazette.dev/core/consumer/protocol"
	recoverylog "go.gazette.dev/core/consumer/recoverylog"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type DeriveTxnState int32

const (
	// IDLE indicates the transaction has not yet begun.
	// It transitions to EXTEND.
	DeriveTxnState_IDLE DeriveTxnState = 0
	// EXTEND extends the derive transaction with additional
	// source or derived collection documents.
	//
	// * The flow consumer sends any number of EXTEND TxnRequests,
	//   containing source collection documents.
	// * Concurrently, the derive worker responds with any number of
	//   EXTEND TxnResponses, each having documents to be added to
	//   the collection being derived.
	// * The flow consumer is responsible for publishing each derived
	//   document to the appropriate collection & partition.
	// * Note TxnRequest and TxnResponse EXTEND messages are _not_ 1:1.
	DeriveTxnState_EXTEND DeriveTxnState = 1
	// FLUSH indicates the transacton pipeline is to flush.
	//
	// * The flow consumer issues FLUSH when its consumer transaction begins to
	//   close.
	// * The derive worker responds with FLUSH to indicate that all source
	//   documents have been processed and all derived documents emitted.
	// * The flow consumer awaits the response FLUSH, while continuing to begin
	//   publish operations for all derived documents seen in the meantime.
	// * On seeing FLUSH, the flow consumer is assured it's sequenced all messages
	//   of the transaction, and can build its consumer.Checkpoint.
	DeriveTxnState_FLUSH DeriveTxnState = 2
	// PREPARE begins a commit of the transaction.
	//
	// * The Flow Consumer sends PREPARE with a consumer.Checkpoint.
	// * On receipt, the derive worker queues an atomic recoverylog.Recorder
	//   block that's conditioned on an (unresolved) "commit" future. Within
	//   this block underlying stores commits (SQLite COMMIT / writing RocksDB
	//   WriteBatch) are issued to persist all state changes of the transaction,
	//   along with the consumer.Checkpoint.
	// * The derive worker responds with PREPARE once all local commits have
	//   completed, and recoverylog writes have been queued (but not started,
	//   awaiting COMMIT).
	// * On receipt, the Flow Consumer arranges to invoke COMMIT on the completion
	//   of all outstanding journal writes -- this the OpFuture passed to the
	//   Store.StartCommit interface. It returns a future which will resolve only
	//   after reading COMMIT from this transaction -- the OpFuture returned by
	//   that interface.
	//
	// It's an error if a prior transaction is still running at the onset of
	// PREPARE. However at the completion of PREPARE, a new & concurrent
	// Transaction may begin, though it itself cannot PREPARE until this
	// Transaction fully completes.
	DeriveTxnState_PREPARE DeriveTxnState = 3
	// COMMIT commits the transaction by resolving the "commit" future created
	// during PREPARE, allowing the atomic commit block created in PREPARE
	// to flush to the recovery log. The derive worker responds with COMMIT
	// when the commit barrier has fully resolved.
	DeriveTxnState_COMMIT DeriveTxnState = 4
)

var DeriveTxnState_name = map[int32]string{
	0: "IDLE",
	1: "EXTEND",
	2: "FLUSH",
	3: "PREPARE",
	4: "COMMIT",
}

var DeriveTxnState_value = map[string]int32{
	"IDLE":    0,
	"EXTEND":  1,
	"FLUSH":   2,
	"PREPARE": 3,
	"COMMIT":  4,
}

func (x DeriveTxnState) String() string {
	return proto.EnumName(DeriveTxnState_name, int32(x))
}

func (DeriveTxnState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_9dd0e7f3bbc7f41f, []int{0}
}

// TxnRequest is the streamed message of a Transaction RPC.
type TxnRequest struct {
	State DeriveTxnState `protobuf:"varint,1,opt,name=state,proto3,enum=flow.DeriveTxnState" json:"state,omitempty"`
	// Collection from which source documents are drawn. Set iff state == EXTEND.
	ExtendSource string `protobuf:"bytes,2,opt,name=extend_source,json=extendSource,proto3" json:"extend_source,omitempty"`
	// Documents of the collection. Set iff state == EXTEND.
	ExtendDocuments [][]byte `protobuf:"bytes,3,rep,name=extend_documents,json=extendDocuments,proto3" json:"extend_documents,omitempty"`
	// Checkpoint to commit. Set iff state == PREPARE.
	PrepareCheckpoint    *protocol.Checkpoint `protobuf:"bytes,4,opt,name=prepare_checkpoint,json=prepareCheckpoint,proto3" json:"prepare_checkpoint,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *TxnRequest) Reset()         { *m = TxnRequest{} }
func (m *TxnRequest) String() string { return proto.CompactTextString(m) }
func (*TxnRequest) ProtoMessage()    {}
func (*TxnRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9dd0e7f3bbc7f41f, []int{0}
}
func (m *TxnRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxnRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxnRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxnRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxnRequest.Merge(m, src)
}
func (m *TxnRequest) XXX_Size() int {
	return m.ProtoSize()
}
func (m *TxnRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TxnRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TxnRequest proto.InternalMessageInfo

func (m *TxnRequest) GetState() DeriveTxnState {
	if m != nil {
		return m.State
	}
	return DeriveTxnState_IDLE
}

func (m *TxnRequest) GetExtendSource() string {
	if m != nil {
		return m.ExtendSource
	}
	return ""
}

func (m *TxnRequest) GetExtendDocuments() [][]byte {
	if m != nil {
		return m.ExtendDocuments
	}
	return nil
}

func (m *TxnRequest) GetPrepareCheckpoint() *protocol.Checkpoint {
	if m != nil {
		return m.PrepareCheckpoint
	}
	return nil
}

// TxnResponse is the streamed response message of a Transaction RPC.
type TxnResponse struct {
	State DeriveTxnState `protobuf:"varint,1,opt,name=state,proto3,enum=flow.DeriveTxnState" json:"state,omitempty"`
	// Documents derived from request documents. Set iff state == EXTEND.
	ExtendDocuments [][]byte `protobuf:"bytes,2,rep,name=extend_documents,json=extendDocuments,proto3" json:"extend_documents,omitempty"`
	// Logical partition labels of these documents in the derived collection.
	// Set iff state == EXTEND.
	ExtendLabels         *protocol1.LabelSet `protobuf:"bytes,3,opt,name=extend_labels,json=extendLabels,proto3" json:"extend_labels,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *TxnResponse) Reset()         { *m = TxnResponse{} }
func (m *TxnResponse) String() string { return proto.CompactTextString(m) }
func (*TxnResponse) ProtoMessage()    {}
func (*TxnResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9dd0e7f3bbc7f41f, []int{1}
}
func (m *TxnResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxnResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxnResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxnResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxnResponse.Merge(m, src)
}
func (m *TxnResponse) XXX_Size() int {
	return m.ProtoSize()
}
func (m *TxnResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TxnResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TxnResponse proto.InternalMessageInfo

func (m *TxnResponse) GetState() DeriveTxnState {
	if m != nil {
		return m.State
	}
	return DeriveTxnState_IDLE
}

func (m *TxnResponse) GetExtendDocuments() [][]byte {
	if m != nil {
		return m.ExtendDocuments
	}
	return nil
}

func (m *TxnResponse) GetExtendLabels() *protocol1.LabelSet {
	if m != nil {
		return m.ExtendLabels
	}
	return nil
}

func init() {
	proto.RegisterEnum("flow.DeriveTxnState", DeriveTxnState_name, DeriveTxnState_value)
	proto.RegisterType((*TxnRequest)(nil), "flow.TxnRequest")
	proto.RegisterType((*TxnResponse)(nil), "flow.TxnResponse")
}

func init() { proto.RegisterFile("go/protocol/flow.proto", fileDescriptor_9dd0e7f3bbc7f41f) }

var fileDescriptor_9dd0e7f3bbc7f41f = []byte{
	// 487 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0x41, 0x8b, 0xd3, 0x40,
	0x18, 0x86, 0x9d, 0xb6, 0x5b, 0x77, 0xa7, 0xeb, 0x9a, 0x0e, 0xb5, 0x94, 0x1e, 0x4a, 0x58, 0x41,
	0xe2, 0x1e, 0x52, 0xa9, 0xa0, 0x07, 0x4f, 0xbb, 0x6d, 0x96, 0x5d, 0x69, 0x75, 0x49, 0x22, 0x88,
	0x97, 0x92, 0x26, 0xdf, 0xc6, 0xb0, 0x69, 0xbe, 0x38, 0x33, 0x59, 0xdb, 0x1f, 0x23, 0xf8, 0x53,
	0x3c, 0x7a, 0x11, 0xfc, 0x0d, 0xeb, 0x1f, 0x91, 0x64, 0x9a, 0xb6, 0x42, 0x15, 0xbc, 0x94, 0x77,
	0xde, 0x79, 0xe8, 0x7c, 0xef, 0x9b, 0x8f, 0xb6, 0x43, 0xec, 0xa7, 0x1c, 0x25, 0xfa, 0x18, 0xf7,
	0xaf, 0x63, 0xfc, 0x6c, 0x16, 0x27, 0x56, 0xcb, 0x75, 0xb7, 0x37, 0xe3, 0x78, 0x03, 0x7c, 0x43,
	0x94, 0x42, 0x51, 0x5d, 0xdd, 0xc7, 0x44, 0x64, 0xf3, 0x7f, 0x10, 0x4f, 0xd6, 0x04, 0x07, 0x1f,
	0x6f, 0x81, 0x2f, 0x63, 0x0c, 0x0b, 0xcd, 0x03, 0x08, 0xa6, 0x98, 0xae, 0xb8, 0x4e, 0x2a, 0x97,
	0x29, 0x88, 0x3e, 0xcc, 0x53, 0xb9, 0x54, 0xbf, 0xab, 0x9b, 0x56, 0x88, 0x21, 0x16, 0xb2, 0x9f,
	0x2b, 0xe5, 0x1e, 0xff, 0x20, 0x94, 0xba, 0x8b, 0xc4, 0x86, 0x4f, 0x19, 0x08, 0xc9, 0x4e, 0xe8,
	0x9e, 0x90, 0x9e, 0x84, 0x0e, 0xd1, 0x89, 0x71, 0x34, 0x68, 0x99, 0x45, 0x94, 0x11, 0xf0, 0xe8,
	0x16, 0xdc, 0x45, 0xe2, 0xe4, 0x77, 0xb6, 0x42, 0xd8, 0x63, 0xfa, 0x00, 0x16, 0x12, 0x92, 0x60,
	0x2a, 0x30, 0xe3, 0x3e, 0x74, 0x2a, 0x3a, 0x31, 0x0e, 0xec, 0x43, 0x65, 0x3a, 0x85, 0xc7, 0x9e,
	0x52, 0x6d, 0x05, 0x05, 0xe8, 0x67, 0x73, 0x48, 0xa4, 0xe8, 0x54, 0xf5, 0xaa, 0x71, 0x68, 0x3f,
	0x54, 0xfe, 0xa8, 0xb4, 0xd9, 0x90, 0xb2, 0x94, 0x43, 0xea, 0x71, 0x98, 0xfa, 0x1f, 0xc1, 0xbf,
	0x49, 0x31, 0x4a, 0x64, 0xa7, 0xa6, 0x13, 0xa3, 0x31, 0x68, 0x99, 0x65, 0x7e, 0x73, 0xb8, 0xbe,
	0xb3, 0x9b, 0x2b, 0x7e, 0x63, 0x1d, 0x7f, 0x21, 0xb4, 0x51, 0xe4, 0x11, 0x29, 0x26, 0x02, 0xfe,
	0x2b, 0xd0, 0xae, 0x59, 0x2b, 0xbb, 0x67, 0x7d, 0xb9, 0xce, 0x1e, 0x7b, 0x33, 0x88, 0xf3, 0x4c,
	0xf9, 0x98, 0xcc, 0x5c, 0x7f, 0xb6, 0x71, 0xee, 0x3b, 0x20, 0xcb, 0x3e, 0x8a, 0xb3, 0x38, 0x79,
	0x4d, 0x8f, 0xfe, 0x7c, 0x9c, 0xed, 0xd3, 0xda, 0xe5, 0x68, 0x6c, 0x69, 0xf7, 0x18, 0xa5, 0x75,
	0xeb, 0xbd, 0x6b, 0xbd, 0x19, 0x69, 0x84, 0x1d, 0xd0, 0xbd, 0xf3, 0xf1, 0x3b, 0xe7, 0x42, 0xab,
	0xb0, 0x06, 0xbd, 0x7f, 0x65, 0x5b, 0x57, 0xa7, 0xb6, 0xa5, 0x55, 0x73, 0x66, 0xf8, 0x76, 0x32,
	0xb9, 0x74, 0xb5, 0xda, 0xe0, 0x1b, 0xa1, 0x75, 0xf5, 0x67, 0xec, 0x94, 0x36, 0x6d, 0x10, 0x12,
	0xb7, 0xbb, 0x60, 0x6d, 0x33, 0x44, 0x0c, 0x63, 0x50, 0x43, 0xcd, 0xb2, 0x6b, 0xd3, 0xca, 0xf7,
	0xa1, 0xbb, 0xb3, 0x4c, 0xf6, 0x82, 0x36, 0x5c, 0xee, 0x25, 0xc2, 0xf3, 0x65, 0x84, 0x09, 0xd3,
	0x54, 0x53, 0x9b, 0xdd, 0xe8, 0x36, 0xb7, 0x1c, 0xd5, 0xae, 0x41, 0x9e, 0x11, 0xf6, 0x8a, 0xd2,
	0xb3, 0x2c, 0x8a, 0x83, 0x8b, 0x28, 0x2f, 0xe6, 0x6f, 0x6f, 0x3e, 0x32, 0xb7, 0xf6, 0xd6, 0x3c,
	0x77, 0x26, 0x05, 0x7e, 0xd6, 0xfe, 0x7e, 0xd7, 0x23, 0x3f, 0xef, 0x7a, 0xe4, 0xeb, 0xaf, 0x1e,
	0xf9, 0xb0, 0x5f, 0xb6, 0x37, 0xab, 0x17, 0xea, 0xf9, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x3a,
	0xc8, 0xcb, 0x93, 0x57, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DeriveClient is the client API for Derive service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DeriveClient interface {
	// RestoreCheckpoint recovers the most recent Checkpoint previously committed
	// to the Store. It is called just once, at Shard start-up. If an external
	// system is used, it should install a transactional "write fence" to ensure
	// that an older Store instance of another process cannot successfully
	// StartCommit after this RestoreCheckpoint returns.
	RestoreCheckpoint(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*protocol.Checkpoint, error)
	// Transaction begins a pipelined derive-worker transaction, following the
	// state machine detailed in DeriveTxnState.
	Transaction(ctx context.Context, opts ...grpc.CallOption) (Derive_TransactionClient, error)
	// BuildHints returns FSMHints which may be played back to fully reconstruct
	// the local filesystem state produced by this derive worker. It may block
	// while pending operations sync to the recovery log.
	BuildHints(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*recoverylog.FSMHints, error)
}

type deriveClient struct {
	cc *grpc.ClientConn
}

func NewDeriveClient(cc *grpc.ClientConn) DeriveClient {
	return &deriveClient{cc}
}

func (c *deriveClient) RestoreCheckpoint(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*protocol.Checkpoint, error) {
	out := new(protocol.Checkpoint)
	err := c.cc.Invoke(ctx, "/flow.Derive/RestoreCheckpoint", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deriveClient) Transaction(ctx context.Context, opts ...grpc.CallOption) (Derive_TransactionClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Derive_serviceDesc.Streams[0], "/flow.Derive/Transaction", opts...)
	if err != nil {
		return nil, err
	}
	x := &deriveTransactionClient{stream}
	return x, nil
}

type Derive_TransactionClient interface {
	Send(*TxnRequest) error
	Recv() (*TxnResponse, error)
	grpc.ClientStream
}

type deriveTransactionClient struct {
	grpc.ClientStream
}

func (x *deriveTransactionClient) Send(m *TxnRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *deriveTransactionClient) Recv() (*TxnResponse, error) {
	m := new(TxnResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *deriveClient) BuildHints(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*recoverylog.FSMHints, error) {
	out := new(recoverylog.FSMHints)
	err := c.cc.Invoke(ctx, "/flow.Derive/BuildHints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeriveServer is the server API for Derive service.
type DeriveServer interface {
	// RestoreCheckpoint recovers the most recent Checkpoint previously committed
	// to the Store. It is called just once, at Shard start-up. If an external
	// system is used, it should install a transactional "write fence" to ensure
	// that an older Store instance of another process cannot successfully
	// StartCommit after this RestoreCheckpoint returns.
	RestoreCheckpoint(context.Context, *empty.Empty) (*protocol.Checkpoint, error)
	// Transaction begins a pipelined derive-worker transaction, following the
	// state machine detailed in DeriveTxnState.
	Transaction(Derive_TransactionServer) error
	// BuildHints returns FSMHints which may be played back to fully reconstruct
	// the local filesystem state produced by this derive worker. It may block
	// while pending operations sync to the recovery log.
	BuildHints(context.Context, *empty.Empty) (*recoverylog.FSMHints, error)
}

// UnimplementedDeriveServer can be embedded to have forward compatible implementations.
type UnimplementedDeriveServer struct {
}

func (*UnimplementedDeriveServer) RestoreCheckpoint(ctx context.Context, req *empty.Empty) (*protocol.Checkpoint, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RestoreCheckpoint not implemented")
}
func (*UnimplementedDeriveServer) Transaction(srv Derive_TransactionServer) error {
	return status.Errorf(codes.Unimplemented, "method Transaction not implemented")
}
func (*UnimplementedDeriveServer) BuildHints(ctx context.Context, req *empty.Empty) (*recoverylog.FSMHints, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuildHints not implemented")
}

func RegisterDeriveServer(s *grpc.Server, srv DeriveServer) {
	s.RegisterService(&_Derive_serviceDesc, srv)
}

func _Derive_RestoreCheckpoint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeriveServer).RestoreCheckpoint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/flow.Derive/RestoreCheckpoint",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeriveServer).RestoreCheckpoint(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Derive_Transaction_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DeriveServer).Transaction(&deriveTransactionServer{stream})
}

type Derive_TransactionServer interface {
	Send(*TxnResponse) error
	Recv() (*TxnRequest, error)
	grpc.ServerStream
}

type deriveTransactionServer struct {
	grpc.ServerStream
}

func (x *deriveTransactionServer) Send(m *TxnResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *deriveTransactionServer) Recv() (*TxnRequest, error) {
	m := new(TxnRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Derive_BuildHints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeriveServer).BuildHints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/flow.Derive/BuildHints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeriveServer).BuildHints(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Derive_serviceDesc = grpc.ServiceDesc{
	ServiceName: "flow.Derive",
	HandlerType: (*DeriveServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RestoreCheckpoint",
			Handler:    _Derive_RestoreCheckpoint_Handler,
		},
		{
			MethodName: "BuildHints",
			Handler:    _Derive_BuildHints_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Transaction",
			Handler:       _Derive_Transaction_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "go/protocol/flow.proto",
}

func (m *TxnRequest) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxnRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.ProtoSize()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxnRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.PrepareCheckpoint != nil {
		{
			size, err := m.PrepareCheckpoint.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintFlow(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if len(m.ExtendDocuments) > 0 {
		for iNdEx := len(m.ExtendDocuments) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ExtendDocuments[iNdEx])
			copy(dAtA[i:], m.ExtendDocuments[iNdEx])
			i = encodeVarintFlow(dAtA, i, uint64(len(m.ExtendDocuments[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.ExtendSource) > 0 {
		i -= len(m.ExtendSource)
		copy(dAtA[i:], m.ExtendSource)
		i = encodeVarintFlow(dAtA, i, uint64(len(m.ExtendSource)))
		i--
		dAtA[i] = 0x12
	}
	if m.State != 0 {
		i = encodeVarintFlow(dAtA, i, uint64(m.State))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TxnResponse) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxnResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.ProtoSize()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxnResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.ExtendLabels != nil {
		{
			size, err := m.ExtendLabels.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintFlow(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ExtendDocuments) > 0 {
		for iNdEx := len(m.ExtendDocuments) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ExtendDocuments[iNdEx])
			copy(dAtA[i:], m.ExtendDocuments[iNdEx])
			i = encodeVarintFlow(dAtA, i, uint64(len(m.ExtendDocuments[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if m.State != 0 {
		i = encodeVarintFlow(dAtA, i, uint64(m.State))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintFlow(dAtA []byte, offset int, v uint64) int {
	offset -= sovFlow(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TxnRequest) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.State != 0 {
		n += 1 + sovFlow(uint64(m.State))
	}
	l = len(m.ExtendSource)
	if l > 0 {
		n += 1 + l + sovFlow(uint64(l))
	}
	if len(m.ExtendDocuments) > 0 {
		for _, b := range m.ExtendDocuments {
			l = len(b)
			n += 1 + l + sovFlow(uint64(l))
		}
	}
	if m.PrepareCheckpoint != nil {
		l = m.PrepareCheckpoint.ProtoSize()
		n += 1 + l + sovFlow(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TxnResponse) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.State != 0 {
		n += 1 + sovFlow(uint64(m.State))
	}
	if len(m.ExtendDocuments) > 0 {
		for _, b := range m.ExtendDocuments {
			l = len(b)
			n += 1 + l + sovFlow(uint64(l))
		}
	}
	if m.ExtendLabels != nil {
		l = m.ExtendLabels.ProtoSize()
		n += 1 + l + sovFlow(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovFlow(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozFlow(x uint64) (n int) {
	return sovFlow(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TxnRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFlow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TxnRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxnRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFlow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= DeriveTxnState(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExtendSource", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFlow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFlow
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFlow
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExtendSource = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExtendDocuments", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFlow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthFlow
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthFlow
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExtendDocuments = append(m.ExtendDocuments, make([]byte, postIndex-iNdEx))
			copy(m.ExtendDocuments[len(m.ExtendDocuments)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PrepareCheckpoint", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFlow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFlow
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFlow
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PrepareCheckpoint == nil {
				m.PrepareCheckpoint = &protocol.Checkpoint{}
			}
			if err := m.PrepareCheckpoint.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFlow(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthFlow
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthFlow
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TxnResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFlow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TxnResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxnResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFlow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= DeriveTxnState(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExtendDocuments", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFlow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthFlow
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthFlow
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExtendDocuments = append(m.ExtendDocuments, make([]byte, postIndex-iNdEx))
			copy(m.ExtendDocuments[len(m.ExtendDocuments)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExtendLabels", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFlow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFlow
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFlow
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ExtendLabels == nil {
				m.ExtendLabels = &protocol1.LabelSet{}
			}
			if err := m.ExtendLabels.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFlow(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthFlow
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthFlow
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipFlow(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFlow
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowFlow
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowFlow
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthFlow
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupFlow
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthFlow
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthFlow        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFlow          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupFlow = fmt.Errorf("proto: unexpected end of group")
)
