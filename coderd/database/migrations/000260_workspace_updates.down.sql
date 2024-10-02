DROP TYPE agent_id_name_pair;

DROP TRIGGER IF EXISTS new_workspace_notify ON workspaces;
DROP FUNCTION IF EXISTS new_workspace_notify;

DROP TRIGGER IF EXISTS new_agent_notify ON workspace_agents;
DROP FUNCTION IF EXISTS new_agent_notify;
