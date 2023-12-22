import fetch from 'node-fetch';
import * as dotenv from 'dotenv';
const fs = require('fs');
const jwt = require('jsonwebtoken');
dotenv.config();

type JwtPayload = {
	iat: number;
	exp: number;
	iss: string;
};

(async (): Promise<void> => {
	if (!process.env.APP_ID) return;
	const now: number = Math.floor(Date.now() / 1000);
	const jwtPayload: JwtPayload = {
		iat: now,
		exp: now + 60 * 10 - 30,
		iss: process.env.APP_ID,
	};
	const jwtSecret = fs.readFileSync('private-key.pem');
	const jwtOptions = {
		algorithm: 'RS256',
	};

	const token = jwt.sign(jwtPayload, jwtSecret, jwtOptions);

	const accessTokenRes = await fetch(
		`https://api.github.com/app/installations/${process.env.INSTALL_ID}/access_tokens`,
		{
			method: 'POST',
			headers: {
				Authorization: `Bearer ${token}`,
				Accept: 'application/vnd.github+json',
			},
		},
	).then(res => res.json());
	console.log(accessTokenRes.token);
})();
