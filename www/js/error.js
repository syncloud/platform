function onError(xhr, textStatus, errorThrown) {
    if (xhr.status === 401) {
        window.location.href = "/server/html/login.html";
    } else {
        window.location.href = "/server/html/error.html";
    }
}