package distribution

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const GoogleOauth2Uri string = "https://accounts.google.com/o/oauth2/auth"
const GoogleOauth2TokenUri string = "https://accounts.google.com/o/oauth2/token"

const BloggerClientRedirectUri string = "urn:ietf:wg:oauth:2.0:oob"
const BloggerClientScope string = "https://www.googleapis.com/auth/blogger"
const BloggerClientPostUri string = "https://www.googleapis.com/blogger/v3/blogs/%s/posts/"

type BloggerClientConfig struct {
	BlogId       string
	ClientId     string
	ClientSecret string
	AuthCode     string
	TokenCache   oauth.Cache
}

type BloggerClient struct {
	BlogId    string
	authCode  string
	postUri   string
	transport *oauth.Transport
}

func NewBloggerClient(config *BloggerClientConfig) *BloggerClient {
	postUri := fmt.Sprintf(BloggerClientPostUri, config.BlogId)
	oauthConfig := &oauth.Config{
		ClientId:     config.ClientId,
		ClientSecret: config.ClientSecret,
		Scope:        BloggerClientScope,
		AuthURL:      GoogleOauth2Uri,
		TokenURL:     GoogleOauth2TokenUri,
		RedirectURL:  BloggerClientRedirectUri,
		TokenCache:   config.TokenCache,
	}
	transport := &oauth.Transport{Config: oauthConfig}
	return &BloggerClient{BlogId: config.BlogId, postUri: postUri, transport: transport}
}

func (b *BloggerClient) Client() (*http.Client, error) {
	if b.transport.Token == nil {
		token, err := b.transport.Config.TokenCache.Token()
		if err != nil {
			_, err = b.transport.Exchange(b.authCode)
			if err != nil {
				return nil, err
			}
		} else {
			b.transport.Token = token
		}
	}
	return b.transport.Client(), nil
}

func (b *BloggerClient) Send(post io.Reader) (string, error) {
	client, err := b.Client()
	if err != nil {
		return "", err
	}
	postBytes, err := ioutil.ReadAll(post)
	if err != nil {
		return "", err
	}
	values := make(url.Values)
	values.Add("fields", "url")
	payload := fmt.Sprintf(`{
		"kind": "blogger#post",
		"blog": {
			"id": "%s"
		},
		"title": "Message",
		"content": "%s,
	}`, b.BlogId, string(postBytes))
	response, err := client.Post(b.postUri+"?"+values.Encode(), "application/json", strings.NewReader(payload))
	if err != nil {
		return "", err
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var f interface{}
	err = json.Unmarshal(bodyBytes, &f)
	if err != nil {
		return "", err
	}
	jsonData := f.(map[string]interface{})
	url, ok := jsonData["url"].(string)
	if !ok {
		return "", errors.New(fmt.Sprintf("Got wrong response: %v\n", jsonData))
	}
	return url, nil
}

func (b *BloggerClient) Receive() ([]io.Reader, error) {
	client, err := b.Client()
	if err != nil {
		return nil, err
	}

	values := make(url.Values)
	values.Add("fields", "items/content")

	// Make the request.
	response, err := client.Get(b.postUri + "?" + values.Encode())
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var f interface{}
	err = json.Unmarshal(bodyBytes, &f)
	if err != nil {
		return nil, err
	}
	jsonData := f.(map[string]interface{})
	items := jsonData["items"].([]interface{})
	fmt.Println(items)
	posts := make([]io.Reader, len(items))
	for i, item := range items {
		postContent := item.(map[string]interface{})
		posts[i] = strings.NewReader(postContent["content"].(string))
	}
	return posts, nil
}
