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

import './common.js'
import './backend/common.js'
import './backend/menu.js'
import './backend/index.js'
import template from './apps.template.js'

function display_apps(data) {
	let html = _.template(template)(data);
	$("#block_apps").html(html);
}

$( document ).ready(() => {
	backend.installed_apps(
    display_apps,
    UiCommon.ui_display_error);
});
