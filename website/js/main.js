window.addEventListener('load', () => {
    main()
})



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
        setupD3(obj.data)
    }
}

function setupD3(data) {
    var svg = d3.select("svg");
    const margin = {top: 20, right: 20, bottom: 30, left: 50};
    const width = +svg.attr("width") - margin.left - margin.right;
    const height = +svg.attr("height") - margin.top - margin.bottom;
    const g = svg.append("g").attr("transform", "translate(" + margin.left + "," + margin.top + ")");
    svg.style('background', '#333')

    const operations = ["GET", "POST", "PUT", "DELETE"];
    var color = d3.scaleOrdinal()
        .domain(operations)
        .range(["rgba(249, 208, 87, 0.7)", "rgba(54, 174, 175, 0.65)"]);

    var x = d3.scaleTime().range([0, width]),
        y = d3.scaleLinear().range([height, 0]),
        z = color;

    var area = d3.area()
        .curve(d3.curveLinear)
        .x(function(d) {
            const k = new Date(d.timestamp * 1000)
            return x(k);
        })
        .y0(y(0))
        .y1(function(d) { return y(d.request_time); });
    data.forEach(function(d) {
        d.request_time = +d.request_time;
    });
    var sources = operations.map(function(op) {
        return {
            id: op,
            values: data.filter(d => d.operation === op).map(function(d) {
                return {timestamp: d.timestamp, request_time: d.request_time};
            })
        };
    });

    x.domain(d3.extent(data, function(d) {
        const k = new Date(d.timestamp * 1000);
        return k;
    }));
    y.domain([
        0,
        d3.max(sources, function(c) { return d3.max(c.values, function(d) { return d.request_time; }); })
    ]);
    z.domain(sources.map(function(c) { return c.id; }));

    g.append("g")
        .attr("class", "axis axis--x")
        .attr("transform", "translate(0," + height + ")")
        .call(d3.axisBottom(x));

    g.append("g")
        .attr("class", "axis axis--y")
        .call(d3.axisLeft(y))
        .append("text")
        .attr("transform", "rotate(-90)")
        .attr("y", 6)
        .attr("dy", "0.71em")
        .attr("fill", "#000")
        .text("Power, kW");

    var source = g.selectAll(".area")
        .data(sources)
        .enter().append("g")
        .attr("class", function(d) { return `area ${d.id}`; })

    source.append("path")
        .attr("d", function(d) { return area(d.values); })
        .style("fill", function(d) { return z(d.id); });
}