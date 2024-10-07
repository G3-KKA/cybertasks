CREATE TABLE tasks(
    id UUID PRIMARY KEY,
    header VARCHAR(255) NOT NULL,
    description TEXT,
    created_at timestamp (2) with time zone NOT NULL,
    status boolean NOT NULL
); 