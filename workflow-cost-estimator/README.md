# Workflow Cost Estimator

The Workflow Cost Estimator is a lambda function with the goal of estimating a given workflow cost. This is not anything more than an estimate. This is not officially supported by CircleCI in any form. The current endpoint is https://ce-scripts.circleci-support.com/api/workflow/cost-estimate.


Query parameters:
| Parameter | Description |
|--|--|
| circle_token | CircleCI API token |
| workflow_id | The workflow ID to estimate the cost of |
| project_name | The repository name corresponding to the workflow_id |
| project_vcs | The name of the VCS provider, IE gh or github |
| project_user | The user or org name of the project corresponding to the workflow_id |

Example response:
```{
    "total_cost": 12.4,
    "total_credits": 20659,
    "total_runtime": "",
    "disclaimer": "This is a cost estimate. This is not an official CircleCI endpoint. Please contact jacobjohnston@circleci.com for questions."
}
