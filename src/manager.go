/**
 * manager.go - manages servers
 *
 * @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
 */
package main

import (
	"./config"
	"./logging"
	"./server"
)

/**
 * Run and control server
 */
func Start(cfg config.Config) {

	log := logging.For("manager")
	log.Info("Starting up...")

	// Go through config and start servers for each server
	for name, s := range cfg.Servers {
		config := prepareConfig(name, s, cfg.Defaults)
		server := server.New(name, config)
		go server.Start()
	}

	log.Info("Start up complete")

	// block forever
	<-(chan string)(nil)
}

/**
 */
func prepareConfig(name string, server config.Server, defaults config.ConnectionOptions) config.Server {

	log := logging.For("manager")

	/* ----- Prerequisites ----- */

	if server.Discovery == nil {
		log.Fatal("No [.discovery] section for", name, "specified")
	}

	if server.Healthcheck == nil {
		log.Fatal("No [.healthcheck] section for", name, "specified. Will allow it later :-)")
	}

	/* ----- Connections params and overrides ----- */

	/* Balance */
	switch server.Balance {
	case
		"weight",
		"leastconn",
		"roundrobin",
		"iphash":
		server.Balance = server.Balance
	case "":
		server.Balance = "weight"
	default:
		log.Fatal("Not supported balance type", server.Balance)
	}

	/* TODO: Still need to decide how to get rid of this */

	if defaults.MaxConnections == nil {
		defaults.MaxConnections = new(int)
	}
	if server.MaxConnections == nil {
		server.MaxConnections = defaults.MaxConnections
	}

	if defaults.ClientIdleTimeout == nil {
		defaults.ClientIdleTimeout = &config.MyDuration{}
	}
	if server.ClientIdleTimeout == nil {
		server.ClientIdleTimeout = defaults.ClientIdleTimeout
	}

	if defaults.BackendIdleTimeout == nil {
		defaults.BackendIdleTimeout = &config.MyDuration{}
	}
	if server.BackendIdleTimeout == nil {
		server.BackendIdleTimeout = defaults.BackendIdleTimeout
	}

	if defaults.BackendConnectionTimeout == nil {
		defaults.BackendConnectionTimeout = &config.MyDuration{}
	}
	if server.BackendConnectionTimeout == nil {
		server.BackendConnectionTimeout = defaults.BackendConnectionTimeout
	}

	return server
}
