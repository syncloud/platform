function onError(xhr, textStatus, errorThrown) {
    if (xhr.status === 401) {
        window.location.href = "/login.html";
    } else {
        window.location.href = "/error.html";
    }
}