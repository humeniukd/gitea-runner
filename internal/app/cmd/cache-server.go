// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"gitea.com/gitea/act_runner/act/artifactcache"
	"gitea.com/gitea/act_runner/internal/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type cacheServerArgs struct {
	Dir  string
	Host string
	Port uint16
}

func runCacheServer(configFile *string, cacheArgs *cacheServerArgs) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadDefault(*configFile)
		if err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}

		initLogging(cfg)

		var (
			dir  = cfg.Cache.Dir
			host = cfg.Cache.Host
			port = cfg.Cache.Port
		)

		// cacheArgs has higher priority
		if cacheArgs.Dir != "" {
			dir = cacheArgs.Dir
		}
		if cacheArgs.Host != "" {
			host = cacheArgs.Host
		}
		if cacheArgs.Port != 0 {
			port = cacheArgs.Port
		}

		// openMode=true: this standalone server has no channel to learn
		// per-job ACTIONS_RUNTIME_TOKENs from Gitea, so bearer auth is
		// skipped and repo scoping collapses. Operators must expose it
		// only on a trusted network.
		cacheHandler, err := artifactcache.StartHandler(
			dir,
			host,
			port,
			true,
			log.StandardLogger().WithField("module", "cache_request"),
		)
		if err != nil {
			return err
		}

		log.Warn("cache-server runs without request authentication; expose it only on a trusted network")
		log.Infof("cache server is listening on %v", cacheHandler.ExternalURL())

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		return nil
	}
}
