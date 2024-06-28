CREATE TABLE IF NOT EXISTS images (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES auth.users,
    status INT NOT NULL DEFAULT 1,
    prompt TEXT NOT NUll,
    deleted BOOLEAN NOT NULL DEFAULT 'false',
    image_location TEXT,
    batch_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);