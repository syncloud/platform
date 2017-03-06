QUnit.test( "sam after complete", function( assert ) {

    assert.expect( 1 );
    var done = assert.async();

    $.mockjax({
        url: "/rest/settings/sam_status",
        responseText: {
            is_running: false
        }
    });

    run_after_sam_is_complete(function() {
        assert.ok( true, 'complete');
        done();
    });

});