import * as Common from './common.js'

test( "show apps success", () => {

  $.ajaxSetup({ async: false });

  Common.check_activation_status();

});

