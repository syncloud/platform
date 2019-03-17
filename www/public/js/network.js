import * as _ from 'underscore';

import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-switch';
import 'bootstrap-switch/dist/css/bootstrap3/bootstrap-switch.css';
import 'font-awesome/css/font-awesome.css'
import '../css/site.css'

import './ui/font.js'
import UiCommon from './ui/common.js'
import './ui/menu.js'
import * as Common from './common.js'
import './backend/menu.js'
import Templates from './network.templates.js'

function check_access(on_complete, on_error) {
    $.get('/rest/access/access').done(on_complete).fail(on_error);
}

function set_access(data, on_complete, on_error) {
    $.get('/rest/access/set_access', data).done(on_complete).fail(on_error);
}

function network_interfaces(on_complete, on_error) {
    $.get('/rest/access/network_interfaces').done(on_complete).fail(on_error);
}

function port_mappings(on_complete, on_error) {
    $.get('/rest/access/port_mappings').done(on_complete).fail(on_error);
}

function ui_display_toggles() {
	$("[type='checkbox']").each(function() {
		$(this).bootstrapSwitch();
	});
};

function ui_display_network(data) {
    $("#block_network").html(_.template(Templates.NetworkTemplate)(data.data));
}

function ui_display_port_mappings(data) {

    var certificate_port_mapping = _.find(data.port_mappings, function(mapping) {return mapping.local_port == 80;});
    var certificate_port = 0;
    if (certificate_port_mapping) {
        certificate_port = certificate_port_mapping.external_port;
    }
    $("#certificate_port").val(certificate_port);
    if (certificate_port != 80) {
        $("#certificate_port_warning").show('slow');
    } else {
        $("#certificate_port_warning").hide('slow');
    }

    var access_port_mapping = _.find(data.port_mappings, function(mapping) {return mapping.local_port == 443;});
    var access_port = 0;
    if (access_port_mapping) {
        access_port = access_port_mapping.external_port
    }
    $("#access_port").val(access_port);
    if (access_port != 443) {
        $("#access_port_warning").show('slow');
    } else {
        $("#access_port_warning").hide('slow');
    }

}

function ui_display_access(data) {
    var access_data = data.data;

    $('#tgl_external').bootstrapSwitch('disabled', false);
    $('#tgl_external').bootstrapSwitch('state', access_data.external_access, true);
    $("#tgl_external_loading").removeClass('opacity-visible');
   
    $('#tgl_ip_autodetect').bootstrapSwitch('disabled', false);
    
    $("#tgl_ip_autodetect_loading").removeClass('opacity-visible');
    var ip_autodetect_enabled;
    if (access_data.hasOwnProperty('public_ip')) {
        ip_autodetect_enabled = false;
        $("#public_ip").val(access_data.public_ip);
    } else {
        ip_autodetect_enabled = true;
        $("#public_ip").val('');
    }

    $('#tgl_ip_autodetect').bootstrapSwitch('state', ip_autodetect_enabled);

    $('#tgl_upnp').bootstrapSwitch('disabled', false);
    $('#tgl_upnp').bootstrapSwitch('state', access_data.upnp_enabled);
    $("#tgl_upnp_loading").removeClass('opacity-visible');

    if (access_data.upnp_available) {
        $("#upnp_warning").hide('slow');
    } else {
        $("#upnp_warning").show('slow');
    }
    
    if (access_data.upnp_available) {
        $("#label_upnp").css("color", "black");
    } else {
        $("#label_upnp").css("color", "red");
    }
    
    $("#btn_save").button('reset');

    $('#tgl_external').bootstrapSwitch('disabled', false);
    $('#tgl_upnp').bootstrapSwitch('disabled', false);

    ui_prepare_external_access();
    ui_prepare_address();
    ui_upnp();

}

function disable_access_controls(disabled) {
    $('#tgl_external').bootstrapSwitch('disabled', disabled);
    $('#tgl_ip_autodetect').bootstrapSwitch('disabled', disabled);
    $('#tgl_upnp').bootstrapSwitch('disabled', disabled);
    $('#public_ip').prop('disabled', disabled);
    $('#certificate_port').prop('disabled', disabled);
    $('#access_port').prop('disabled', disabled);
}

