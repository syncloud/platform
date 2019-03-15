const backend = require('./index');

const mock = require('../__mocks__/jquery.mockjax');

test('backend', () => { 
	expect(backend).toBeDefined(); 
});
