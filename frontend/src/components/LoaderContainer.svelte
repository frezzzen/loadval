<script>
    import Loader from "./Loader.svelte";
</script>

<div class="loader-container">
    <div class="loader-content">
        <h1 class="loader-title">Waiting for Valorant to start...</h1>
        <div class="loader-wrapper">
            <Loader />
        </div>
        <p class="loader-subtitle">
            Please ensure Valorant is running and try again
        </p>
    </div>
</div>

<style lang="scss">
    .loader-container {
        display: flex;
        align-items: center;
        justify-content: center;
        min-height: 100vh;
        width: 100%;
        background: linear-gradient(135deg, #070707 0%, #0a0a0a 100%);
        padding: 2rem;
    }

    .loader-content {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        text-align: center;
        max-width: 40rem;
        width: 100%;
    }

    .loader-title {
        font-size: 2.4rem;
        font-weight: 600;
        color: #cccaca;
        margin-bottom: 2rem;
        text-align: center;
        background: linear-gradient(135deg, #cccaca 0%, #a6adc8 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
        animation: fadeInUp 0.8s ease-out;
    }

    .loader-wrapper {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 2rem;
        padding: 2rem;
        background: linear-gradient(
            135deg,
            rgba(255, 255, 255, 0.05),
            rgba(255, 255, 255, 0.02)
        );
        backdrop-filter: blur(10px);
        border-radius: 2rem;
        border: 1px solid rgba(255, 255, 255, 0.1);
        box-shadow:
            0 8px 32px rgba(0, 0, 0, 0.3),
            inset 0 1px 0 rgba(255, 255, 255, 0.1);
        animation: fadeInUp 0.8s ease-out 0.2s both;

        :global svg {
            width: 100%;
            height: 100%;
        }
    }

    .loader-subtitle {
        font-size: 1.4rem;
        color: #a6adc8;
        opacity: 0.8;
        margin: 0;
        animation: fadeInUp 0.8s ease-out 0.4s both;
    }

    @keyframes fadeInUp {
        from {
            opacity: 0;
            transform: translateY(2rem);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    @media (max-width: 768px) {
        .loader-container {
            padding: 1.5rem;
        }

        .loader-title {
            font-size: 2rem;
            margin-bottom: 1.5rem;
        }

        .loader-wrapper {
            padding: 1.5rem;
            margin-bottom: 1.5rem;
        }

        .loader-subtitle {
            font-size: 1.2rem;
        }
    }

    @media (max-width: 480px) {
        .loader-container {
            padding: 1rem;
        }

        .loader-title {
            font-size: 1.8rem;
            margin-bottom: 1rem;
        }

        .loader-wrapper {
            padding: 1rem;
            margin-bottom: 1rem;
        }

        .loader-subtitle {
            font-size: 1.1rem;
        }
    }
</style>
