package config

type MirrorConfig struct {
	Name      string `mapstructure:"name"`
	RemoteUri string `mapstructure:"remote_uri"`
	MirrorUri string `mapstructure:"mirror_uri"`
}
