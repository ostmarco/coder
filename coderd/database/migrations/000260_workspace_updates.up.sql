CREATE TYPE agent_id_name_pair AS (
	id uuid,
	name text
);

CREATE FUNCTION new_workspace_notify() RETURNS trigger
	LANGUAGE plpgsql
AS $$
DECLARE
BEGIN
	-- Write to the notification channel `new_workspace:owner_id`
	-- with the workspace id as the payload.
	PERFORM pg_notify('new_workspace:' || NEW.owner_id, NEW.id);
	RETURN NEW;
END;
$$;

CREATE TRIGGER new_workspace_notify
	AFTER INSERT ON workspaces
	FOR EACH ROW
EXECUTE FUNCTION new_workspace_notify();


CREATE FUNCTION new_agent_notify() RETURNS trigger
	LANGUAGE plpgsql
AS $$
DECLARE
	workspace_owner_id uuid;
BEGIN
	SELECT workspaces.owner_id
	INTO workspace_owner_id
	FROM
		workspaces
	WHERE
		workspaces.id = (
			SELECT
				workspace_id
			FROM
				workspace_builds
			WHERE
				workspace_builds.job_id = (
					SELECT
						job_id
					FROM
						workspace_resources
					WHERE
						workspace_resources.id = (
							SELECT
								resource_id
							FROM
								workspace_agents
							WHERE
								workspace_agents.id = NEW.id
						)
				)
		);
	-- Agents might not belong to a workspace (template imports)
	IF workspace_owner_id IS NOT NULL THEN
		-- Write to the notification channel `new_agent:workspace_owner_id`
		PERFORM pg_notify('new_agent:' || workspace_owner_id, '');
	END IF;
	RETURN NEW;
END;
$$;


CREATE TRIGGER new_agent_notify
	AFTER INSERT ON workspace_agents
	FOR EACH ROW
EXECUTE FUNCTION new_agent_notify();

