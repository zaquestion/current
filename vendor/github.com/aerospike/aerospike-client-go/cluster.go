// Copyright 2013-2016 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/aerospike/aerospike-client-go/logger"

	. "github.com/aerospike/aerospike-client-go/types"
	. "github.com/aerospike/aerospike-client-go/types/atomic"
)

type partitionMap map[string][][]*Node

// String implements stringer interface for partitionMap
func (pm partitionMap) String() string {
	res := bytes.Buffer{}
	for ns, replicaArray := range pm {
		for i, nodeArray := range replicaArray {
			for j, node := range nodeArray {
				res.WriteString(ns)
				res.WriteString(",")
				res.WriteString(strconv.Itoa(i))
				res.WriteString(",")
				res.WriteString(strconv.Itoa(j))
				res.WriteString(",")
				if node != nil {
					res.WriteString(node.String())
				} else {
					res.WriteString("NIL")
				}
				res.WriteString("\n")
			}
		}
	}
	return res.String()
}

// Cluster encapsulates the aerospike cluster nodes and manages
// them.
type Cluster struct {
	// Initial host nodes specified by user.
	seeds *SyncVal //[]*Host

	// All aliases for all nodes in cluster.
	aliases *SyncVal //map[Host]*Node

	// Active nodes in cluster.
	nodes *SyncVal //[]*Node

	// Hints for best node for a partition
	partitionWriteMap atomic.Value //partitionMap

	// Random node index.
	nodeIndex *AtomicInt

	// Random partition replica index.
	replicaIndex *AtomicInt

	clientPolicy ClientPolicy

	mutex       sync.Mutex
	wgTend      sync.WaitGroup
	tendChannel chan struct{}
	closed      AtomicBool

	// Aerospike v3.6.0+
	supportsFloat, supportsBatchIndex, supportsReplicasAll, supportsGeo *AtomicBool
	requestProleReplicas                                                *AtomicBool

	// User name in UTF-8 encoded bytes.
	user string

	// Password in hashed format in bytes.
	password *SyncVal // []byte
}

// NewCluster generates a Cluster instance.
func NewCluster(policy *ClientPolicy, hosts []*Host) (*Cluster, error) {
	newCluster := &Cluster{
		clientPolicy: *policy,
		nodeIndex:    NewAtomicInt(0),
		replicaIndex: NewAtomicInt(0),
		tendChannel:  make(chan struct{}),

		seeds:   NewSyncVal(hosts),
		aliases: NewSyncVal(make(map[Host]*Node)),
		nodes:   NewSyncVal([]*Node{}),

		password: NewSyncVal(nil),

		supportsFloat:        NewAtomicBool(false),
		supportsBatchIndex:   NewAtomicBool(false),
		supportsReplicasAll:  NewAtomicBool(false),
		supportsGeo:          NewAtomicBool(false),
		requestProleReplicas: NewAtomicBool(policy.RequestProleReplicas),
	}

	newCluster.partitionWriteMap.Store(make(partitionMap))

	// setup auth info for cluster
	if policy.RequiresAuthentication() {
		newCluster.user = policy.User
		hashedPass, err := hashPassword(policy.Password)
		if err != nil {
			return nil, err
		}
		newCluster.password = NewSyncVal(hashedPass)
	}

	// try to seed connections for first use
	newCluster.waitTillStabilized(policy.FailIfNotConnected)

	// apply policy rules
	if policy.FailIfNotConnected && !newCluster.IsConnected() {
		return nil, fmt.Errorf("Failed to connect to host(s): %v. The network connection(s) to cluster nodes may have timed out, or the cluster may be in a state of flux.", hosts)
	}

	// start up cluster maintenance go routine
	newCluster.wgTend.Add(1)
	go newCluster.clusterBoss(&newCluster.clientPolicy)

	Logger.Debug("New cluster initialized and ready to be used...")
	return newCluster, nil
}

// String implements the stringer interface
func (clstr *Cluster) String() string {
	return fmt.Sprintf("%v", clstr.nodes)
}

