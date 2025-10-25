<script lang="ts">
	import { onMount, onDestroy } from "svelte";
	import {
		QuitApp,
		MaximiseApp,
		MinimiseApp,
		ReloadApp,
		LogInfo,
		LogError,
	} from "../wailsjs/go/app/App";
	import { GetMainData } from "../wailsjs/go/main/ValorantAPI";
	import type { main } from "../wailsjs/go/models";
	import AppContainer from "./components/AppContainer.svelte";
	import LoaderContainer from "./components/LoaderContainer.svelte";
	import UpdateNotification from "./components/UpdateNotification.svelte";
	import ConfirmModal from "./components/core/ConfirmModal.svelte";
	import { useModalManager } from "./managers/modal-manager.svelte";

	const modalManager = useModalManager();

	function QuitButton() {
		QuitApp();
	}

	function MaximiseButton() {
		MaximiseApp();
	}

	function MinimiseButton() {
		MinimiseApp();
	}

	function ReloadButton() {
		LogInfo("Reloading application...");
		ReloadApp();
	}

	let loadout = $state<main.PlayerLoadoutResponse | null>(null);
	let ownedSkins = $state<main.OwnedItemsResponseEntitlement[] | null>(null);
	let ownedSkinVariants = $state<main.OwnedItemsResponseEntitlement[] | null>(
		null,
	);
	let ownedAgents = $state<main.OwnedItemsResponseEntitlement[] | null>(null);
	let ownedCards = $state<main.OwnedItemsResponseEntitlement[] | null>(null);
	let windowSize = $state({
		width: window.innerWidth,
		height: window.innerHeight,
	});

	function handleResize() {
		windowSize = { width: window.innerWidth, height: window.innerHeight };
	}

	onMount(() => {
		LogInfo("Frontend application mounted");

		window.addEventListener("resize", handleResize);

		(async () => {
			try {
				LogInfo("Fetching main data from Valorant API...");
				const mainData = await GetMainData();
				loadout = mainData.PlayerLoadout || null;
				ownedSkins = mainData.OwnedSkins[0].Entitlements || [];
				ownedSkinVariants =
					mainData.OwnedSkinVariants[0].Entitlements || [];
				ownedAgents = mainData.OwnedAgents[0].Entitlements || [];
				ownedCards = mainData.OwnedCards[0].Entitlements || [];
				LogInfo("Main data loaded successfully");
			} catch (error) {
				LogError("Failed to load main data: " + error);
				console.error("Failed to load main data:", error);
			}
		})();
	});

	onDestroy(() => {
		window.removeEventListener("resize", handleResize);
	});
</script>

