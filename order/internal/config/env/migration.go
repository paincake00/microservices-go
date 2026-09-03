package env

import "github.com/caarlos0/env/v11"

type migrationEnvConfig struct {
	MigrationDir string `env:"MIG_DIR,required"`
}

type MigrationConfig struct {
	raw migrationEnvConfig
}

func NewMigrationConfig() (*MigrationConfig, error) {
	var raw migrationEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}
	return &MigrationConfig{raw: raw}, nil
}

func (mc *MigrationConfig) MigrationDir() string {
	return mc.raw.MigrationDir
}
