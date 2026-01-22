-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items (
    id CHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(1000),
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_items_title ON items(title);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS item_properties (
    id CHAR(36) PRIMARY KEY,
    item_id CHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    value VARCHAR(1000) NOT NULL,
    CONSTRAINT fk_item_properties_item
        FOREIGN KEY (item_id) REFERENCES items(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_item_properties_item_id ON item_properties(item_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_item_properties_name ON item_properties(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS item_properties;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
