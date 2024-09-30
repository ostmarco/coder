package tailnet_test

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/coder/coder/v2/coderd/database"
	"github.com/coder/coder/v2/coderd/database/pubsub"
	"github.com/coder/coder/v2/codersdk"
	"github.com/coder/coder/v2/tailnet"
	"github.com/coder/coder/v2/tailnet/proto"
)

func TestWorkspaceUpdates(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	peerID := uuid.New()

	ws1id := uuid.New()
	uuid1slice := tailnet.UUIDToByteSlice(ws1id)
	agent1id := uuid.New()
	agent1idslice := tailnet.UUIDToByteSlice(agent1id)
	ws2id := uuid.New()
	ws2idslice := tailnet.UUIDToByteSlice(ws2id)
	ws3id := uuid.New()
	ws3idslice := tailnet.UUIDToByteSlice(ws3id)
	ownerid := uuid.New()
	agent2id := uuid.New()
	agent2idslice := tailnet.UUIDToByteSlice(agent2id)
	ws4id := uuid.New()
	ws4idslice := tailnet.UUIDToByteSlice(ws4id)

	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		db := &mockWorkspaceStore{
			orderedRows: []database.GetWorkspacesAndAgentsRow{
				// Gains a new agent
				{
					ID:         ws1id,
					Name:       "ws1",
					OwnerID:    ownerid,
					JobStatus:  database.ProvisionerJobStatusRunning,
					Transition: database.WorkspaceTransitionStart,
					Agents: []database.AgentIDNamePair{
						{
							ID:   agent1id,
							Name: "agent1",
						},
					},
				},
				// Changes status
				{
					ID:         ws2id,
					Name:       "ws2",
					OwnerID:    ownerid,
					JobStatus:  database.ProvisionerJobStatusRunning,
					Transition: database.WorkspaceTransitionStart,
				},
				// Is deleted
				{
					ID:         ws3id,
					Name:       "ws3",
					OwnerID:    ownerid,
					JobStatus:  database.ProvisionerJobStatusSucceeded,
					Transition: database.WorkspaceTransitionStop,
				},
			},
		}

		ps := &mockPubsub{
			cbs: map[string]pubsub.Listener{},
		}

		updateProvider, err := tailnet.NewUpdatesProvider(ctx, db, ps)
		require.NoError(t, err)

		ch, err := updateProvider.Subscribe(peerID, ownerid)
		require.NoError(t, err)

		update, ok := <-ch
		require.True(t, ok)
		slices.SortFunc(update.UpsertedWorkspaces, func(a, b *proto.Workspace) int {
			return strings.Compare(a.Name, b.Name)
		})
		require.Equal(t, &proto.WorkspaceUpdate{
			UpsertedWorkspaces: []*proto.Workspace{
				{
					Id:     uuid1slice,
					Name:   "ws1",
					Status: proto.Workspace_STARTING,
				},
				{
					Id:     ws2idslice,
					Name:   "ws2",
					Status: proto.Workspace_STARTING,
				},
				{
					Id:     ws3idslice,
					Name:   "ws3",
					Status: proto.Workspace_STOPPED,
				},
			},
			UpsertedAgents: []*proto.Agent{
				{
					Id:          agent1idslice,
					Name:        "agent1",
					WorkspaceId: uuid1slice,
				},
			},
			DeletedWorkspaces: []*proto.Workspace{},
			DeletedAgents:     []*proto.Agent{},
		}, update)

		// Update the database
		db.orderedRows = []database.GetWorkspacesAndAgentsRow{
			{
				ID:         ws1id,
				Name:       "ws1",
				OwnerID:    ownerid,
				JobStatus:  database.ProvisionerJobStatusRunning,
				Transition: database.WorkspaceTransitionStart,
				Agents: []database.AgentIDNamePair{
					{
						ID:   agent1id,
						Name: "agent1",
					},
					{
						ID:   agent2id,
						Name: "agent2",
					},
				},
			},
			{
				ID:         ws2id,
				Name:       "ws2",
				OwnerID:    ownerid,
				JobStatus:  database.ProvisionerJobStatusRunning,
				Transition: database.WorkspaceTransitionStop,
			},
			{
				ID:         ws4id,
				Name:       "ws4",
				OwnerID:    ownerid,
				JobStatus:  database.ProvisionerJobStatusRunning,
				Transition: database.WorkspaceTransitionStart,
			},
		}
		ps.Publish(codersdk.AllWorkspacesNotifyChannel, nil)

		update, ok = <-ch
		require.True(t, ok)
		slices.SortFunc(update.UpsertedWorkspaces, func(a, b *proto.Workspace) int {
			return strings.Compare(a.Name, b.Name)
		})
		require.Equal(t, &proto.WorkspaceUpdate{
			UpsertedWorkspaces: []*proto.Workspace{
				{
					// Changed status
					Id:     ws2idslice,
					Name:   "ws2",
					Status: proto.Workspace_STOPPING,
				},
				{
					// New workspace
					Id:     ws4idslice,
					Name:   "ws4",
					Status: proto.Workspace_STARTING,
				},
			},
			UpsertedAgents: []*proto.Agent{
				{
					Id:          agent2idslice,
					Name:        "agent2",
					WorkspaceId: uuid1slice,
				},
			},
			DeletedWorkspaces: []*proto.Workspace{
				{
					Id:     ws3idslice,
					Name:   "ws3",
					Status: proto.Workspace_STOPPED,
				},
			},
			DeletedAgents: []*proto.Agent{},
		}, update)
	})

	t.Run("Resubscribe", func(t *testing.T) {
		t.Parallel()

		db := &mockWorkspaceStore{
			orderedRows: []database.GetWorkspacesAndAgentsRow{
				{
					ID:         ws1id,
					Name:       "ws1",
					OwnerID:    ownerid,
					JobStatus:  database.ProvisionerJobStatusRunning,
					Transition: database.WorkspaceTransitionStart,
					Agents: []database.AgentIDNamePair{
						{
							ID:   agent1id,
							Name: "agent1",
						},
					},
				},
			},
		}

		ps := &mockPubsub{
			cbs: map[string]pubsub.Listener{},
		}

		updateProvider, err := tailnet.NewUpdatesProvider(ctx, db, ps)
		require.NoError(t, err)

		ch, err := updateProvider.Subscribe(peerID, ownerid)
		require.NoError(t, err)

		update, ok := <-ch
		require.True(t, ok)
		slices.SortFunc(update.UpsertedWorkspaces, func(a, b *proto.Workspace) int {
			return strings.Compare(a.Name, b.Name)
		})
		require.Equal(t, &proto.WorkspaceUpdate{
			UpsertedWorkspaces: []*proto.Workspace{
				{
					Id:     uuid1slice,
					Name:   "ws1",
					Status: proto.Workspace_STARTING,
				},
			},
			UpsertedAgents: []*proto.Agent{
				{
					Id:          agent1idslice,
					Name:        "agent1",
					WorkspaceId: uuid1slice,
				},
			},
			DeletedWorkspaces: []*proto.Workspace{},
			DeletedAgents:     []*proto.Agent{},
		}, update)

		// Unsubscribe
		// err = updateProvider.Unsubscribe(ownerid)
		// require.NoError(t, err)
	})
}

