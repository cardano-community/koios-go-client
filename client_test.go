// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios

// func TestNewDefaults(t *testing.T) {
// 	api, err := New()
// 	assert.NoError(t, err)
// 	if assert.NotNil(t, api) {
// 		raw := fmt.Sprintf(
// 			"%s://%s/api/%s/",
// 			DefaultScheme,
// 			MainnetHost,
// 			DefaultAPIVersion,
// 		)
// 		u, err := url.ParseRequestURI(raw)
// 		assert.NoError(t, err, "default url can not be constructed")
// 		assert.Equal(t, u.String(), api.BaseURL(), "invalid default base url")
// 	}
// }

// func TestOptions(t *testing.T) {
// 	api, err := New(
// 		Host("localhost"),
// 		APIVersion("v1"),
// 		Port(8080),
// 		Scheme("http"),
// 		RateLimit(100),
// 		Origin("http://localhost.localdomain"),
// 		CollectRequestsStats(true),
// 	)
// 	assert.NoError(t, err)
// 	if assert.NotNil(t, api) {
// 		assert.Equal(t, "http://localhost:8080/api/v1/", api.BaseURL(), "invalid default base url")
// 	}

// 	api2, err2 := New(Scheme("ws"))
// 	assert.EqualError(t, err2, "scheme must be http or https")
// 	assert.Nil(t, api2)
// }

// func TestOptionErrs(t *testing.T) {
// 	client, _ := New()
// 	assert.Error(t, HTTPClient(http.DefaultClient).apply(client),
// 		"should not allow changing http client.")
// 	assert.Error(t, RateLimit(0).apply(client),
// 		"should not unlimited requests p/s")
// 	assert.Error(t, Origin("localhost").apply(client),
// 		"origin should be valid http origin")
// 	_, err := New(Origin("localhost.localdomain"))
// 	assert.Error(t, err, "New should return err when option is invalid")
// }

// func TestHTTPClient(t *testing.T) {
// 	client, err := New(HTTPClient(http.DefaultClient))
// 	assert.NotEqual(t, http.DefaultClient, client)
// 	assert.Error(t, err)
// }

// func TestReadResponseBody(t *testing.T) {
// 	// enure that readResponseBody behaves consistently
// 	nil1, nil2 := ReadResponseBody(nil)
// 	assert.Nil(t, nil1)
// 	assert.Nil(t, nil2)
// }

// func TestNewClients(t *testing.T) {
// 	c1, err := New()
// 	assert.NoError(t, err)
// 	assert.NotNil(t, c1)

// 	client := http.DefaultClient
// 	client.Timeout = 10 * time.Second
// 	c2, err2 := New(HTTPClient(client))
// 	assert.NoError(t, err2)
// 	assert.NotNil(t, c2)

// 	c3, err := c1.WithOptions(HTTPClient(client))
// 	assert.NoError(t, err)
// 	assert.NotNil(t, c3)
// }

// var errApply = errors.New("apply error")

// func TestApplyError(t *testing.T) {
// 	res := &Response{}

// 	res.applyError([]byte(
// 		"{hint:\"the-hint\",details:\"the-details\",code:101,message:\"the-message\"}"),
// 		errApply,
// 	)
// 	assert.Equal(t, "apply error: invalid character 'h' looking for beginning of object key string", res.Error.Message)

// 	res.applyError(
// 		[]byte("{\"hint\":\"the-hint\",\"details\":\"the-details\",\"code\":\"101\",\"message\":\"the-message\"}"),
// 		errApply,
// 	)
// 	assert.Equal(t, "the-hint", res.Error.Hint)
// 	assert.Equal(t, "apply error: the-message", res.Error.Message)
// 	assert.Equal(t, ErrorCode("101"), res.Error.Code)
// 	assert.Equal(t, "the-details", res.Error.Details)
// }

// func TestBaseURL(t *testing.T) {
// 	api, _ := New()
// 	err := api.setBaseURL("http", "localhost", "v2", 9000)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "http://localhost:9000/api/v2/", api.url.String())

// 	err2 := api.setBaseURL("http", "localhost\\invalid", "v2", 9000)
// 	assert.EqualError(
// 		t,
// 		err2,
// 		"parse \"http://localhost\\\\invalid:9000/api/v2/\": invalid character \"\\\\\" in host name",
// 	)
// }
