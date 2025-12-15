INSERT INTO roles (roleId, roleName)
VALUES
    (1, 'user'),
    (2, 'admin')
ON CONFLICT (roleId) DO NOTHING;