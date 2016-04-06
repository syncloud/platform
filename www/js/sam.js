function run_after_sam_is_complete(on_complete) {

    var recheck_function = function () { run_after_sam_is_complete(on_complete); }

    var recheck_timeout = 2000;
    $.get('/rest/settings/sam_status')
            .done(function(sam) {
                if (sam.is_running)
                    setTimeout(recheck_function, recheck_timeout);
                else
                    on_complete();
            })
            .fail(function() {
                setTimeout(recheck_function, recheck_timeout);
            })

}
