const Network = require('./network');
const mock = require('../__mocks__/jquery.mockjax');

test('network save access', () => {
  $.ajaxSetup({ async: false });
  var response;
  Network.set_access(true,
        true,
        true,
        '1.1.1.1',
        80,
        443,
        function (data) {
            response = data;
        },
        function () {}
      );
  expect(response).toBeDefined();
});
