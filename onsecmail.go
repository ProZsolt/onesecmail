package onesecmail

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Mail struct {
	ID          int          `json:"id"`
	From        string       `json:"from"`
	Subject     string       `json:"subject"`
	Date        string       `json:"date"`
	Attachments []Attachment `json:"attachments"`
	Body        string       `json:"body"`
	TextBody    string       `json:"textBody"`
	HTMLBody    string       `json:"htmlBody"`
}

type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Size        int    `json:"size"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

var DefaultClient = &Client{
	baseURL:    "https://www.1secmail.com/api/v1/",
	httpClient: http.DefaultClient,
}

func (c Client) doQuery(query url.Values, ret any) error {
	req, err := http.NewRequest(http.MethodGet, c.baseURL, nil)
	if err != nil {
		return err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %v", resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(ret)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) GenRandomMailbox(count int) ([]string, error) {
	q := url.Values{}
	q.Add("action", "genRandomMailbox")
	if count > 0 {
		q.Add("count", strconv.Itoa(count))
	}

	var ret []string
	err := c.doQuery(q, &ret)
	return ret, err
}

func (c Client) GetDomainList() ([]string, error) {
	q := url.Values{}
	q.Add("action", "getDomainList")

	var ret []string
	err := c.doQuery(q, &ret)
	return ret, err
}

func (c Client) GetMessages(login string, domain string) ([]Mail, error) {
	q := url.Values{}
	q.Add("action", "getMessages")
	q.Add("login", login)
	q.Add("domain", domain)

	var ret []Mail
	err := c.doQuery(q, &ret)
	return ret, err
}

func (c Client) ReadMessage(login string, domain string, id int) (Mail, error) {
	q := url.Values{}
	q.Add("action", "readMessage")
	q.Add("login", login)
	q.Add("domain", domain)
	q.Add("id", strconv.Itoa(id))

	var ret Mail
	err := c.doQuery(q, &ret)
	return ret, err
}

func (c Client) DownloadAttachment(login string, domain string, id int, file string) (io.ReadCloser, error) {
	q := url.Values{}
	q.Add("action", "download")
	q.Add("login", login)
	q.Add("domain", domain)
	q.Add("id", strconv.Itoa(id))
	q.Add("file", file)

	req, err := http.NewRequest(http.MethodGet, c.baseURL, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected HTTP status: %v", resp.Status)
	}

	return resp.Body, nil
}

type Mailbox struct {
	client *Client
	login  string
	domain string
}

func splitAddress(addr string) (string, string) {
	sp := strings.Split(addr, "@")
	if len(sp) != 2 {
		panic("invalid address")
	}
	return sp[0], sp[1]
}

func GenerateRandomMailboxes(count int) ([]Mailbox, error) {
	addrs, err := DefaultClient.GenRandomMailbox(count)
	if err != nil {
		return nil, err
	}
	var mbs []Mailbox
	for _, addr := range addrs {
		l, d := splitAddress(addr)
		mbs = append(mbs, Mailbox{
			client: DefaultClient,
			login:  l,
			domain: d,
		})
	}
	return mbs, nil
}

func GenerateRandomMailbox() (Mailbox, error) {
	mbs, err := GenerateRandomMailboxes(1)
	if err != nil {
		return Mailbox{}, err
	}
	return mbs[0], nil
}

func (mb Mailbox) Address() string {
	return fmt.Sprintf("%s@%s", mb.login, mb.domain)
}

func (mb Mailbox) GetMessages() ([]Mail, error) {
	return mb.client.GetMessages(mb.login, mb.domain)
}

func (mb Mailbox) ReadMessage(id int) (Mail, error) {
	return mb.client.ReadMessage(mb.login, mb.domain, id)
}

func (mb Mailbox) DownloadAttachment(id int, file string) (io.ReadCloser, error) {
	return mb.client.DownloadAttachment(mb.login, mb.domain, id, file)
}
