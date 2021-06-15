#!/bin/bash

${APP_DIR}/app/mustash/mo ${APP_DIR}/app/mustash/templates/env_template.mo 	        >>   ${APP_DIR}/app/.env
${APP_DIR}/app/mustash/mo ${APP_DIR}/app/mustash/templates/database_template.mo		>>   ${APP_DIR}/app/config/database.yml
${APP_DIR}/app/mustash/mo ${APP_DIR}/app/mustash/templates/datacollectors_template.mo	>>   ${APP_DIR}/app/config/datacollectors.yml
#  ${APP_DIR}/app/mustash/mo ${APP_DIR}/app/mustash/templates/editors_template.mo 	>>   ${APP_DIR}/app/config/editors.yml
#  ${APP_DIR}/app/mustash/mo ${APP_DIR}/app/mustash/templates/inference_template.mo 	>>   ${APP_DIR}/app/config/inference.yml
${APP_DIR}/app/mustash/mo ${APP_DIR}/app/mustash/templates/secrets_template.mo 		>>   ${APP_DIR}/app/config/secrets.yml
#  ${APP_DIR}/app/mustash/mo ${APP_DIR}/app/mustash/templates/spectra_template.mo 	>>   ${APP_DIR}/app/config/spectra.yml
${APP_DIR}/app/mustash/mo ${APP_DIR}/app/mustash/templates/storage_template.mo 		>>   ${APP_DIR}/app/config/storage.yml
${APP_DIR}/app/mustash/mo ${APP_DIR}/app/mustash/templates/user_props_template.mo 	>>   ${APP_DIR}/app/config/user_props.yml


