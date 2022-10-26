package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"math/rand"
	//	"bytes"
	"sync"
	"sync/atomic"
	"time"

	//	"6.824/labgob"
	"../labrpc"
)


//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in part 2D you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh, but set CommandValid to false for these
// other uses.
//
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int

	// For 2D:
	SnapshotValid bool
	Snapshot      []byte
	SnapshotTerm  int
	SnapshotIndex int
}

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	// 2A CODE
	role string
	latestHeartBeatTime int64
	term int
	termTicket map[int]bool
	termVoteNums int
	//termStarTime float64
	//termEndTime float64
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	var isleader bool
	// Your code here (2A).
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if rf.role == "leader"{
		isleader = true
	}
	return rf.term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}


//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}


//
// A service wants to switch to snapshot.  Only do so if Raft hasn't
// have more recent info since it communicate the snapshot on applyCh.
//
func (rf *Raft) CondInstallSnapshot(lastIncludedTerm int, lastIncludedIndex int, snapshot []byte) bool {

	// Your code here (2D).

	return true
}

// the service says it has created a snapshot that has
// all info up to and including index. this means the
// service no longer needs the log through (and including)
// that index. Raft should now trim its log as much as possible.
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// Your code here (2D).

}


//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	// 2A code
	Term int
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	// 2A CODE
	Option bool
}


type HeartBeatArgs struct {
	Term int
}

type HeartBeatReply struct {
}

func (rf *Raft) HeartBeat(args *HeartBeatArgs, reply *HeartBeatReply){
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if args.Term >= rf.term{
		rf.term = args.Term
		if rf.role != "follower"{  // 此时接收到的RPC任期如果大于等于自己的，就说明发送者已经当选leader。
			rf.role = "follower"
		}
		rf.latestHeartBeatTime = time.Now().UnixNano()
	}
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	// 2A CODE
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if rf.term > args.Term{
		return
	}
	if !rf.termTicket[args.Term]{
		reply.Option = true
		rf.termTicket[args.Term] = true
	}
}


func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

func (rf *Raft) sendHeartBeat(){
	for {
		rf.mu.Lock()
		if rf.role == "leader" {
			for i,_ := range rf.peers{
				go func(n int) {
					rf.mu.Lock()
					args := &HeartBeatArgs{Term: rf.term}
					rf.mu.Unlock()
					reply := &HeartBeatReply{}
					if n != rf.me{
						rf.peers[n].Call("Raft.HeartBeat", args, reply)
					}
				}(i)
			}
		}
		time.Sleep(100*time.Millisecond)
		rf.mu.Unlock()
	}
}



func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (2B).


	return index, term, isLeader
}


func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}


func (rf *Raft) startElection(){
	rf.mu.Lock()
	rf.term ++
	rf.role = "candidate"
	rf.termVoteNums = 0
	if !rf.termTicket[rf.term]{
		rf.termTicket[rf.term] = true
		rf.termVoteNums ++
	}
	rf.mu.Unlock()
	for i,_ := range rf.peers{
		if i != rf.me{
			go func(n int) {
				rf.mu.Lock()
				args := &RequestVoteArgs{Term: rf.term}
				rf.mu.Unlock()
				resp := &RequestVoteReply{}
				if ok:= rf.sendRequestVote(n, args, resp);ok{
					if resp.Option{
						rf.mu.Lock()
						rf.termVoteNums ++
						rf.mu.Unlock()
					}
				}
			}(i)
		}
	}
	go func() {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(500)
		time.Sleep(time.Millisecond*time.Duration(n + 500))
		rf.mu.Lock()
		if rf.role != "leader"{
			rf.role = "follower"
		}
		rf.mu.Unlock()
	}()

	for {
		rf.mu.Lock()
		if rf.role == "candidate"{
			if rf.termVoteNums * 2 > len(rf.peers){
				rf.role = "leader"
				rf.mu.Unlock()
				break
			}
		}else{
			rf.mu.Unlock()
			break
		}
		rf.mu.Unlock()
		time.Sleep(1*time.Millisecond)
	}
}


// The ticker go routine starts a new election if this peer hasn't received
// heartsbeats recently.
func (rf *Raft) ticker() {
	for rf.killed() == false {

		// Your code here to check if a leader election should
		// be started and to randomize sleeping time using
		// time.Sleep().
		// 2A CODE
		rf.mu.Lock()
		if (time.Now().UnixNano() - rf.latestHeartBeatTime) / (1000*1000) > 1000 && rf.role == "follower"{
			rf.mu.Unlock()
			rf.startElection()
		}else{
			rf.mu.Unlock()
			time.Sleep(50*time.Millisecond)
		}
	}
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.role = "follower"
	rf.termTicket = make(map[int]bool, 0)
	// Your initialization code here (2A, 2B, 2C).


	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	// start ticker goroutine to start elections
	go rf.ticker()  // 选举检测

	go rf.sendHeartBeat() // 心跳发送

	return rf
}
