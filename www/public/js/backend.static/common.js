window.onerror = function(msg, url, linenumber) {
    alert('Error message: '+msg+'\nURL: '+url+'\nLine Number: '+linenumber);
    return true;
};

function success_callbacks(parameters, data) {
    if (parameters.hasOwnProperty("done")) {
        parameters.done(data);
    }
    if (parameters.hasOwnProperty("always")) {
        parameters.always();
    }
}