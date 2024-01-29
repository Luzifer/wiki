package main

import "embed"

//go:embed frontend/* default_files/*
var assets embed.FS
