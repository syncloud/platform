function on_error(status, error) {
    if (status === 401) {
        window.location.href = "/login.html";
    } else if (status === 0) {
        console.log('user navigated away from the page');
    } else {
        window.location.href = "/error.html";
    }
}
