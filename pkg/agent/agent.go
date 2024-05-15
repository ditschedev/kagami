package agent

import (
	"github.com/ditschedev/kagami/pkg/config"
	"github.com/ditschedev/kagami/pkg/log"
	"github.com/ditschedev/kagami/pkg/mirror"
	"github.com/spf13/viper"
)

func Start() {
	var repositories []config.MirrorConfig

	err := viper.UnmarshalKey("repositories", &repositories)
	if err != nil {
		log.Fatal("Could not load config file")
		return
	}

	for _, repo := range repositories {
		m := mirror.New(&repo)
		m.Mirror()
	}
}
