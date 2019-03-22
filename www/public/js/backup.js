import * as _ from 'underscore';
import $ from 'jquery';
import jQuery from 'jquery';
import dateFormat from 'dateformat';
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
    columnDefs: [
        {
            headerName: 'App',
            field: 'app',
            cellStyle: { 'text-align': "left" }
        },
        {
            headerName: 'Date',
            field: 'date',
            cellRenderer: (data) => {
                return dateFormat(data.date, 'mm/dd/yyyy HH:MM')
            },
            cellStyle: { 'text-align': "left" }
        }
    ],
    columnTypes: {
        "dateColumn": {
            filter: 'agDateColumnFilter',
            suppressMenu:true
        }
    },
    rowData: [
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },
        { app: 'Nextcloud', date:  1294862756114  },

    ]
};

$( document ).ready(function () {
  let eGridDiv = document.querySelector('#backupGrid');
  
  let grid = new Grid(eGridDiv, gridOptions);
  gridOptions.api.sizeColumnsToFit();
});
