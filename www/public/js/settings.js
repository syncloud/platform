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
import './backend/menu.js'

function check_versions(on_complete, on_error) {

    backend.check_versions(function () {
        Common.run_after_sam_is_complete(
            backend.job_status,
            setTimeout,
            function () {
                backend.get_versions(
                    on_complete,
                    on_error);
            }, on_error);
        }, on_error);
}

function platform_upgrade(on_complete, on_error) {

    backend.platform_upgrade(function (data) {
        Common.check_for_service_error(data, function () {
            Common.run_after_sam_is_complete(
                backend.job_status,
                setTimeout,
                function () {
                    backend.get_versions(
                         on_complete,
                         on_error);
                 }, on_error);
        }, on_error);
    }, on_error);
    
}

function sam_upgrade(on_complete, on_error) {

    backend.sam_upgrade(function (data) {
        Common.check_for_service_error(data, function () {
            Common.run_after_sam_is_complete(
                backend.job_status,
                setTimeout,
                function () {
                    backend.get_versions(
                        on_complete,
                        on_error);
                }, on_error);
        }, on_error);
    }, on_error);
    
}

module.exports = {
	check_versions,
	platform_upgrade,
	sam_upgrade
};