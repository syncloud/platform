var Mobile;

function Desktop() {}
Desktop.prototype.getApiUrl = function() {};
Desktop.prototype.getDomain = function() {};
Desktop.prototype.getRedirectLogin = function() {};
Desktop.prototype.getRedirectPassword = function() {};
Desktop.prototype.saveCredentials = function(mac_address, user, pass) {};


if (typeof Android !== 'undefined') {
    Mobile = new AndroidMobile();
} else {
    Mobile = new Desktop();
}