function on_error(status, error) {
    if (status === 401) {
        window.location.href = "/login.html";
    } else {
        window.location.href = "/error.html";
    }
}
