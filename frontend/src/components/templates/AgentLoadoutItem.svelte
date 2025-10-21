<script lang="ts">
    import type { Agent } from "src/types/agent.types";

    type Props = {
        agent: Agent;
        hasLoadout: boolean;
        onAgentClick: (agent: Agent) => void;
    };

    let { agent, hasLoadout, onAgentClick }: Props = $props();

    function handleClick() {
        onAgentClick(agent);
    }
</script>

<button
    class="agent-loadout-item"
    class:has-loadout={hasLoadout}
    onclick={handleClick}
>
    <div class="agent-image-container">
        <img
            src={agent.displayIcon}
            alt={agent.displayName}
            class="agent-icon"
        />
        {#if hasLoadout}
            <div class="loadout-indicator">
                <svg
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="3"
                >
                    <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
            </div>
        {/if}
    </div>
    <span class="agent-name">{agent.displayName}</span>
</button>

<style lang="scss">
    .agent-loadout-item {
        position: relative;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.75rem;
        padding: 1rem;
        background: linear-gradient(
            135deg,
            rgba(173, 64, 255, 0.03),
            rgba(122, 40, 203, 0.01)
        );
        border: 2px solid rgba(173, 64, 255, 0.2);
        border-radius: 12px;
        cursor: pointer;
        transition: all 0.3s ease;

        &:hover {
            transform: translateY(-4px);
            border-color: rgba(173, 64, 255, 0.6);
            box-shadow: 0 8px 24px rgba(173, 64, 255, 0.3);

            .agent-icon {
                transform: scale(1.05);
            }
        }

        &.has-loadout {
            border-color: rgba(70, 243, 85, 0.4);
            background: linear-gradient(
                135deg,
                rgba(70, 243, 85, 0.08),
                rgba(70, 243, 85, 0.02)
            );

            &:hover {
                border-color: rgba(70, 243, 85, 0.6);
                box-shadow: 0 8px 24px rgba(70, 243, 85, 0.3);
            }
        }

        .agent-image-container {
            position: relative;
            width: 80px;
            height: 80px;
            border-radius: 50%;
            overflow: hidden;
            background: linear-gradient(
                135deg,
                rgba(173, 64, 255, 0.15),
                rgba(122, 40, 203, 0.1)
            );
            border: 2px solid rgba(173, 64, 255, 0.3);
            transition: all 0.3s ease;

            .agent-icon {
                width: 100%;
                height: 100%;
                object-fit: cover;
                transition: transform 0.3s ease;
            }

            .loadout-indicator {
                position: absolute;
                bottom: -2px;
                right: -2px;
                width: 28px;
                height: 28px;
                background: #46f355;
                border: 2px solid rgba(0, 0, 0, 0.8);
                border-radius: 50%;
                display: flex;
                align-items: center;
                justify-content: center;
                color: rgba(0, 0, 0, 0.8);
                animation: pulse 2s ease-in-out infinite;
            }
        }

        .agent-name {
            font-size: 0.875rem;
            font-weight: 600;
            color: rgba(255, 255, 255, 0.9);
            text-align: center;
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }
    }

    @keyframes pulse {
        0%,
        100% {
            box-shadow: 0 0 0 0 rgba(70, 243, 85, 0.7);
        }
        50% {
            box-shadow: 0 0 0 8px rgba(70, 243, 85, 0);
        }
    }
</style>
