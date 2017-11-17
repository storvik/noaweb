package noaweb

import (
	"flag"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Instance variable that holds the instance
var noawebinst *Instance

// Instance is the main instance, contains all configurations.
type Instance struct {
	router        *mux.Router
	subrouters    Subrouters
	configfile    string
	workingdir    string
	commandFlags  CommandFlags
	dbConnections dbConnections
	mailguns      mailgunConfigs
	Port          int               `json:"Port"`
	Hostname      string            `json:"Hostname"`
	TLS           bool              `json:"TLS"`
	TLSCache      string            `json:"TLSCache"`
	TLSEmail      string            `json:"TLSEmail"`
	AssetsDir     string            `json:"AssetsDir"`
	UserConfig    map[string]string `json:"UserConfig"`
}

type CommandFlags struct {
	devMode       bool
	migrateAction string
}

// NewInstance creates new stweb instance which listen on 8080.
func NewInstance() *Instance {
	dir, _ := os.Getwd()
	instance := &Instance{
		router:        mux.NewRouter().StrictSlash(true),
		subrouters:    make(Subrouters),
		workingdir:    dir,
		dbConnections: make(dbConnections),
		mailguns:      make(mailgunConfigs),
		Hostname:      "",
		Port:          8080,
		TLS:           false,
		TLSCache:      "",
		TLSEmail:      "",
		AssetsDir:     dir,
	}

	instance.parseFlags()

	noawebinst = instance
	return instance
}

// NewInstanceParse parses given config file.
func NewInstanceParse(configFile string) *Instance {
	dir, _ := os.Getwd()
	instance := &Instance{
		router:        mux.NewRouter().StrictSlash(true),
		subrouters:    make(Subrouters),
		configfile:    configFile,
		workingdir:    dir,
		dbConnections: make(dbConnections),
		mailguns:      make(mailgunConfigs),
		AssetsDir:     dir,
	}

	instance.parseFlags()
	parseConfig(configFile, instance)

	noawebinst = instance
	return instance
}

func (i *Instance) parseFlags() {
	dev := flag.Bool("dev", false, "run in development mode, without tls")
	str := flag.String("migrate", "up", "migrate database (up|down|refresh|false)")

	flag.Parse()
	i.commandFlags.devMode = *dev
	i.commandFlags.migrateAction = *str
}

// HandlerFunc is the type for handler functions.
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Subrouters type
type Subrouters map[string]*mux.Router
