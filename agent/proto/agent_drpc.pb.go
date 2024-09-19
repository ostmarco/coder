// Code generated by protoc-gen-go-drpc. DO NOT EDIT.
// protoc-gen-go-drpc version: v0.0.33
// source: agent/proto/agent.proto

package proto

import (
	context "context"
	errors "errors"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	drpc "storj.io/drpc"
	drpcerr "storj.io/drpc/drpcerr"
)

type drpcEncoding_File_agent_proto_agent_proto struct{}

func (drpcEncoding_File_agent_proto_agent_proto) Marshal(msg drpc.Message) ([]byte, error) {
	return proto.Marshal(msg.(proto.Message))
}

func (drpcEncoding_File_agent_proto_agent_proto) MarshalAppend(buf []byte, msg drpc.Message) ([]byte, error) {
	return proto.MarshalOptions{}.MarshalAppend(buf, msg.(proto.Message))
}

func (drpcEncoding_File_agent_proto_agent_proto) Unmarshal(buf []byte, msg drpc.Message) error {
	return proto.Unmarshal(buf, msg.(proto.Message))
}

func (drpcEncoding_File_agent_proto_agent_proto) JSONMarshal(msg drpc.Message) ([]byte, error) {
	return protojson.Marshal(msg.(proto.Message))
}

func (drpcEncoding_File_agent_proto_agent_proto) JSONUnmarshal(buf []byte, msg drpc.Message) error {
	return protojson.Unmarshal(buf, msg.(proto.Message))
}

type DRPCAgentClient interface {
	DRPCConn() drpc.Conn

	GetManifest(ctx context.Context, in *GetManifestRequest) (*Manifest, error)
	GetServiceBanner(ctx context.Context, in *GetServiceBannerRequest) (*ServiceBanner, error)
	UpdateStats(ctx context.Context, in *UpdateStatsRequest) (*UpdateStatsResponse, error)
	UpdateLifecycle(ctx context.Context, in *UpdateLifecycleRequest) (*Lifecycle, error)
	BatchUpdateAppHealths(ctx context.Context, in *BatchUpdateAppHealthRequest) (*BatchUpdateAppHealthResponse, error)
	UpdateStartup(ctx context.Context, in *UpdateStartupRequest) (*Startup, error)
	BatchUpdateMetadata(ctx context.Context, in *BatchUpdateMetadataRequest) (*BatchUpdateMetadataResponse, error)
	BatchCreateLogs(ctx context.Context, in *BatchCreateLogsRequest) (*BatchCreateLogsResponse, error)
	GetAnnouncementBanners(ctx context.Context, in *GetAnnouncementBannersRequest) (*GetAnnouncementBannersResponse, error)
	MarkWorkspacePrebuildReady(ctx context.Context, in *MarkWorkspacePrebuildReadyRequest) (*MarkWorkspacePrebuildReadyResponse, error)
}

type drpcAgentClient struct {
	cc drpc.Conn
}

func NewDRPCAgentClient(cc drpc.Conn) DRPCAgentClient {
	return &drpcAgentClient{cc}
}

func (c *drpcAgentClient) DRPCConn() drpc.Conn { return c.cc }

