# Recovery

## Database

SQLite file path: `~/Library/Application Support/route-deck/route-deck.db` (`os.UserConfigDir()/route-deck/route-deck.db`).

GORM `AutoMigrate` adds tables and columns. It does not drop columns and does not reliably rename them.

Before a release that changes GORM models in a breaking way:

1. Copy the SQLite file (and `-wal` / `-shm` if present) to a timestamped backup.
2. Run the new build.
3. If the app cannot open the DB, restore the backup and stay on the previous build.

MUST NOT delete the live database to “fix” a migration.

## ProxyBridge left dirty

`kill -9` on the CLI can leave extension rules in place (CLI docs). Recovery: start the CLI and stop it with SIGTERM, or quit the ProxyBridge app. RouteDeck MUST NOT report `Active` while that state is unknown.

## Secrets

Passwords are in Keychain, not in the SQLite backup. Restoring a database does not restore deleted Keychain items. Export of profiles omits secrets unless a future task explicitly adds an opt-in.

## Diagnostics

A diagnostic bundle (later task) includes versions, state, redacted settings, and recent errors. It MUST NOT include passwords or raw profile JSON that contains passwords.