// Maintains the cluster on intervals.
// All clean up code for cluster is here as well.
func (clstr *Cluster) clusterBoss(policy *ClientPolicy) {
	defer clstr.wgTend.Done()

	tendInterval := policy.TendInterval
	if tendInterval <= 10*time.Millisecond {
		tendInterval = 10 * time.Millisecond
	}

Loop:
	for {
		select {
		case <-clstr.tendChannel:
			// tend channel closed
			break Loop
		case <-time.After(tendInterval):
			if err := clstr.tend(policy.FailIfNotConnected); err != nil {
				Logger.Warn(err.Error())
			}
		}
	}

	// cleanup code goes here
	clstr.closed.Set(true)

	// close the nodes
	nodeArray := clstr.GetNodes()
	for _, node := range nodeArray {
		node.Close()
	}
}

// AddSeeds adds new hosts to the cluster.
// They will be added to the cluster on next tend call.
func (clstr *Cluster) AddSeeds(hosts []*Host) {
	clstr.seeds.Update(func(val interface{}) (interface{}, error) {
		seeds := val.([]*Host)
		seeds = append(seeds, hosts...)
		return seeds, nil
	})
}

// Updates cluster state
func (clstr *Cluster) tend(failIfNotConnected bool) error {
	nodes := clstr.GetNodes()
	nodes_before_tend := len(nodes)

	// All node additions/deletions are performed in tend goroutine.
	// If active nodes don't exist, seed cluster.
	if len(nodes) == 0 {
		Logger.Info("No connections available; seeding...")
		if _, err := clstr.seedNodes(failIfNotConnected); err != nil {
			return err
		}

		// refresh nodes list after seeding
		nodes = clstr.GetNodes()
	}

	// Refresh all known nodes.
	friendList := []*Host{}
	refreshCount := 0

	floatSupport := true
	batchIndexSupport := true
	replicasAllSupport := true
	geoSupport := true

	for _, node := range nodes {
		// Clear node reference counts.
		node.referenceCount.Set(0)

		if node.IsActive() {
			if friends, err := node.Refresh(friendList); err != nil {
				node.failures.IncrementAndGet()
				Logger.Warn("Node `%s` refresh failed: %s", node, err)
			} else {
				node.failures.Set(0)
				refreshCount++

				// make sure ALL nodes support float

				floatSupport = floatSupport && node.supportsFloat.Get()
				batchIndexSupport = batchIndexSupport && node.supportsBatchIndex.Get()
				replicasAllSupport = replicasAllSupport && node.supportsReplicasAll.Get()
				geoSupport = geoSupport && node.supportsGeo.Get()

				if friends != nil {
					friendList = append(friendList, friends...)
				}
			}
		}

	}

	if !floatSupport {
		Logger.Warn("Some cluster nodes do not support float type. Disabling native float support in the client library...")
	}

	// Disable prole requests if some nodes don't support it.
	if clstr.clientPolicy.RequestProleReplicas && !replicasAllSupport {
		Logger.Warn("Some nodes don't support 'replicas-all'. Will use 'replicas-master' for all nodes.")
	}

	// set the cluster supported features
	clstr.supportsFloat.Set(floatSupport)
	clstr.supportsBatchIndex.Set(batchIndexSupport)
	clstr.supportsReplicasAll.Set(replicasAllSupport)
	clstr.requestProleReplicas.Set(clstr.clientPolicy.RequestProleReplicas && replicasAllSupport)
	clstr.supportsGeo.Set(geoSupport)

	// Add nodes in a batch.
	if addList := clstr.findNodesToAdd(friendList); len(addList) > 0 {
		clstr.addNodes(addList)
	}

	// IMPORTANT: Remove must come after add to remove aliases
	// Handle nodes changes determined from refreshes.
	// Remove nodes in a batch.
	if removeList := clstr.findNodesToRemove(refreshCount); len(removeList) > 0 {
		clstr.removeNodes(removeList)
	}

	// only log if node count is changed
	if nodes_before_tend != len(clstr.GetNodes()) {
		Logger.Info("Tend finished. Live node count: %d", len(clstr.GetNodes()))
	}
	return nil
}

