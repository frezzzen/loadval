<script lang="ts">
    type Props = {
        isOpen: boolean;
        title: string;
        message: string;
        confirmText?: string;
        cancelText?: string;
        type?: "danger" | "warning" | "info";
        onConfirm: () => void;
        onCancel: () => void;
    };

    let {
        isOpen = $bindable(false),
        title,
        message,
        confirmText = "Confirm",
        cancelText = "Cancel",
        type = "info",
        onConfirm,
        onCancel,
    }: Props = $props();

    function handleConfirm() {
        onConfirm();
        isOpen = false;
    }

    function handleCancel() {
        onCancel();
        isOpen = false;
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Escape") {
            handleCancel();
        } else if (e.key === "Enter" && e.ctrlKey) {
            handleConfirm();
        }
    }

    function handleBackdropClick(e: MouseEvent) {
        if (e.target === e.currentTarget) {
            handleCancel();
        }
    }
</script>

{#if isOpen}
    <div
        class="modal-backdrop"
        onclick={handleBackdropClick}
        onkeydown={handleKeydown}
        role="presentation"
    >
        <div
            class="modal-container"
            role="dialog"
            aria-modal="true"
            aria-labelledby="modal-title"
        >
            <div class="modal-header">
                <h3 id="modal-title" class="modal-title">{title}</h3>
            </div>

            <div class="modal-body">
                <p class="modal-message">{message}</p>
            </div>

            <div class="modal-footer">
                {#if cancelText}
                    <button class="btn btn-secondary" onclick={handleCancel}>
                        {cancelText}
                    </button>
                {/if}
                <button
                    class="btn btn-primary"
                    class:btn-danger={type === "danger"}
                    class:btn-warning={type === "warning"}
                    onclick={handleConfirm}
                >
                    {confirmText}
                </button>
            </div>
        </div>
    </div>
{/if}

<style lang="scss">
    .modal-backdrop {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.75);
        backdrop-filter: blur(8px);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 9999;
        animation: fadeIn 0.2s ease-out;
    }

    @keyframes fadeIn {
        from {
            opacity: 0;
        }
        to {
            opacity: 1;
        }
    }

    .modal-container {
        background: linear-gradient(
            135deg,
            rgba(30, 30, 40, 0.98),
            rgba(20, 20, 30, 0.98)
        );
        border: 2px solid rgba(173, 64, 255, 0.3);
        border-radius: 16px;
        box-shadow:
            0 20px 60px rgba(0, 0, 0, 0.5),
            0 0 40px rgba(173, 64, 255, 0.2);
        min-width: 400px;
        max-width: 500px;
        animation: slideUp 0.3s ease-out;
        overflow: hidden;
    }

    @keyframes slideUp {
        from {
            transform: translateY(20px);
            opacity: 0;
        }
        to {
            transform: translateY(0);
            opacity: 1;
        }
    }

    .modal-header {
        padding: 1.5rem 2rem;
        border-bottom: 1px solid rgba(173, 64, 255, 0.2);
        background: linear-gradient(
            135deg,
            rgba(173, 64, 255, 0.1),
            rgba(122, 40, 203, 0.05)
        );
    }

    .modal-title {
        margin: 0;
        font-size: 1.5rem;
        font-weight: 700;
        color: #fff;
        letter-spacing: -0.025em;
    }

    .modal-body {
        padding: 2rem;
    }

    .modal-message {
        margin: 0;
        font-size: 1.1rem;
        line-height: 1.6;
        color: rgba(255, 255, 255, 0.9);
    }

    .modal-footer {
        padding: 1.5rem 2rem;
        display: flex;
        justify-content: flex-end;
        gap: 1rem;
        border-top: 1px solid rgba(173, 64, 255, 0.2);
        background: linear-gradient(
            135deg,
            rgba(173, 64, 255, 0.05),
            rgba(122, 40, 203, 0.02)
        );
    }

    .btn {
        padding: 0.75rem 1.5rem;
        border-radius: 8px;
        font-weight: 600;
        font-size: 1rem;
        cursor: pointer;
        transition: all 0.2s ease;
        border: none;
        outline: none;

        &:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
        }

        &:active {
            transform: translateY(0);
        }

        &:focus-visible {
            outline: 2px solid rgba(173, 64, 255, 0.6);
            outline-offset: 2px;
        }
    }

    .btn-primary {
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
            box-shadow: 0 4px 16px rgba(173, 64, 255, 0.4);
        }
    }

    .btn-danger {
        background: linear-gradient(
            135deg,
            rgba(255, 70, 85, 0.8),
            rgba(230, 57, 70, 0.8)
        );
        color: white;
        border: 1px solid rgba(255, 70, 85, 0.4);

        &:hover {
            background: linear-gradient(
                135deg,
                rgba(255, 70, 85, 1),
                rgba(230, 57, 70, 1)
            );
            box-shadow: 0 4px 16px rgba(255, 70, 85, 0.4);
        }
    }

    .btn-warning {
        background: linear-gradient(
            135deg,
            rgba(255, 193, 7, 0.8),
            rgba(255, 160, 0, 0.8)
        );
        color: #000;
        border: 1px solid rgba(255, 193, 7, 0.4);

        &:hover {
            background: linear-gradient(
                135deg,
                rgba(255, 193, 7, 1),
                rgba(255, 160, 0, 1)
            );
            box-shadow: 0 4px 16px rgba(255, 193, 7, 0.4);
        }
    }

    .btn-secondary {
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
            color: #ffffff;
        }
    }
</style>
