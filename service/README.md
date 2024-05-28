# Back end for fast storage service

env info:

```json
{
    "DATABASE_USERNAME": "postgres",
    "DATABASE_PASSWORD": "postgres",
    "DATABASE_HOST": "localhost",
    "DATABASE_PORT": "5432",
    "DATABASE_NAME": "fast-storage-service",
    "DATABASE_MIGRATION": "true",
    "DATABASE_INITIALIZATION_DATA": "false",
    "KEYCLOAK_CLIENT_ID": "fast-storage-service",
    "KEYCLOAK_CLIENT_SECRET": "Wog91caQCn39MZwgRH9TxM9MZ0oqu4GG",
    "KEYCLOAK_API_URL": "https://localhost:8443",
    "KEYCLOAK_ADMIN_USERNAME": "admin",
    "KEYCLOAK_ADMIN_PASSWORD": "admin",
    "ACCOUNT_ACTIVE_HOST": "http://localhost:8080",
    "NFS_HOST": "/dev/nvme0n1p7",
    "MOUNT_FOLDER": "/mount",
    "OUTLOOK_USERNAME": "",
    "OUTLOOK_PASSWORD": ""
}
```

## Running on local with docker

Use docker incase you don't have Go installed on your machine

To run on local with docker. Running the following command:

```shell
docker compose build
docker compose up -d
```

To stop app from running. Running the following command:

```shell
docker compose down
```
