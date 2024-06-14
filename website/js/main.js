window.addEventListener('load', () => {
    main()
})

const OPERATIONS = {
    GET: 'GET',
    POST: 'POST',
    PUT: 'PUT',
    DELETE: 'DELETE',
}

function main() {
    const wsUrl = "ws://127.0.0.1:8081/ws"
    const ws = new WebSocket(wsUrl)
    ws.onerror = err => {
        console.log('an error is occurred while doing that', err)
    }
    ws.onopen = () => {
        ws.send('getData')
    }
    ws.onmessage = (event) => {
        handleWebsocketMessage(event.data)
    }
    
}

function handleWebsocketMessage(data) {
    const obj = JSON.parse(data)
    if (obj.action === 'allLogs') {
        setupGraph(obj.data)
    } else if (obj.action === 'singleLog') {
        globalData.push(obj.data);
        updateChart()
    } else {
        console.warn(`unkown`, obj)
    }
    
}

// since data can be updated, define it in this scope
const globalData = [];

let chartInstance;

function setupGraph(data) {
    const chartDom = document.getElementById('main');
    chartInstance = echarts.init(chartDom, 'dark');

    // Initialize the chart
    const option = {
        title: {
            text: 'Request Time'
        },
        tooltip: {
            trigger: 'axis'
        },
        legend: {
            data: ['GET', 'POST', 'PUT', 'DELETE']
        },
        xAxis: {
            type: 'time',
            boundaryGap: false
        },
        yAxis: {
            type: 'value'
        },
    };
    option && chartInstance.setOption(option);
    
    // push the initial data
    globalData.push(...data);
    
    // update chart
    updateChart()
}

function updateChart() {
    const chartData = {
        [OPERATIONS.GET]: {
            timestamps: [],
            requestTimes: [],
        }, 
        [OPERATIONS.POST]: {
            timestamps: [],
            requestTimes: [],
        },
        [OPERATIONS.PUT]: {
            timestamps: [],
            requestTimes: [],
        },
        [OPERATIONS.DELETE]: {
            timestamps: [],
            requestTimes: [],
        }
    }
    // add data
    globalData.forEach(item => {
        const date = new Date(item.timestamp * 1000);
        chartData[item.operation].timestamps.push(date)
        chartData[item.operation].requestTimes.push(item.request_time);
    });
    
    // generate series from chart data
    const series = Object.keys(chartData).reduce((acc, curr) => {
        const singleData = {
            name: curr,
            data: chartData[curr].timestamps.map(
                (time, index) => [time, chartData[curr].requestTimes[index]]
            ),
            type: 'line',
            areaStyle: {},
            smooth: true
        }
        acc.push(singleData)
        return acc
    }, [])
    
    // set series
    chartInstance.setOption(
        {
            series
        }
    )
}