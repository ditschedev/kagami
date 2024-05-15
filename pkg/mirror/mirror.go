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

type Mirror interface {
	Mirror()
}

type mirror struct {
	config  *config.MirrorConfig
	spinner *spinner.Spinner
	tmpDir  string
}

func New(cfg *config.MirrorConfig) Mirror {
	return &mirror{
		config:  cfg,
		spinner: spinner.New(spinner.CharSets[40], 100*time.Millisecond),
	}
}

func (m *mirror) Mirror() {
	defer m.cleanup()

	color.Set(color.FgHiBlack)
	_ = m.spinner.Color("fgHiBlack")
	color.Set(color.FgHiBlack)
	m.spinner.Suffix = fmt.Sprintf(" Mirroring repository %s", m.config.Name)
	m.spinner.FinalMSG = fmt.Sprintf("✓ Successfully mirrored %s\n", m.config.Name)
	m.spinner.Start()

	// check if already cloned
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

	//// fetch
	//err = repo.Fetch(&git.FetchOptions{
	//	RemoteName: "origin",
	//})
	//if err != nil {
	//
	//	if !errors.Is(err, git.NoErrAlreadyUpToDate) {
	//		log.Fatal("Could not fetch latest changes")
	//		return
	//	}
	//
	//	if !initialClone {
	//		log.Write(fmt.Sprintf("No changes detected for %s", m.config.Name), color.FgHiGreen)
	//		return
	//	}
	//}
	//
	//w, err := repo.Worktree()
	//if err != nil {
	//	fmt.Println(err)
	//	log.Fatal("Could not get worktree")
	//}
	//
	//// Pull the latest changes from the origin remote and merge into the current branch
	//err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	//if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
	//	fmt.Println(err)
	//	log.Fatal("Could not pull latest changes")
	//}

	// push
	err = repo.Push(&git.PushOptions{
		RemoteName: "mirror",
	})
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			//log.Write(fmt.Sprintf("✓ No changes detected for %s", m.config.Name), color.FgHiGreen)
			return
		}

		log.Error(err, "Could not push changes")
		return
	}

	//log.Write(fmt.Sprintf("✓ Successfully mirrored %s", m.config.Name), color.FgHiGreen)
}

func (m *mirror) cleanup() {
	if _, err := os.Stat(m.tmpDir); !os.IsNotExist(err) {
		err := os.RemoveAll(m.tmpDir)
		if err != nil {
			log.Error(err, "Could not remove cache directory")
		}
	}

	color.Set(color.FgGreen)
	m.spinner.Stop()
}
