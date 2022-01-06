/* These scripts run the first time the Postgres container starts. The file is mounted inside the container. */

SELECT 'CREATE DATABASE nsx_api'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'nsx_api')\gexec
