package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tmknom/cross/internal/errlib"
)

type Repository struct {
	repo     *git.Repository
	worktree *git.Worktree
}

func OpenRepo(path string) (*Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, errlib.Wrapf(err, "cannot open repository: %s", path)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, errlib.Wrapf(err, "invalid repository: %#v", repo)
	}

	result := &Repository{repo: repo, worktree: worktree}
	return result, nil
}

func (r *Repository) Add() error {
	opts := &git.AddOptions{
		All: true,
	}
	err := r.worktree.AddWithOptions(opts)
	if err != nil {
		return errlib.Wrapf(err, "cannot git add: %#v", opts)
	}
	return nil
}

func (r *Repository) Commit(message string) (string, error) {
	opts := &git.CommitOptions{
		AllowEmptyCommits: true,
	}
	hash, err := r.worktree.Commit(message, opts)
	if err != nil {
		return "", errlib.Wrapf(err, "cannot git commit: %#v", opts)
	}
	return hash.String(), nil
}

func (r *Repository) SwitchOrCreate(branch string) error {
	err := r.Switch(branch)
	if err != nil {
		return r.createBranchAndSwitch(branch)
	}
	return nil
}

func (r *Repository) Switch(branch string) error {
	return r.checkout(branch, false)
}

func (r *Repository) createBranchAndSwitch(branch string) error {
	return r.checkout(branch, true)
}

func (r *Repository) checkout(branch string, create bool) error {
	opts := &git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Create: create,
		Keep:   true,
	}
	err := r.worktree.Checkout(opts)
	if err != nil {
		return errlib.Wrapf(err, "cannot git checkout: %#v", opts)
	}

	return nil
}
