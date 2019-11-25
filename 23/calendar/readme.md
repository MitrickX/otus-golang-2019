Simple calendar with http interface <br>
In project used cobra (with viper) and zap <br>

**Usage:** <br>

Simple call (will use default config): <br>
For run http service <br>
**calendar run http** <br>
For run grpc service <br>
**calendar run grpc** <br>

If you want set own custom config: <br>
**calendar --config <path_to_config> run [http|grpc]** <br>

In config 'db' is DB connection settings (DB is PostgreSQL)<br>
If you want off DB storage just don't have 'db' key in config <br><br>

For run tests:<br>
**go test -v -race ./...**

If you want test sql.Storage, you must CREATE test database first and connection settings must be declared in test.yaml config within 'db' key AND tests must be running from directory internal/storage/sql<br>
**cd internal/storage/sql && go test -v -race . && cd ../../../**
