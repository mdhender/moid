#!/bin/bash
##############################################################################
# This script rebuilds the localhost database.
##############################################################################
# Confirm that we're running from the root of the repository.
[ -f go.mod -a -f bin/rebuild_db.sh ] || {
  echo error: must run from the root of the repository
  exit 2
}
##############################################################################
#
db=testdata/localhost/moid.db
##############################################################################
# In theory, the DDL contains all the commands needed to rebuild the database
# without removing the file and starting from scratch. In a better world,
# we'd be able to do this without having to drop the tables.
for ddl in \
  internal/generators/sqlc/202502110915_initial.sql \
; do
  echo " info: running '${ddl}..."
  [ -f "${ddl}" ] || {
    echo "error: '${ddl}' does not exist"
    exit 2
  }
  sqlite3 "${db}" ".read ${ddl}" || {
    echo "error: unable to rebuild the database"
    exit 2
  }
done
##############################################################################
# Dump the schema, as a courtesy.
sqlite3 "${db}" .schema || {
  echo "error: unable to dump the schema"
  exit 2
}
##############################################################################
#
echo " info: rebuilt '${db}"
exit 0
