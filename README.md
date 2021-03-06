# Message broker

Implementation of the server-client solution for storing key-value data.

### Prerequisites

The Docker should be pre-installed to run the project.

### Installing

First of all you should clone this repo on your local machine.

```
git clone git@github.com:alexkylt/message_broker.git
```

You need run the below command to up and run the project using Makefile from the repo folder:

```
make all
```

The above command will build and run the full project:
1) Docker container for a build purpose("go" entry point);
2) Run the "check" rule(goimports, govet, golint);
3) Create docker network;
4) Build server part;
5) Build client part;
6) Build and run docker container for Postgres server;
7) Build and run docker container for Server;
8) Build and run docker container for Client;

After all steps from the ```make all``` command finished you will see the CLI in your terminal:




```
$ 
```

You need to specify the appropriate command along with key/value values you want to get, set or delete.

The possible commands are:
  * **GET** {key_name}                - get a KV pair from DB/map  
  * **SET** {key_name} {key_value}    - insert/update a KV pair to DB/map
  * **DEL** {key_name}                - delete a KV pair from DB/map
  * **EXIT** - exit from the CLI

CLI usage sample:
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

### Configuration

Server uses Postgres database as storage by defualt. The defaul server port is 9090.  If you want to use Go map as a storage instead of Postgres database or to use other Server port - you may specify appropriate values in the Makefile:

```
SERVER_HOST_NAME ?= server
STORAGE_MODE ?= "map"
SERVER_PORT ?= 9090
```

or you may pass these parameters to ```make``` command:
```
make all SERVER_HOST_NAME=messagebroker SERVER_PORT=7788 STORAGE_MODE=map
```

NOTE: You should not specify the ```SERVER_HOST_NAME``` as ```localhost``` or ```127.0.0.1```. Server/Client/Database containers use docker network for communication.

### Using

If you already install project on your local machine using ```make all``` command you can rebuild the project. You should run the ```make all``` command again to rebuild the project if you want to use other storage or port.

After exiting CLI you  can run the project by the following command: ```make run_client``` - this command will run the Server too.
