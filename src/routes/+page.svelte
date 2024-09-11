<script>
	import { onMount, tick } from 'svelte'; // Import lifecycle hooks and utilities from Svelte
	import Chart from 'chart.js/auto'; // Import Chart.js for chart rendering
	import { Node, Svelvet } from 'svelvet'; // Import Svelvet components for the resizable and draggable nodes

	// References to the canvas elements for different chart types
	let pieChartCanvas;
	let lineChartCanvas;
	let barChartCanvas;
	let radarChartCanvas;

	// Chart instances for different chart types
	let pieChart = null;
	let lineChart = null;
	let barChart = null;
	let radarChart = null;

	// Default dimensions for each chart node
	let node1Width = 500;
	let node1Height = 400;
	let node2Width = 500;
	let node2Height = 400;
	let node3Width = 500;
	let node3Height = 400;
	let node4Width = 500;
	let node4Height = 400;
	let nodeTreeWidth = 500;
	let nodeTreeHeight = 500;

	// Window size to make the layout responsive
	let windowWidth = 1000;
	let windowHeight = 1000;

	// Handles resizing of nodes, updates the width and height of a node based on its ID
	function handleResize({ width, height }, nodeId) {
		if (nodeId === 1) {
			node1Width = width;
			node1Height = height;
		} else if (nodeId === 2) {
			node2Width = width;
			node2Height = height;
		} else if (nodeId === 3) {
			node3Width = width;
			node3Height = height;
		} else if (nodeId === 4) {
			node4Width = width;
			node4Height = height;
		} else if (nodeId === 'tree') {
			nodeTreeWidth = width;
			nodeTreeHeight = height;
		}
	}

	// Mock data representing different hosts for a treeview
	let hosts = [
		{
			name: 'vigilant-vino-iamr-02',
			id: 'host-1',
			expanded: false,
			status: 'Unknown',
		},
		{
			name: 'made-up-host-01',
			id: 'host-2',
			expanded: false,
			status: 'Unknown',
		},
		{
			name: 'made-up-host-02',
			id: 'host-3',
			expanded: false,
			status: 'Unknown',
		},
	];

	// Stores the currently selected host for the dashboard view
	let selectedHost = null;

	// Toggles the expanded state of the host in the treeview
	function toggleHost(host) {
		host.expanded = !host.expanded;
	}

	// Loads the dashboard for a selected host and fetches its status
	function loadHostDashboard(host) {
		selectedHost = host;
		if (host.name === 'vigilant-vino-iamr-02') {
			fetch('http://localhost:8080/api/graphite')
				.then((res) => res.json())
				.then((data) => {
					host.status = 'UP';
				})
				.catch(() => {
					host.status = 'DOWN';
				});
		} else {
			host.status = 'UP';
		}
	}

	// Fetches graphite data from an API endpoint
	async function fetchGraphiteData() {
		try {
			const response = await fetch('http://localhost:8080/api/graphite');
			if (!response.ok) throw new Error(`API request failed with status ${response.status}`);
			const data = await response.json();
			return data;
		} catch (error) {
			console.error('Error fetching graphite data:', error);
			return null;
		}
	}

	// Updates the line chart with new data
	function updateLineChart(data) {
		if (lineChart) {
			lineChart.data.labels = data.map(item => new Date(item[0] * 1000).toLocaleString());
			lineChart.data.datasets[0].data = data.map(item => item[1]);
			lineChart.update();
		}
	}

	// Updates the bar chart with new data
	function updateBarChart(data) {
		if (barChart) {
			barChart.data.labels = data.map(item => new Date(item[0] * 1000).toLocaleString());
			barChart.data.datasets[0].data = data.map(item => item[1]);
			barChart.update();
		}
	}

	// Updates the radar chart with new data
	function updateRadarChart(data) {
		if (radarChart) {
			const transformedData = data.map(item => item[1]);
			radarChart.data.datasets[0].data = transformedData;
			radarChart.update();
		}
	}

	// Lifecycle hook that runs after the component is mounted
	onMount(async () => {
		// Set the initial window size
		windowWidth = window.innerWidth;
		windowHeight = window.innerHeight;

		// Update window size on resize events
		window.addEventListener('resize', () => {
			windowWidth = window.innerWidth;
			windowHeight = window.innerHeight;
		});

		await tick(); // Wait for DOM updates

		// Fetch graphite data from the API
		const graphiteData = await fetchGraphiteData();

		// Initialize the pie chart if the canvas is available
		if (pieChartCanvas instanceof HTMLCanvasElement) {
			pieChart = new Chart(pieChartCanvas, {
				type: 'pie',
				data: {
					labels: ['Red', 'Blue', 'Yellow'],
					datasets: [{
						data: [300, 50, 100],
						backgroundColor: ['#FF6384', '#36A2EB', '#FFCE56']
					}]
				},
				options: {
					responsive: true,
					maintainAspectRatio: false
				}
			});
		}

		// Initialize the line chart and update it with data
		if (lineChartCanvas instanceof HTMLCanvasElement) {
			lineChart = new Chart(lineChartCanvas, {
				type: 'line',
				data: {
					labels: [],
					datasets: [{
						label: 'Host Alive',
						data: [],
						borderColor: '#4BC0C0',
						tension: 0.1
					}]
				},
				options: {
					responsive: true,
					maintainAspectRatio: false
				}
			});
			if (graphiteData && Array.isArray(graphiteData[0]?.datapoints)) {
				updateLineChart(graphiteData[0].datapoints);
			}
		}

		// Initialize the bar chart and update it with data
		if (barChartCanvas instanceof HTMLCanvasElement) {
			barChart = new Chart(barChartCanvas, {
				type: 'bar',
				data: {
					labels: [],
					datasets: [{
						label: 'Bar Chart Data',
						data: [],
						backgroundColor: '#FF6384'
					}]
				},
				options: {
					responsive: true,
					maintainAspectRatio: false
				}
			});
			if (graphiteData && Array.isArray(graphiteData[0]?.datapoints)) {
				updateBarChart(graphiteData[0].datapoints);
			}
		}

		// Initialize the radar chart and update it with data
		if (radarChartCanvas instanceof HTMLCanvasElement) {
			radarChart = new Chart(radarChartCanvas, {
				type: 'radar',
				data: {
					labels: ['Metric 1', 'Metric 2', 'Metric 3', 'Metric 4', 'Metric 5'],
					datasets: [{
						label: 'Radar Chart Data',
						data: [],
						backgroundColor: 'rgba(34, 202, 236, 0.2)',
						borderColor: 'rgba(34, 202, 236, 1)'
					}]
				},
				options: {
					responsive: true,
					maintainAspectRatio: false
				}
			});
			if (graphiteData && Array.isArray(graphiteData[0]?.datapoints)) {
				updateRadarChart(graphiteData[0].datapoints);
			}
		}
	});
