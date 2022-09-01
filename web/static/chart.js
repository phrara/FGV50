

var option = {
    "animation": true,
    "animationThreshold": 2000,
    "animationDuration": 1000,
    "animationEasing": "cubicOut",
    "animationDelay": 0,
    "animationDurationUpdate": 300,
    "animationEasingUpdate": "cubicOut",
    "animationDelayUpdate": 0,
    "color": [
        "#c23531",
        "#2f4554",
        "#61a0a8",
        "#d48265",
        "#749f83",
        "#ca8622",
        "#bda29a",
        "#6e7074",
        "#546570",
        "#c4ccd3",
        "#f05b72",
        "#ef5b9c",
        "#f47920",
        "#905a3d",
        "#fab27b",
        "#2a5caa",
        "#444693",
        "#726930",
        "#b2d235",
        "#6d8346",
        "#ac6767",
        "#1d953f",
        "#6950a1",
        "#918597"
    ],
    "series": [
        {
            "type": "tree",
            "data": [
                {
                    "name": "",
                    "children": []
                }
            ],
            "symbol": "emptyCircle",
            "emphasis": {
                "focus": "descendant"
            },
            "top": '18%',
            "bottom": "14%",
            "symbolSize": 7,
            "edgeShape": "polyline",
            "edgeForkPosition": "63%",
            "roam": true,
            "expandAndCollapse": true,
            "orient": "LR",
            "label": {
                "show": true,
                "position": "top",
                "margin": 8
            },
            lineStyle: {
                width: 5
            },
            "leaves": {
                "label": {
                    "show": true,
                    "position": "top",
                    "margin": 8
                }
            }
        }
    ],
    "legend": [
        {
            "data": [],
            "selected": {}
        }
    ],
    "tooltip": {
        "show": true,
        "trigger": "item",
        "triggerOn": "mousemove",
        "axisPointer": {
            "type": "line"
        },
        "showContent": true,
        "alwaysShowContent": false,
        "showDelay": 0,
        "hideDelay": 100,
        "textStyle": {
            "fontSize": 14
        },
        "borderWidth": 0,
        "padding": 5
    }
};

function getOption() {
    return option
}

function setRoot(root) {
    option.series[0].data[0].name = root
    return option
}

function addIP(ip) {
    for (let v of option.series[0].data[0].children) {
        if (ip === v.name) {
            return
        }
    }
    var c = {
        "name": ip,
        "children": []
    }
    option.series[0].data[0].children.push(c)
}

function addPortInfo(ip, port, proto, banner) {
    for (let v of option.series[0].data[0].children) {
        if (ip === v.name) {
            var c = {
                "name": "port:" + port + "-" + proto + "-" + banner
            }
            v.children.push(c)
        }
    }
}

function addHDInfo(ip, mac, dev) {
  for (let v of option.series[0].data[0].children) {
    if (ip === v.name) {
        var c = {
            "name": "mac:" + mac + "-" + "dev:" + dev
        }
        v.children.push(c)
    }
  }
}