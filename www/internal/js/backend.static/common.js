function success_callbacks(parameters, data) {
    if (parameters.hasOwnProperty("done")) {
        parameters.done(data);
    }
    if (parameters.hasOwnProperty("always")) {
        parameters.always();
    }
}