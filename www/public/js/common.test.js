const Common = require('./common');

test('job status', () => { 
    var checker_count = 0;
    var checker_on_complete;
    var checker_job;
    var checker_on_error;
    var checker = function (job, on_complete, on_error) {
        checker_count += 1;
        checker_job = job;
        checker_on_complete = on_complete;
        checker_on_error = on_error;
    };
    
    Common.run_after_job_is_complete(checker, function(func, timeout) { func(); }, function() {}, function(a, b, c) {}, 'test');

    checker_on_complete({is_running: true});

    expect(checker_job).toEqual('test');
    expect(checker_count).toEqual(2);
});