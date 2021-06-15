production:
  :services:
    - :name: 'mailcollector'
      :every: 5 # minutes
    - :name: 'folderwatchersftp'
      :every: 5 # minutes
    - :name: 'folderwatcherlocal'
      :every: 5 # minutes
    - :name: 'filewatchersftp'
      :every: 2 # minutes
    - :name: 'filewatcherlocal'
      :every: 2 # minutes

  :mailcollector:
    :server: 'imap.server.de'
    :mail_address: 'service@mail'
    :password: 'password'
    # :port: 993 default
    # :ssl: true default
    :aliases:
      - 'alias_one@kit.edu'
      - 'alias_two@kit.edu'

  :sftpusers:
    - :user: 'user1'
      :password: 'pass'
    - :user: 'user2'
      :password: 'pass'

  :localcollectors:
    - :path: '/home/ftpuser'
    - :path: '/home/eln/public'
