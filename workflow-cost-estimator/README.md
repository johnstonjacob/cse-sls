
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
`{
    "total_cost": 12.4,
    "total_credits": 20659,
    "total_runtime": "",
    "disclaimer": "This is a cost estimate. This is not an official CircleCI endpoint. Please contact jacobjohnston@circleci.com for questions."
}`


## Limitations
AWS API Gateways are limited to 30 seconds of execution time, including cold start time. This script is written in Go, so the uploaded binary doesn't require much startup time. The highest job number I have tried this script with is 92, and the response time ( not execution time ) was <10 seconds. 99.9% of users should be fine with this limitation. If you are consistently seeing timeouts, and there is no [statuspage](https://status.circleci.com/) about an outage, please raise an issue with the parameters you passed into the function.

## Deploying

If you don't want to use the provided endpoint for any reason

If you want a custom domain:
Prerequisites: Go, NPM, and serverless installed on your system.

 1. `npm install`
 2. Edit `serverless.yml` with your new domain information
 3. `sls create_domain`
 4. Wait 40~ minutes
 5. `make deploy`

If you don't want a custom domain:
Prerequisites: Go and serverless installed on your system.

 1. Remove plugins and custom sections in `serverless.yml`
 2. `make deploy`

## I don't want to use Lambda for {xzy}
Not currently a supported solution. Decoupling the logic from the Lambda specifics should be fairly straightforward. 

## I want to deploy to $FaaS provider
Not currently a supported solution
 1. Update `serverless.yml` with your chosen provider
 2. Remove Lambda logic and replace with $FaaS provider logic
