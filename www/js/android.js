function AndroidMobile() {}

AndroidMobile.prototype.getRedirectLogin = function() {
    return Android.getRedirectLogin();
};

AndroidMobile.prototype.getRedirectPassword = function() {
    return Android.getRedirectPassword();
};

AndroidMobile.prototype.saveCredentials = function(mac_address, name, password) {
    Android.saveCredentials(mac_address, name, password);
};