package main

import (
	"bytes"
	"fmt"

	"html/template"

	"gopkg.in/mailgun/mailgun-go.v1"
)

func sendVerifyMail(name, email, verifyKey string) bool {
	type Mail struct {
		Name       string
		SiteName   string
		ActiveLink string
	}
	var doc bytes.Buffer
	var templateString = `{{.Name}} 您好,
    这是里是{{.SiteName}}, 请您点击下面的链接确认您的邮箱:
    {{.ActiveLink}}
-------------
如需帮助请联系: support@linkown.com`
	t := template.New("")
	t, _ = t.Parse(templateString)
	p := Mail{Name: name, SiteName: gconf.SiteName, ActiveLink: gconf.SiteLink + `/verifyKey?k=` + verifyKey}
	t.Execute(&doc, p)

	if gconf.MailGun.ApiKey != "" && gconf.MailGun.Domain != "" && gconf.MailGun.PublicApiKey != "" && gconf.MailGun.Sender != "" {
		mg := mailgun.NewMailgun(gconf.MailGun.Domain, gconf.MailGun.ApiKey, gconf.MailGun.PublicApiKey)

		message := mailgun.NewMessage(
			gconf.MailGun.Sender,
			gconf.SiteName+" 邮箱确认",
			doc.String(),
			email)
		resp, id, err := mg.Send(message)
		if err != nil {
			return false
		}
		fmt.Printf("ID: %s Resp: %s\n", id, resp)
		return true
	}
	return false
}
