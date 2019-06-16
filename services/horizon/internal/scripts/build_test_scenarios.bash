#! /usr/bin/env bash
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GOTOP="$( cd "$DIR/../../../../../../../.." && pwd )"
PACKAGES=$(find $GOTOP/src/github.com/kinecosystem/go/services/horizon/internal/test/scenarios -iname '*.rb' -not -name '_common_accounts.rb')
# PACKAGES=$(find $GOTOP/src/github.com/kinecosystem/go/services/horizon/internal/test/scenarios -iname 'kahuna.rb')

go install github.com/kinecosystem/go/services/horizon



dropdb hayashi_scenarios --if-exists
createdb hayashi_scenarios

export STELLAR_CORE_DATABASE_URL="postgres://localhost/hayashi_scenarios?sslmode=disable"
export DATABASE_URL="postgres://localhost/horizon_scenarios?sslmode=disable"
export NETWORK_PASSPHRASE="Test SDF Network ; September 2015"
export STELLAR_CORE_URL="http://localhost:8080"
export SKIP_CURSOR_UPDATE="true"

# run all scenarios
for i in $PACKAGES; do
  
  CORE_SQL="${i%.rb}-core.sql"
  HORIZON_SQL="${i%.rb}-horizon.sql"
  echo "package $i"
  echo "core_sql: $CORE_SQL"
  scc -r $i --dump-root-db > $CORE_SQL --stellar-core-bin ./stellar-core
  echo "done with scc"
  # load the core scenario
  psql $STELLAR_CORE_DATABASE_URL < $CORE_SQL

  # recreate horizon dbs
  echo "recreating horizon db"
  dropdb horizon_scenarios --if-exists
  createdb horizon_scenarios

  # import the core data into horizon
  ./horizon db init
  ./horizon db init-asset-stats
  ./horizon db rebase

  # write horizon data to sql file
  pg_dump $DATABASE_URL \
    --clean --if-exists --no-owner --no-acl --inserts \
    | sed '/SET idle_in_transaction_session_timeout/d' \
    | sed '/SET row_security/d' \
    > $HORIZON_SQL
done

echo 'genrating bindata'
# commit new sql files to bindata
go generate github.com/kinecosystem/go/services/horizon/internal/test/scenarios
# go test github.com/kinecosystem/go/services/horizon/internal/ingest
