#!/bin/bash

checkFolderExists(){
    if [[ ! -d $1 ]]; then
        echo "    Folder $1 does not exist. Please create it."
        return 1
    else
        echo "    Found folder $1."
        return 0
    fi
}

checkSymlinkExists(){
    if [[ ! -d $1 ]]; then
        echo "    Symlink $1 does not exist. Please create it."
        return 1
    else
        echo "    Found symlink $1."
        return 0
    fi
}

checkFolderIsWritable(){
    if [[ ! -w $1 ]]; then
        echo "    Folder $1 has no write permission. Please grant it."
        return 1
    else
        echo "    Folder $1 has write permission."
        return 0
    fi
}

generateNewSecrets(){
    SECRET_KEY="$(bundle exec rake secret)"
    cat <<EOF > config/secrets.yml
production:
  secret_key_base: $SECRET_KEY
EOF
}

echo "---- Init script ----"

echo "Checks for the file system:"
# check accessibility of shared folders
if ! mount | grep 'on /shared' 2>&1 1>/dev/null; then
    echo "    The shared folder is not correctly connected as volume. Please make sure that a folder shared/ is available next to the docker-compose.yml file."
    exit 1
else
    echo "    Found folder /shared."
fi

if ! checkFolderExists "/shared/eln"        ; then exit 1; fi
if ! checkFolderExists "/shared/eln/config" ; then exit 1; fi
if ! checkFolderExists "/shared/eln/log"    ; then exit 1; fi
if ! checkFolderExists "/shared/eln/public" ; then exit 1; fi
if ! checkFolderExists "/shared/eln/tmp"    ; then exit 1; fi
if ! checkFolderExists "/shared/eln/uploads"; then exit 1; fi

if ! checkSymlinkExists "/chemotion/app/log"                      ; then exit 1; fi
if ! checkSymlinkExists "/chemotion/app/tmp"                      ; then exit 1; fi
if ! checkSymlinkExists "/chemotion/app/uploads"                  ; then exit 1; fi
if ! checkSymlinkExists "/chemotion/app/config/database.yml"      ; then exit 1; fi
if ! checkSymlinkExists "/chemotion/app/config/datacollector.yml" ; then exit 1; fi
if ! checkSymlinkExists "/chemotion/app/config/editors.yml"       ; then exit 1; fi
if ! checkSymlinkExists "/chemotion/app/config/secrets.yml"       ; then exit 1; fi
if ! checkSymlinkExists "/chemotion/app/config/storage.yml"       ; then exit 1; fi
if ! checkSymlinkExists "/chemotion/app/public/assets"            ; then exit 1; fi
if ! checkSymlinkExists "/chemotion/app/public/packs"             ; then exit 1; fi
if ! checkSymlinkExists "/chemotion/app/public/welcome-message.md"; then exit 1; fi

# check write permissions in folder
if ! checkFolderIsWritable "/shared/eln"        ; then exit 1; fi
if ! checkFolderIsWritable "/shared/eln/config" ; then exit 1; fi
if ! checkFolderIsWritable "/shared/eln/log"    ; then exit 1; fi
if ! checkFolderIsWritable "/shared/eln/public" ; then exit 1; fi
if ! checkFolderIsWritable "/shared/eln/tmp"    ; then exit 1; fi
if ! checkFolderIsWritable "/shared/eln/uploads"; then exit 1; fi

[[ -z $USERID ]] && export USERID=1000
chown -R $USERID:$USERID /shared

# check existance of certain files?
# copy files - still needed with new Dockerfile?

echo "Checks for the database:"
# check accessibility of DB:
db_profile="production"
db_configfile="/shared/eln/config/database.yml"
if [ -f $db_configfile ]; then
    for variable in $(python3 /etc/scripts/parseYML.py read --upper --prefix=DB_ $db_configfile $db_profile); do
        if [[ $variable == *"=<"* ]]; then
            echo "You are using environment definitions for the database configuration. This is obsolete."
            exit 1
        fi    
    done
 else
    echo "Cannot find the configuration file $db_configfile."
    exit 1
