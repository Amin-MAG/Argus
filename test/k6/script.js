import http from 'k6/http';
import {check, sleep} from 'k6';

export const options = {
    stages: [
        {duration: '25s', target: 10},
        {duration: '35s', target: 30},
        {duration: '25s', target: 20},
        {duration: '20s', target: 0},
    ],
};

function getRandomIp() {
    const segments = [];
    for (let i = 0; i < 4; i++) {
        segments.push(Math.floor(Math.random() * 256));
    }
    return segments.join('.');
}

function testCreateAgent() {
    const url = 'https://argus.ghasvari.com/api/v1/agents';
    const payload = JSON.stringify({
        ip_address: getRandomIp(),
    });
    const params = {
        headers: {
            'Content-Type': 'application/json',
            'API-Key': "test_api_key",
        },
    };
    const createRes = http.post(url, payload, params);
    check(createRes, {
        'create agent status is 201': (r) => r.status === 201,
    });

    console.log(createRes.body)
    const agentId = JSON.parse(createRes.body)["agent"].id;
    const getRes = http.get(`${url}/${agentId}`, params);
    check(getRes, {
        'get agent status is 200': (r) => r.status === 200,
    });
}

function getRandomQueryParams() {
    const page = Math.floor(Math.random() * 10) + 1; // Random page number between 1 and 10
    const pageSize = Math.floor(Math.random() * 20) + 1; // Random page size between 1 and 20
    const ipAddress = getRandomIp();
    const sortBy = 'id';
    const order = Math.random() > 0.5 ? 'asc' : 'desc'; // Random order 'asc' or 'desc'

    return `page=${page}&page_size=${pageSize}&ip_address=${ipAddress}&sort_by=${sortBy}&order=${order}`;
}

function testGetAgents() {
    const baseUrl = 'https://argus.ghasvari.com/api/v1/agents';
    const queryParams = getRandomQueryParams();
    const url = `${baseUrl}?${queryParams}`;
    const params = {
        headers: {
            'Content-Type': 'application/json',
            'API-Key': "test_api_key",
        },
    };

    // Send GET request to retrieve list of agents
    const res = http.get(url, params);

    // Check the response status and content
    check(res, {
        'get all agents status is 200': (r) => r.status === 200,
        'get all agents response contains agents': (r) => JSON.parse(r.body).length > 0,
    });
}

export default function () {
    let res = http.get('https://argus.ghasvari.com/api/v1/health/ping', {
        headers: {
            'Content-Type': 'application/json',
            'API-Key': "test_api_key",
        },
    });
    console.log(`Ping Response status: ${res.status}`);

    testCreateAgent();

    testGetAgents();

    sleep(1);
}
