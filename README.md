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
- thinking of when we have to update the fingers tables of nodes (when 2 nodes enter the network : risk of update the table twice) -> search in the litterature and discuss solution

TOCHANGE : 
- Finger table should be :
example for node "05" 3bit space

| Index | Key | Successor | 
| ----- | --- | --------- |
|     0 |  06 |        07 | 
|     1 |  08 |        08 | 
|     2 |  06 |        07 | 

