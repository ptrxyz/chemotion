class SendSharingNotificationJob < ActiveJob::Base
  queue_as :default

  def perform(user, data)
    SharingNotificationMailer.send_mail(user, data).deliver_now
  end
end
