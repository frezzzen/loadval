<script lang="ts">
    import { onMount } from "svelte";
    import {
        CheckForUpdates,
        DownloadUpdate,
        InstallUpdate,
        OpenDownloadFolder,
    } from "../../wailsjs/go/main/UpdateAPI";
    import type { main } from "../../wailsjs/go/models";

    let updateInfo = $state<main.UpdateInfo | null>(null);
    let showNotification = $state(false);
    let isDownloading = $state(false);
    let downloadPath = $state("");
    let error = $state("");
    let checking = $state(false);

    onMount(() => {
        checkForUpdates();

        const interval = setInterval(
            () => {
                checkForUpdates();
            },
            6 * 60 * 60 * 1000,
        );

        return () => clearInterval(interval);
    });

    async function checkForUpdates() {
        if (checking) return;

        checking = true;
        error = "";

        try {
            const info = await CheckForUpdates();
            updateInfo = info;

            if (info.available) {
                showNotification = true;
            }
        } catch (err) {
            console.error("Failed to check for updates:", err);
            error = err instanceof Error ? err.message : String(err);
        } finally {
            checking = false;
        }
    }

    async function handleDownload() {
        if (!updateInfo?.downloadUrl) return;

        isDownloading = true;
        error = "";

        try {
            const path = await DownloadUpdate(updateInfo.downloadUrl);
            downloadPath = path;
        } catch (err) {
            console.error("Failed to download update:", err);
            error = err instanceof Error ? err.message : String(err);
        } finally {
            isDownloading = false;
        }
    }

    async function handleInstall() {
        if (!downloadPath) return;

        try {
            await InstallUpdate(downloadPath);
        } catch (err) {
            console.error("Failed to install update:", err);
            error = err instanceof Error ? err.message : String(err);
        }
    }

    async function handleOpenFolder() {
        try {
            await OpenDownloadFolder();
        } catch (err) {
            console.error("Failed to open folder:", err);
        }
    }

    function formatFileSize(bytes: number): string {
        if (bytes === 0) return "0 Bytes";
        const k = 1024;
        const sizes = ["Bytes", "KB", "MB", "GB"];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return (
            Math.round((bytes / Math.pow(k, i)) * 100) / 100 + " " + sizes[i]
        );
    }

    function dismiss() {
        showNotification = false;
    }
</script>

