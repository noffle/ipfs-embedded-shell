package embeddedShell

import (
	"io"

	"gopkg.in/errgo.v1"

	chunker "github.com/ipfs/go-ipfs-chunker"
	dag "github.com/ipfs/go-merkledag"
	unixfs "github.com/ipfs/go-unixfs"
	importer "github.com/ipfs/go-unixfs/importer"
)

func (s *Shell) Add(r io.Reader) (string, error) {
	dag, err := importer.BuildDagFromReader(
		s.node.DAG,
		chunker.DefaultSplitter(r),
	)
	if err != nil {
		return "", errgo.Notef(err, "add: importing DAG failed.")
	}
	return dag.Cid().String(), nil
}

// AddLink creates a unixfs symlink and returns its hash
func (s *Shell) AddLink(target string) (string, error) {
	d, _ := unixfs.SymlinkData(target)
	nd := dag.NodeWithData(d)
	err := s.node.DAG.Add(s.ctx, nd)
	if err != nil {
		return "", err
	}
	return nd.Cid().String(), nil
}
