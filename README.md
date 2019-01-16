# Message broked

Implementation of the server-client solution for storing key-value data.

### Prerequisites

The Docker should be pre-installed to run the project.

```
Give examples
```

### Installing

First of all you should clone this repo on your local machine.

```
git clone git@github.com:alexkylt/message_broker.git
```

Actuall after this you are able to up and run the project using Makefile from the repo folder:

```
make all
```

The above command will build and run the full project:
1) Docker container for a build purposer("go" entry point);
2) Run the "check" command(goimports, govet, golint);
3) Create docker network;
4) Build server part;
5) Build client part;
6) Build and run docker container for Postgres server;
7) Build and run docker container for Server;
8) Build and run docker container for Client;

After all steps from the "make all" command finished you will see the CLI in your terminal:




```
$ 
```

You need to specify the appropriate command along with key/value values you want to get, set or delete, e.g.:
```

$ GET key7
(key7, There is no key7 key in th db.)
$ SET key7 qwerty
$ GET key7
(key7, qwerty)
$ DEL key7
$ GET key7
(key7, There is no key7 key in th db.)
$ EXIT
```
The possible commands are:
  * **GET** {key_name}                - get a KV pair from DB/map  
  * **SET** {key_name} {key_value}    - insert/update a KV pair to DB/map
  * **DEL** {key_name}                - delete a KV pair from DB/map
  * **EXIT**                          - exit from the CLI
  
  Server uses Postgres database as storage by defualt
  
STORAGE_MODE ?= "db"
SERVER_PORT ?= 9090
SERVER_HOST ?= $(DOCKER_IMAGE_SERVER)