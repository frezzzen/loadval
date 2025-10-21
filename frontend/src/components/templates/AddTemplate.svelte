<script lang="ts">
    type Props = {
        onAddTemplate: () => void;
    };

    const { onAddTemplate }: Props = $props();

    const handleAddTemplate = () => {
        onAddTemplate();
    };
</script>

<button class="add-template" onclick={handleAddTemplate}>
    <div class="add-template-content">
        <div class="add-icon">
            <svg
                width="48"
                height="48"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
            >
                <line x1="12" y1="5" x2="12" y2="19"></line>
                <line x1="5" y1="12" x2="19" y2="12"></line>
            </svg>
        </div>
        <span class="add-label">Create Template</span>
        <span class="add-description">Save a new loadout template</span>
    </div>
</button>

<style lang="scss">
    .add-template {
        display: flex;
        align-items: center;
        justify-content: center;
        min-height: 180px;
        padding: 2rem;
        background: linear-gradient(
            135deg,
            rgba(173, 64, 255, 0.03),
            rgba(122, 40, 203, 0.01)
        );
        border: 2px dashed rgba(173, 64, 255, 0.3);
        border-radius: 12px;
        cursor: pointer;
        transition: all 0.3s ease;

        &:hover {
            border-color: rgba(173, 64, 255, 0.6);
            background: linear-gradient(
                135deg,
                rgba(173, 64, 255, 0.08),
                rgba(122, 40, 203, 0.04)
            );
            transform: translateY(-2px);
            box-shadow: 0 8px 24px rgba(173, 64, 255, 0.25);

            .add-icon {
                transform: scale(1.1) rotate(90deg);
                background: linear-gradient(
                    135deg,
                    rgba(173, 64, 255, 0.4),
                    rgba(122, 40, 203, 0.3)
                );
            }

            .add-label {
                color: #fff;
            }
        }

        .add-template-content {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 1rem;

            .add-icon {
                display: flex;
                align-items: center;
                justify-content: center;
                width: 72px;
                height: 72px;
                background: linear-gradient(
                    135deg,
                    rgba(173, 64, 255, 0.2),
                    rgba(122, 40, 203, 0.15)
                );
                border-radius: 50%;
                transition: all 0.3s ease;
                color: rgba(173, 64, 255, 0.9);

                svg {
                    display: block;
                }
            }

            .add-label {
                font-size: 1.125rem;
                font-weight: 600;
                color: rgba(255, 255, 255, 0.8);
                text-transform: uppercase;
                letter-spacing: 0.05em;
                transition: color 0.3s ease;
            }

            .add-description {
                font-size: 0.875rem;
                color: rgba(255, 255, 255, 0.5);
                font-weight: 400;
                text-align: center;
            }
        }
    }
</style>
