module.exports = {
	env: {
		es6: true
	},

	parserOptions: {
		sourceType: "module",
		ecmaVersion: 2020,
		ecmaFeatures: {
			jsx: false
		}
	},

	rules: {
		"no-var": "error",
		"space-in-parens": "error",
		"no-multiple-empty-lines": "error",
		"prefer-const": "error",
		"no-use-before-define": "error"
	}
}