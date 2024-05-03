-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Cars (
    id SERIAL PRIMARY KEY,
                                     regNum VARCHAR(255) Not NULL,
                                     mark VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    year INTEGER NULL,
    owner_name VARCHAR(255) NOT NULL,
    owner_surname VARCHAR(255) NOT NULL,
    owner_patronymic VARCHAR(255) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Cars;
-- +goose StatementEnd
