/*
Copyright 2023 The Bestchains Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package auth

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/bestchains/bc-cli/pkg/common"
	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/int128/oauth2cli"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
	"k8s.io/klog/v2"
)

const (
	localServerSuccessHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>已认证</title>
	<script>
		window.close()
	</script>
	<style>
		body {
			background-color: #eee;
			margin: 0;
			padding: 0;
			font-family: sans-serif;
		}
		.placeholder {
			margin: 2em;
			padding: 2em;
			background-color: #fff;
			border-radius: 1em;
		}
	</style>
</head>
<body>
	<div class="placeholder">
		<h1>已认证</h1>
		<p>现在您可以关闭该窗口。</p>
	</div>
</body>
</html>
`
)

func Auth(ctx context.Context, config *common.AuthConfig) (authConfig *common.AuthConfig, err error) {
	if !config.Enable {
		return config, nil
	}
	enableAuth = true
	var client *client
	client, err = newClient(ctx, *config)
	if err != nil {
		return nil, err
	}
	if config.Expiry != 0 && config.IDToken != "" {
		if time.Now().Before(time.Unix(config.Expiry, 0)) {
			klog.V(2).Infoln("Parse ID token from config file and try to verify it is valid...")
			err = client.verifyIDToken(ctx)
		} else {
			klog.V(2).Infoln("ID token has expired, try to refresh it...")
			err = client.refresh(ctx)
		}

		if err == nil {
			idToken = config.IDToken
			return &client.AuthConfig, nil
		}
		klog.Errorf("failed to verify or refresh ID token: %v", err)
	}
	return client.newAuthReq(ctx)
}

func newClient(ctx context.Context, config common.AuthConfig) (*client, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, // ignore tls verify, because each has its own cert
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	provider, err := gooidc.NewProvider(ctx, config.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("oidc discovery error: %w", err)
	}
	oauth2Config := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.IssuerURL + "/auth",
			TokenURL: config.IssuerURL + "/token",
		},
		Scopes: []string{"openid", "email", "groups", "profile", "offline_access"},
	}
	return &client{
		httpClient:   httpClient,
		oauth2Config: oauth2Config,
		provider:     provider,
		AuthConfig:   config,
	}, nil
}

type client struct {
	AuthConfig   common.AuthConfig
	oauth2Config oauth2.Config
	provider     *gooidc.Provider
	httpClient   *http.Client
}

func (c *client) newAuthReq(ctx context.Context) (config *common.AuthConfig, err error) {
	ready := make(chan string, 1)
	defer close(ready)
	cfg := oauth2cli.Config{
		OAuth2Config:           c.oauth2Config,
		LocalServerReadyChan:   ready,
		Logf:                   klog.V(2).Infof,
		LocalServerBindAddress: []string{common.LocalBindPort},
		LocalServerSuccessHTML: localServerSuccessHTML,
	}
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case url := <-ready:
			klog.V(2).Infof("Open %s in default browser...", url)
			err := browser.OpenURL(url)
			if err != nil {
				klog.Errorf("could not open the browser: %w", err)
			}
			return err
		case <-ctx.Done():
			return fmt.Errorf("context done while waiting for authorization: %w", ctx.Err())
		}
	})
	eg.Go(func() error {
		ctx = c.wrapContext(ctx)
		token, err := oauth2cli.GetToken(ctx, cfg)
		if err != nil {
			return fmt.Errorf("could not get a token: %w", err)
		}
		klog.V(2).Infof("You got a valid token, will expiry in %s", time.Until(token.Expiry))
		c.AuthConfig.RefreshToken = token.RefreshToken
		c.AuthConfig.Expiry = token.Expiry.Unix()
		return c.verifyToken(ctx, token)
	})
	if err := eg.Wait(); err != nil {
		klog.Errorf("authorization error: %s", err)
		return &c.AuthConfig, err
	}
	idToken = c.AuthConfig.IDToken
	return &c.AuthConfig, nil
}

func (c *client) verifyToken(ctx context.Context, token *oauth2.Token) error {
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return fmt.Errorf("id_token is missing in the token response: %v", token)
	}
	c.AuthConfig.IDToken = idToken
	if err := c.verifyIDToken(ctx); err != nil {
		return err
	}
	return nil
}

func (c *client) verifyIDToken(ctx context.Context) error {
	verifier := c.provider.Verifier(&gooidc.Config{ClientID: c.AuthConfig.ClientID})
	idToken, err := verifier.Verify(ctx, c.AuthConfig.IDToken)
	if err != nil {
		return fmt.Errorf("could not verify the ID token: %w", err)
	}
	oUser := make(map[string]interface{})
	err = idToken.Claims(&oUser)
	if err != nil {
		return fmt.Errorf("could not parse the ID token: %w", err)
	}
	c.AuthConfig.Username = oUser["preferred_username"].(string)
	return nil
}

func (c *client) wrapContext(ctx context.Context) context.Context {
	if c.httpClient != nil {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, c.httpClient)
	}
	return ctx
}

// refresh sends a refresh token request and returns a token set.
func (c *client) refresh(ctx context.Context) error {
	ctx = c.wrapContext(ctx)
	currentToken := &oauth2.Token{
		Expiry:       time.Now(),
		RefreshToken: c.AuthConfig.RefreshToken,
	}
	source := c.oauth2Config.TokenSource(ctx, currentToken)
	token, err := source.Token()
	if err != nil {
		klog.Errorf("could not refresh the token: %s", err)
		return fmt.Errorf("could not refresh the token: %w", err)
	}
	c.AuthConfig.RefreshToken = token.RefreshToken
	c.AuthConfig.Expiry = token.Expiry.Unix()
	return c.verifyToken(ctx, token)
}

var idToken string
var enableAuth bool

func AddAuthHeader(req *http.Request) {
	if enableAuth && idToken != "" {
		req.Header.Add("Authorization", "Bearer "+idToken)
	}
}
