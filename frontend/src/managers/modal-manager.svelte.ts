type ModalConfig = {
    title: string;
    message: string;
    confirmText?: string;
    cancelText?: string;
    type?: "danger" | "warning" | "info";
};

export class ModalManager {
    isOpen = $state(false);
    config = $state<ModalConfig>({
        title: "",
        message: "",
        confirmText: "Confirm",
        cancelText: "Cancel",
        type: "info",
    });

    private resolvePromise: ((value: boolean) => void) | null = null;

    confirm(config: ModalConfig): Promise<boolean> {
        this.config = {
            confirmText: "Confirm",
            cancelText: "Cancel",
            type: "info",
            ...config,
        };
        this.isOpen = true;

        return new Promise((resolve) => {
            this.resolvePromise = resolve;
        });
    }

    handleConfirm() {
        if (this.resolvePromise) {
            this.resolvePromise(true);
            this.resolvePromise = null;
        }
        this.isOpen = false;
    }

    handleCancel() {
        if (this.resolvePromise) {
            this.resolvePromise(false);
            this.resolvePromise = null;
        }
        this.isOpen = false;
    }
}

let modalManager: ModalManager;

export const useModalManager = () => {
    if (!modalManager) {
        modalManager = new ModalManager();
    }
    return modalManager;
};

