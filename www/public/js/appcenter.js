import * as _ from 'underscore';
import $ from 'jquery';
import jQuery from 'jquery';

import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-switch';
import 'bootstrap-switch/dist/css/bootstrap3/bootstrap-switch.css';
import 'font-awesome/css/font-awesome.css'
import '../css/site.css'

import './ui/font.js'
import UiCommon from './ui/common.js'
import './ui/menu.js'
import Common from './common.js'
import './backend/common.js'
import './backend/menu.js'
import './backend/appcenter.js'

import Templates from './appcenter.templates.js'

function available_apps(on_complete, on_error) {

    backend.available_apps(
         function (data) {
            Common.check_for_service_error(
                data,
                function() {
                    on_complete(data);
                },
                on_error);
         },
         on_error
    );
}
function display_apps(data) {
		$("#block_apps").html(_.template(Templates.AppsTemplate)(data));
}

$( document ).ready(function() {
  if (typeof mock !== 'undefined') { console.log("backend mock") };
  available_apps(
    display_apps,
    UiCommon.ui_display_error
		);
});


module.exports = {
  available_apps
};
