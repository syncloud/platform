import './backend/index.js'
import * as _ from 'underscore'
import './apps.template.js'
import '../css/site.css'
import '../css/bootstrap.css'
import '../css/bootstrap-switch.css'
import '../css/font-awesome.css'

function display_apps(data) {                   	$("#block_apps").html(_.template(apps_template)(data));                              }                                                                                               $( document ).ready(function() {                                backend.installed_apps(                                         display_apps,                                   ui_display_error                );                              });
