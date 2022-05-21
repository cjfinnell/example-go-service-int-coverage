import http from 'k6/http';
import exec from 'k6/execution';
import { Httpx } from 'https://jslib.k6.io/httpx/0.0.6/index.js';
import { describe } from 'https://jslib.k6.io/expect/0.0.5/index.js';

export let options = {
	vus: 100,
	duration: '10s',
	thresholds: {
		http_req_duration: [
			'p(50)<80',
			'p(95)<200',
		],
	},
};

let session = new Httpx({
	baseURL: 'http://rediswrapper:8080',
	headers: {
		'Authorization': 'whatever'
	}
})

export default function () {
	let key = exec.vu.idInTest;
	let value = 'content';
	describe('set-get-del', (t) => {
		let res = session.get(`/${key}`);
		t.expect(res.status).as(`miss`).toEqual(404);

		res = session.post(`/${key}/${value}`);
		t.expect(res.status).as(`set`).toEqual(201);

		res = session.get(`/${key}`);
		t.expect(res.status).as(`hit`).toEqual(200);

		res = session.delete(`/${key}`);
		t.expect(res.status).as(`del`).toEqual(200);
	});
}
