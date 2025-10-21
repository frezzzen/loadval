<script lang="ts">
    import AddTemplate from "./AddTemplate.svelte";
    import TemplateCreationView from "./TemplateCreationView.svelte";
    import TemplateItem from "./TemplateItem.svelte";
    import AgentLoadoutItem from "./AgentLoadoutItem.svelte";
    import AgentLoadoutView from "./AgentLoadoutView.svelte";
    import Loader from "../Loader.svelte";
    import type { Agent } from "../../types/agent.types";
    import type { Template } from "../../types/template.type";
    import { useTemplateManager } from "../../managers/template-manager.svelte";
    import { useValorantManager } from "../../managers/valorant-manager.svelte";
    import type { main } from "../../../wailsjs/go/models";

    type Props = {
        loadout: main.PlayerLoadoutResponse;
    };

    let { loadout }: Props = $props();

    const templateManager = useTemplateManager();
    const valorantManager = useValorantManager();

    let currentView = $state<
        "templates" | "template-creation" | "template-edit" | "agent-loadout"
    >("templates");
    let selectedAgent = $state<Agent | null>(null);
    let selectedTemplate = $state<Template | null>(null);

    $effect(() => {
        if (valorantManager.agents.length === 0) {
            valorantManager.GetAgents();
        }
    });

    const playableAgents = $derived(
        valorantManager.agents.filter((agent) => agent.isPlayableCharacter),
    );

    function onAddTemplateClick() {
        currentView = "template-creation";
    }

    function onAgentClick(agent: Agent) {
        selectedAgent = agent;
        currentView = "agent-loadout";
    }

    function onTemplateEdit(template: Template) {
        selectedTemplate = template;
        currentView = "template-edit";
    }

    function onBack() {
        currentView = "templates";
        selectedAgent = null;
        selectedTemplate = null;
    }
</script>

