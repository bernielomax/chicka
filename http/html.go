package http

var indexHTML = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width,initial-scale=1">
    <title>Chicka Results</title>
    <style type="text/css">
        body {
            font-family: sans-serif;
        }
        table {
            table-layout: fixed;
            width: 100%;
        }

        table th {
            text-align: left;
        }

        table thead tr th {
            border-bottom: 2px solid #999;
        }
        table th,
        table td {
            border-right: 1px solid #DDD;
            padding: 5px;
        }
        table tbody tr:nth-child(even) td {
            background: #F1F1F1;
        }

        table .short-col {
            width: 110px;
        }

        table .date-col {
            width: 400px;
        }

        .pass {
            color: green;
        }
        .fail {
            color: red;
        }
    </style>
</head>
<body>

<div id="root"></div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.3.1/react.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.3.1/react-dom.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/axios/0.16.2/axios.js"></script>
<script>

"use strict";

var _extends = Object.assign || function (target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i]; for (var key in source) { if (Object.prototype.hasOwnProperty.call(source, key)) { target[key] = source[key]; } } } return target; };

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _possibleConstructorReturn(self, call) { if (!self) { throw new ReferenceError("this hasn't been initialised - super() hasn't been called"); } return call && (typeof call === "object" || typeof call === "function") ? call : self; }

function _inherits(subClass, superClass) { if (typeof superClass !== "function" && superClass !== null) { throw new TypeError("Super expression must either be null or a function, not " + typeof superClass); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, enumerable: false, writable: true, configurable: true } }); if (superClass) Object.setPrototypeOf ? Object.setPrototypeOf(subClass, superClass) : subClass.__proto__ = superClass; }

var App = function (_React$Component) {
	_inherits(App, _React$Component);

	function App(props) {
		_classCallCheck(this, App);

		var _this2 = _possibleConstructorReturn(this, _React$Component.call(this, props));

		_this2.getTableData = _this2.getTableData.bind(_this2);

		_this2.state = {
			loadingError: false,
			results: {}
		};
		return _this2;
	}

	App.prototype.componentDidMount = function componentDidMount() {
		var _this3 = this;

		_this3.getTableData();

		setInterval(function () {
			_this3.getTableData();
		}, 5000);
	};

	App.prototype.getTableData = function getTableData() {
		var _this = this;

		console.log('getTableData');

		axios.get("http://127.0.0.1:9090").then(function (response) {
			_this.setState({
				results: response.data
			});
		}).catch(function (error) {
			console.log(error);
			this.setState({
				loadingError: true
			});
		});
	};

	App.prototype.render = function render() {
		var resultsObj = this.state.results;
		var unsortedRows = [];
		var rows = [];
		for (var key in resultsObj) {
			if (resultsObj.hasOwnProperty(key)) {

				unsortedRows.push(_extends({
					"date": key
				}, resultsObj[key]));
			}
		}

		unsortedRows.sort(function (a, b) {
			return parseInt(b.date) - parseInt(a.date);
		});

		unsortedRows.map(function (row) {
			var newDate = new Date(parseInt(row.date) / 1000000);
			rows.push(React.createElement(
				"tr",
				{ key: row.date },
				React.createElement(
					"td",
					{ className: "date-col" },
					newDate.toString()
				),
				React.createElement(
					"td",
					null,
					row.description
				),
				React.createElement(
					"td",
					{ className: "short-col" },
					row.expect.toString()
				),
				React.createElement(
					"td",
					{ className: "short-col" },
					row.result.toString()
				),
				React.createElement(
					"td",
					{ className: "short-col" },
					row.data
				),
				React.createElement(
					"td",
					{ className: "short-col" },
					React.createElement(
						"span",
						{ className: row.expect === row.result ? "pass" : "fail" },
						row.expect === row.result ? "✔" : "✖"
					)
				)
			));
		});

		return React.createElement(
			"div",
			null,
			React.createElement(
				"table",
				null,
				React.createElement(
					"thead",
					null,
					React.createElement(
						"tr",
						null,
						React.createElement(
							"th",
							{ className: "date-col" },
							"Date"
						),
						React.createElement(
							"th",
							null,
							"Description"
						),
						React.createElement(
							"th",
							{ className: "short-col" },
							"Expect"
						),
						React.createElement(
							"th",
							{ className: "short-col" },
							"Result"
						),
						React.createElement(
							"th",
							{ className: "short-col" },
							"Data"
						),
						React.createElement("th", { className: "short-col" })
					)
				),
				React.createElement(
					"tbody",
					null,
					rows
				)
			)
		);
	};

	return App;
}(React.Component);

ReactDOM.render(React.createElement(App, null), document.getElementById("root"));

</script>
</body>


</html>
`
