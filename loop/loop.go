package loop

import (
	"sync"

	"github.com/Artic67/architecture-lab4/parser"
)

// Structure of trusted handler
type trustedHandler struct {
	loop *EventLoop
}

// Post Method of trusted handler
func (th *trustedHandler) Post(cmd parser.Command) {
	th.loop.commands.push(cmd)
}

// Structure of commands queue
type cmdQueue struct {
	commands    []parser.Command
	hasElements sync.Cond
}

// IsEmpty Method of commands queue
func (queue *cmdQueue) isEmpty() bool {
	queue.hasElements.L.Lock()
	defer queue.hasElements.L.Unlock()
	return len(queue.commands) == 0
}

// Push Method of commands queue
func (queue *cmdQueue) push(cmd parser.Command) {
	queue.hasElements.L.Lock()
	queue.commands = append(queue.commands, cmd)
	queue.hasElements.L.Unlock()
	queue.hasElements.Broadcast()
}

// Pull Method of commands queue
func (queue *cmdQueue) pull() parser.Command {
	queue.hasElements.L.Lock()
	for len(queue.commands) == 0 {
		queue.hasElements.Wait()
	}
	defer queue.hasElements.L.Unlock()
	cmd := queue.commands[0]
	queue.commands[0] = nil
	queue.commands = queue.commands[1:]
	return cmd
}

// Structure of Event Loop
type EventLoop struct {
	commands *cmdQueue
	stopCond sync.Cond

	stopLocker    sync.Mutex
	stopRequested bool
	stopped       bool
	wasStarted    bool
}

// Initialization Method of Event Loop
func (loop *EventLoop) init() {
	loop.commands = &cmdQueue{
		hasElements: *sync.NewCond(&sync.Mutex{}),
	}
	loop.stopCond = *sync.NewCond(&sync.Mutex{})
	loop.stopLocker = sync.Mutex{}
	loop.wasStarted = true
	loop.stopRequested = false
	loop.stopped = false
}

// Dispose Method of Event Loop
func (loop *EventLoop) dispose() {
	loop.commands = nil
}

// Post Method of Event Loop
func (loop *EventLoop) Post(cmd parser.Command) {
	loop.verifyRunning()
	if loop.isStopRequested() || loop.isStopped() {
		return
	}
	loop.commands.push(cmd)
}

// Await Finish Method of Event Loop
func (loop *EventLoop) AwaitClose() {
	loop.verifyRunning()
	if loop.isStopped() {
		return
	}
	if !loop.isStopRequested() {
		loop.stopLocker.Lock()
		loop.stopRequested = true
		loop.stopLocker.Unlock()
	}
	loop.stopCond.L.Lock()
	loop.stopCond.Wait()
	loop.stopCond.L.Unlock()
}

// Start Method of Event Loop
func (loop *EventLoop) Start() {
	loop.init()
	go loop.listen()
}

// Method of Event Loop to check is stop requested
func (loop *EventLoop) isStopRequested() bool {
	loop.stopLocker.Lock()
	defer loop.stopLocker.Unlock()
	return loop.stopRequested
}

// Method of Event Loop to check is stopped
func (loop *EventLoop) isStopped() bool {
	loop.stopLocker.Lock()
	defer loop.stopLocker.Unlock()
	return loop.stopped
}

// Stop Method of Event Loop
func (loop *EventLoop) stop() {
	loop.stopLocker.Lock()
	loop.stopped = true
	loop.stopRequested = false
	loop.wasStarted = true // it means loop was started at once
	loop.stopLocker.Unlock()
	loop.stopCond.Broadcast()
	loop.dispose()
}

// Listen Method of Event Loop
func (loop *EventLoop) listen() {
	for !loop.isStopRequested() || !loop.commands.isEmpty() {
		cmd := loop.commands.pull()
		cmd.Execute(&trustedHandler{loop})
	}
	loop.stop()
}

// Method of Event Loop to verify running
func (loop *EventLoop) verifyRunning() {
	if !loop.wasStarted && !loop.stopped {
		panic("Unable to perform an action. Loop was not started")
	}
}