// Tend the cluster until it has stabilized and return control.
// This helps avoid initial database request timeout issues when
// a large number of threads are initiated at client startup.
//
// If the cluster has not stabilized by the timeout, return
// control as well.  Do not return an error since future
// database requests may still succeed.
func (clstr *Cluster) waitTillStabilized(failIfNotConnected bool) {
	count := -1

	doneCh := make(chan bool, 1)

	// will run until the cluster is stabilized
	go func() {
		for {
			if err := clstr.tend(failIfNotConnected); err != nil {
				Logger.Warn(err.Error())
			}

			// Check to see if cluster has changed since the last Tend().
			// If not, assume cluster has stabilized and return.
			if count == len(clstr.GetNodes()) {
				break
			}

			time.Sleep(time.Millisecond)

			count = len(clstr.GetNodes())
		}
		doneCh <- true
	}()

	// returns either on timeout or on cluster stabilization
	timeout := time.After(clstr.clientPolicy.Timeout)
	select {
	case <-timeout:
		return
	case <-doneCh:
		return
	}
}

func (clstr *Cluster) findAlias(alias *Host) *Node {
	res, _ := clstr.aliases.GetSyncedVia(func(val interface{}) (interface{}, error) {
		aliases := val.(map[Host]*Node)
		return aliases[*alias], nil
	})

	return res.(*Node)
}

func (clstr *Cluster) setPartitions(partMap partitionMap) {
	clstr.partitionWriteMap.Store(partMap)
}

func (clstr *Cluster) getPartitions() partitionMap {
	return clstr.partitionWriteMap.Load().(partitionMap)
}

func (clstr *Cluster) updatePartitions(conn *Connection, node *Node) (int, error) {
	parser, err := newPartitionParser(conn, node, clstr.getPartitions(), _PARTITIONS, clstr.requestProleReplicas.Get())
	if err != nil {
		return -1, err
	}

	// update partition write map
	if parser.isPartitionMapCopied() {
		clstr.setPartitions(parser.getPartitionMap())
	}

	return parser.getGeneration(), nil
}

// Adds seeds to the cluster
func (clstr *Cluster) seedNodes(failIfNotConnected bool) (bool, error) {
	// Must copy array reference for copy on write semantics to work.
	seedArrayIfc, _ := clstr.seeds.GetSyncedVia(func(val interface{}) (interface{}, error) {
		seeds := val.([]*Host)
		seeds_copy := make([]*Host, len(seeds))
		copy(seeds_copy, seeds)

		return seeds_copy, nil
	})
	seedArray := seedArrayIfc.([]*Host)

	errorList := []error{}

	Logger.Info("Seeding the cluster. Seeds count: %d", len(seedArray))

	// Add all nodes at once to avoid copying entire array multiple times.
	list := []*Node{}

	for _, seed := range seedArray {
		// Check if seed already exists in cluster.
		if clstr.findAlias(seed) != nil {
			continue
		}

		seedNodeValidator, err := newNodeValidator(clstr, seed, clstr.clientPolicy.Timeout)
		if err != nil {
			Logger.Warn("Seed %s failed: %s", seed.String(), err.Error())
			if failIfNotConnected {
				errorList = append(errorList, err)
			}
			continue
		}

		var nv *nodeValidator
		// Seed host may have multiple aliases in the case of round-robin dns configurations.
		for _, alias := range seedNodeValidator.aliases {

			if *alias == *seed {
				nv = seedNodeValidator
			} else {
				nv, err = newNodeValidator(clstr, alias, clstr.clientPolicy.Timeout)
				if err != nil {
					Logger.Warn("Seed %s failed: %s", seed.String(), err.Error())
					if failIfNotConnected {
						errorList = append(errorList, err)
					}
					continue
				}
			}

			if !clstr.findNodeName(list, nv.name) {
				node := clstr.createNode(nv)
				clstr.addAliases(node)
				list = append(list, node)
			}
		}
	}

	var err error
	if len(list) > 0 {
		clstr.addNodesCopy(list)
		return true, nil
	} else if failIfNotConnected {
		msg := "Failed to connect to host(s): "
		for _, err := range errorList {
			msg += "\n" + err.Error()
		}
		err = NewAerospikeError(INVALID_NODE_ERROR, msg)
	}

	return false, err
}

