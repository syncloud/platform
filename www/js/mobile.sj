var Mobile;

function Desktop() {}
Desktop.prototype.getRedirectLogin = function() {};
Desktop.prototype.getRedirectPassword = function() {};


if (typeof Android !== 'undefined') {
    Mobile = new AndroidMobile();
} else {
    Mobile = new Desktop();
}