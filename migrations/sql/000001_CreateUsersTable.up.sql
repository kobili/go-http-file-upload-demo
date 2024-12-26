CREATE TABLE users (
    user_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(140) UNIQUE
);

CREATE TABLE profile_pictures (
    id BIGSERIAL PRIMARY KEY,
    file_path VARCHAR(250) UNIQUE,
    user_id uuid REFERENCES users (user_id)
)
