<script lang="ts">
	import { onMount } from 'svelte';
	import { type PageProps } from './$types';
	import { toast } from '@zerodevx/svelte-toast';

	let { params }: PageProps = $props();

	let notifications = $state<string[]>([]);
	let eventHistory = $state<any[]>([]);
	let socket: WebSocket | null = null;
	let message = $state('User logged in');

	async function sendMessage(event: Event) {
		event.preventDefault();
		const payload = {
			subject: `tenant.${params.tenant_id}`,
			data: {
				id: crypto.randomUUID(),
				tenant_id: params.tenant_id,
				message: message,
				timestamp: new Date().toISOString()
			}
		};

		try {
			const response = await fetch('http://localhost:4001/events', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(payload)
			});

			if (!response.ok) {
				console.error('Failed to send message:', response.statusText);
			}
			toast.push(`Message sent: ${message}`, {
				theme: {
					'--toastBackground': '#4caf50',
					'--toastColor': '#fff'
				},
				duration: 3000
			});
		} catch (error) {
			console.error('Error sending message:', error);
		}
	}

	async function fetchEventHistory() {
		try {
			const response = await fetch(`http://localhost:4001/events/all?tenant_id=${params.tenant_id}`);
			if (response.ok) {
				const data = await response.json();
				eventHistory = data || [];
			} else {
				console.error('Failed to fetch event history:', response.statusText);
			}
		} catch (error) {
			console.error('Error fetching event history:', error);
		}
	}

	onMount(() => {
		fetchEventHistory();
		const subject = `tenant.${params.tenant_id}`;
		const wsUrl = `ws://localhost:4001/subscribe?subject=${subject}`;
		socket = new WebSocket(wsUrl);

		socket.onopen = () => {
			console.log('WebSocket connection established');
		};

		socket.onmessage = (event) => {
			try {
				const data = JSON.parse(event.data);
				notifications = [...notifications, JSON.stringify(data, null, 2)];
				toast.push(`Notification received`, {
					theme: {
						'--toastBackground': '#2196f3',
						'--toastColor': '#fff'
					},
					duration: 3000
				});
			} catch (e) {
				console.error('Error parsing message:', e);
				notifications = [...notifications, event.data];
			}
		};

		socket.onerror = (error) => {
			console.error('WebSocket error:', error);
		};

		socket.onclose = () => {
			console.log('WebSocket connection closed');
		};

		return () => {
			if (socket) {
				socket.close();
			}
		};
	});
</script>

<div class="container mx-auto p-4">
	<h1 class="text-3xl font-bold mb-4">Tenant: {params.tenant_id}</h1>

	<div class="card bg-base-100 shadow-xl">
		<div class="card-body">
			<h2 class="card-title">Send Notification</h2>
			<form onsubmit={sendMessage} class="form-control">
				<textarea bind:value={message} class="textarea textarea-bordered" rows="3"></textarea>
				<div class="card-actions justify-end mt-4">
					<button type="submit" class="btn btn-primary">Send</button>
				</div>
			</form>
		</div>
	</div>

	{#if notifications.length > 0}
		<div class="mt-8">
			<h2 class="text-2xl font-bold mb-4">Live Notifications</h2>
			<div class="mockup-code">
				{#each notifications as notification, i (i)}
					<pre data-prefix={i + 1}><code>{notification}</code></pre>
				{/each}
			</div>
		</div>
	{/if}

	{#if eventHistory.length > 0}
		<div class="mt-8">
			<h2 class="text-2xl font-bold mb-4">Event History</h2>
			<div class="mockup-code">
				{#each eventHistory as event, i (i)}
					<pre data-prefix={i + 1}><code>{JSON.stringify(event.data, null, 2)}</code></pre>
				{/each}
			</div>
		</div>
	{/if}
</div>
