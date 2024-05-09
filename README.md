# architect

Architect is a tool for creating web applications on Golang. It consists of two parts: bin for generating frame for new microservice and lib for running this microservice.

## FUTURE

### General
Architect is the tool that is responsible for the fundament and architecture of application. 
In the best case architect allow to construct artcitecture of application on the fly. 

### Goals 
- Simplify mecroservice development.
- Speed up devlopmnet.
- Unify architecture of all Golang applications.
- Move all infrastructure for app so that developer can focus only on logic.
- Instill golang writing style.

### Details
Command scathces:
- add - base for other add commands. [x]
    - client - adds clinet for other service + connection base. In best case - generate form proto client code. [x]
    - postgres - adds code for connection to postgres. Can consider other DBs. 
    - repository - adds base repo code. [x]
    - manager - adds base business-unit code + interfaces. [x]
    - sub_manager - adds base business-unit code + interfaces. [x]
    - kafka (consumer/producer)kafka message broker - adds code for . Need to think more, mb some supestructure above kafka/rabbiMQ.
    - cron - adds base code for cron.
    - proto-service - adds proto service. [x]
- show - base for other show commands.
    - structure - shows architecture of application. [x]

## TODO
- Move from env to yaml config ([viper](https://github.com/rakyll/statik)).
- Move swagger from apps to architect ([statik](https://github.com/rakyll/statik))
- Write own protoc generate service layer and to delete dirty hack for services. 
- Make update command to apdate architect as lib.
- Make upgarade command to apdate architect as a binary.
- Move from latest to specific verions of deps.
- Keep records every CHANGELOG (separate feature merge request before release).
- Generate config from values.yaml files. 
- Update existing files, not recreate it. 
- Docker (with add commands or separate).
- Kafka/RabbitMQ, clients, repos - to adapter layer.
- Interactive mode.
- Documentation in README.md.
- Look in FUTURE section.
- Change work with proto deps - now it's a little bit complicated (buf is a choice, but not in Russia).
- Check build of application.
- Add real gitlab CI/CD with building and pushing to docker-registry.
- Add integration tests.
- Support for windows & macOS
- Client add is very dirty. Need to rethink approach. Need checks and clean code. 
- Short flag for command architecture.
- Add help for makefiles. [x]
- Add work with metrics.
- Fix business_error. Now it throws Internal errors mo matter if higher was thwon other gRPC error. 
- Add work with flags (ports & etc).
- Add work with middleware. User can define middleware in architect.NewApp() func. Mb some default.
- Fix work with CORS. Now it works only for GET methods.