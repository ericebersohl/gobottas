package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ericebersohl/gobottas/discussion"
	"log"
)

// Contains Gobottas functions and data
type Registry struct {
	Interceptors  map[CommandType]Interceptor // all built-in interceptors
	Executors     map[CommandType]Executor    // all built-in executors
	DirPath       string                      // path to local data
	CommandPrefix uint8                       // character that precedes all Gobottas commands
	DQueue        *discussion.Queue           // the data structure that holds discussion queue data
}

// Function to call all Interceptors on a message
func (r *Registry) Intercept(msg *Message) error {
	for _, i := range r.Interceptors {
		err := i(msg)
		if err != nil {
			log.Printf("Registry.Intercept: %v", err)
			return err
		}
	}
	return nil
}

// Calls the Executor to which the Registry points for the Message CommandType
func (r *Registry) Execute(msg *Message, s *discordgo.Session) error {

	// check if the executor is in the map; if not do nothing
	if e, ok := r.Executors[msg.CommandType]; ok {
		err := e(s, r, msg)
		if err != nil {
			log.Printf("Registry.Execute: %v", err)
			return err
		}
	}

	return nil
}

// Declaration of built-in Interceptors
var BuiltinInterceptors = map[CommandType]Interceptor{
	Help:  HelpInterceptor,
	Meme:  MemeInterceptor,
	Queue: QueueInterceptor,
}

// Declaration of built-in Executors
var BuiltinExecutors = map[CommandType]Executor{
	Help:  HelpExecutor,
	Meme:  MemeExecutor,
	Queue: QueueExecutor,
}

type RegistryOpt func(*Registry)

// Get a new registry, optionally pass RegistryOpts for custom functionality
func NewRegistry(opts ...RegistryOpt) *Registry {
	r := Registry{
		Interceptors:  BuiltinInterceptors,
		Executors:     BuiltinExecutors,
		DirPath:       "/tmp",
		CommandPrefix: '&',
		DQueue:        discussion.NewQueue(),
	}

	for _, o := range opts {
		o(&r)
	}

	return &r
}

// Add a custom Interceptor to the Registry
func WithCustomInterceptor(i Interceptor, t CommandType) RegistryOpt {
	return func(r *Registry) {
		r.Interceptors[t] = i
	}
}

// Add a custom Executor to the Registry
func WithCustomExecutor(e Executor, t CommandType) RegistryOpt {
	return func(r *Registry) {
		r.Executors[t] = e
	}
}

// Set the path to local files
func WithDirPath(s string) RegistryOpt {
	return func(r *Registry) {
		r.DirPath = s
	}
}

// Set the CommandPrefix
func WithCommandPrefix(c uint8) RegistryOpt {
	return func(r *Registry) {
		r.CommandPrefix = c
	}
}

// Add an existing or externally defined queue to the registry
func WithDiscussionQueue(q *discussion.Queue) RegistryOpt {
	return func(r *Registry) {
		r.DQueue = q
	}
}
