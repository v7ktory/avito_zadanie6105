CREATE TABLE tender (
  id UUID PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  description TEXT NOT NULL,
  service_type VARCHAR(50) NOT NULL,
  status VARCHAR(50) NOT NULL,
  version INT NOT NULL DEFAULT 1,
  organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
  creator_username VARCHAR(50) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);