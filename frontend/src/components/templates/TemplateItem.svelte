<script lang="ts">
    import { useTemplateManager } from "../../managers/template-manager.svelte";
    import type { Template } from "../../types/template.type";

    type Props = {
        template: Template;
        onEdit: (template: Template) => void;
    };

    let { template, onEdit }: Props = $props();

    const templateManager = useTemplateManager();

    function handleDelete(e: Event) {
        e.stopPropagation();
        if (confirm(`Are you sure you want to delete "${template.name}"?`)) {
            templateManager.removeTemplate(template.id);
        }
    }

    function handleClick() {
        templateManager.setActiveTemplate(template.id);
        templateManager.saveLoadout(template.loadout);
    }

    function handleEdit(e: Event) {
        e.stopPropagation();
        templateManager.setActiveTemplate(null);
        onEdit(template);
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Enter" || e.key === " ") {
            e.preventDefault();
            onEdit(template);
        }
    }
</script>

<div
    class="template-item"
    class:active={templateManager.activeTemplate === template.id}
    onclick={handleClick}
    onkeydown={handleKeydown}
    role="button"
    tabindex="0"
>
    <div class="template-content">
        <div class="template-header">
            <h3 class="template-name">{template.name}</h3>
            {#if template.agent}
                <span class="agent-badge">{template.agent}</span>
            {/if}
        </div>

        <div class="action-buttons">
            <button
                class="btn-edit"
                title="Edit template"
                onclick={handleEdit}
                aria-label="Edit template"
            >
                <svg
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                >
                    <path
                        d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"
                    ></path>
                    <path
                        d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"
                    ></path>
                </svg>
            </button>
            <button
                class="btn-delete"
                title="Delete template"
                onclick={handleDelete}
                aria-label="Delete template"
            >
                <svg
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                >
                    <path
                        d="M3 6h18M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
                    ></path>
                </svg>
            </button>
        </div>
    </div>
</div>

<style lang="scss">
    .template-item {
        &.active {
            background: linear-gradient(
                135deg,
                rgba(173, 64, 255, 0.15),
                rgba(122, 40, 203, 0.1)
            );
            border-color: rgba(173, 64, 255, 0.6);
            box-shadow: 0 8px 20px rgba(173, 64, 255, 0.3);
            .action-buttons {
                opacity: 1;
            }
        }

        background: linear-gradient(
            135deg,
            rgba(173, 64, 255, 0.08),
            rgba(122, 40, 203, 0.04)
        );
        border: 1px solid rgba(173, 64, 255, 0.2);
        border-radius: 12px;
        padding: 1.25rem;
        transition: all 0.3s ease;
        cursor: pointer;
        outline: none;

        &:hover,
        &:focus {
            transform: translateY(-2px);
            box-shadow: 0 8px 20px rgba(173, 64, 255, 0.3);
            border-color: rgba(173, 64, 255, 0.6);

            .action-buttons {
                opacity: 1;
            }
        }

        &:focus {
            box-shadow: 0 0 0 3px rgba(173, 64, 255, 0.4);
        }

        .template-content {
            display: flex;
            justify-content: space-between;
            align-items: center;
            gap: 1rem;

            .template-header {
                display: flex;
                align-items: center;
                gap: 1rem;
                flex: 1;
                min-width: 0;

                .template-name {
                    margin: 0;
                    font-size: 1.4rem;
                    font-weight: 600;
                    color: #fff;
                    white-space: nowrap;
                    overflow: hidden;
                    text-overflow: ellipsis;
                }

                .agent-badge {
                    display: inline-flex;
                    align-items: center;
                    padding: 0.25rem 0.75rem;
                    background: rgba(255, 255, 255, 0.1);
                    border: 1px solid rgba(255, 255, 255, 0.2);
                    border-radius: 6px;
                    font-size: 0.75rem;
                    font-weight: 600;
                    color: rgba(255, 255, 255, 0.8);
                    text-transform: uppercase;
                    letter-spacing: 0.05em;
                    white-space: nowrap;
                }
            }

            &:hover {
                .action-buttons {
                    opacity: 1;
                }
            }

            .action-buttons {
                display: flex;
                gap: 0.5rem;
                opacity: 0;
                transition: opacity 0.3s ease;
            }

            .btn-edit {
                display: flex;
                align-items: center;
                justify-content: center;
                width: 40px;
                height: 40px;
                padding: 0;
                background: rgba(64, 173, 255, 0.15);
                border: 1px solid rgba(64, 173, 255, 0.3);
                border-radius: 8px;
                cursor: pointer;
                transition: all 0.3s ease;
                color: #40adff;

                &:hover {
                    transform: scale(1.1);
                    background: rgba(64, 173, 255, 0.25);
                    border-color: rgba(64, 173, 255, 0.5);
                }

                svg {
                    flex-shrink: 0;
                }
            }

            .btn-delete {
                display: flex;
                align-items: center;
                justify-content: center;
                width: 40px;
                height: 40px;
                padding: 0;
                background: rgba(255, 70, 85, 0.15);
                border: 1px solid rgba(255, 70, 85, 0.3);
                border-radius: 8px;
                cursor: pointer;
                transition: all 0.3s ease;
                color: #ff4655;

                &:hover {
                    transform: scale(1.1);
                    background: rgba(255, 70, 85, 0.25);
                    border-color: rgba(255, 70, 85, 0.5);
                }

                svg {
                    flex-shrink: 0;
                }
            }
        }
    }
</style>
