if (navigator.userAgent.indexOf("PhantomJS") !== -1) {
    String.prototype.format = function () {
        var args = arguments;
        return this.replace(/{(\d+)}/g, function (match, number) {
            return typeof args[number] != 'undefined'
      ? args[number]
      : '{' + number + '}';
        });
    };

    String.prototype.teamCityEscape = function () {
        return this.replace(/['\n\r\|\[\]]/g, function (match) {
            switch (match) {
                case "'":
                    return "|'";
                case "\n":
                    return "|n";
                case "\r":
                    return "|r";
                case "|":
                    return "||";
                case "[":
                    return "|[";
                case "]":
                    return "|]";
                default:
                    return match;
            }
        });
    };
    
    /* TODO (dw): Have this passed through as a param to PhantonJS.exe? */
    var suiteName = "QUnit Tests";
    var currentTestName = "";
    var hasBegun = false;

    qunitBegin = function () {
        // TODO (dw): Should be able to use QUnit.begin() - but that doesn't seem to fire.
        console.log("##teamcity[testSuiteStarted name='{0}']".format(suiteName.teamCityEscape()));
    };

    /* QUnit.testStart({ name }) */
    QUnit.testStart = function (args) {
        if (!hasBegun) {
            qunitBegin();
            hasBegun = true;
        }

        currentTestName = args.name;
    };

    /* QUnit.log({ result, actual, expected, message }) */
    QUnit.log = function (args) {
        var currentAssertion = "{0} > {1}".format(currentTestName, args.message).teamCityEscape();

        console.log("##teamcity[testStarted name='{0}']".format(currentAssertion));

        if (!args.result) {
            console.log("##teamcity[testFailed type='comparisonFailure' name='{0}' details='expected={1}, actual={2}' expected='{1}' actual='{2}']".format(currentAssertion, args.expected.teamCityEscape(), args.actual.teamCityEscape()));
        }

        console.log("##teamcity[testFinished name='{0}']".format(currentAssertion));
    };

    /* QUnit.done({ failed, passed, total, runtime }) */
    QUnit.done = function (args) {
        console.log("##teamcity[testSuiteFinished name='{0}']".format(suiteName.teamCityEscape()));
    };
}