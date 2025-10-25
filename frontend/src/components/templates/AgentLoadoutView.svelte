<script lang="ts">
    import type { Agent } from "../../types/agent.types";
    import LoadoutView from "../loadout/LoadoutView.svelte";
    import type { main } from "../../../wailsjs/go/models";
    import { useTemplateManager } from "../../managers/template-manager.svelte";
    import { useModalManager } from "../../managers/modal-manager.svelte";

    type Props = {
        agent: Agent;
        initialLoadout: main.PlayerLoadoutResponse;
        onBack: () => void;
    };

    let { agent, initialLoadout, onBack }: Props = $props();

    let templateManager = useTemplateManager();
    let modalManager = useModalManager();

    let loadout = $state<main.PlayerLoadoutResponse>(initialLoadout);

    function handleSave() {
        templateManager.setAgentLoadout(agent.uuid, loadout);
        onBack();
    }

    async function handleClear() {
        const confirmed = await modalManager.confirm({
            title: "Clear Agent Loadout",
            message: `Are you sure you want to clear the loadout for ${agent.displayName}? This action cannot be undone.`,
            confirmText: "Clear",
            cancelText: "Cancel",
            type: "danger",
        });

        if (confirmed) {
            templateManager.removeAgentLoadout(agent.uuid);
            onBack();
        }
    }
</script>

<div class="agent-loadout-view">
    <div class="header">
        <div class="header-top">
            <div class="agent-info">
                <img
                    src={agent.displayIcon}
                    alt={agent.displayName}
                    class="agent-avatar"
                />
                <div class="agent-details">
                    <h2 class="title">{agent.displayName} Loadout</h2>
                    <p class="agent-role">{agent.role.displayName}</p>
                </div>
            </div>
            <div class="actions">
                <button class="btn-secondary" onclick={onBack}>
                    <span>Cancel</span>
                </button>
                {#if templateManager.hasAgentLoadout(agent.uuid)}
                    <button class="btn-danger" onclick={handleClear}>
                        <span>Clear</span>
                    </button>
                {/if}
                <button class="btn-primary" onclick={handleSave}>
                    <span>Save Loadout</span>
                </button>
            </div>
        </div>
    </div>
    <LoadoutView bind:loadout />
</div>

<style lang="scss">
    .agent-loadout-view {
        display: flex;
        flex-direction: column;
        gap: 2rem;

        .header {
            display: flex;
            flex-direction: column;
            gap: 1.5rem;
            padding: 1.5rem;
            background: linear-gradient(
                135deg,
                rgba(173, 64, 255, 0.08),
                rgba(122, 40, 203, 0.04)
            );
            border-radius: 12px;
            border: 1px solid rgba(173, 64, 255, 0.2);
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);

            .header-top {
                display: flex;
                justify-content: space-between;
                align-items: center;
                flex-wrap: wrap;
                gap: 1.5rem;

                .agent-info {
                    display: flex;
                    align-items: center;
                    gap: 1.25rem;

                    .agent-avatar {
                        width: 64px;
                        height: 64px;
                        border-radius: 50%;
                        border: 2px solid rgba(173, 64, 255, 0.4);
                        object-fit: cover;
                        background: linear-gradient(
                            135deg,
                            rgba(173, 64, 255, 0.15),
                            rgba(122, 40, 203, 0.1)
                        );
                    }

                    .agent-details {
                        display: flex;
                        flex-direction: column;
                        gap: 0.25rem;

                        .title {
                            margin: 0;
                            font-size: 1.75rem;
                            font-weight: 700;
                            color: #fff;
                            text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
                            letter-spacing: -0.025em;
                        }

                        .agent-role {
                            margin: 0;
                            font-size: 0.875rem;
                            color: rgba(255, 255, 255, 0.6);
                            font-weight: 500;
                            text-transform: uppercase;
                            letter-spacing: 0.05em;
                        }
                    }
                }

                .actions {
                    display: flex;
                    gap: 0.75rem;
                    flex-wrap: wrap;
                }
            }
        }

        button {
            padding: 1rem 2rem;
            border-radius: 10px;
            font-weight: 600;
            font-size: 1.05rem;
            cursor: pointer;
            transition: all 0.2s ease;
            border: none;
            outline: none;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
            white-space: nowrap;

            span {
                display: inline-block;
            }

            &:hover {
                transform: translateY(-2px);
                box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
            }

            &:active {
                transform: translateY(0);
                box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
            }

            &.btn-primary {
                background: linear-gradient(
                    135deg,
                    rgba(173, 64, 255, 0.8),
                    rgba(122, 40, 203, 0.8)
                );
                color: white;
                border: 1px solid rgba(173, 64, 255, 0.4);

                &:hover {
                    background: linear-gradient(
                        135deg,
                        rgba(173, 64, 255, 1),
                        rgba(122, 40, 203, 1)
                    );
                    border-color: rgba(173, 64, 255, 0.6);
                    box-shadow: 0 4px 16px rgba(173, 64, 255, 0.4);
                }
            }

            &.btn-secondary {
                background: linear-gradient(
                    135deg,
                    rgba(255, 255, 255, 0.08),
                    rgba(200, 200, 200, 0.05)
                );
                color: rgba(255, 255, 255, 0.9);
                border: 1px solid rgba(255, 255, 255, 0.2);

                &:hover {
                    background: linear-gradient(
                        135deg,
                        rgba(255, 255, 255, 0.15),
                        rgba(200, 200, 200, 0.1)
                    );
                    border-color: rgba(255, 255, 255, 0.4);
                    color: #ffffff;
                }
            }

            &.btn-danger {
                background: linear-gradient(
                    135deg,
                    rgba(255, 70, 85, 0.2),
                    rgba(230, 57, 70, 0.15)
                );
                color: #ff6b7a;
                border: 1px solid rgba(255, 70, 85, 0.4);

                &:hover {
                    background: linear-gradient(
                        135deg,
                        rgba(255, 70, 85, 0.35),
                        rgba(230, 57, 70, 0.3)
                    );
                    border-color: rgba(255, 70, 85, 0.6);
                    color: #ff8590;
                    box-shadow: 0 4px 16px rgba(255, 70, 85, 0.3);
                }
            }
        }
    }
</style>
