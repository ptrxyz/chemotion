Rails.application.configure do
  if ENV["GIMME_THE_MEAT"].present?
    config.log_level = :warn
    Rails.logger.level = 0
  end
end

