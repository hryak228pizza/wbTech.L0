-- +goose Up
-- +goose StatementBegin
alter table orders add column simplex TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table orders drop column simplex;
-- +goose StatementEnd
