package handlers

import (
	"log"

	"F22/config"

	"github.com/globalsign/mgo"
	"github.com/gorilla/sessions"

)

// Provider holds application wide variables
type Provider struct {
	log     *log.Logger
	cfg     *config.Config
	db      *mgo.Session
	session *sessions.CookieStore
}

func NewProvider(log *log.Logger, cfg *config.Config, db *mgo.Session, session *sessions.CookieStore) *Provider {
	return &Provider{
		log:     log,
		cfg:     cfg,
		db:      db,
		session: session,
	}
}

func (p *Provider) Logger() *log.Logger { return p.log }

func (p *Provider) DB() *mgo.Session { return p.db }

func (p *Provider) Session() *sessions.CookieStore { return p.session }

func (p *Provider) Config() *config.Config { return p.cfg }

