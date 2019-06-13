
# Workflow Cost Estimator

The Workflow Cost Estimator is a lambda function with the goal of estimating a given workflow cost. This is not anything more than an estimate. This is not officially supported by CircleCI in any form. The current endpoint is https://ce-scripts.circleci-support.com/api/workflow/cost-estimate/$WORKFLOW_ID.


Query parameters:

| Parameter | Description |
|--|--|
| circle_token | CircleCI API token |
| circle_url | (Optional) Location of your CircleCI install (if applicable) Default: https://circleci.com |

**NOTE: This script does not currently support CircleCI server.**

Example response:
```
{
    "total_cost": 0.01,
    "total_credits": 8,
    "total_runtime": "46s",
    "disclaimer": "This is a cost estimate. This is not an official CircleCI endpoint. Please contact jacobjohnston@circleci.com for questions.",
    "jobs": [
        {
            "job_name": "Blocked - test_parallel",
            "total_cost": 0,
            "total_credits": 0,
            "total_runtime": "0s"
        },
        {
            "job_name": "Blocked - approve_deploy",
            "total_cost": 0,
            "total_credits": 0,
            "total_runtime": "0s"
        },
        {
            "job_name": "Blocked - some_other_job",
            "total_cost": 0,
            "total_credits": 0,
            "total_runtime": "0s"
        },
        {
            "job_name": "Blocked - deploy",
            "total_cost": 0,
            "total_credits": 0,
            "total_runtime": "0s"
        },
        {
            "job_name": "test",
            "total_cost": 0.01,
            "total_credits": 7.666666666666667,
            "total_runtime": "46s"
        }
    ]
}
```

## Limitations
AWS API Gateways are limited to 30 seconds of execution time, including cold start time. This script is written in Go, so the uploaded binary doesn't require much startup time. The highest job number I have tried this script with is 92, and the response time ( not execution time ) was <10 seconds. 99.9% of users should be fine with this limitation. If you are consistently seeing timeouts, and there is no [statuspage](https://status.circleci.com/) about an outage, please raise an issue with the parameters you passed into the function.
