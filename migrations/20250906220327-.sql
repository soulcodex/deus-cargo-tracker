-- +migrate Up
CREATE TABLE cargoes_tracking
(
    id            VARCHAR(50) PRIMARY KEY,
    cargo_id      VARCHAR(50)  NOT NULL,
    entry_type    VARCHAR(100) NOT NULL,
    status_before VARCHAR(50) DEFAULT NULL,
    status_after  VARCHAR(50) DEFAULT NULL,
    created_at    TIMESTAMPTZ  NOT NULL
);
CREATE INDEX cargoes_tracking_idx_cargo_id ON cargoes_tracking (cargo_id);
CREATE INDEX cargoes_tracking_idx_entry_type ON cargoes_tracking (entry_type);
-- +migrate Down
DROP INDEX IF EXISTS cargoes_tracking_idx_cargo_id;
DROP INDEX IF EXISTS cargoes_tracking_idx_entry_type;
DROP TABLE IF EXISTS cargoes_tracking;