// Finds a node by name in a list of nodes
func (clstr *Cluster) findNodeName(list []*Node, name string) bool {
	for _, node := range list {
		if node.GetName() == name {
			return true
		}
	}
	return false
}

func (clstr *Cluster) addAlias(host *Host, node *Node) {
	if host != nil && node != nil {
		clstr.aliases.Update(func(val interface{}) (interface{}, error) {
			aliases := val.(map[Host]*Node)
			aliases[*host] = node
			return aliases, nil
		})
	}
}

func (clstr *Cluster) removeAlias(alias *Host) {
	if alias != nil {
		clstr.aliases.Update(func(val interface{}) (interface{}, error) {
			aliases := val.(map[Host]*Node)
			delete(aliases, *alias)
			return aliases, nil
		})
	}
}

func (clstr *Cluster) findNodesToAdd(hosts []*Host) []*Node {
	list := make([]*Node, 0, len(hosts))

	for _, host := range hosts {
		if nv, err := newNodeValidator(clstr, host, clstr.clientPolicy.Timeout); err != nil {
			Logger.Warn("Add node %s failed: %s", host.Name, err.Error())
		} else {
			node := clstr.findNodeByName(nv.name)
			// make sure node is not already in the list to add
			if node == nil {
				for _, n := range list {
					if n.GetName() == nv.name {
						node = n
						break
					}
				}
			}

			if node != nil {
				// Duplicate node name found.  This usually occurs when the server
				// services list contains both internal and external IP addresses
				// for the same node.  Add new host to list of alias filters
				// and do not add new node.
				node.referenceCount.IncrementAndGet()
				node.AddAlias(host)
				clstr.addAlias(host, node)
				continue
			}
			node = clstr.createNode(nv)
			list = append(list, node)
		}
	}
	return list
}

func (clstr *Cluster) createNode(nv *nodeValidator) *Node {
	return newNode(clstr, nv)
}

func (clstr *Cluster) findNodesToRemove(refreshCount int) []*Node {
	nodes := clstr.GetNodes()

	removeList := []*Node{}

	for _, node := range nodes {
		if !node.IsActive() {
			// Inactive nodes must be removed.
			removeList = append(removeList, node)
			continue
		}

		switch len(nodes) {
		case 1:
			// Single node clusters rely on whether it responded to info requests.
			if node.failures.Get() >= 5 {
				// Remove node.  Seeds will be tried in next cluster tend iteration.
				removeList = append(removeList, node)
			}

		case 2:
			// Two node clusters require at least one successful refresh before removing.
			if refreshCount == 1 && node.referenceCount.Get() == 0 && node.failures.Get() > 0 {
				// Node is not referenced nor did it respond.
				removeList = append(removeList, node)
			}

		default:
			// Multi-node clusters require two successful node refreshes before removing.
			if refreshCount >= 2 && node.referenceCount.Get() == 0 {
				// Node is not referenced by other nodes.
				// Check if node responded to info request.
				if node.failures.Get() == 0 {
					// Node is alive, but not referenced by other nodes.  Check if mapped.
					if !clstr.findNodeInPartitionMap(node) {
						// Node doesn't have any partitions mapped to it.
						// There is no point in keeping it in the cluster.
						removeList = append(removeList, node)
					}
				} else {
					// Node not responding. Remove it.
					removeList = append(removeList, node)
				}
			}
		}
	}
	return removeList
}

func (clstr *Cluster) findNodeInPartitionMap(filter *Node) bool {
	partitions := clstr.getPartitions()

	for _, replicaArray := range partitions {
		for _, nodeArray := range replicaArray {
			for _, node := range nodeArray {
				// Use reference equality for performance.
				if node == filter {
					return true
				}
			}
		}
	}
	return false
}

func (clstr *Cluster) addNodes(nodesToAdd []*Node) {
	// Add all nodes at once to avoid copying entire array multiple times.
	for _, node := range nodesToAdd {
		clstr.addAliases(node)
	}
	clstr.addNodesCopy(nodesToAdd)
}

