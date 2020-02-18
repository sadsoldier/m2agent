
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUFFIXES = .go .c
OBJEXT= none

EXTRA_PROGRAMS = builder
builder_SOURCES = builder.go

GOFLAGS = -ldflags="-s -w"

builder$(EXEEXT): $(builder_SOURCES)
	$(GO) build -o builder$(EXEEXT) $(builder_SOURCES)

sbin_PROGRAMS = a2srv a2pwd

a2srv_SOURCES = a2srv.go
a2pwd_SOURCES = a2pwd.go

EXTRA_a2pwd_SOURCES = \
	config/config.go \
	server/user-model/user_model.go

EXTRA_a2srv_SOURCES = \
	config/config.go \
	daemon/daemon.go \
	tools/file.go \
	bundle/public.go \
	\
	server/server.go \
	server/db-controller/db_controller.go \
	server/db-models/pgdb-model/pgdb_model.go \
	server/dump-controller/dump_controller.go \
	server/dump-models/pgdump-model/pgdump_model.go \
	server/dbuser-controller/dbuser_controller.go \
	server/dbuser-models/pguser-model/pguser_model.go \
	server/status-controller/status_controller.go \
	server/user-controller/user_controller.go \
	server/user-model/user_model.go


bundle/public.go: builder$(EXEEXT) public/index.html
	test -d bundle || $(MKDIR_P) bundle
	./builder$(EXEEXT) --package=bundle -o bundle/public.go public/

a2srv$(EXEEXT): $(a2srv_SOURCES) $(EXTRA_a2srv_SOURCES)
	$(GO) build $(GOFLAGS) -o a2srv$(EXEEXT) $(a2srv_SOURCES)

a2pwd$(EXEEXT): $(a2pwd_SOURCES) $(EXTRA_a2pwd_SOURCES)
	$(GO) build $(GOFLAGS) -o a2pwd$(EXEEXT) $(a2pwd_SOURCES)

EXTRA_DIST = README.md public/ \
	go.mod \
	go.sum \
	\
	front/src/ \
	front/public/ \
	\
	front/.gitignore \
	front/babel.config.js \
	front/package-lock.json \
	front/package.json \
	front/postcss.config.js \
	front/README.md \
	\
	front/webpack.config.devel.js \
	front/webpack.config.prod.js

clean-local:
	rm -f builder$(EXEEXT)
	rm -rf autom4te.cache

install-data-local:
	test -z $(DESTDIR)$(APP_LIBDIR) || $(MKDIR_P) $(DESTDIR)$(APP_LIBDIR)
#	cp -r public $(DESTDIR)$(APP_LIBDIR)
	test -z $(DESTDIR)$(APP_CONFDIR) || $(MKDIR_P) $(DESTDIR)$(APP_CONFDIR)
	test -z $(DESTDIR)$(APP_LOGDIR) || $(MKDIR_P) $(DESTDIR)$(APP_LOGDIR)
	test -z $(DESTDIR)$(APP_RUNDIR) || $(MKDIR_P) $(DESTDIR)$(APP_RUNDIR)
	test -z $(DESTDIR)$(APP_DATABASEDIR) || $(MKDIR_P) $(DESTDIR)$(APP_DATABASEDIR)
#EOF
