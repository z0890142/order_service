CREATE TABLE IF NOT EXISTS Patient (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    age INTEGER NOT NULL,
    gender SMALLINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_patient_timestamp
BEFORE UPDATE ON Patient
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();


INSERT INTO Patient (name, age, gender) VALUES
    ('John', 35, 1),
    ('Alice', 28, 0),
    ('Michael', 40, 1),
    ('Emily', 32, 0),
    ('David', 45, 1);
