#!/bin/bash
set -euo pipefail

cd /chemotion/app || exit 200
echo "u = User.find_by(type: 'Admin'); u.password='chemotion'; u.account_active=true; u.save" | bundle exec rails c
