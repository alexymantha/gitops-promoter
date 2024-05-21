package github

import (
	"context"
	promoterv1alpha1 "github.com/argoproj/promoter/api/v1alpha1"
	"github.com/argoproj/promoter/internal/scms"
	"github.com/google/go-github/v61/github"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"
)

type CommitStatus struct {
	client *github.Client
}

var _ scms.CommitStatusProvider = &CommitStatus{}

func NewGithubCommitStatusProvider(secret v1.Secret) (*CommitStatus, error) {
	client, err := GetClient(secret)
	if err != nil {
		return nil, err
	}

	return &CommitStatus{
		client: client,
	}, nil
}

func (cs CommitStatus) Set(ctx context.Context, commitStatus *promoterv1alpha1.CommitStatus) (*promoterv1alpha1.CommitStatus, error) {
	logger := log.FromContext(ctx)
	logger.Info("Setting Commit Status")

	repoStatus, response, err := cs.client.Repositories.CreateStatus(ctx, commitStatus.Spec.RepositoryReference.Owner, commitStatus.Spec.RepositoryReference.Name, commitStatus.Spec.Sha, &github.RepoStatus{
		State:       github.String(commitStatus.Spec.State),
		TargetURL:   github.String(commitStatus.Spec.Url),
		Description: github.String(commitStatus.Spec.Description),
		Context:     github.String(commitStatus.Spec.Name),
	})
	if err != nil {
		return nil, err
	}
	logger.Info("github rate limit",
		"limit", response.Rate.Limit,
		"remaining", response.Rate.Remaining,
		"reset", response.Rate.Remaining)
	logger.Info("github response status",
		"status", response.Status)

	commitStatus.Status.Id = strconv.FormatInt(*repoStatus.ID, 10)
	commitStatus.Status.State = *repoStatus.State
	commitStatus.Status.Sha = commitStatus.Spec.Sha
	commitStatus.Status.ObservedGeneration = commitStatus.Generation
	return commitStatus, nil
}
