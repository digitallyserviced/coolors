package runtime

import (
	"embed"
	"net/http"

	"github.com/pgavlin/femto"
)

//go:embed colorschemes/* syntax/*
var embedFs embed.FS

// emb
var Files = femto.NewRuntimeFiles(http.FS(embedFs))
