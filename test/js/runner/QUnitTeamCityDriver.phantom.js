/**
* Wait until the test condition is true or a timeout occurs. Useful for waiting
* on a server response or for a ui change (fadeIn, etc.) to occur.
*
* @param testFx javascript condition that evaluates to a boolean,
* it can be passed in as a string (e.g.: '1 == 1' or '$("#bar").is(":visible")' or
* as a callback function.
* @param onReady what to do when testFx condition is fulfilled,
* it can be passed in as a string (e.g.: '1 == 1' or '$("#bar").is(":visible")' or
* as a callback function.
* @param timeOutMillis the max amount of time to wait. If not specified, 3 sec is used.
*/
var system = require('system');
var webpageUrl = 'http:///www.example.com';

function waitFor(testFx, onReady, timeOutMillis) {
    var maxtimeOutMillis = timeOutMillis ? timeOutMillis : 3001, //< Default Max Timout is 3s
        start = new Date().getTime(),
        condition = false,
        interval = setInterval(function () {
            if ((new Date().getTime() - start < maxtimeOutMillis) && !condition) {
                // If not time-out yet and condition not yet fulfilled
                condition = (typeof (testFx) === 'string' ? eval(testFx) : testFx()); //< defensive code
            } else {
                if (!condition) {
                    // If condition still not fulfilled (timeout but condition is 'false')
                    console.log('waitFor() timeout');
                    phantom.exit(1);
                } else {
                    // Condition fulfilled (timeout and/or condition is 'true')
                    typeof (onReady) === 'string' ? eval(onReady) : onReady(); //< Do what it's supposed to do once the condition is fulfilled
                    clearInterval(interval); //< Stop this interval
                }
            }
        }, 100); //< repeat check every 250ms
};

// phantoms.args only available in PhantomJS 1.x http://phantomjs.org/api/phantom/property/args.html
if (phantom.args) {
	if (phantom.args.length === 0 || phantom.args.length > 2) {
		console.log('Usage: run-qunit.js URL');
		phantom.exit(1);
	} else {
		webpageUrl = phantom.args[0];
	}
} else if (system.args) {
	if (system.args.length === 0 || system.args.length < 2) {
		console.log('Usage: run-qunit.js URL');
		phantom.exit(1);
	} else {
		webpageUrl = system.args[1];
		webpageUrl = webpageUrl.replace(/\\/g, '/');
	}	
} else {
	console.log('Unable to access args from command line.  Check PhantomJS docs');
}

var page = new WebPage();

// Route 'console.log()' calls from within the Page context to the main Phantom context (i.e. current 'this')
page.onConsoleMessage = function (msg) {
    console.log(msg);
};

page.onError = function(error) {
	console.error('(OnError) ' + error.url + ': ' + error.errorString);
}

page.onResourceError = function(resourceError) {
	console.error('(OnResourceError) ' + resourceError.url + ': ' + resourceError.errorString);
}

//page.onResourceRequested = function (request) {
//    console.log('(OnRequest)' + JSON.stringify(request, undefined, 4));
//};

page.open(webpageUrl, function (status) {
    if (status !== 'success') {
        console.log('Unable to access network');
        phantom.exit(1);
    } else {
        waitFor(function () {
            return page.evaluate(function () {
                var el = document.getElementById('qunit-testresult');
                if (el && el.innerText.match('completed')) {
                    return true;
                }
                return false;
            });
        }, function () {
            phantom.exit();
        }, 60000);
    }
});