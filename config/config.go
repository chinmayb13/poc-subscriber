package config

import "github.com/spf13/viper"

type AppConfig struct {
	DB     DBConfig
	Kafka  KafkaCfg
	PubSub PubSubCfg
}

type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
}

type KafkaCfg struct {
	Topic   string `mapstructure:"KAFKA_TOPIC"`
	Broker  string `mapstructure:"KAFKA_BROKER"`
	GroupID string `mapstructure:"KAFKA_GROUP_ID"`
}

type PubSubCfg struct {
	ProjectID string `mapstructure:"PUBSUB_PROJECT_ID"`
	SubID     string `mapstructure:"PUBSUB_SUB_ID"`
}

func LoadConfig(path string) (config AppConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config.local")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	var db DBConfig
	var kafka KafkaCfg
	var pubsub PubSubCfg
	err = viper.Unmarshal(&db)
	if err != nil {
		return
	}
	err = viper.Unmarshal(&kafka)
	if err != nil {
		return
	}
	err = viper.Unmarshal(&pubsub)
	if err != nil {
		return
	}
	config.DB = db
	config.Kafka = kafka
	config.PubSub = pubsub
	return
}