{#if showNotification && updateInfo?.available}
    <div class="update-notification">
        <div class="notification-content">
            <div class="notification-header">
                <div class="header-left">
                    <svg
                        width="24"
                        height="24"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        class="update-icon"
                    >
                        <path
                            d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2"
                        ></path>
                    </svg>
                    <div class="title-section">
                        <h3>Update Available</h3>
                        <p class="version-info">
                            v{updateInfo.latestVersion}
                            {#if updateInfo.publishedAt}
                                · {updateInfo.publishedAt}
                            {/if}
                        </p>
                    </div>
                </div>
                <button class="close-btn" onclick={dismiss} aria-label="Close">
                    <svg
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <line x1="18" y1="6" x2="6" y2="18"></line>
                        <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                </button>
            </div>

            {#if updateInfo.releaseNotes}
                <div class="release-notes">
                    <p>
                        {updateInfo.releaseNotes.substring(0, 200)}{updateInfo
                            .releaseNotes.length > 200
                            ? "..."
                            : ""}
                    </p>
                </div>
            {/if}

            {#if error}
                <div class="error-message">
                    <svg
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <circle cx="12" cy="12" r="10"></circle>
                        <line x1="12" y1="8" x2="12" y2="12"></line>
                        <line x1="12" y1="16" x2="12.01" y2="16"></line>
                    </svg>
                    <span>{error}</span>
                </div>
            {/if}

            <div class="notification-actions">
                {#if downloadPath}
                    <button class="btn btn-success" onclick={handleInstall}>
                        <svg
                            width="18"
                            height="18"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                        >
                            <polyline points="20 6 9 17 4 12"></polyline>
                        </svg>
                        <span>Install Now</span>
                    </button>
                    <button
                        class="btn btn-secondary"
                        onclick={handleOpenFolder}
                    >
                        <span>Open Folder</span>
                    </button>
                {:else if isDownloading}
                    <button class="btn btn-primary" disabled>
                        <div class="spinner"></div>
                        <span>Downloading...</span>
                    </button>
                {:else}
                    <button class="btn btn-primary" onclick={handleDownload}>
                        <svg
                            width="18"
                            height="18"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                        >
                            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"
                            ></path>
                            <polyline points="7 10 12 15 17 10"></polyline>
                            <line x1="12" y1="15" x2="12" y2="3"></line>
                        </svg>
                        <span
                            >Download ({formatFileSize(
                                updateInfo.downloadSize,
                            )})</span
                        >
                    </button>
                {/if}
                <button class="btn btn-text" onclick={dismiss}>
                    <span>Later</span>
                </button>
            </div>
        </div>
    </div>
{/if}

<style lang="scss">
    .update-notification {
        position: fixed;
        top: 20px;
        right: 20px;
        z-index: 9999;
        max-width: 450px;
        animation: slideIn 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }

    @keyframes slideIn {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }

    .notification-content {
        background: linear-gradient(
            135deg,
            rgba(173, 64, 255, 0.15),
            rgba(122, 40, 203, 0.1)
        );
        border: 1px solid rgba(173, 64, 255, 0.3);
        border-radius: 12px;
        padding: 1.5rem;
        box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
        backdrop-filter: blur(10px);
    }

    .notification-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 1rem;

        .header-left {
            display: flex;
            align-items: flex-start;
            gap: 0.75rem;

            .update-icon {
                color: #ad40ff;
                flex-shrink: 0;
                margin-top: 2px;
            }

            .title-section {
                h3 {
                    margin: 0;
                    font-size: 1.125rem;
                    font-weight: 700;
                    color: #fff;
                    letter-spacing: -0.025em;
                }

                .version-info {
                    margin: 0.25rem 0 0 0;
                    font-size: 0.875rem;
                    color: rgba(255, 255, 255, 0.6);
                    font-weight: 500;
                }
            }
        }

        .close-btn {
            background: none;
            border: none;
            color: rgba(255, 255, 255, 0.6);
            cursor: pointer;
            padding: 0.25rem;
            border-radius: 4px;
            transition: all 0.2s ease;
            flex-shrink: 0;

            &:hover {
                color: rgba(255, 255, 255, 0.9);
                background: rgba(255, 255, 255, 0.1);
            }
        }
    }

    .release-notes {
        margin: 0 0 1rem 0;
        padding: 0.75rem;
        background: rgba(0, 0, 0, 0.2);
        border-radius: 8px;
        border-left: 3px solid #ad40ff;

        p {
            margin: 0;
            font-size: 0.875rem;
            color: rgba(255, 255, 255, 0.8);
            line-height: 1.5;
        }
    }

    .error-message {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        margin: 0 0 1rem 0;
        padding: 0.75rem;
        background: rgba(255, 70, 85, 0.15);
        border: 1px solid rgba(255, 70, 85, 0.3);
        border-radius: 8px;
        color: #ff6b7a;
        font-size: 0.875rem;

        svg {
            flex-shrink: 0;
        }

        span {
            flex: 1;
        }
    }

    .notification-actions {
        display: flex;
        gap: 0.75rem;
        flex-wrap: wrap;
    }

    .btn {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1.25rem;
        border-radius: 8px;
        font-weight: 600;
        font-size: 0.9375rem;
        cursor: pointer;
        transition: all 0.2s ease;
        border: none;
        outline: none;
        white-space: nowrap;

        svg {
            flex-shrink: 0;
        }

        &:disabled {
            opacity: 0.6;
            cursor: not-allowed;
        }

        &.btn-primary {
            background: linear-gradient(
                135deg,
                rgba(173, 64, 255, 0.8),
                rgba(122, 40, 203, 0.8)
            );
            color: white;
            border: 1px solid rgba(173, 64, 255, 0.4);

            &:hover:not(:disabled) {
                background: linear-gradient(
                    135deg,
                    rgba(173, 64, 255, 1),
                    rgba(122, 40, 203, 1)
                );
                box-shadow: 0 4px 16px rgba(173, 64, 255, 0.4);
                transform: translateY(-1px);
            }
        }

        &.btn-success {
            background: linear-gradient(
                135deg,
                rgba(70, 243, 85, 0.8),
                rgba(40, 203, 55, 0.8)
            );
            color: rgba(0, 0, 0, 0.9);
            border: 1px solid rgba(70, 243, 85, 0.4);

            &:hover:not(:disabled) {
                background: linear-gradient(
                    135deg,
                    rgba(70, 243, 85, 1),
                    rgba(40, 203, 55, 1)
                );
                box-shadow: 0 4px 16px rgba(70, 243, 85, 0.4);
                transform: translateY(-1px);
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

            &:hover:not(:disabled) {
                background: linear-gradient(
                    135deg,
                    rgba(255, 255, 255, 0.15),
                    rgba(200, 200, 200, 0.1)
                );
                border-color: rgba(255, 255, 255, 0.4);
            }
        }

        &.btn-text {
            background: none;
            color: rgba(255, 255, 255, 0.7);
            border: none;
            padding: 0.75rem 1rem;

            &:hover:not(:disabled) {
                color: rgba(255, 255, 255, 0.9);
                background: rgba(255, 255, 255, 0.05);
            }
        }
    }

    .spinner {
        width: 16px;
        height: 16px;
        border: 2px solid rgba(255, 255, 255, 0.3);
        border-top-color: white;
        border-radius: 50%;
        animation: spin 0.6s linear infinite;
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
</style>
