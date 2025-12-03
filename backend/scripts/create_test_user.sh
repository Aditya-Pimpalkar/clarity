#!/bin/bash
echo "🔧 Creating test user in llm_observability database..."

docker exec -i llm-obs-clickhouse clickhouse-client --database=llm_observability --multiquery << 'SQL'
-- Create test organization (ignore if exists)
INSERT INTO organizations (id, name, created_at, updated_at) 
SELECT 'org-test-123', 'Test Organization', now(), now()
WHERE NOT EXISTS (SELECT 1 FROM organizations WHERE id = 'org-test-123');

-- Create test user (ignore if exists) - WITHOUT password_hash
INSERT INTO users (
    id, 
    organization_id,
    email, 
    name, 
    role, 
    created_at,
    updated_at
)
SELECT 
    'user-test-123',
    'org-test-123',
    'diksha@example.com',
    'Diksha Sahare',
    'admin',
    now(),
    now()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'diksha@example.com');

-- Create test project (ignore if exists)
INSERT INTO projects (
    id,
    name,
    organization_id,
    created_at,
    updated_at
)
SELECT 
    'proj-test-123',
    'Test Project',
    'org-test-123',
    now(),
    now()
WHERE NOT EXISTS (SELECT 1 FROM projects WHERE id = 'proj-test-123');

SELECT '✅ Setup complete!' as status;
SELECT count() as user_count FROM users;
SELECT count() as org_count FROM organizations;
SELECT count() as project_count FROM projects;

-- Show the created user
SELECT id, email, name, organization_id, role FROM users WHERE email = 'diksha@example.com' FORMAT Pretty;
SQL

echo ""
echo "✅ Test data created!"
echo ""
echo "📧 Login credentials:"
echo "   Email: diksha@example.com"
echo "   Password: anything (mock auth enabled)"
echo ""
