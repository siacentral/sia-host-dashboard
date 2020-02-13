async function sendJSONRequest(url, method, data) {
	let headers = {};

	const r = await fetch(url, {
			method: method,
			mode: 'cors',
			cache: 'no-cache',
			headers: headers,
			body: data ? JSON.stringify(data) : null
		}),
		resp = await r.json();

	return { statusCode: r.status >= 200 && r.status < 300 ? 200 : r.status, body: resp };
}

export async function getAverageSettings() {
	const resp = await sendJSONRequest(`https://api.siacentral.com/v2/hosts/settings/average`, 'GET', null, true);

	if (resp.statusCode !== 200)
		throw new Error(resp.body.message);

	return resp.body.settings;
}

export async function getCoinPrice() {
	const resp = await sendJSONRequest('https://api.siacentral.com/v2/market/exchange-rate', 'GET', null, true);

	if (resp.statusCode !== 200)
		throw new Error(resp.body.message);

	return resp.body;
}

export async function getSnapshots(end) {
	if (!end)
		end = new Date();

	const resp = await sendJSONRequest(`${process.env.VUE_APP_API_BASE_URL}/api/snapshots?end=${Math.round(end.getTime() / 1000)}`, 'GET', null, true);

	if (resp.statusCode !== 200)
		throw new Error(resp.body.message);

	if (!Array.isArray(resp.body.snapshots))
		return [];

	return resp.body.snapshots;
}

export async function getStatus() {
	const resp = await sendJSONRequest(`${process.env.VUE_APP_API_BASE_URL}/api/status`, 'GET', null, true);

	if (resp.statusCode !== 200)
		throw new Error(resp.body.message);

	return resp.body;
}

export async function getTotals(end) {
	if (!end)
		end = new Date();

	const resp = await sendJSONRequest(`${process.env.VUE_APP_API_BASE_URL}/api/totals?date=${Math.round(end.getTime() / 1000)}`, 'GET', null, true);

	if (resp.statusCode !== 200)
		throw new Error(resp.body.message);

	return resp.body;
}