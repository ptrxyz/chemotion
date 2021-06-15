production:
  adapter: postgresql
  encoding: unicode
  database: {{DB_NAME}}
  pool: 5
  username: {{DB_ROLE}}
  password: {{DB_PW}}
  host: {{DB_HOST}}
  port: {{DB_PORT}}
