-- +migrate Up
CREATE TABLE cargoes
(
    id         VARCHAR(50) PRIMARY KEY,
    vessel_id  VARCHAR(50) NOT NULL,
    items      JSONB       NOT NULL,
    status     VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX cargoes_idx_vessel_id ON cargoes (vessel_id);
CREATE INDEX cargoes_idx_status ON cargoes (status);
-- +migrate Down
DROP INDEX IF EXISTS cargoes_idx_vessel_id;
DROP INDEX IF EXISTS cargoes_idx_status;
DROP TABLE IF EXISTS cargoes;
