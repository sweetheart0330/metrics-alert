CREATE TABLE metrics (
                         id SERIAL PRIMARY KEY,
                         metric_id VARCHAR(255) NOT NULL,
                         metric_type VARCHAR(10) NOT NULL CHECK (metric_type IN ('counter', 'gauge')),
                         timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                         delta_value BIGINT,
                         gauge_value DOUBLE PRECISION,
                         hash VARCHAR(64),
                         created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX ux_metrics_metric_id ON metrics(metric_id);
CREATE INDEX idx_metrics_id_timestamp ON metrics (metric_id, timestamp DESC);
CREATE INDEX idx_metrics_type_timestamp ON metrics (metric_type, timestamp DESC);
CREATE INDEX idx_metrics_timestamp_brin ON metrics USING BRIN (timestamp);
