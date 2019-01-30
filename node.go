package embeddedShell

import (
	"fmt"
	"io/ioutil"

	"github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"golang.org/x/net/context"
	"gopkg.in/errgo.v1"
)

func NewDefaultNodeWithFSRepo(ctx context.Context, repoPath string) (*core.IpfsNode, error) {
	r, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, errgo.Notef(err, "opening fsrepo failed")
	}
	node, err := core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Repo:   r,
	})
	if err != nil {
		return nil, errgo.Notef(err, "ipfs NewNode() failed.")
	}
	// TODO: can we bootsrap localy/mdns first and fall back to default?
	err = node.Bootstrap(core.DefaultBootstrapConfig)
	if err != nil {
		return nil, errgo.Notef(err, "ipfs Bootstrap() failed.")
	}
	return node, nil
}

func NewTmpDirNode(ctx context.Context) (*core.IpfsNode, error) {
	dir, err := ioutil.TempDir("", "ipfs-shell")
	if err != nil {
		return nil, fmt.Errorf("failed to get temp dir: %s", err)
	}

	cfg, err := config.Init(ioutil.Discard, 1024)
	if err != nil {
		return nil, err
	}

	err = fsrepo.Init(dir, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	repo, err := fsrepo.Open(dir)
	if err != nil {
		return nil, err
	}

	return core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Repo:   repo,
	})
}
