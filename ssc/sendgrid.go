package ssc

import (
	"os"

	"github.com/ecsaas/sendgrid-sendmail-config/DEFINE_VARIABLES/SENDGRID"
	"github.com/ecsaas/sendgrid-sendmail-config/DEFINE_VARIABLES/SENDGRID_ENV"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendGridRequestSendmail(
	subject string,
	modelContent map[string]string,
	modelFrom map[string]string,
	modelTo map[string]string,
	modelReplyTo map[string]string,
	sendAtSinceTimeUnix int,
) (stt int, xMsgId string, err error) {
	if modelContent == nil || modelFrom == nil || modelTo == nil ||
		os.Getenv(SENDGRID_ENV.SENDGRID_API_KEY) == "" {
		return
	}
	m := mail.NewV3Mail()

	//subject = "Your Example Order Confirmation"
	m.Subject = subject

	//modelContent[SENDGRID.CONTENT_TYPE] = SENDGRID.VALUE_TEXT_HTML
	//modelContent[SENDGRID.CONTENT_VALUE] = "<p>Hello from Twilio SendGrid!</p><p>Sending with the email service trusted by developers and marketers for <strong>time-savings</strong>, <strong>scalability</strong>, and <strong>delivery expertise</strong>.</p><p>%open-track%</p>"
	m.AddContent(mail.NewContent(
		modelContent[SENDGRID.CONTENT_TYPE],
		modelContent[SENDGRID.CONTENT_VALUE],
	))

	//modelFrom[SENDGRID.NAME] = "Win Autoketing Sales Team"
	//modelFrom[SENDGRID.ADDRESS] = "win@autoketing.com"
	m.SetFrom(mail.NewEmail(modelFrom[SENDGRID.NAME], modelFrom[SENDGRID.ADDRESS]))

	personalization := mail.NewPersonalization()
	//modelTo[SENDGRID.NAME] = "MANH"
	//modelTo[SENDGRID.ADDRESS] = "manhbv@hellomedia.vn"
	personalization.AddTos([]*mail.Email{
		mail.NewEmail(modelTo[SENDGRID.NAME], modelTo[SENDGRID.ADDRESS]),
	}...)
	m.AddPersonalizations(personalization)

	if modelReplyTo != nil {
		//modelReplyTo[SENDGRID.NAME] = "Support Service Team"
		//modelReplyTo[SENDGRID.ADDRESS] = "manhbv@hellomedia.vn"
		m.SetReplyTo(mail.NewEmail(modelReplyTo[SENDGRID.NAME], modelReplyTo[SENDGRID.ADDRESS]))
	}

	if sendAtSinceTimeUnix > 0 {
		m.SetSendAt(sendAtSinceTimeUnix)
	}

	//trackingSettings := mail.NewTrackingSettings()
	//clickTrackingSetting := mail.NewClickTrackingSetting()
	//clickTrackingSetting.SetEnable(true)
	//clickTrackingSetting.SetEnableText(false)
	//trackingSettings.SetClickTracking(clickTrackingSetting)
	//openTrackingSetting := mail.NewOpenTrackingSetting()
	//openTrackingSetting.SetEnable(true)
	//openTrackingSetting.SetSubstitutionTag("%open-track%")
	//trackingSettings.SetOpenTracking(openTrackingSetting)
	//subscriptionTrackingSetting := mail.NewSubscriptionTrackingSetting()
	//subscriptionTrackingSetting.SetEnable(false)
	//trackingSettings.SetSubscriptionTracking(subscriptionTrackingSetting)
	//m.SetTrackingSettings(trackingSettings)

	sgRequest := sendgrid.GetRequest(
		os.Getenv(SENDGRID_ENV.SENDGRID_API_KEY),
		SENDGRID.SENDGRID_ENDPOINT,
		SENDGRID.SENDGRID_HOST,
	)
	sgRequest.Method = SENDGRID.POST
	sgRequest.Body = mail.GetRequestBody(m)

	sgResponse, err := sendgrid.API(sgRequest)
	if err == nil {
		stt = sgResponse.StatusCode
		if len(sgResponse.Headers[SENDGRID.X_MESSAGE_ID]) > 0 {
			xMsgId = sgResponse.Headers[SENDGRID.X_MESSAGE_ID][0]
		}
	}
	return
}
