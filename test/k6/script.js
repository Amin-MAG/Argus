import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
    stages: [
        { duration: '25s', target: 10 },
        { duration: '35s', target: 30 },
        { duration: '25s', target: 20 },
        { duration: '20s', target: 0 },
    ],
};

export default function () {
    const res = http.get('https://argus.ghasvari.com/api/v1/health/ping');
    sleep(1);
}
