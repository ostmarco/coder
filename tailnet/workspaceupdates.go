package tailnet

import (
	"context"
	"database/sql"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"github.com/coder/coder/v2/coderd/database"
	"github.com/coder/coder/v2/coderd/database/pubsub"
	"github.com/coder/coder/v2/coderd/util/slice"
	"github.com/coder/coder/v2/codersdk"
	"github.com/coder/coder/v2/tailnet/proto"
)

type workspacesByOwner map[uuid.UUID]workspacesByID

type workspacesByID map[uuid.UUID]ownedWorkspace

type ownedWorkspace struct {
	WorkspaceName string
	JobStatus     database.ProvisionerJobStatus
	Transition    database.WorkspaceTransition
	Agents        []database.AgentIDNamePair
}

// Equal does not compare agents
func (w ownedWorkspace) Equal(other ownedWorkspace) bool {
	return w.WorkspaceName == other.WorkspaceName &&
		w.JobStatus == other.JobStatus &&
		w.Transition == other.Transition
}

func convertRows(v []database.GetWorkspacesAndAgentsRow) workspacesByOwner {
	m := make(map[uuid.UUID]workspacesByID)
	for _, ws := range v {
		owned := ownedWorkspace{
			WorkspaceName: ws.Name,
			JobStatus:     ws.JobStatus,
			Transition:    ws.Transition,
		}
		if byID, exists := m[ws.OwnerID]; !exists {
			m[ws.OwnerID] = map[uuid.UUID]ownedWorkspace{ws.ID: owned}
		} else {
			byID[ws.ID] = owned
			m[ws.OwnerID] = byID
		}
	}
	return workspacesByOwner(m)
}

func convertStatus(status database.ProvisionerJobStatus, trans database.WorkspaceTransition) proto.Workspace_Status {
	wsStatus := codersdk.ConvertWorkspaceStatus(codersdk.ProvisionerJobStatus(status), codersdk.WorkspaceTransition(trans))
	return WorkspaceStatusToProto(wsStatus)
}

type sub struct {
	mu     sync.Mutex
	userID uuid.UUID
	tx     chan<- *proto.WorkspaceUpdate
	prev   workspacesByID
}

func (s *sub) send(all workspacesByOwner) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Filter to only the workspaces owned by the user
	own := all[s.userID]
	update := produceUpdate(s.prev, own)
	s.prev = own
	s.tx <- update
}

type WorkspaceUpdatesProvider interface {
	Subscribe(userID uuid.UUID) (<-chan *proto.WorkspaceUpdate, error)
}

type WorkspaceStore interface {
	GetWorkspaceByAgentID(ctx context.Context, agentID uuid.UUID) (database.GetWorkspaceByAgentIDRow, error)
	GetWorkspacesAndAgents(ctx context.Context) ([]database.GetWorkspacesAndAgentsRow, error)
}

type updatesProvider struct {
	mu       sync.RWMutex
	db       WorkspaceStore
	ps       pubsub.Pubsub
	subs     []*sub
	cancelFn func()
}

var _ WorkspaceUpdatesProvider = (*updatesProvider)(nil)

func (u *updatesProvider) Start() error {
	cancel, err := u.ps.Subscribe(codersdk.AllWorkspacesNotifyChannel, u.handleUpdate)
	if err != nil {
		return err
	}
	u.cancelFn = cancel
	return nil
}

func (u *updatesProvider) Stop() {
	u.cancelFn()
}

func (u *updatesProvider) handleUpdate(ctx context.Context, _ []byte) {
	rows, err := u.db.GetWorkspacesAndAgents(ctx)
	if err != nil && !xerrors.Is(err, sql.ErrNoRows) {
		// TODO: Log
		return
	}

	wg := &sync.WaitGroup{}
	latest := convertRows(rows)
	u.mu.RLock()
	defer u.mu.RUnlock()
	for _, sub := range u.subs {
		sub := sub
		wg.Add(1)
		go func() {
			sub.send(latest)
			defer wg.Done()
		}()
	}
	wg.Wait()
}

func NewUpdatesProvider(db WorkspaceStore, ps pubsub.Pubsub) WorkspaceUpdatesProvider {
	return &updatesProvider{
		db:   db,
		ps:   ps,
		subs: make([]*sub, 0),
	}
}

func (u *updatesProvider) Subscribe(userID uuid.UUID) (<-chan *proto.WorkspaceUpdate, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	tx := make(chan *proto.WorkspaceUpdate, 1)
	sub := &sub{
		userID: userID,
		tx:     tx,
		prev:   make(workspacesByID),
	}
	u.subs = append(u.subs, sub)
	return tx, nil
}

func produceUpdate(old, new workspacesByID) *proto.WorkspaceUpdate {
	out := &proto.WorkspaceUpdate{
		UpsertedWorkspaces: []*proto.Workspace{},
		UpsertedAgents:     []*proto.Agent{},
		DeletedWorkspaces:  []*proto.Workspace{},
		DeletedAgents:      []*proto.Agent{},
	}

	for wsID, newWorkspace := range new {
		oldWorkspace, exists := old[wsID]
		// Upsert both workspace and agents if the workspace is new
		if !exists {
			out.UpsertedWorkspaces = append(out.UpsertedWorkspaces, &proto.Workspace{
				Id:     UUIDToByteSlice(wsID),
				Name:   newWorkspace.WorkspaceName,
				Status: convertStatus(newWorkspace.JobStatus, newWorkspace.Transition),
			})
			for _, agent := range newWorkspace.Agents {
				out.UpsertedAgents = append(out.UpsertedAgents, &proto.Agent{
					Id:          UUIDToByteSlice(agent.ID),
					Name:        agent.Name,
					WorkspaceId: UUIDToByteSlice(wsID),
				})
			}
			continue
		}
		// Upsert workspace if the workspace is updated
		if !newWorkspace.Equal(oldWorkspace) {
			out.UpsertedWorkspaces = append(out.UpsertedWorkspaces, &proto.Workspace{
				Id:     UUIDToByteSlice(wsID),
				Name:   newWorkspace.WorkspaceName,
				Status: convertStatus(newWorkspace.JobStatus, newWorkspace.Transition),
			})
		}

		add, remove := slice.SymmetricDifference(oldWorkspace.Agents, newWorkspace.Agents)
		for _, agent := range add {
			out.UpsertedAgents = append(out.UpsertedAgents, &proto.Agent{
				Id:          UUIDToByteSlice(agent.ID),
				Name:        agent.Name,
				WorkspaceId: UUIDToByteSlice(wsID),
			})
		}
		for _, agent := range remove {
			out.DeletedAgents = append(out.DeletedAgents, &proto.Agent{
				Id:          UUIDToByteSlice(agent.ID),
				Name:        agent.Name,
				WorkspaceId: UUIDToByteSlice(wsID),
			})
		}
	}

	// Delete workspace and agents if the workspace is deleted
	for wsID, oldWorkspace := range old {
		if _, exists := new[wsID]; !exists {
			out.DeletedWorkspaces = append(out.DeletedWorkspaces, &proto.Workspace{
				Id:     UUIDToByteSlice(wsID),
				Name:   oldWorkspace.WorkspaceName,
				Status: convertStatus(oldWorkspace.JobStatus, oldWorkspace.Transition),
			})
			for _, agent := range oldWorkspace.Agents {
				out.DeletedAgents = append(out.DeletedAgents, &proto.Agent{
					Id:          UUIDToByteSlice(agent.ID),
					Name:        agent.Name,
					WorkspaceId: UUIDToByteSlice(wsID),
				})
			}
		}
	}

	return out
}
