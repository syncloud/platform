import * as _ from 'underscore';
import $ from 'jquery';
import jQuery from 'jquery';
import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-switch';
import 'bootstrap-switch/dist/css/bootstrap3/bootstrap-switch.css';
import 'font-awesome/css/font-awesome.css'
import '../css/site.css'
import UiCommon from './ui/common.js'
import './ui/menu.js'

import Common from './common.js'
import {Grid} from "ag-grid-community";
import "ag-grid-community/dist/styles/ag-grid.css";
import "ag-grid-community/dist/styles/ag-theme-balham.css";

const gridOptions = {
    defaultColDef: {
        cellStyle: { 'text-align': "left" },
        
    },
    columnDefs: [
        {
            headerName: 'File',
            field: 'file',
            resizable: true,
            sortable: true,
            filter: true,
            filter: 'agTextColumnFilter'
        },
        {
            headerName: 'Actions',
            width: 100,
            resizable: false,
            cellRenderer: (params) => { 
              var div = document.createElement('div');
              div.innerHTML = `
                <i class='fa fa-undo' style='padding-left: 20px'></i>
                <i class='fa fa-trash' style='padding-left: 20px'></i>
             `;
              var buttons = div.querySelectorAll('i');
              buttons[0].addEventListener('click', () => { 
                $.post(
                 '/rest/backup/restore', 
                 { file: params.data.file },
                 (data) => {
                   alert("done");
                   params.api.redrawRows();
                 }
                );
              }); 
              buttons[1].addEventListener('click', () => { 
               $.post(
                 '/rest/backup/remove', 
                 { file: params.data.file }
                 
                )
                .fail(function(a,b,c) {alert("failed");})
                .done(function(data) {
                   alert(data);
                   params.api.redrawRows();});
              }); 
              return div;
            }
        },
      
    ],
    suppressDragLeaveHidesColumns: true,
    floatingFilter: true,
    domLayout: 'autoHeight'
};

$( document ).ready(function () {
  if (typeof mock !== 'undefined') { console.log("backend mock") };

  let eGridDiv = document.querySelector('#backupGrid');
  
  let grid = new Grid(eGridDiv, gridOptions);
  $.getJSON('/rest/backup/list')
      .done((response) => { 
        gridOptions.api.setRowData(response.data);
        gridOptions.api.sizeColumnsToFit();
      });
});
