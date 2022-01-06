/* These scripts run the first time the Postgres container starts. The file is mounted inside the container. */

SELECT 'CREATE DATABASE nsx-api'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'nsx-api')\gexec