func (clstr *Cluster) addAliases(node *Node) {
	// Add node's aliases to global alias set.
	// Aliases are only used in tend goroutine, so synchronization is not necessary.
	nodeAliases := node.GetAliases()

	clstr.aliases.Update(func(val interface{}) (interface{}, error) {
		aliases := val.(map[Host]*Node)

		for _, alias := range nodeAliases {
			aliases[*alias] = node
		}

		return aliases, nil
	})
}

func (clstr *Cluster) addNodesCopy(nodesToAdd []*Node) {
	clstr.nodes.Update(func(val interface{}) (interface{}, error) {
		nodes := val.([]*Node)
		nodes = append(nodes, nodesToAdd...)
		return nodes, nil
	})
}

func (clstr *Cluster) removeNodes(nodesToRemove []*Node) {
	// There is no need to delete nodes from partitionWriteMap because the nodes
	// have already been set to inactive. Further connection requests will result
	// in an exception and a different node will be tried.

	// Cleanup node resources.
	for _, node := range nodesToRemove {
		// Remove node's aliases from cluster alias set.
		// Aliases are only used in tend goroutine, so synchronization is not necessary.
		for _, alias := range node.GetAliases() {
			Logger.Debug("Removing alias ", alias)
			clstr.removeAlias(alias)
		}
		node.Close()
	}

	// Remove all nodes at once to avoid copying entire array multiple times.
	clstr.removeNodesCopy(nodesToRemove)
}

func (clstr *Cluster) setNodes(nodes []*Node) {
	// Replace nodes with copy.
	clstr.nodes.Set(nodes)
}

func (clstr *Cluster) removeNodesCopy(nodesToRemove []*Node) {
	// Create temporary nodes array.
	// Since nodes are only marked for deletion using node references in the nodes array,
	// and the tend goroutine is the only goroutine modifying nodes, we are guaranteed that nodes
	// in nodesToRemove exist.  Therefore, we know the final array size.
	nodes := clstr.GetNodes()
	nodeArray := []*Node{}
	count := 0

	// Add nodes that are not in remove list.
	for _, node := range nodes {
		if clstr.nodeExists(node, nodesToRemove) {
			Logger.Info("Removed node `%s`", node)
		} else {
			nodeArray = append(nodeArray, node)
			count++
		}
	}

	clstr.setNodes(nodeArray)
}

func (clstr *Cluster) nodeExists(search *Node, nodeList []*Node) bool {
	for _, node := range nodeList {
		if node.Equals(search) {
			return true
		}
	}
	return false
}

// IsConnected returns true if cluster has nodes and is not already closed.
func (clstr *Cluster) IsConnected() bool {
	// Must copy array reference for copy on write semantics to work.
	nodeArray := clstr.GetNodes()
	return (len(nodeArray) > 0) && !clstr.closed.Get()
}

func (clstr *Cluster) getReadNode(partition *Partition, replica ReplicaPolicy) (*Node, error) {
	switch replica {
	case MASTER:
		return clstr.getMasterNode(partition)
	case MASTER_PROLES:
		return clstr.getMasterProleNode(partition)
	default:
		// includes case RANDOM:
		return clstr.GetRandomNode()
	}
}

func (clstr *Cluster) getMasterNode(partition *Partition) (*Node, error) {
	pmap := clstr.getPartitions()
	replicaArray := pmap[partition.Namespace]

	if replicaArray != nil {
		node := replicaArray[0][partition.PartitionId]
		if node != nil && node.IsActive() {
			return node, nil
		}
	}

	return clstr.GetRandomNode()
}

func (clstr *Cluster) getMasterProleNode(partition *Partition) (*Node, error) {
	pmap := clstr.getPartitions()
	replicaArray := pmap[partition.Namespace]

	if replicaArray != nil {
		for range replicaArray {
			index := int(math.Abs(float64(clstr.replicaIndex.IncrementAndGet() % len(replicaArray))))
			node := replicaArray[index][partition.PartitionId]
			if node != nil && node.IsActive() {
				return node, nil
			}
		}
	}

	return clstr.GetRandomNode()
}

