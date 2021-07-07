function createRequest(url, body) {
	return new Request(url, {
		method: "POST",
		credentials: "same-origin",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(body),
	});
}

// This is a multi-line comment using double slashes
// and should render appropriately on output.
export async function getAccountDetails(accountEmail) {
	const res = await fetch(createRequest("/rpc/account.Accounts/GetAccountDetails", { "accountEmail": accountEmail }));
	const jsonBody = await res.json();
	if (res.ok) {
		return jsonBody;
	}
	throw new Error(jsonBody.msg);
}

// This is a multi-line block comment, which is
//also expected to render properly.
export async function logout(accountId, token) {
	const res = await fetch(createRequest("/rpc/account.Accounts/Logout", { "accountId": accountId, "token": token }));
	const jsonBody = await res.json();
	if (res.ok) {
		return jsonBody;
	}
	throw new Error(jsonBody.msg);
}

// Creates a checkout session for the given item.
export async function createCheckoutSession(itemId) {
	const res = await fetch(createRequest("/rpc/shop.Shop/CreateCheckoutSession", { "itemId": itemId }));
	const jsonBody = await res.json();
	if (res.ok) {
		return jsonBody;
	}
	throw new Error(jsonBody.msg);
}