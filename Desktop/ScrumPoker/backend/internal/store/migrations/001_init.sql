CREATE TABLE IF NOT EXISTS rooms (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    code CHAR(6) NOT NULL UNIQUE,
    host_member_id BIGINT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'lobby',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS members (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    room_id BIGINT NOT NULL,
    display_name VARCHAR(60) NOT NULL,
    is_host BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    session_token VARCHAR(128) NOT NULL UNIQUE,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMP NULL,
    CONSTRAINT fk_members_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

ALTER TABLE rooms
    ADD CONSTRAINT fk_rooms_host_member FOREIGN KEY (host_member_id) REFERENCES members(id) ON DELETE SET NULL;

CREATE TABLE IF NOT EXISTS tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    room_id BIGINT NOT NULL,
    title VARCHAR(180) NOT NULL,
    status VARCHAR(32) NOT NULL,
    average_label VARCHAR(64) NULL,
    average_numeric DECIMAL(10,2) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revealed_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    CONSTRAINT fk_tasks_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

CREATE INDEX idx_tasks_room_status ON tasks(room_id, status);

CREATE TABLE IF NOT EXISTS votes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id BIGINT NOT NULL,
    member_id BIGINT NOT NULL,
    value VARCHAR(32) NOT NULL,
    numeric_value DECIMAL(10,2) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uniq_task_member_vote (task_id, member_id),
    CONSTRAINT fk_votes_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_votes_member FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE
);