func (c *drpcAgentClient) GetManifest(ctx context.Context, in *GetManifestRequest) (*Manifest, error) {
	out := new(Manifest)
	err := c.cc.Invoke(ctx, "/coder.agent.v2.Agent/GetManifest", drpcEncoding_File_agent_proto_agent_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcAgentClient) GetServiceBanner(ctx context.Context, in *GetServiceBannerRequest) (*ServiceBanner, error) {
	out := new(ServiceBanner)
	err := c.cc.Invoke(ctx, "/coder.agent.v2.Agent/GetServiceBanner", drpcEncoding_File_agent_proto_agent_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcAgentClient) UpdateStats(ctx context.Context, in *UpdateStatsRequest) (*UpdateStatsResponse, error) {
	out := new(UpdateStatsResponse)
	err := c.cc.Invoke(ctx, "/coder.agent.v2.Agent/UpdateStats", drpcEncoding_File_agent_proto_agent_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcAgentClient) UpdateLifecycle(ctx context.Context, in *UpdateLifecycleRequest) (*Lifecycle, error) {
	out := new(Lifecycle)
	err := c.cc.Invoke(ctx, "/coder.agent.v2.Agent/UpdateLifecycle", drpcEncoding_File_agent_proto_agent_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcAgentClient) BatchUpdateAppHealths(ctx context.Context, in *BatchUpdateAppHealthRequest) (*BatchUpdateAppHealthResponse, error) {
	out := new(BatchUpdateAppHealthResponse)
	err := c.cc.Invoke(ctx, "/coder.agent.v2.Agent/BatchUpdateAppHealths", drpcEncoding_File_agent_proto_agent_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcAgentClient) UpdateStartup(ctx context.Context, in *UpdateStartupRequest) (*Startup, error) {
	out := new(Startup)
	err := c.cc.Invoke(ctx, "/coder.agent.v2.Agent/UpdateStartup", drpcEncoding_File_agent_proto_agent_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcAgentClient) BatchUpdateMetadata(ctx context.Context, in *BatchUpdateMetadataRequest) (*BatchUpdateMetadataResponse, error) {
	out := new(BatchUpdateMetadataResponse)
	err := c.cc.Invoke(ctx, "/coder.agent.v2.Agent/BatchUpdateMetadata", drpcEncoding_File_agent_proto_agent_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcAgentClient) BatchCreateLogs(ctx context.Context, in *BatchCreateLogsRequest) (*BatchCreateLogsResponse, error) {
	out := new(BatchCreateLogsResponse)
	err := c.cc.Invoke(ctx, "/coder.agent.v2.Agent/BatchCreateLogs", drpcEncoding_File_agent_proto_agent_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcAgentClient) GetAnnouncementBanners(ctx context.Context, in *GetAnnouncementBannersRequest) (*GetAnnouncementBannersResponse, error) {
	out := new(GetAnnouncementBannersResponse)
	err := c.cc.Invoke(ctx, "/coder.agent.v2.Agent/GetAnnouncementBanners", drpcEncoding_File_agent_proto_agent_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcAgentClient) MarkWorkspacePrebuildReady(ctx context.Context, in *MarkWorkspacePrebuildReadyRequest) (*MarkWorkspacePrebuildReadyResponse, error) {
	out := new(MarkWorkspacePrebuildReadyResponse)
	err := c.cc.Invoke(ctx, "/coder.agent.v2.Agent/MarkWorkspacePrebuildReady", drpcEncoding_File_agent_proto_agent_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type DRPCAgentServer interface {
	GetManifest(context.Context, *GetManifestRequest) (*Manifest, error)
	GetServiceBanner(context.Context, *GetServiceBannerRequest) (*ServiceBanner, error)
	UpdateStats(context.Context, *UpdateStatsRequest) (*UpdateStatsResponse, error)
	UpdateLifecycle(context.Context, *UpdateLifecycleRequest) (*Lifecycle, error)
	BatchUpdateAppHealths(context.Context, *BatchUpdateAppHealthRequest) (*BatchUpdateAppHealthResponse, error)
	UpdateStartup(context.Context, *UpdateStartupRequest) (*Startup, error)
	BatchUpdateMetadata(context.Context, *BatchUpdateMetadataRequest) (*BatchUpdateMetadataResponse, error)
	BatchCreateLogs(context.Context, *BatchCreateLogsRequest) (*BatchCreateLogsResponse, error)
	GetAnnouncementBanners(context.Context, *GetAnnouncementBannersRequest) (*GetAnnouncementBannersResponse, error)
	MarkWorkspacePrebuildReady(context.Context, *MarkWorkspacePrebuildReadyRequest) (*MarkWorkspacePrebuildReadyResponse, error)
}

type DRPCAgentUnimplementedServer struct{}

func (s *DRPCAgentUnimplementedServer) GetManifest(context.Context, *GetManifestRequest) (*Manifest, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCAgentUnimplementedServer) GetServiceBanner(context.Context, *GetServiceBannerRequest) (*ServiceBanner, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCAgentUnimplementedServer) UpdateStats(context.Context, *UpdateStatsRequest) (*UpdateStatsResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCAgentUnimplementedServer) UpdateLifecycle(context.Context, *UpdateLifecycleRequest) (*Lifecycle, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCAgentUnimplementedServer) BatchUpdateAppHealths(context.Context, *BatchUpdateAppHealthRequest) (*BatchUpdateAppHealthResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCAgentUnimplementedServer) UpdateStartup(context.Context, *UpdateStartupRequest) (*Startup, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCAgentUnimplementedServer) BatchUpdateMetadata(context.Context, *BatchUpdateMetadataRequest) (*BatchUpdateMetadataResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCAgentUnimplementedServer) BatchCreateLogs(context.Context, *BatchCreateLogsRequest) (*BatchCreateLogsResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCAgentUnimplementedServer) GetAnnouncementBanners(context.Context, *GetAnnouncementBannersRequest) (*GetAnnouncementBannersResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCAgentUnimplementedServer) MarkWorkspacePrebuildReady(context.Context, *MarkWorkspacePrebuildReadyRequest) (*MarkWorkspacePrebuildReadyResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

type DRPCAgentDescription struct{}

func (DRPCAgentDescription) NumMethods() int { return 10 }

func (DRPCAgentDescription) Method(n int) (string, drpc.Encoding, drpc.Receiver, interface{}, bool) {
	switch n {
	case 0:
		return "/coder.agent.v2.Agent/GetManifest", drpcEncoding_File_agent_proto_agent_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCAgentServer).
					GetManifest(
						ctx,
						in1.(*GetManifestRequest),
					)
			}, DRPCAgentServer.GetManifest, true
	case 1:
		return "/coder.agent.v2.Agent/GetServiceBanner", drpcEncoding_File_agent_proto_agent_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCAgentServer).
					GetServiceBanner(
						ctx,
						in1.(*GetServiceBannerRequest),
					)
			}, DRPCAgentServer.GetServiceBanner, true
	case 2:
		return "/coder.agent.v2.Agent/UpdateStats", drpcEncoding_File_agent_proto_agent_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCAgentServer).
					UpdateStats(
						ctx,
						in1.(*UpdateStatsRequest),
					)
			}, DRPCAgentServer.UpdateStats, true
	case 3:
		return "/coder.agent.v2.Agent/UpdateLifecycle", drpcEncoding_File_agent_proto_agent_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCAgentServer).
					UpdateLifecycle(
						ctx,
						in1.(*UpdateLifecycleRequest),
					)
			}, DRPCAgentServer.UpdateLifecycle, true
	case 4:
		return "/coder.agent.v2.Agent/BatchUpdateAppHealths", drpcEncoding_File_agent_proto_agent_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCAgentServer).
					BatchUpdateAppHealths(
						ctx,
						in1.(*BatchUpdateAppHealthRequest),
					)
			}, DRPCAgentServer.BatchUpdateAppHealths, true
	case 5:
		return "/coder.agent.v2.Agent/UpdateStartup", drpcEncoding_File_agent_proto_agent_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCAgentServer).
					UpdateStartup(
						ctx,
						in1.(*UpdateStartupRequest),
					)
			}, DRPCAgentServer.UpdateStartup, true
	case 6:
		return "/coder.agent.v2.Agent/BatchUpdateMetadata", drpcEncoding_File_agent_proto_agent_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCAgentServer).
					BatchUpdateMetadata(
						ctx,
						in1.(*BatchUpdateMetadataRequest),
					)
			}, DRPCAgentServer.BatchUpdateMetadata, true
	case 7:
		return "/coder.agent.v2.Agent/BatchCreateLogs", drpcEncoding_File_agent_proto_agent_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCAgentServer).
					BatchCreateLogs(
						ctx,
						in1.(*BatchCreateLogsRequest),
					)
			}, DRPCAgentServer.BatchCreateLogs, true
	case 8:
		return "/coder.agent.v2.Agent/GetAnnouncementBanners", drpcEncoding_File_agent_proto_agent_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCAgentServer).
					GetAnnouncementBanners(
						ctx,
						in1.(*GetAnnouncementBannersRequest),
					)
			}, DRPCAgentServer.GetAnnouncementBanners, true
	case 9:
		return "/coder.agent.v2.Agent/MarkWorkspacePrebuildReady", drpcEncoding_File_agent_proto_agent_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCAgentServer).
					MarkWorkspacePrebuildReady(
						ctx,
						in1.(*MarkWorkspacePrebuildReadyRequest),
					)
			}, DRPCAgentServer.MarkWorkspacePrebuildReady, true
	default:
		return "", nil, nil, nil, false
	}
}

