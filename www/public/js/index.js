import * as _ from 'underscore';
import $ from 'jquery';

import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap-switch';
import 'bootstrap-switch/dist/css/bootstrap3/bootstrap-switch.min.css';

import '../css/site.css'
import '../css/font-awesome.css'

import './ui/font.js'
import './ui/common.js'
import './ui/menu.js'

import './common.js'
import './backend.static/common.js'
import './backend.static/menu.js'
import './backend.static/index.js'
import './apps.template.js'

function display_apps(data) {
	$("#block_apps").html(_.template(apps_template)(data));                              

}

$( document ).ready(function() {
	backend.installed_apps(display_apps, ui_display_error);
});
