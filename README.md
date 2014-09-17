GrandRepositorySky
==================

Grand Repository in the Sky : POC implementation of _Distributed Hash Table_ (DHT) based on _Chord_ protocol (see [This paper ](http://pdos.csail.mit.edu/papers/chord:sigcomm01/chord_sigcomm.pdf))

:construction: Current objective :construction: : **1** :sweat_smile:
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

TODO : 
- implement the function to create finger table (the primitive are already done)
- Upgrade the lookup function to take advantage of the fingers table (+test)
- Add a predecessor node pointer in DHTNode
- Extract the next real node from the fingers table (fingers[0] became a new atribut of DHTNode)
- Finger table is now :
example for node "05" 3bit space

| Index | Key | Successor | 
| ----- | --- | --------- |
|     0 |  06 |        07 | 
|     1 |  08 |        08 | 
|     2 |  06 |        07 | 
make sure that all the test and existing code take account of this modification

### Finger table calculation 

2->5

New node 4 !

2->4->5

When a node enter the ring, we initialize its fingers table (to be sure that its succesor is *(5)*) and update the fingers table of its predecessor *(2)*. 

### Ring stabilization
To avoid the depreciation of all the fingers table after some new nodes join the ring, there's a mecanism to update the fingers table each 5min (pifometric value).
**Goroutine** ?

## 2. Network communication

SPRINT : 
- Do a sample project to understand how to send/receive and parse messages in GO
- Make a library to send/receive messages
- Do logs (with timestamps, colors)
- Make tests
- Modify the existing code to make stuffs work
- Handle when a node quit the ring (not gracefuly)


TODO : 
- All methods which return a `*DHTNode` must return an **id** and we have to perform a `lookup(id)`