fi   

source <( python3 /etc/scripts/parseYML.py read --upper --prefix=DB_ $db_configfile $db_profile )

# simply waits for the DB to be up and done booting
echo "    Evaluated configuration file: $db_configfile"
echo "    Imported profile: $db_profile"
echo "    Connecting to host: $DB_HOST ..."
iterator=1
while ! pg_isready -h $DB_HOST 1>/dev/null 2>&1; do
    ((iterator++))
    echo "    Database instance not ready. Waiting ..."
    sleep 10
    if [ $iterator -eq 5 ]; then
        echo "    Database cannot be reached on $DB_HOST, please check the connection!"
        exit  1
    fi
done
echo "    Database instance ready."

# check correct setup of the DB and initialize DB
echo "    Creating database ..."
if ! (echo "\q" | psql -d $DB_DATABASE -h $DB_HOST -U $DB_USERNAME 2>/dev/null); then
    echo "    Can not connect to database or database needs to be initialized."
    read -e -p "
    Do you want to initialize the database ? [yes/N] " YN
    if [[ $YN == "yes" ]]; then
        sleep 3
        echo "    Dropping database $DB_DATABASE if it exists ..."
        psql --host="$DB_HOST" --username 'postgres' -c "
            DROP DATABASE IF EXISTS $DB_DATABASE;"
        echo "    Dropping role $DB_USERNAME if it exists ..."
        echo "    Creating role $DB_USERNAME with password $DB_PASSWORD ..."
        psql --host="$DB_HOST" --username 'postgres' -c "
            DROP ROLE IF EXISTS $DB_USERNAME;
            CREATE ROLE $DB_USERNAME LOGIN CREATEDB NOSUPERUSER PASSWORD '$DB_PASSWORD';"
        echo "    Creating database $DB_DATABASE for owner $DB_USERNAME ..."
        psql --host="$DB_HOST" --username 'postgres' -c "            
            CREATE DATABASE $DB_DATABASE OWNER $DB_USERNAME;
        " || {
            echo "    Could not create database. PSQL returned [$?]."
            exit 1
        }
        psql --host="$DB_HOST" --username="$DB_USERNAME" -c "
            CREATE EXTENSION IF NOT EXISTS \"pg_trgm\";
            CREATE EXTENSION IF NOT EXISTS \"hstore\";
            CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";
            ALTER USER $DB_USERNAME PASSWORD '$DB_PASSWORD';
        " || {
            echo "    Failed to set password for database user. PSQL returned [$?]."
            exit 1
        }
    else
        exit 1
    fi
fi

echo "    Database up and running."

# block PubChem
echo "Database setup:"
echo "    Blocking access to PubChem server..."
echo "127.0.0.1    pubchem.ncbi.nlm.nih.gov" >> /etc/hosts

cd /chemotion/app/

read -e -p "
Do you want to generate a new secrets.yml ? [yes/N] " YN
[[ $YN == "yes" ]] && generateNewSecrets

echo "    Initializing database schemas..."
bundle exec rake db:create
echo "    Database created."
bundle exec rake db:migrate
echo "    Database migrated."

read -e -p "
Do you want to seed database data ? [yes/N] " YN
[[ $YN == "yes" ]] && bundle exec rake db:seed && echo "    Database seeded."

echo "    Creating sprites..."
bundle exec rake ketcherails:import:common_templates
rm -rf /chemotion/app/app/public/images/ketcherails/icons/original/*
bundle exec rails r 'MakeKetcherailsSprites.perform_now'

# unblock PubChem
# do not use -i here. Docker prevents it from working...
echo "    Unblocking access to PubChem server..."
sed '/pubchem.ncbi.nlm.nih.gov/d' /etc/hosts > /etc/hosts

bundle exec rake assets:precompile