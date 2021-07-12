#!/bin/bash

export mustash=${APP_DIR}/app/mustash/mo
export templates=${APP_DIR}/app/mustash/templates

${mustash} ${templates}/env_template.mo 		> ${APP_DIR}/app/.env
${mustash} ${templates}/database_template.mo		> ${APP_DIR}/app/config/database.yml    
${mustash} ${templates}/datacollectors_template.mo	> ${APP_DIR}/app/config/datacollectors.yml
#${mustash} ${templates}/editors_template.mo 		> ${APP_DIR}/app/config/editors.yml
#${mustash} ${templates}/inference_template.mo 		> ${APP_DIR}/app/config/inference.yml
${mustash} ${templates}/secrets_template.mo 		> ${APP_DIR}/app/config/secrets.yml
#${mustash} ${templates}/spectra_template.mo 		> ${APP_DIR}/app/config/spectra.yml
${mustash} ${templates}/storage_template.mo 		> ${APP_DIR}/app/config/storage.yml
${mustash} ${templates}/user_props_template.mo 		> ${APP_DIR}/app/config/user_props.yml

