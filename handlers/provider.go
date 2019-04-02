package handlers

import (
	"log"

	"F22/config"

	"github.com/globalsign/mgo"
)

// Provider holds application wide variables
type Provider struct {
	log     *log.Logger
	cfg     *config.Config
	db      *mgo.Session
}

func NewProvider(log *log.Logger, cfg *config.Config, db *mgo.Session) *Provider {
	return &Provider{
		log:     log,
		cfg:     cfg,
		db:      db,
	}
}

func (p *Provider) Logger() *log.Logger { return p.log }

func (p *Provider) DB() *mgo.Session { return p.db }

func (p *Provider) Config() *config.Config { return p.cfg }

