package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/SubhashBose/go-repro/lib"
)

func parseCommandline() (cfg lib.Config, err error) {
	var (
		mappingDefs, rewriteDefs string
		sslAllowInsecure         bool
		noLogging                bool
		showVersion              bool
		configFile               string
	)

	flag.StringVar(&mappingDefs, "mappings", "", "mapping definitions, format: local=remote,[local=remote,...]")
	flag.StringVar(&rewriteDefs, "rewrite", "", "comma-separated list of regexes indetifying routes whose response will be rewritten")
	flag.BoolVar(&sslAllowInsecure, "allow-insecure", false, "accept insecure upstream connections")
	flag.BoolVar(&noLogging, "no-logging", false, "disable logging via x-go-repro-log headers")
	flag.StringVar(&configFile, "config", "", "read YAML config from file (all other options are ignored)")
	flag.BoolVar(&showVersion, "version", false, "display version")

	flag.Usage = func() {
		fmt.Fprint(os.Stdout, "usage: go-repro [options]\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if showVersion {
		fmt.Printf("version: %s\n", lib.Version())
		os.Exit(0)
	}

	if configFile != "" {
		var yamlConfig YamlConfig

		yamlConfig, err = UnmarshalYamlConfigFile(configFile)

		if err == nil {
			cfg, err = yamlConfig.createReproConfig()
		}
	} else {
		cfg = lib.NewConfig()
		cfg.SetSSLAllowInsecure(sslAllowInsecure)
		cfg.SetNoLogging(noLogging)

		err = addMappings(mappingDefs, &cfg)

		if err == nil {
			err = addRewrites(rewriteDefs, &cfg)
		}
	}

	return
}

func addMappings(def string, cfg *lib.Config) (err error) {
	if def == "" {
		return
	}

	for _, definition := range strings.Split(def, ",") {
		parts := strings.Split(definition, "=")

		if len(parts) != 3 {
			err = errors.New(fmt.Sprintf("syntax error in mapping: %s", def))
		} else {
			err = cfg.AddMapping(parts[0], parts[1], parts[2])
		}

		if err != nil {
			return
		}
	}

	return
}

func addRewrites(def string, cfg *lib.Config) (err error) {
	if def == "" {
		return
	}

	for _, definition := range strings.Split(def, ",") {
		err = cfg.AddRewriteRoute(definition)

		if err != nil {
			return
		}
	}

	return
}

func main() {
	var err error

	cfg, err := parseCommandline()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	if cfg.CountMappings() == 0 {
		fmt.Println("nothing to do")
		return
	}

	r, err := lib.NewRepro(cfg)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err = <-r.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
