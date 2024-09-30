package coderd_test

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/coder/coder/v2/coderd"
	"github.com/coder/coder/v2/coderd/database"
	"github.com/coder/coder/v2/coderd/database/pubsub"
	"github.com/coder/coder/v2/coderd/rbac"
	"github.com/coder/coder/v2/codersdk"
	"github.com/coder/coder/v2/tailnet"
	"github.com/coder/coder/v2/tailnet/proto"
	"github.com/coder/coder/v2/testutil"
)

func TestWorkspaceUpdates(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	peerID := uuid.New()

	ws1ID := uuid.New()
	ws1IDSlice := tailnet.UUIDToByteSlice(ws1ID)
	agent1ID := uuid.New()
	agent1IDSlice := tailnet.UUIDToByteSlice(agent1ID)
	ws2ID := uuid.New()
	ws2IDSlice := tailnet.UUIDToByteSlice(ws2ID)
	ws3ID := uuid.New()
	ws3IDSlice := tailnet.UUIDToByteSlice(ws3ID)
	ownerID := uuid.New()
	agent2ID := uuid.New()
	agent2IDSlice := tailnet.UUIDToByteSlice(agent2ID)
	ws4ID := uuid.New()
	ws4IDSlice := tailnet.UUIDToByteSlice(ws4ID)

	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		db := &mockWorkspaceStore{
			orderedRows: []database.GetWorkspacesAndAgentsRow{
				// Gains a new agent
				{
					ID:         ws1ID,
					Name:       "ws1",
					OwnerID:    ownerID,
					JobStatus:  database.ProvisionerJobStatusRunning,
					Transition: database.WorkspaceTransitionStart,
					Agents: []database.AgentIDNamePair{
						{
							ID:   agent1ID,
							Name: "agent1",
						},
					},
				},
				// Changes status
				{
					ID:         ws2ID,
					Name:       "ws2",
					OwnerID:    ownerID,
					JobStatus:  database.ProvisionerJobStatusRunning,
					Transition: database.WorkspaceTransitionStart,
				},
				// Is deleted
				{
					ID:         ws3ID,
					Name:       "ws3",
					OwnerID:    ownerID,
					JobStatus:  database.ProvisionerJobStatusSucceeded,
					Transition: database.WorkspaceTransitionStop,
				},
			},
		}

		ps := &mockPubsub{
			cbs: map[string]pubsub.Listener{},
		}

		updateProvider, err := coderd.NewUpdatesProvider(ctx, db, ps)
		defer updateProvider.Stop()
		require.NoError(t, err)

		ch, err := updateProvider.Subscribe(peerID, ownerID)
		require.NoError(t, err)

		update, ok := <-ch
		require.True(t, ok)
		slices.SortFunc(update.UpsertedWorkspaces, func(a, b *proto.Workspace) int {
			return strings.Compare(a.Name, b.Name)
		})
		require.Equal(t, &proto.WorkspaceUpdate{
			UpsertedWorkspaces: []*proto.Workspace{
				{
					Id:     ws1IDSlice,
					Name:   "ws1",
					Status: proto.Workspace_STARTING,
				},
				{
					Id:     ws2IDSlice,
					Name:   "ws2",
					Status: proto.Workspace_STARTING,
				},
				{
					Id:     ws3IDSlice,
					Name:   "ws3",
					Status: proto.Workspace_STOPPED,
				},
			},
			UpsertedAgents: []*proto.Agent{
				{
					Id:          agent1IDSlice,
					Name:        "agent1",
					WorkspaceId: ws1IDSlice,
				},
			},
			DeletedWorkspaces: []*proto.Workspace{},
			DeletedAgents:     []*proto.Agent{},
		}, update)

		// Update the database
		db.orderedRows = []database.GetWorkspacesAndAgentsRow{
			{
				ID:         ws1ID,
				Name:       "ws1",
				OwnerID:    ownerID,
				JobStatus:  database.ProvisionerJobStatusRunning,
				Transition: database.WorkspaceTransitionStart,
				Agents: []database.AgentIDNamePair{
					{
						ID:   agent1ID,
						Name: "agent1",
					},
					{
						ID:   agent2ID,
						Name: "agent2",
					},
				},
			},
			{
				ID:         ws2ID,
				Name:       "ws2",
				OwnerID:    ownerID,
				JobStatus:  database.ProvisionerJobStatusRunning,
				Transition: database.WorkspaceTransitionStop,
			},
			{
				ID:         ws4ID,
				Name:       "ws4",
				OwnerID:    ownerID,
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
					Id:     ws2IDSlice,
					Name:   "ws2",
					Status: proto.Workspace_STOPPING,
				},
				{
					// New workspace
					Id:     ws4IDSlice,
					Name:   "ws4",
					Status: proto.Workspace_STARTING,
				},
			},
			UpsertedAgents: []*proto.Agent{
				{
					Id:          agent2IDSlice,
					Name:        "agent2",
					WorkspaceId: ws1IDSlice,
				},
			},
			DeletedWorkspaces: []*proto.Workspace{
				{
					Id:     ws3IDSlice,
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
					ID:         ws1ID,
					Name:       "ws1",
					OwnerID:    ownerID,
					JobStatus:  database.ProvisionerJobStatusRunning,
					Transition: database.WorkspaceTransitionStart,
					Agents: []database.AgentIDNamePair{
						{
							ID:   agent1ID,
							Name: "agent1",
						},
					},
				},
			},
		}

		ps := &mockPubsub{
			cbs: map[string]pubsub.Listener{},
		}

		updateProvider, err := coderd.NewUpdatesProvider(ctx, db, ps)
		defer updateProvider.Stop()
		require.NoError(t, err)

		ch, err := updateProvider.Subscribe(peerID, ownerID)
		require.NoError(t, err)

		expected := &proto.WorkspaceUpdate{
			UpsertedWorkspaces: []*proto.Workspace{
				{
					Id:     ws1IDSlice,
					Name:   "ws1",
					Status: proto.Workspace_STARTING,
				},
			},
			UpsertedAgents: []*proto.Agent{
				{
					Id:          agent1IDSlice,
					Name:        "agent1",
					WorkspaceId: ws1IDSlice,
				},
			},
			DeletedWorkspaces: []*proto.Workspace{},
			DeletedAgents:     []*proto.Agent{},
		}

		update := testutil.RequireRecvCtx(ctx, t, ch)
		slices.SortFunc(update.UpsertedWorkspaces, func(a, b *proto.Workspace) int {
			return strings.Compare(a.Name, b.Name)
		})
		require.Equal(t, expected, update)

		updateProvider.Unsubscribe(ownerID)
		require.NoError(t, err)
		ch, err = updateProvider.Subscribe(peerID, ownerID)
		require.NoError(t, err)

		update = testutil.RequireRecvCtx(ctx, t, ch)
		slices.SortFunc(update.UpsertedWorkspaces, func(a, b *proto.Workspace) int {
			return strings.Compare(a.Name, b.Name)
		})
		require.Equal(t, expected, update)
	})
}

type mockWorkspaceStore struct {
	orderedRows []database.GetWorkspacesAndAgentsRow
}

// GetWorkspaceRBACByAgentID implements tailnet.UpdateQuerier.
func (*mockWorkspaceStore) GetWorkspaceRBACByAgentID(context.Context, uuid.UUID) (rbac.Objecter, error) {
	panic("unimplemented")
}

// GetWorkspacesAndAgents implements tailnet.UpdateQuerier.
func (m *mockWorkspaceStore) GetWorkspacesAndAgents(context.Context) ([]database.GetWorkspacesAndAgentsRow, error) {
	return m.orderedRows, nil
}

var _ coderd.UpdateQuerier = (*mockWorkspaceStore)(nil)

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
