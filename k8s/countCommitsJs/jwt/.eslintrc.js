module.exports = {
	env: {
		node: true,
		es2016: true,
	},
	extends: ['standard', 'prettier'],
	parser: '@typescript-eslint/parser',
	parserOptions: {
		ecmaVersion: 'latest',
		sourceType: 'module',
	},
	plugins: ['@typescript-eslint'],
	rules: {},
};
