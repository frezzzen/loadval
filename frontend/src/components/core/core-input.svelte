<script lang="ts">
    type Props = {
        label: string;
        value: string;
        placeholder?: string;
    };

    let { label, value = $bindable(), placeholder = "" }: Props = $props();
</script>

<div class="core-input">
    <input type="text" bind:value {placeholder} />
</div>

<style lang="scss">
    .core-input {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }

    label {
        font-size: 0.95rem;
        font-weight: 600;
        color: rgba(255, 255, 255, 0.8);
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    input {
        font-size: 1rem;
        font-weight: 400;
        color: #ffffff;
        font-family: "DM Sans", sans-serif;
        background: rgba(255, 255, 255, 0.05);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 8px;
        padding: 0.875rem 1rem;
        transition: all 0.2s ease;
        outline: none;

        &::placeholder {
            color: rgba(255, 255, 255, 0.4);
        }

        &:hover:not(:focus) {
            border-color: rgba(255, 255, 255, 0.2);
            background: rgba(255, 255, 255, 0.08);
        }

        &:focus {
            border-color: rgba(173, 64, 255, 0.6);
            background: rgba(255, 255, 255, 0.08);
            box-shadow: 0 0 0 3px rgba(173, 64, 255, 0.1);
        }
    }
</style>
