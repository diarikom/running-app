package nmailgun

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"html/template"
	"time"
)

type Config struct {
	PrivateApiKey string `mapstructure:"private_api_key"`
	Domain        string `mapstructure:"domain"`
	TemplatePath  string `mapstructure:"template_path"`
	Region        string `mapstructure:"region"`
	DefaultSender string `mapstructure:"default_sender"`
}

type SendOpt struct {
	Sender       string
	Recipients   []string
	Subject      string
	Text         string
	TemplateFile string
	TemplateData interface{}
}

type Mailer struct {
	Domain        string
	DefaultSender string
	TemplatePath  string
	Timeout       time.Duration
	client        mailgun.Mailgun
}

func (m *Mailer) GetSender(name string) string {
	return fmt.Sprintf("%s@%s", name, m.Domain)
}

func (m *Mailer) GetDefaultSender() string {
	return fmt.Sprintf("%s@%s", m.DefaultSender, m.Domain)
}

func (m *Mailer) Send(opt SendOpt) error {
	// Parse template file
	html, err := parseTemplate(m.TemplatePath+opt.TemplateFile, opt.TemplateData)
	if err != nil {
		return err
	}

	// Create message
	message := m.client.NewMessage(opt.Sender, opt.Subject, opt.Text, opt.Recipients...)
	message.SetHtml(html)

	// Init context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send html
	if _, _, err := m.client.Send(ctx, message); err != nil {
		return err
	}
	return nil
}

// ParseTemplate returns parsed template with data in string format
// If there is an error, it will return response with error data
func parseTemplate(templateFilePath string, data interface{}) (content string, err error) {
	// ParseFiles creates a new Template and parses the template definitions from
	// the named files. The returned template's name will have the (base) name and
	// (parsed) contents of the first file. There must be at least one file.
	// If an error occurs, parsing stops and the returned *Template is nil.
	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return "", err
	}

	// A Buffer is a variable-sized buffer of bytes with Read and Write methods.
	// The zero value for Buffer is an empty buffer ready to use.
	buf := new(bytes.Buffer)

	// Execute applies a parsed template to the specified data object,
	// writing the output to wr.
	// If an error occurs executing the template or writing its output,
	// execution stops, but partial results may already have been written to
	// the output writer.
	// A template may be executed safely in parallel.
	if err := tmpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func New(config Config) *Mailer {
	// Init client
	cl := mailgun.NewMailgun(config.Domain, config.PrivateApiKey)

	// If Region is europe, set base api url
	if config.Region == "eu" {
		cl.SetAPIBase(mailgun.APIBaseEU)
	}

	// If default sender not found, change to no-reply
	if config.DefaultSender == "" {
		config.DefaultSender = "no-reply"
	}

	return &Mailer{
		TemplatePath:  config.TemplatePath,
		client:        cl,
		Domain:        config.Domain,
		DefaultSender: config.DefaultSender,
	}
}
