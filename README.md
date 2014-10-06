GrandRepositorySky
==================

Grand Repository in the Sky : POC implementation of _Distributed Hash Table_ (DHT) based on _Chord_ protocol (see [This paper ](http://pdos.csail.mit.edu/papers/chord:sigcomm01/chord_sigcomm.pdf))

:construction: Current objective :construction: : **2** :suspect:
---------------------

### 6 objectives : 
1. Chord DHT playground 
	- organize node in ring
	- lookup function
	- finger table
2. :wrench:Fully working chord DHT
	- separate nodes and makes them communicate
3. Replication
	- Add data to nodes
	- Manage data when nodes leave or enter the network
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

###TODO : 
- Upgrade the lookup function to take advantage of the fingers table (***+test***)

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
Call the method ``gimmeGraph`` on any node, export the result to a file and process with your best graphviz, I recommend circo. ``circo graph.gv -Tsvg -o viz`` for svg² output

## 2. Network communication
### Ring stabilization
To avoid the depreciation of all the fingers table after some new nodes join the ring, there's a mecanism to update the fingers table each 5min (pifometric value).
**Goroutine** ?

SPRINT : 
- Do a sample project to understand how to send/receive and parse messages in GO
- Make a library to send/receive messages 		Tim
- Do logs (with timestamps, colors ?) 		Tim
- Make tests
- Modify the existing code to make stuffs work
- Handle when a node quit the ring (not gracefuly)
- What lookup can/should/actually return

Work In progress (Tim):
- logs : almost done, easy part
- send/receive message library : 
	- marshal/unmarshal -> easy with stdlib
	- use of channel : the idea is to add in the message struct a channel for the anwer
	- critical ressource, to be protected:
		- successor
		- predecessor
		- data
		

TODO : 
- All methods which return a `*DHTNode` must return an **id** and we have to perform a `lookup(id)`

###launch parameters
- IP, PORT
- (IP, PORT of an existing node)
- (ID)
- (Data already owned by the node (dump of json), to simulate the case when a node already holding data quit temporarily the ring)


###CLI interface
Useful ? **Objective 4** provide a web service, we just have to do more things than just thoses (and make a page to visualize and act on the node)
abilities : 
- lookup
- updateFingersTable
- printInformation (fingerTable)
- printGraphViz
- areYouAlive

###message format
- origin (first emiter of the message)
	- IP
	- Port
- Destination
	- IP
	- Port
- type of operation :
	- IAmNewHere
	- update successor/predecessor
	- updateFingerTable
	- lookup / finded
	- gimmeInfo / IAmTheNSA
	- AreYouAlive / IAmAlive


Work in progress (Al):
- Get requests and show brutal infos
- Web client
- Web controls


TODO :
- Multiple instances server
- Graph interface



## 3. Data & Replication
facebook duplicate data of one node on 2 other nodes : enough ?

###Choice of underlying database
- Redis
- mongodb 

### new messages type
- DataReplication
