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

For run tests:<br>
**go test -v -race ./...**
