package main

import "tp/client/helpers"

type Client struct {
	Config            helpers.Config
	UploadedFileCount int
}

func newClient() *Client {
	config := helpers.LoadConfig()
	helpers.InitStore(*config)
	return &Client{
		Config:            *config,
		UploadedFileCount: 0,
	}
}
