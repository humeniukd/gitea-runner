package exprparser

import "gitea.com/gitea/act_runner/pkg/model"

func (impl *interperterImpl) getNeedsTransitive(job *model.Job) []string {
	needs := job.Needs()

	for _, need := range needs {
		parentNeeds := impl.getNeedsTransitive(impl.config.Run.Workflow.GetJob(need))
		needs = append(needs, parentNeeds...)
	}

	return needs
}

func (impl *interperterImpl) always() (any, error) {
	return true, nil
}

func (impl *interperterImpl) jobSuccess() (any, error) {
	jobs := impl.config.Run.Workflow.Jobs
	jobNeeds := impl.getNeedsTransitive(impl.config.Run.Job())

	for _, needs := range jobNeeds {
		if jobs[needs].Result != "success" {
			return false, nil
		}
	}

	return true, nil
}

func (impl *interperterImpl) stepSuccess() (any, error) {
	return impl.env.Job.Status == "success", nil
}

func (impl *interperterImpl) jobFailure() (any, error) {
	jobs := impl.config.Run.Workflow.Jobs
	jobNeeds := impl.getNeedsTransitive(impl.config.Run.Job())

	for _, needs := range jobNeeds {
		if jobs[needs].Result == "failure" {
			return true, nil
		}
	}

	return false, nil
}

func (impl *interperterImpl) stepFailure() (any, error) {
	return impl.env.Job.Status == "failure", nil
}

func (impl *interperterImpl) cancelled() (any, error) {
	return impl.env.Job.Status == "cancelled", nil
}
