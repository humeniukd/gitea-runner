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

		// If cache.external_secret is set, run with bearer auth on the cache
		// API and expose the /_internal control-plane so remote runners can
		// pre-register their per-job ACTIONS_RUNTIME_TOKENs against this
		// server. Otherwise stay in legacy openMode (unauthenticated, intended
		// for trusted networks only).
		secret := cfg.Cache.ExternalSecret
		openMode := secret == ""
		cacheHandler, err := artifactcache.StartHandler(
			dir,
			host,
			port,
			openMode,
			secret,
			log.StandardLogger().WithField("module", "cache_request"),
		)
		if err != nil {
			return err
		}

		if openMode {
			log.Warn("cache-server runs without request authentication; expose it only on a trusted network or set cache.external_secret on both ends")
		} else {
			log.Info("cache-server: bearer auth + control-plane enabled via cache.external_secret")
		}
		log.Infof("cache server is listening on %v", cacheHandler.ExternalURL())

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		return nil
	}
}
