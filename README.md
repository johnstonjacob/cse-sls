# cse-sls

Serverless functions for the CCI CSE org.


Current functions:

[Workflow Cost Estimator](https://github.com/johnstonjacob/cse-sls/tree/master/src/workflow-cost-estimator)


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

## I don't want to use Lambda for {xzy} reason.
Not currently a supported solution. However, decoupling the logic from the Lambda specifics should be fairly straightforward. 

## I want to deploy to $FaaS provider
Not currently a supported solution
 1. Update `serverless.yml` with your chosen provider
 2. Remove Lambda logic and replace with $FaaS provider logic
