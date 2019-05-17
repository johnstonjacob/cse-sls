#! /usr/local/bin/node
const axios = require('axios');
const credit_price = '0.0006';
const resource_classes = {
    docker: {
        small: 5,
        medium: 10,
        'medium+': 15,
        large: 20,
        xlarge: 40,
        '2xlarge': 80,
        '2xlarge+': 100,
        '3xlarge': 160,
        '4xlarge': 320,
    },
    machine: {
        small: 5,
        medium: 10,
        large: 20,
        xlarge: 40,
        '2xlarge': 80,
        '3xlarge': 120,
    },
    macOS: {},
    GPU: {},
    windows: {},
};

let circle_url = 'https://circleci.com';
const v1_url = circle_url + '/api/v1/';
const v2_url = circle_url + '/api/v2/';

//TODO: retrieve this information from API
let project_name, project_user;

let circle_token;
let workflowID;


async function get_workflow_status() {
    const workflow_params = {
        method: 'GET',
        url: v2_url + 'workflow/' + workflowID,
        auth: {
            username: circle_token,
        },
    };
    try {
        const res = await axios(workflow_params);
        const { data } = res;

        if (data.status !== 'success' && data.status !== 'failed') {
            console.error(`Workflow status is ${data.status}. Exiting..`);
            console.error('Workflow status must be "success" or "failed" to estimate cost.');
            return {
                statusCode: 202,
                body: {
                    message: `Workflow status is ${
                        data.status
                    }. Workflow status must be "success" or "failed to estimate cost."`,
                },
            };
        }

        console.log(`Workflow status is ${data.status}. Getting jobs ids..`);
    } catch (error) {
        console.error(error.response.data);
        return {
            statusCode: 500,
            body: {
                error: error.response.data,
            },
        };
    }
}

async function get_workflow_jobs() {
    const job_id_params = {
        method: 'GET',
        url: v2_url + 'workflow/' + workflowID + '/jobs',
        auth: {
            username: circle_token,
        },
    };

    try {
        const res = await axios(job_id_params);
        const { data } = res;
        return { jobs: data['jobs'], ok: true };
    } catch (error) {
        return {
            res: {
                statusCode: 500,
                body: {
                    error: error.response.data,
                },
            },
            ok: false,
        };
    }
}

async function tallyJobCost(jobs) {
    const totalCost = 0

    return {
        ok: true,
        totalCost,
    }
}

async function main() {
    if (!circle_token || !workflowID || !project_name || !project_user) {
        return {
            statusCode: 401,
            body: {
                error: 'Please provide circle_token AND workflow_id parameters.',
            },
        };
    }

    const workflowRes = await get_workflow_status();
    if (workflowRes !== undefined) return workflowRes;

    const { jobs, ok, res } = await get_workflow_jobs();

    if (!ok) return res;

    const = await tallyJobCost(jobs);


    return {
        statusCode: 200,
        body: {
            jobs,
        },
    };
}

module.exports.workflow_cost_estimate = async (event) => {
    if (event['queryStringParameters']) {
        ({ workflowID, circle_token, project_name, project_user } = event['queryStringParameters']);
    }

    const res = await main();
    res.body['disclaimer'] =
        'NOTICE - THIS IS A COST ESTIMATE. THERE IS NO GUARANTEE THIS ESTIMATE IS CORRECT.';
    res.body = JSON.stringify(res.body, null, 2);
    return res;
};
