package tailnet

import (
	"slices"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/coder/coder/v2/coderd/database"
	"github.com/coder/coder/v2/tailnet/proto"
)

func TestProduceUpdate(t *testing.T) {
	t.Parallel()

	uuid1 := uuid.New()
	uuid1slice := UUIDToByteSlice(uuid1)
	uuid2 := uuid.New()
	uuid2slice := UUIDToByteSlice(uuid2)
	uuid3 := uuid.New()
	uuid3slice := UUIDToByteSlice(uuid3)
	uuid4 := uuid.New()
	uuid4slice := UUIDToByteSlice(uuid4)
	uuid5 := uuid.New()
	uuid5slice := UUIDToByteSlice(uuid5)
	tt := []struct {
		name     string
		old, new workspacesByID
		expected *proto.WorkspaceUpdate
	}{
		{
			name: "InitialUpdate",
			old:  map[uuid.UUID]ownedWorkspace{},
			new: map[uuid.UUID]ownedWorkspace{
				uuid1: {
					WorkspaceName: "ws1",
					JobStatus:     database.ProvisionerJobStatusRunning,
					Transition:    database.WorkspaceTransitionStart,
					Agents: []database.AgentIDNamePair{
						{
							ID:   uuid2,
							Name: "agent1",
						},
						{
							ID:   uuid3,
							Name: "agent2",
						},
					},
				},
			},
			expected: &proto.WorkspaceUpdate{
				UpsertedWorkspaces: []*proto.Workspace{
					{
						Id:     uuid1slice,
						Name:   "ws1",
						Status: proto.Workspace_STARTING,
					},
				},
				UpsertedAgents: []*proto.Agent{
					{
						Id:          uuid2slice,
						Name:        "agent1",
						WorkspaceId: uuid1slice,
					},
					{
						Id:          uuid3slice,
						Name:        "agent2",
						WorkspaceId: uuid1slice,
					},
				},
				DeletedWorkspaces: []*proto.Workspace{},
				DeletedAgents:     []*proto.Agent{},
			},
		},
		{
			name: "Complex",
			old: map[uuid.UUID]ownedWorkspace{
				// Gains a new agent
				uuid1: {
					WorkspaceName: "ws1",
					JobStatus:     database.ProvisionerJobStatusSucceeded,
					Transition:    database.WorkspaceTransitionStart,
					Agents: []database.AgentIDNamePair{
						{
							ID:   uuid2,
							Name: "agent1",
						},
					},
				},
				// Changes status
				uuid2: {
					WorkspaceName: "ws2",
					JobStatus:     database.ProvisionerJobStatusRunning,
					Transition:    database.WorkspaceTransitionStart,
				},
				// Is deleted
				uuid3: {
					WorkspaceName: "ws3",
					JobStatus:     database.ProvisionerJobStatusSucceeded,
					Transition:    database.WorkspaceTransitionStop,
				},
			},
			new: map[uuid.UUID]ownedWorkspace{
				uuid1: {
					WorkspaceName: "ws1",
					JobStatus:     database.ProvisionerJobStatusSucceeded,
					Transition:    database.WorkspaceTransitionStart,
					Agents: []database.AgentIDNamePair{
						{
							ID:   uuid2,
							Name: "agent1",
						},
						// New agent
						{
							ID:   uuid4,
							Name: "agent2",
						},
					},
				},
				uuid2: {
					WorkspaceName: "ws2",
					JobStatus:     database.ProvisionerJobStatusRunning,
					Transition:    database.WorkspaceTransitionStop,
					Agents:        []database.AgentIDNamePair{},
				},
				// New workspace
				uuid5: {
					WorkspaceName: "ws4",
					JobStatus:     database.ProvisionerJobStatusRunning,
					Transition:    database.WorkspaceTransitionStart,
					Agents:        []database.AgentIDNamePair{},
				},
			},
			expected: &proto.WorkspaceUpdate{
				UpsertedWorkspaces: []*proto.Workspace{
					{
						// Changed status
						Id:     uuid2slice,
						Name:   "ws2",
						Status: proto.Workspace_STOPPING,
					},
					{
						// New workspace
						Id:     uuid5slice,
						Name:   "ws4",
						Status: proto.Workspace_STARTING,
					},
				},
				UpsertedAgents: []*proto.Agent{
					{
						Id:          uuid4slice,
						Name:        "agent2",
						WorkspaceId: uuid1slice,
					},
				},
				DeletedWorkspaces: []*proto.Workspace{
					{
						Id:     uuid3slice,
						Name:   "ws3",
						Status: proto.Workspace_STOPPED,
					},
				},
				DeletedAgents: []*proto.Agent{},
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := produceUpdate(tc.old, tc.new)
			slices.SortFunc(actual.UpsertedWorkspaces, func(a, b *proto.Workspace) int {
				return strings.Compare(a.Name, b.Name)
			})
			slices.SortFunc(actual.DeletedWorkspaces, func(a, b *proto.Workspace) int {
				return strings.Compare(a.Name, b.Name)
			})
			slices.SortFunc(actual.UpsertedAgents, func(a, b *proto.Agent) int {
				return strings.Compare(a.Name, b.Name)
			})
			slices.SortFunc(actual.DeletedAgents, func(a, b *proto.Agent) int {
				return strings.Compare(a.Name, b.Name)
			})
			require.Equal(t, tc.expected, actual)
		})
	}
}
