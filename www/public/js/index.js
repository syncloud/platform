import * as _ from 'underscore';
import $ from 'jquery';
import jQuery from 'jquery';

import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-switch';
import 'bootstrap-switch/dist/css/bootstrap3/bootstrap-switch.css';
import 'font-awesome/css/font-awesome.css'
import '../css/site.css'
import * as UiCommon from './ui/common.js'
import './ui/menu.js'
import Common from './common.js'
import Templates from './index.templates.js'

function installed_apps(on_complete, on_error) {
    $.get('/rest/installed_apps').done(on_complete).fail(on_error);
}

function display_apps(data) {
	let html = _.template(Templates.IndexTemplate)(data);
	$("#block_apps").html(html);
}

$( document ).ready(() => {
 if (typeof mock !== 'undefined') { console.log(mock) };
 installed_apps(
    display_apps,
    UiCommon.ui_display_error);
});
