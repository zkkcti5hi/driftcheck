// Package config provides loading and validation of driftcheck configuration
// files written in YAML.
//
// A typical configuration file looks like:
//
//	manifest:
//	  path: docker-compose.yml
//	  type: compose        # compose | kubernetes
//	output:
//	  format: table        # text | json | table
//	  quiet: false
//	docker:
//	  host: unix:///var/run/docker.sock
//
// Use [Load] to read a file from disk, or construct a [Config] directly
// when integrating driftcheck as a library.
package config
