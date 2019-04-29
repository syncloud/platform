import * as Common from './common.js'
import * as Mock from '../__mocks__/jquery.mockjax.js'

test('job status', () => {

   $.ajaxSetup({ async: false });

    var on_complete_count = 0;
    Common.run_after_job_is_complete(
        function(func, timeout) { func(); },
        function() {
            on_complete_count += 1;
        },
        function(a, b, c) {},
        Common.INSTALLER_UPDATE_URL,
        (resp) => { return false; }
        );

    expect(on_complete_count).toEqual(1);
});

test('find app', () => { 
    let app = Common.find_app(Mock.versions_data.data, 'platform');
    expect(app.app.name).toEqual('Platform');
});
