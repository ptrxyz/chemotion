#!/bin/bash
set -euo pipefail

cd /chemotion/app
echo "u = User.find_by(type: 'Admin'); u.password='chemotion'; u.account_active=true; u.save" | rails c