type mockWorkspaceStore struct {
	orderedRows []database.GetWorkspacesAndAgentsRow
}

var _ tailnet.WorkspaceStore = (*mockWorkspaceStore)(nil)

func (*mockWorkspaceStore) GetWorkspaceByAgentID(context.Context, uuid.UUID) (database.GetWorkspaceByAgentIDRow, error) {
	return database.GetWorkspaceByAgentIDRow{}, nil
}
func (db *mockWorkspaceStore) GetWorkspacesAndAgents(context.Context) ([]database.GetWorkspacesAndAgentsRow, error) {
	return db.orderedRows, nil
}

type mockPubsub struct {
	cbs map[string]pubsub.Listener
}

// Close implements pubsub.Pubsub.
func (*mockPubsub) Close() error {
	panic("unimplemented")
}

// Publish implements pubsub.Pubsub.
func (m *mockPubsub) Publish(event string, message []byte) error {
	cb, ok := m.cbs[event]
	if !ok {
		return nil
	}
	cb(context.Background(), message)
	return nil
}

// Subscribe implements pubsub.Pubsub.
func (m *mockPubsub) Subscribe(event string, listener pubsub.Listener) (cancel func(), err error) {
	m.cbs[event] = listener
	return func() {}, nil
}

// SubscribeWithErr implements pubsub.Pubsub.
func (*mockPubsub) SubscribeWithErr(string, pubsub.ListenerWithErr) (func(), error) {
	panic("unimplemented")
}

var _ pubsub.Pubsub = (*mockPubsub)(nil)
