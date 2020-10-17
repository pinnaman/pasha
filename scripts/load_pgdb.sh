#!/bin/bash

cd scripts
docker exec -it datadoor psql -h localhost -U postgres -f /scripts/db_setup.sql
#cd cmd/news/pg_data
echo "loading news database from backup...."
docker exec -i datadoor psql -U postgres -d pasha_ddoor_db < pg_backup.sql