func DRPCRegisterAgent(mux drpc.Mux, impl DRPCAgentServer) error {
	return mux.Register(impl, DRPCAgentDescription{})
}

type DRPCAgent_GetManifestStream interface {
	drpc.Stream
	SendAndClose(*Manifest) error
}

type drpcAgent_GetManifestStream struct {
	drpc.Stream
}

func (x *drpcAgent_GetManifestStream) SendAndClose(m *Manifest) error {
	if err := x.MsgSend(m, drpcEncoding_File_agent_proto_agent_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCAgent_GetServiceBannerStream interface {
	drpc.Stream
	SendAndClose(*ServiceBanner) error
}

type drpcAgent_GetServiceBannerStream struct {
	drpc.Stream
}

func (x *drpcAgent_GetServiceBannerStream) SendAndClose(m *ServiceBanner) error {
	if err := x.MsgSend(m, drpcEncoding_File_agent_proto_agent_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCAgent_UpdateStatsStream interface {
	drpc.Stream
	SendAndClose(*UpdateStatsResponse) error
}

type drpcAgent_UpdateStatsStream struct {
	drpc.Stream
}

func (x *drpcAgent_UpdateStatsStream) SendAndClose(m *UpdateStatsResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_agent_proto_agent_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCAgent_UpdateLifecycleStream interface {
	drpc.Stream
	SendAndClose(*Lifecycle) error
}

type drpcAgent_UpdateLifecycleStream struct {
	drpc.Stream
}

func (x *drpcAgent_UpdateLifecycleStream) SendAndClose(m *Lifecycle) error {
	if err := x.MsgSend(m, drpcEncoding_File_agent_proto_agent_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCAgent_BatchUpdateAppHealthsStream interface {
	drpc.Stream
	SendAndClose(*BatchUpdateAppHealthResponse) error
}

type drpcAgent_BatchUpdateAppHealthsStream struct {
	drpc.Stream
}

func (x *drpcAgent_BatchUpdateAppHealthsStream) SendAndClose(m *BatchUpdateAppHealthResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_agent_proto_agent_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCAgent_UpdateStartupStream interface {
	drpc.Stream
	SendAndClose(*Startup) error
}

type drpcAgent_UpdateStartupStream struct {
	drpc.Stream
}

func (x *drpcAgent_UpdateStartupStream) SendAndClose(m *Startup) error {
	if err := x.MsgSend(m, drpcEncoding_File_agent_proto_agent_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCAgent_BatchUpdateMetadataStream interface {
	drpc.Stream
	SendAndClose(*BatchUpdateMetadataResponse) error
}

type drpcAgent_BatchUpdateMetadataStream struct {
	drpc.Stream
}

func (x *drpcAgent_BatchUpdateMetadataStream) SendAndClose(m *BatchUpdateMetadataResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_agent_proto_agent_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCAgent_BatchCreateLogsStream interface {
	drpc.Stream
	SendAndClose(*BatchCreateLogsResponse) error
}

type drpcAgent_BatchCreateLogsStream struct {
	drpc.Stream
}

func (x *drpcAgent_BatchCreateLogsStream) SendAndClose(m *BatchCreateLogsResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_agent_proto_agent_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCAgent_GetAnnouncementBannersStream interface {
	drpc.Stream
	SendAndClose(*GetAnnouncementBannersResponse) error
}

type drpcAgent_GetAnnouncementBannersStream struct {
	drpc.Stream
}

func (x *drpcAgent_GetAnnouncementBannersStream) SendAndClose(m *GetAnnouncementBannersResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_agent_proto_agent_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCAgent_MarkWorkspacePrebuildReadyStream interface {
	drpc.Stream
	SendAndClose(*MarkWorkspacePrebuildReadyResponse) error
}

type drpcAgent_MarkWorkspacePrebuildReadyStream struct {
	drpc.Stream
}

func (x *drpcAgent_MarkWorkspacePrebuildReadyStream) SendAndClose(m *MarkWorkspacePrebuildReadyResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_agent_proto_agent_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}
