Rails.application.configure do
  # Settings specified here will take precedence over those in config/application.rb.

  # Code is not reloaded between requests.
  config.cache_classes = true

  # Eager load code on boot. This eager loads most of Rails and
  # your application in memory, allowing both threaded web servers
  # and those relying on copy on write to perform better.
  # Rake tasks automatically ignore this option for performance.
  config.eager_load = true

  # Full error reports are disabled and caching is turned on.
  config.consider_all_requests_local       = false
  config.action_controller.perform_caching = true

  # Enable Rack::Cache to put a simple HTTP cache in front of your application
  # Add `rack-cache` to your Gemfile before enabling this.
  # For large-scale production use, consider using a caching reverse proxy like
  # NGINX, varnish or squid.
  # config.action_dispatch.rack_cache = true

  # Disable serving static files from the `/public` folder by default since
  # Apache or NGINX already handles this.
  # config.serve_static_files = ENV['RAILS_SERVE_STATIC_FILES'].present?
  # UPDATE we have to enable it for engines with assets
  config.serve_static_files = true

  # Compress JavaScripts and CSS.
  config.assets.js_compressor = Uglifier.new(harmony: true, mangle: false, compress: false)
  config.assets.css_compressor = :sass

  # Do not fallback to assets pipeline if a precompiled asset is missed.
  config.assets.compile = true

  # Asset digests allow you to set far-future HTTP expiration dates on all assets,
  # yet still be able to expire them through the digest params.
  config.assets.digest = true

  # `config.assets.precompile` and `config.assets.version` have moved to config/initializers/assets.rb

  TurboSprockets.configure do |config|
    config.precompiler.enabled = false
    config.preloader.enabled = false
    config.preloader.worker_count = 2
    config.precompiler.worker_count = 2
  end

  # Specifies the header that your server uses for sending files.
  # config.action_dispatch.x_sendfile_header = 'X-Sendfile' # for Apache
  # config.action_dispatch.x_sendfile_header = 'X-Accel-Redirect' # for NGINX

  # Force all access to the app over SSL, use Strict-Transport-Security, and use secure cookies.
  # config.force_ssl = true

  # Use the lowest log level to ensure availability of diagnostic information
  # when problems arise.
  config.log_level = :info

  # Prepend all log lines with the following tags.
  # config.log_tags = [ :subdomain, :uuid ]

  # Use a different logger for distributed setups.
  # config.logger = ActiveSupport::TaggedLogging.new(SyslogLogger.new)

  # Use a different cache store in production.
  # config.cache_store = :mem_cache_store

  # Enable serving of images, stylesheets, and JavaScripts from an asset server.
  # config.action_controller.asset_host = 'http://assets.example.com'

  # Ignore bad email addresses and do not raise email delivery errors.
  # Set this to true and configure the email server for immediate delivery to raise delivery errors.
  # config.action_mailer.raise_delivery_errors = false

  # Enable locale fallbacks for I18n (makes lookups for any locale fall back to
  # the I18n.default_locale when a translation cannot be found).
  config.i18n.fallbacks = true

  # Send deprecation notices to registered listeners.
  config.active_support.deprecation = :notify

  # Use default logging formatter so that PID and timestamp are not suppressed.
  #config.log_formatter = ::Logger::Formatter.new
  config.log_formatter = OnelineLogFormatter.new


  config.action_mailer.raise_delivery_errors = true
  config.action_mailer.delivery_method = :smtp
  config.action_mailer.default_url_options = { host: ENV['SMTP_HOST']}
  config.action_mailer.smtp_settings = {
    :address              => ENV["SMTP_ADDRESS"],
    :port                 => ENV["SMTP_PORT"],
    :user_name            => ENV['SMTP_USERNAME'],
    :domain               => ENV['SMTP_DOMAIN'],
    :password             => ENV['SMTP_PASSWORD'],
    :authentication       => ENV['SMTP_AUTH'] && ENV['SMTP_AUTH'].to_sym,
    :enable_starttls_auto => ENV['SMTP_TLS'] && ENV['SMTP_TLS'].match(/true/),
    :openssl_verify_mode  => ENV['SMTP_SSL_MODE']
  }

  config.browserify_rails.commandline_options = "-t [ babelify --presets [ @babel/preset-env  @babel/preset-react ] --plugins [ @babel/plugin-proposal-object-rest-spread ] ] -g uglifyify "
end
