QUnit.test( "get redirect login", function( assert ) {
  Mobile = new Desktop();
  assert.ok( !Mobile.getRedirectLogin());
});

QUnit.test( "get redirect login on Mobile", function( assert ) {
  Mobile = new TestMobile();
  assert.equal( Mobile.getRedirectLogin(), "login");
});

QUnit.test( "get redirect password", function( assert ) {
  Mobile = new Desktop();
  assert.ok( !Mobile.getRedirectPassword());
});

QUnit.test( "get redirect password on Mobile", function( assert ) {
  Mobile = new TestMobile();
  assert.equal( Mobile.getRedirectPassword(), "pass");
});

function TestMobile() {}

TestMobile.prototype.getRedirectLogin = function() {
  return "login";
};

TestMobile.prototype.getRedirectPassword = function() {
  return "pass";
};