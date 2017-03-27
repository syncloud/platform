QUnit.test( "job status", function( assert ) {

    var checker_count = 0;
    var checker_on_complete;
    var checker_job;
    var checker = function (job, on_complete, on_error) {
        checker_count += 1;
        checker_job = job;
        checker_on_complete = on_complete;
        this.on_error = on_error;
    };
    
    run_after_job_is_complete(checker, function(func, timeout) { func(); }, function() {}, function(a, b, c) {}, 'test');

    checker_on_complete({is_running: true});

    assert.deepEqual( checker_job, 'test');
    assert.deepEqual( checker_count, 2);
});