</script>

<!-- HTML layout and structure of the dashboard -->
<h1>Cool Charts</h1>

<div id="app">
<Svelvet {windowWidth} {windowHeight} initialLocation={{ x: 250, y: 250 }} background>
	<!-- Pie Chart Node -->
	<Node 
		id="1" 
		position={{ x: 50, y: 10 }} 
		width={node1Width} 
		height={node1Height} 
		bgColor="white" 
		aria-label="Pie Chart Node"
		role="region"
		resizable 
		on:resizeEnd={(e) => handleResize(e.detail, 1)}
	>
		<canvas bind:this={pieChartCanvas} width={node1Width} height={node1Height}></canvas>
	</Node>

	<!-- Line Chart Node -->
	<Node 
		id="2" 
		position={{ x: 50, y: 460 }} 
		width={node2Width} 
		height={node2Height} 
		bgColor="white" 
		aria-label="Line Chart Node"
		role="region"
		resizable 
		on:resizeEnd={(e) => handleResize(e.detail, 2)}
	>
		<canvas bind:this={lineChartCanvas} width={node2Width} height={node2Height}></canvas>
	</Node>

	<!-- Bar Chart Node -->
	<Node 
		id="3" 
		position={{ x: 600, y: 10 }} 
		width={node3Width} 
		height={node3Height} 
		bgColor="white" 
		aria-label="Bar Chart Node"
		role="region"
		resizable 
		on:resizeEnd={(e) => handleResize(e.detail, 3)}
	>
		<canvas bind:this={barChartCanvas} width={node3Width} height={node3Height}></canvas>
	</Node>

	<!-- Radar Chart Node -->
	<Node 
		id="4" 
		position={{ x: 600, y: 460 }} 
		width={node4Width} 
		height={node4Height} 
		bgColor="white" 
		aria-label="Radar Chart Node"
		role="region"
		resizable 
		on:resizeEnd={(e) => handleResize(e.detail, 4)}
	>
		<canvas bind:this={radarChartCanvas} width={node4Width} height={node4Height}></canvas>
	</Node>

	<!-- Tree View Node for displaying hosts -->
	<Node
		id="tree"
		position={{ x: 1200, y: 10 }}
		width={nodeTreeWidth}
		height={nodeTreeHeight}
		bgColor="lightgray"
		aria-label="Treeview Node"
		role="region"
		resizable
		on:resizeEnd={(e) => handleResize(e.detail, 'tree')}
	>
		<div>
		  <h2>Hosts</h2>
		  <ul>
			{#each hosts as host (host.id)}
			  <li>
				<button on:click|stopPropagation={() => toggleHost(host)} aria-expanded={host.expanded} aria-controls={host.id} style="background: none; border: none; text-align: left; padding: 0;">
				  <strong>{host.name}</strong> (Status: {host.status})
				  {host.expanded ? '[-]' : '[+]'}
				</button>

				{#if host.expanded}
				  <div style="margin-left: 20px;" id={host.id}>
					<p>Host ID: {host.id}</p>
					<button on:click={() => loadHostDashboard(host)}>
					  Load Dashboard
					</button>
					{#if selectedHost && selectedHost.id === host.id}
					  <div style="margin-top: 10px;">
						<h3>Dashboard for {selectedHost.name}</h3>
						<p>Status: {selectedHost.status}</p>
					  </div>
					{/if}
				  </div>
				{/if}
			  </li>
			{/each}
		  </ul>
		</div>
	</Node>
</Svelvet>
</div>

<!-- Basic styling for the layout -->
<style>
	h1 {
	  margin: 0;
	  padding: 10px;
	}
  
	#app {
	  width: 100vw;
	  height: 100vh;
	}

	ul {
		list-style-type: none;
		padding-left: 0;
	}

	li {
		margin-bottom: 10px;
	}

	button {
		margin-top: 5px;
	}
</style>