<main>
	<UpdateNotification />

	<drag role="button" tabindex="0" style="--wails-draggable:drag">
		<div class="app-title">LOADVAL</div>
		<div class="drag-spacer"></div>
		<button onclick={ReloadButton} title="Reload App">
			<span class="material-icons">refresh</span>
		</button>
		<button onclick={MaximiseButton} title="Maximize">
			<span class="material-icons">fullscreen</span>
		</button>
		<button onclick={MinimiseButton} title="Minimize">
			<span class="material-icons">minimize</span>
		</button>
		<button onclick={QuitButton} title="Close">
			<span class="material-icons"> close </span>
		</button>
	</drag>

	{#if loadout && ownedSkins && ownedSkinVariants && ownedAgents && ownedCards}
		<AppContainer
			{loadout}
			{ownedSkins}
			{ownedSkinVariants}
			{ownedAgents}
			{ownedCards}
		/>
	{:else}
		<div class="loader-container">
			<LoaderContainer />
		</div>
	{/if}
</main>

<ConfirmModal
	bind:isOpen={modalManager.isOpen}
	title={modalManager.config.title}
	message={modalManager.config.message}
	confirmText={modalManager.config.confirmText}
	cancelText={modalManager.config.cancelText}
	type={modalManager.config.type}
	onConfirm={() => modalManager.handleConfirm()}
	onCancel={() => modalManager.handleCancel()}
/>

<style lang="scss">
	main {
		margin-top: 4rem;
		position: relative;
		transition: transform 0.1s ease-out;

		&.dragging {
			cursor: grabbing;
			user-select: none;
		}
	}

	.title {
		font-size: 1.6rem;
		color: white;
		font-weight: 600;
		text-align: center;
		margin-bottom: 1rem;
		text-align: center;
	}

	drag {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		width: 100%;
		z-index: 1000;
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 1rem;
		background: linear-gradient(
			135deg,
			rgba(255, 255, 255, 0.1),
			rgba(255, 255, 255, 0.05)
		);
		backdrop-filter: blur(10px);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-top: none;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
		cursor: move;
		transition: all 0.3s ease;

		&:hover {
			background: linear-gradient(
				135deg,
				rgba(255, 255, 255, 0.15),
				rgba(255, 255, 255, 0.08)
			);
			box-shadow: 0 6px 25px rgba(0, 0, 0, 0.4);
		}

		.app-title {
			font-size: 1.5rem;
			font-weight: 700;
			color: #fff;
			text-transform: uppercase;
			letter-spacing: 0.1em;
			text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
			background: linear-gradient(135deg, #ad40ff, #7a28cb);
			-webkit-background-clip: text;
			-webkit-text-fill-color: transparent;
			background-clip: text;
			pointer-events: none;
		}

		.drag-spacer {
			flex: 1;
		}

		button {
			display: flex;
			align-items: center;
			justify-content: center;
			width: 3.2rem;
			height: 3.2rem;
			border: none;
			border-radius: 0.8rem;
			background: rgba(255, 255, 255, 0.1);
			color: #cccaca;
			cursor: pointer;
			transition: all 0.2s ease;
			backdrop-filter: blur(5px);
			border: 1px solid rgba(255, 255, 255, 0.1);
			pointer-events: auto;

			&:hover {
				background: rgba(255, 255, 255, 0.2);
				transform: translateY(-1px);
				box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
			}

			&:active {
				transform: translateY(0);
				box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
			}

			&:nth-child(3) {
				&:hover {
					background: rgba(139, 92, 246, 0.2);
					color: #8b5cf6;
				}
			}

			&:nth-child(4) {
				&:hover {
					background: rgba(34, 197, 94, 0.2);
					color: #22c55e;
				}
			}

			&:nth-child(5) {
				&:hover {
					background: rgba(251, 191, 36, 0.2);
					color: #fbbf24;
				}
			}

			&:nth-child(6) {
				&:hover {
					background: rgba(239, 68, 68, 0.2);
					color: #ef4444;
				}
			}

			.material-icons {
				font-size: 1.8rem;
				transition: transform 0.2s ease;
			}

			&:hover .material-icons {
				transform: scale(1.1);
			}
		}
	}

	@media (max-width: 1200px) {
		drag {
			padding: 0.9rem;

			.app-title {
				font-size: 1.3rem;
			}

			button {
				width: 3rem;
				height: 3rem;

				.material-icons {
					font-size: 1.7rem;
				}
			}
		}
	}

	@media (max-width: 1000px) {
		drag {
			padding: 0.8rem;

			.app-title {
				font-size: 1.2rem;
			}

			button {
				width: 2.8rem;
				height: 2.8rem;

				.material-icons {
					font-size: 1.6rem;
				}
			}
		}
	}

	@media (max-width: 800px) {
		drag {
			padding: 0.7rem;
			gap: 0.4rem;

			.app-title {
				font-size: 1.1rem;
			}

			button {
				width: 2.6rem;
				height: 2.6rem;

				.material-icons {
					font-size: 1.5rem;
				}
			}
		}

		main {
			margin-top: 3.5rem;
		}
	}
</style>