<div class="templates">
    {#if currentView === "template-creation"}
        <TemplateCreationView {loadout} {onBack} />
    {:else if currentView === "template-edit" && selectedTemplate}
        <TemplateCreationView
            {loadout}
            {onBack}
            editingTemplate={selectedTemplate}
        />
    {:else if currentView === "agent-loadout" && selectedAgent}
        <AgentLoadoutView
            agent={selectedAgent}
            initialLoadout={templateManager.getAgentLoadout(
                selectedAgent.uuid,
            ) || JSON.parse(JSON.stringify(loadout))}
            {onBack}
        />
    {:else}
        <div class="section">
            <div class="section-header">
                <div class="section-title-row">
                    <h2 class="section-title">Agent Loadouts</h2>
                    <div class="switch-container">
                        <label class="switch">
                            <input
                                type="checkbox"
                                checked={templateManager.isAgentLoadoutsEnabled}
                                oninput={() => {
                                    templateManager.setIsAgentLoadoutsEnabled(
                                        !templateManager.isAgentLoadoutsEnabled,
                                    );
                                }}
                            />
                            <span class="slider"></span>
                        </label>
                    </div>
                </div>
                <p class="section-description">
                    Click on an agent to customize their loadout <br />
                    If agent loadouts are enabled, custom loadouts will be overridden
                    when you select an agent.
                </p>
            </div>
            {#if valorantManager.agents.length === 0}{:else}
                <div
                    class="agent-grid"
                    class:disabled={!templateManager.isAgentLoadoutsEnabled}
                >
                    {#each playableAgents as agent}
                        <AgentLoadoutItem
                            {agent}
                            hasLoadout={templateManager.hasAgentLoadout(
                                agent.uuid,
                            )}
                            {onAgentClick}
                            disabled={!templateManager.isAgentLoadoutsEnabled}
                        />
                    {/each}
                </div>
            {/if}
        </div>

        <div class="section">
            <div class="section-header">
                <h2 class="section-title">General Templates</h2>
                <p class="section-description">
                    Create custom templates for quick loadout switching
                </p>
            </div>
            <AddTemplate onAddTemplate={onAddTemplateClick} />
            {#if templateManager.templates.length > 0}
                <div class="template-list">
                    {#each templateManager.templates as template}
                        <TemplateItem {template} onEdit={onTemplateEdit} />
                    {/each}
                </div>
            {:else}
                <div class="empty-state">
                    <div class="empty-icon">
                        <svg
                            width="64"
                            height="64"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="1.5"
                        >
                            <path
                                d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"
                            ></path>
                            <polyline points="17 21 17 13 7 13 7 21"></polyline>
                            <polyline points="7 3 7 8 15 8"></polyline>
                        </svg>
                    </div>
                    <h3 class="empty-title">No Templates Yet</h3>
                    <p class="empty-description">
                        Create your first template to save and manage your
                        favorite weapon loadouts
                    </p>
                </div>
            {/if}
        </div>
    {/if}
</div>

<style lang="scss">
    .templates {
        display: flex;
        flex-direction: column;
        gap: 2rem;
        padding: 1rem;

        .template-list {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
            gap: 1rem;
            margin-top: 0.5rem;

            @media (max-width: 768px) {
                grid-template-columns: 1fr;
            }
        }

        .section {
            display: flex;
            flex-direction: column;
            gap: 1.5rem;

            .section-header {
                display: flex;
                flex-direction: column;
                gap: 0.5rem;

                .section-title-row {
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    gap: 1rem;

                    .section-title {
                        margin: 0;
                        font-size: 1.5rem;
                        font-weight: 700;
                        color: #fff;
                        text-transform: uppercase;
                        letter-spacing: 0.05em;
                    }

                    .switch-container {
                        display: flex;
                        align-items: center;
                        gap: 0.75rem;

                        .switch {
                            position: relative;
                            display: inline-block;
                            width: 52px;
                            height: 28px;

                            input {
                                opacity: 0;
                                width: 0;
                                height: 0;

                                &:checked + .slider {
                                    background: linear-gradient(
                                        135deg,
                                        #ad40ff,
                                        #7a28cb
                                    );
                                    border-color: #ad40ff;

                                    &:before {
                                        transform: translateX(24px);
                                        box-shadow: 0 2px 8px
                                            rgba(173, 64, 255, 0.4);
                                    }
                                }

                                &:focus + .slider {
                                    box-shadow: 0 0 0 3px
                                        rgba(173, 64, 255, 0.2);
                                }
                            }

                            .slider {
                                position: absolute;
                                cursor: pointer;
                                top: 0;
                                left: 0;
                                right: 0;
                                bottom: 0;
                                background: linear-gradient(
                                    135deg,
                                    rgba(255, 255, 255, 0.15),
                                    rgba(255, 255, 255, 0.05)
                                );
                                transition: all 0.3s
                                    cubic-bezier(0.4, 0, 0.2, 1);
                                border-radius: 28px;
                                border: 1px solid rgba(255, 255, 255, 0.2);
                                backdrop-filter: blur(10px);

                                &:before {
                                    position: absolute;
                                    content: "";
                                    height: 22px;
                                    width: 22px;
                                    left: 3px;
                                    bottom: 2px;
                                    background: linear-gradient(
                                        135deg,
                                        #ffffff,
                                        #f8f9fa
                                    );
                                    transition: all 0.3s
                                        cubic-bezier(0.4, 0, 0.2, 1);
                                    border-radius: 50%;
                                    box-shadow:
                                        0 2px 6px rgba(0, 0, 0, 0.15),
                                        0 1px 2px rgba(0, 0, 0, 0.1);
                                }

                                &:hover {
                                    border-color: rgba(255, 255, 255, 0.3);
                                    background: linear-gradient(
                                        135deg,
                                        rgba(255, 255, 255, 0.2),
                                        rgba(255, 255, 255, 0.1)
                                    );
                                }
                            }
                        }

                        .switch-label {
                            font-size: 0.9rem;
                            font-weight: 500;
                            color: rgba(255, 255, 255, 0.8);
                            min-width: 60px;
                        }
                    }
                }

                .section-title {
                    margin: 0;
                    font-size: 1.5rem;
                    font-weight: 700;
                    color: #fff;
                    text-transform: uppercase;
                    letter-spacing: 0.05em;
                }

                .section-description {
                    margin: 0;
                    font-size: 0.95rem;
                    color: rgba(255, 255, 255, 0.6);
                    font-weight: 400;
                }
            }
        }

        .agent-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
            gap: 1rem;
            transition: all 0.3s ease;

            @media (max-width: 768px) {
                grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
            }

            &.disabled {
                opacity: 0.4;
                pointer-events: none;
                position: relative;
            }
        }

        .loader-container {
            display: flex;
            justify-content: center;
            padding: 4rem 2rem;
        }

        .empty-state {
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            padding: 4rem 2rem;
            background: linear-gradient(
                135deg,
                rgba(173, 64, 255, 0.03),
                rgba(122, 40, 203, 0.01)
            );
            border: 2px dashed rgba(173, 64, 255, 0.2);
            border-radius: 12px;
            text-align: center;

            .empty-icon {
                margin-bottom: 1.5rem;
                color: rgba(255, 255, 255, 0.3);

                svg {
                    display: block;
                }
            }

            .empty-title {
                margin: 0 0 0.75rem 0;
                font-size: 1.5rem;
                font-weight: 600;
                color: rgba(255, 255, 255, 0.8);
            }

            .empty-description {
                margin: 0;
                font-size: 1rem;
                color: rgba(255, 255, 255, 0.5);
                max-width: 400px;
            }
        }
    }
</style>
