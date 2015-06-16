function AndroidMobile() {}
AndroidMobile.prototype.getRedirectLogin = function() {
    return Android.getRedirectLogin();
};

AndroidMobile.prototype.getRedirectPassword = function() {
    return Android.getRedirectPassword();
};

AndroidMobile.prototype.saveCredentials = function(name, password) {
    Android.saveCredentials(name, password);
};