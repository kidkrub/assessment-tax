package config

import (
	"os"
	"strconv"
)

type env func(key string) string

type cfg struct {
	getEnv env
}

type Server struct {
	Hostname string
	PORT     int
}

type Database struct {
	DatabaseUrl string
}

const (
	cHostname    = "HOSTNAME"
	cPort        = "PORT"
	cDatabaseUrl = "DATABASE_URL"
)

func New() *cfg {
	return &cfg{getEnv: os.Getenv}
}

func (c *cfg) Server() Server {
	return Server{c.envString(cHostname, ""), c.envInt(cPort, 1323)}
}

func (c *cfg) Database() Database {
	return Database{c.envString(cDatabaseUrl, "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable")}
}

func (c *cfg) envString(key, defaultValue string) string {
	value := c.getEnv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *cfg) envInt(key string, defaultValue int) int {
	value := c.getEnv(key)

	val, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return val
}
