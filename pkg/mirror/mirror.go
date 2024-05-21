package mirror

import (
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/ditschedev/kagami/pkg/config"
	"github.com/ditschedev/kagami/pkg/log"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"os"
	"path"
	"time"
)

const (
	mirrorStatusSuccess = iota
	mirrorStatusNoChanges
	mirrorStatusError
)

type Mirror interface {
	Mirror()
}

type mirror struct {
	config  *config.MirrorConfig
	spinner *spinner.Spinner
	tmpDir  string
	status  int
}

func New(cfg *config.MirrorConfig) Mirror {
	return &mirror{
		config:  cfg,
		spinner: spinner.New(spinner.CharSets[40], 100*time.Millisecond),
		status:  mirrorStatusError,
	}
}

func (m *mirror) Mirror() {
	defer m.cleanup()

	color.Set(color.FgHiBlack)
	_ = m.spinner.Color("fgHiBlack")
	color.Set(color.FgHiBlack)
	m.spinner.Suffix = fmt.Sprintf(" Mirroring repository %s\n", m.config.Name)
	m.spinner.Start()

	cachePath, err := os.MkdirTemp("", "kagami-mirrors")
	localPath := path.Join(cachePath, m.config.Name)

	if err != nil {
		log.Error(err, "Could not create temporary directory")
		return
	}

	m.tmpDir = cachePath

	authMethod, err := ssh.NewSSHAgentAuth("git")
	if err != nil {
		log.Error(err, "Could not create ssh auth method")
		return
	}

	repo, err := git.PlainClone(localPath, false, &git.CloneOptions{
		URL:        m.config.RemoteUri,
		Auth:       authMethod,
		Mirror:     true,
		RemoteName: "origin",
	})

	if err != nil {
		log.Error(err, "Could not clone repository")
		return
	}

	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name:   "mirror",
		URLs:   []string{m.config.MirrorUri},
		Mirror: true,
	})

	if err != nil {
		log.Error(err, "Could not create mirror remote")
		return
	}

	// push
	err = repo.Push(&git.PushOptions{
		RemoteName: "mirror",
	})
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			m.status = mirrorStatusNoChanges
			return
		}

		log.Error(err, "Could not push changes")
		return
	}

	m.status = mirrorStatusSuccess
}

func (m *mirror) cleanup() {
	if _, err := os.Stat(m.tmpDir); !os.IsNotExist(err) {
		err := os.RemoveAll(m.tmpDir)
		if err != nil {
			log.Error(err, "Could not remove cache directory")
		}
	}

	if m.status == mirrorStatusError {
		return
	}

	color.Set(color.FgGreen)
	m.spinner.FinalMSG = fmt.Sprintf("✓ Successfully mirrored %s\n", m.config.Name)
	m.spinner.Stop()
}
