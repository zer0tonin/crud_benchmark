import http from 'k6/http';
import { parseHTML } from 'k6/html';
import { check } from 'k6';

export const options = {
  discardResponseBodies: false,
  scenarios: {
    write_scenario: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '30s', target: 2 },
        { duration: '30s', target: 5 },
        { duration: '30s', target: 0 },
      ],
      gracefulRampDown: '0s',
      exec: 'write',
    },
    read_scenario: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '10s', target: 0 },
        { duration: '20s', target: 18 },
        { duration: '30s', target: 45 },
        { duration: '30s', target: 0 },
      ],
      gracefulRampDown: '0s',
      exec: 'read',
    },
  },
};

export function write() {
  const targetURL = `http://${__ENV.HOSTNAME}`
  const payload = {
    title: 'Lorem Ipsum',
    body: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum',
  }
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  }
  const res = http.post(targetURL + '/', JSON.stringify(payload), params)
  check(res, {
    'is status 201': r => r.status === 201,
  })
}

export function read() {
  const targetURL = `http://${__ENV.HOSTNAME}`
  let res = http.get(targetURL + '/');
  check(res, {
    'is list ok': r =>
      r.status === 200 && r.body.includes('CRUD Benchmark'),
  });

  const choice = Math.floor(Math.random() * 10)
  const link = parseHTML(res.body).find('a').get(choice).getAttribute('href')
  res = http.get(targetURL + link)
  check(res, {
    'is get ok': r => r.status === 200 && r.body.includes('Lorem Ipsum'),
  })
}