function ui_check_access() {
    disable_access_controls(true);

    $("#tgl_external_loading").addClass('opacity-visible');
    $("#tgl_upnp_loading").addClass('opacity-visible');
    $("#tgl_ip_autodetect_loading").addClass('opacity-visible');
    $('#btn_save').button('loading');

    check_access(
        (data) => {
            Common.check_for_service_error(
                data,
                () => ui_display_access(data), 
                UiCommon.ui_display_error);
        },
        UiCommon.ui_display_error);

    port_mappings(ui_display_port_mappings, UiCommon.ui_display_error);
}

function ui_check_network() {
    network_interfaces(ui_display_network, UiCommon.ui_display_error);
}

function ui_prepare_external_access() {
		var toggle = $("#tgl_external");
		var enabled = toggle.bootstrapSwitch('state');
		if (enabled) {
            $("#external_block").show('slow');
        } else {
            $("#external_block").hide('slow');
        }
}

function ui_prepare_address() {
		var toggle = $("#tgl_ip_autodetect");
		var enabled = toggle.bootstrapSwitch('state');
		$('#public_ip').prop('disabled', enabled);
}

function ui_upnp() {
		var toggle = $("#tgl_upnp");
		var enabled = toggle.bootstrapSwitch('state');
		$('#certificate_port').prop('disabled', enabled);
		$('#access_port').prop('disabled', enabled);
}

function isValidPort(port) {
    return Number.isNaN(port) || port < 1 || port > 65535
}
function error(message) {
    return {
        status: 200,
        responseJSON: {
            message: message
        }
    }
}

export function set_access(
        upnp_enabled,
        external_access,
        ip_autodetect,
        public_ip,
        certificate_port,
        access_port,
        on_complete,
        on_error) {

        var request_data = {
           upnp_enabled: upnp_enabled,
           external_access: external_access,
           certificate_port: certificate_port,
           access_port: access_port
        };

        if (!ip_autodetect) {
            request_data.public_ip = public_ip;
        }

        set_access(request_data, on_complete, on_error);
    };
    
$(document).ready(function () {
    if (typeof mock !== 'undefined') { console.log("backend mock") };

    ui_display_toggles();

    $("#tgl_external").on('switchChange.bootstrapSwitch', function (event, state) {
        event.preventDefault();
        ui_prepare_external_access();
    });

    $("#tgl_ip_autodetect").on('switchChange.bootstrapSwitch', function (event, state) {
        event.preventDefault();
        ui_prepare_address();
    });

    $("#tgl_upnp").on('switchChange.bootstrapSwitch', function (event, state) {
        event.preventDefault();
        ui_upnp();
    });

    $("#btn_save").click(function ( event ) {
        event.preventDefault();
        
        var access_toggle = $("#tgl_external");
        var upnp_enabled = $("#tgl_upnp").bootstrapSwitch('state');
        var ip_autodetect = $("#tgl_ip_autodetect").bootstrapSwitch('state');

        var certificate_port_string = $("#certificate_port").val();
        var certificate_port = parseInt(certificate_port_string);

        var access_port_string = $("#access_port").val();
        var access_port = parseInt(access_port_string);

        var public_ip = $("#public_ip").val().trim();

        if (upnp_enabled) {
            certificate_port = 0;
            access_port = 0;
        } else {
            
            if (isValidPort(certificate_port)) {
                UiCommon.ui_display_error(error("certificate port (" + certificate_port_string + ") has to be between 1 and 65535"), {}, {});
                return;
            }
            if (isValidPort(access_port)) {
                UiCommon.ui_display_error(error("access port (" + access_port_string + ") has to be between 1 and 65535"), {}, {});
                return;
            }
        }
        
        disable_access_controls(true);
        var btn = $(this);
        btn.button('loading');

        set_access(
                upnp_enabled,
                access_toggle.bootstrapSwitch('state'),
                ip_autodetect,
                public_ip,
                certificate_port,
                access_port,
                function(data) {
                    Common.check_for_service_error(
                        data,
                        ui_check_access,
                        function (xhr, textStatus, errorThrown) {
                            UiCommon.ui_display_error(xhr, textStatus, errorThrown);
                            ui_check_access();
                        }
                    );
                },
                function (xhr, textStatus, errorThrown) {
                    UiCommon.ui_display_error(xhr, textStatus, errorThrown);
                    ui_check_access();
                }
        );
    });

    ui_check_access();
    ui_check_network();

});