// GetRandomNode returns a random node on the cluster
func (clstr *Cluster) GetRandomNode() (*Node, error) {
	// Must copy array reference for copy on write semantics to work.
	nodeArray := clstr.GetNodes()
	length := len(nodeArray)
	for i := 0; i < length; i++ {
		// Must handle concurrency with other non-tending goroutines, so nodeIndex is consistent.
		index := int(math.Abs(float64(clstr.nodeIndex.IncrementAndGet() % length)))
		node := nodeArray[index]

		if node.IsActive() {
			// Logger.Debug("Node `%s` is active. index=%d", node, index)
			return node, nil
		}
	}
	return nil, NewAerospikeError(INVALID_NODE_ERROR)
}

// GetNodes returns a list of all nodes in the cluster
func (clstr *Cluster) GetNodes() []*Node {
	// Must copy array reference for copy on write semantics to work.
	return clstr.nodes.Get().([]*Node)
}

// GetNodeByName finds a node by name and returns an
// error if the node is not found.
func (clstr *Cluster) GetNodeByName(nodeName string) (*Node, error) {
	node := clstr.findNodeByName(nodeName)

	if node == nil {
		return nil, NewAerospikeError(INVALID_NODE_ERROR)
	}
	return node, nil
}

func (clstr *Cluster) findNodeByName(nodeName string) *Node {
	// Must copy array reference for copy on write semantics to work.
	for _, node := range clstr.GetNodes() {
		if node.GetName() == nodeName {
			return node
		}
	}
	return nil
}

// Close closes all cached connections to the cluster nodes
// and stops the tend goroutine.
func (clstr *Cluster) Close() {
	if !clstr.closed.Get() {
		// send close signal to maintenance channel
		close(clstr.tendChannel)

		// wait until tend is over
		clstr.wgTend.Wait()
	}
}

// MigrationInProgress determines if any node in the cluster
// is participating in a data migration
func (clstr *Cluster) MigrationInProgress(timeout time.Duration) (res bool, err error) {
	if timeout <= 0 {
		timeout = _DEFAULT_TIMEOUT
	}

	done := make(chan bool, 1)

	go func() {
		// this function is guaranteed to return after _DEFAULT_TIMEOUT
		nodes := clstr.GetNodes()
		for _, node := range nodes {
			if node.IsActive() {
				if res, err = node.MigrationInProgress(); res || err != nil {
					done <- true
					return
				}
			}
		}

		res, err = false, nil
		done <- false
	}()

	dealine := time.After(timeout)
	for {
		select {
		case <-dealine:
			return false, NewAerospikeError(TIMEOUT)
		case <-done:
			return res, err
		}
	}
}

// WaitUntillMigrationIsFinished will block until all
// migration operations in the cluster all finished.
func (clstr *Cluster) WaitUntillMigrationIsFinished(timeout time.Duration) (err error) {
	if timeout <= 0 {
		timeout = _NO_TIMEOUT
	}
	done := make(chan error, 1)

	go func() {
		// this function is guaranteed to return after timeout
		// no go routines will be leaked
		for {
			if res, err := clstr.MigrationInProgress(timeout); err != nil || !res {
				done <- err
				return
			}
		}
	}()

	dealine := time.After(timeout)
	select {
	case <-dealine:
		return NewAerospikeError(TIMEOUT)
	case err = <-done:
		return err
	}
}

// Password returns the password that is currently used with the cluster.
func (clstr *Cluster) Password() (res []byte) {
	pass := clstr.password.Get()
	if pass != nil {
		return pass.([]byte)
	}
	return nil
}

func (clstr *Cluster) changePassword(user string, password string, hash []byte) {
	// change password ONLY if the user is the same
	if clstr.user == user {
		clstr.clientPolicy.Password = password
		clstr.password.Set(hash)
	}
}

// ClientPolicy returns the client policy that is currently used with the cluster.
func (clstr *Cluster) ClientPolicy() (res ClientPolicy) {
	return clstr.clientPolicy
}
