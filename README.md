GrandRepositorySky
==================

Grand Repository in the Sky : POC implementation of _Distributed Hash Table_ (DHT) based on _Chord_ protocol (see [This paper ](http://pdos.csail.mit.edu/papers/chord:sigcomm01/chord_sigcomm.pdf))

The final goal of the project is to deploy a node inside a docker container

Installation :
----
Assuming that you already have Go install, with a `GOPATH` variable set : 
```
go get github.com/tgermain/grandRepositorySky 
```

Run the project : 
----

Local demo : 
```bash
#launch the first node with default parameters (localhost:4321)
go run main.go 
#launch the second node on (localhost:4322) and connect to the first one 
go run main.go -p 4322 -d 4321 
```

There are web page interface at `localhost:4321` and `localhost:4322`.

Docker demo:
Assuming you already have docker install :
```bash
# get the latest version of a node container
docker pull tgermain/repo_sky:latest
# run the orchestrator
go run webMaster/server.go
```
The web interface to launch new container is server on `localhost:8080`.


Documentation :
----
see `doc/d7024-lab-report.pdf`

:construction: Current objective :construction: : **5** :godmode:
---------------------

### 6 objectives : 
1. Chord DHT playground 
	- organize node in ring
	- lookup function
	- finger table
2. :wrench:Fully working chord DHT
	- separate nodes and makes them communicate
	- Manage data when nodes leave or enter the network
3. Replication
	- Add data to nodes
4. Web service
	- _POST_ : new key-value pair
	- _GET_ : find value of a key
	- _PUT_ : update value of a key
	- _DELETE_ : delete a key-value pair
5. Virtualization 
	- package a node as a Docker container
	- create a web service to setup a Sky network easily
		- creating openStack VM
		- deploying docker node on it
		- visualization
		- REST API
6. :lock:(optional) Security and Encryption



## 1. Chord DHT playground 

to run the playground test : 
`go test dhtPlayground_test.go dhtPlayground.go -v`

###DONE:
- Extract the next real node from the fingers table (fingers[0] became a new atribut of DHTNode)
- implement the function to create finger table (the primitive are already done)
- Add a predecessor node pointer in DHTNode
- Add a graphviz visualization (see func ``gimmeGraph`` and test ``TestGraph`` in test3bit_test.go)

### Fingers table calculation 

2->5

New node 4 !

2->4->5

When a node enter the ring, we initialize its fingers table (to be sure that its succesor is *(5)*) and update the fingers table of its predecessor *(2)*. 


### Ring visualization
Call the method ``gimmeGraph`` on any node, export the result to a file and process with your best graphviz, I recommend circo. ``circo graph.gv -Tsvg -o viz`` for svg output

## 2. Network communication
### Ring stabilization DONE
To avoid the depreciation of all the fingers table after some new nodes join the ring, there's a mecanism to update the fingers table each 5min (pifometric value).
**Goroutine** ?

Done : 
- Handle when a node quit the ring (not gracefuly)
- semaphore for critical part (predecessor, successor maybe others ...)
- Make a library to send/receive messages 		
- Do logs (with timestamps, colors ?) 		
- Modify the existing code to make stuffs work 
- What lookup can/should/actually return 
- Ring stabilization 
- heartBeat 
- All methods which return a `*DHTNode` must return an **id** and we have to perform a `lookup(id)`

###launch parameters DONE
- IP, PORT 
- (IP, PORT of an existing node)
- (ID)


###message format
- origin (first emiter of the message)
	- IP
	- Port
- Destination
	- IP
	- Port
- type of operation :
	- LOOKUP / LOOKUPRESPONSE
	- UPDATESUCCESSOR
	- UPDATEPREDECESSOR
	- UPDATEFINGERTABLE
	- PRINTRING
	- JOINRING
	- AREYOUALIVE / IAMALIVE
	- GETSUCCESORE / GETSUCCESORERESPONSE


Work in progress (Al):
- Get requests and show brutal infos :neckbeard: :rage4: :boar: yeah BRUTAL ! :)
- Web client
- Web controls


Done :
- Multiple instances server
- Graph interface



## 3. Data & Replication
facebook duplicate data of one node on 2 other nodes : enough !

- Data place on node A is replicated to A.successor and A.predecessor
- Special case when there is:
  - only 1 node -> data are not replicated
  - only 2 node -> data are only replicated to the predecessor

### new messages type
- DataReplication

DONE : 
- methods in node.go to add, retreive, modify and delete a (key-value) data
- new messages : 
	- POSTDATA
	- GETDATA / GETDATARESPONSE
	- MODIFYDATA
	- DELETEDATA


## 5. Virtualization

Problem with openStach API -> we did not use it.

###Docker usage
See the docker file to see how the image is built.

**Important** : due to the dynamic connection of the node on the same virtual machine (it's not a use case), we do not isolate the network interface of containers. They all use the *host network interface* via `NetworkMode: "host"` when starting a container.