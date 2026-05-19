SELECT 'CREATE DATABASE shorten'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'shorten');