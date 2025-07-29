# Running the app

Generate new JWT keys for the project:
```bash
sh scripts/keys.sh
```

Compile the SQL queries using sqlc:
```bash
sh scripts/sqlc.sh
```

Spin up the project:
```bash
DATABASE_PROVIDER=postgres \
POSTGRES_HOST=postgres \
POSTGRES_PORT=5432 \
POSTGRES_DB=postgres \
POSTGRES_USER=postgres \
POSTGRES_PASSWORD=postgres \
docker compose up --build
```

Run the Atlas migrations:
```bash
sh scripts/atlas.sh
```

Access the project at:
```bash
# Client
http://localhost:3000
# Admin
http://localhost:3001
```

For Grafana Monitoring, check the README.md in `/grafana` folder.  
For Kubernetes Deployment + Monitoring, check the README.md in `/kube` folder.
