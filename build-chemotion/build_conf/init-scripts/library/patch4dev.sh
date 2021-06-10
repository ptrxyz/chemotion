#!/bin/bash
# RUNAS: user

if [[ -d "${APP_DIR}/app/config" ]]; then 
	cd "${APP_DIR}/app/config"
	for i in *.example; do
		mv ${i} ${i%.*};
	done
fi

if [[ -f "${APP_DIR}/app/.env.development" ]]; then
	mv ${APP_DIR}/app/.env.development ${APP_DIR}/app/.env
	export RAILS_ENV=development
	export SECRET_KEY_BASE=b9a89bc2ad33

	echo -e "SECRET_KEY_BASE='${SECRET_KEY_BASE}'\n"\
	"DB_NAME='${DB_NAME}'\n"\
	"DB_ROLE='${DB_ROLE}'\n"\
	"DB_PW='${DB_PW}'\n"\
	"DB_HOST='${DB_HOST}'\n"\
	"DB_PORT=${DB_PORT}\n"\
	"RAILS_ENV=${RAILS_ENV}\n" | sed 's/^ //g' >>${APP_DIR}/app/.env

	mkdir -p ${APP_DIR}/app/config

	echo -e "${RAILS_ENV}:\n"\
	"   adapter: postgresql\n"\
	"   encoding: unicode\n"\
	"   database: <%=ENV['DB_NAME']%>\n"\
	"   pool: 5\n"\
	"   username: <%=ENV['DB_ROLE']%>\n"\
	"   password: <%=ENV['DB_PW']%>\n"\
	"   host: <%=ENV['DB_HOST']%>\n"\
	"   port: <%=ENV['DB_PORT']%>\n" | sed 's/^ //g' >${APP_DIR}/app/config/database.yml

	echo -e "${RAILS_ENV}:\n"\
	"   secret_key_base: <%=ENV['SECRET_KEY_BASE']%>\n" | sed 's/^ //g' >${APP_DIR}/app/config/secrets.yml
fi






