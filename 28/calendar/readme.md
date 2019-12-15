Simple calendar with http interface <br>
In project used cobra (with viper) and zap <br>

**Usage:** <br>

Simple call (will use default config): <br>

For run http service <br>
**calendar http** <br>

For run grpc service <br>
**calendar grpc** <br>

For run notification scheduler <br>
**calendar scheduler** <br>

For run notification sender <br>
**calendar sender** <br>

If you want set own custom config: <br>
**calendar --config <path_to_config> [http|grpc|scheduler|sender]** <br>

In config 'db' is DB connection settings (DB is PostgreSQL)<br>
If you want off DB storage just don't have 'db' key in config <br><br>

For run tests:<br>
**go test -v -race ./...**

If you want test sql.Storage, you must CREATE test database first and connection settings must be declared in test.yaml config within 'db' key AND tests must be running from directory internal/storage/sql<br>
**cd internal/storage/sql && go test -v -race . && cd ../../../**
