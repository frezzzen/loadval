<script lang="ts">
    type Props = {
        label: string;
        options: { label: string; value: string }[];
        value: string;
    };

    let { label, options, value = $bindable() }: Props = $props();

    let isOpen = $state(false);
    let selectRef: HTMLDivElement;

    function toggleDropdown() {
        isOpen = !isOpen;
    }

    function selectOption(option: { label: string; value: string }) {
        value = option.value;
        isOpen = false;
    }

    function handleClickOutside(event: MouseEvent) {
        if (selectRef && !selectRef.contains(event.target as Node)) {
            isOpen = false;
        }
    }

    const selectedOption = $derived(
        options.find((option) => option.value === value)
    );

    $effect(() => {
        if (isOpen) {
            document.addEventListener("click", handleClickOutside);
        } else {
            document.removeEventListener("click", handleClickOutside);
        }

        return () => {
            document.removeEventListener("click", handleClickOutside);
        };
    });
</script>

<div class="core-select" bind:this={selectRef}>
    <div class="label">{label}</div>
    <div class="custom-select">
        <button
            type="button"
            class="selected-value"
            class:open={isOpen}
            onclick={toggleDropdown}
        >
            <span>{selectedOption?.label}</span>
            <span class="material-icons" class:rotate={isOpen}
                >arrow_drop_down</span
            >
        </button>

        {#if isOpen}
            <div class="options">
                {#each options as option}
                    <button
                        type="button"
                        class="option"
                        class:selected={option.value === value}
                        onclick={() => selectOption(option)}
                    >
                        {option.label}
                    </button>
                {/each}
            </div>
        {/if}
    </div>
</div>

<style lang="scss">
    .core-select {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        position: relative;
    }

    .label {
        font-size: 1.4rem;
        font-weight: 600;
        color: #cccaca;
        margin-bottom: 0.5rem;
        text-transform: uppercase;
        letter-spacing: 0.1em;
    }

    .custom-select {
        position: relative;
        width: 100%;
    }

    .selected-value {
        width: 100%;
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0.75rem 0.5rem;
        font-size: 1.4rem;
        font-weight: 400;
        color: #cccaca;
        font-family: "DM Sans", sans-serif;
        background-color: transparent;
        border: none;
        border-bottom: 2px solid #3c3c3c;
        cursor: pointer;
        transition:
            border-color 0.3s ease,
            transform 0.3s ease;

        &:hover:not(.open) {
            border-bottom-color: #505050;
        }

        &.open {
            border-bottom-color: #ad40ff;
            transform: translateY(-1px);
        }

        .material-icons {
            transition: transform 0.2s ease;
            font-size: 2.4rem;

            &.rotate {
                transform: rotate(180deg);
            }
        }
    }

    .options {
        position: absolute;
        top: calc(100% + 0.8rem);
        left: 0;
        width: 100%;
        display: flex;
        flex-direction: column;
        background-color: #2c2c2c;
        border: 1px solid #3c3c3c;
        border-radius: 0.8rem;
        overflow: hidden;
        z-index: 1000;
        max-height: 20rem;
        overflow-y: auto;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
    }

    .option {
        width: 100%;
        padding: 0.8rem 1.2rem;
        font-size: 1.4rem;
        font-weight: 400;
        color: #cccaca;
        font-family: "DM Sans", sans-serif;
        background-color: transparent;
        border: none;
        text-align: left;
        cursor: pointer;
        transition: background-color 0.2s ease;

        &:hover {
            background-color: #3c3c3c;
        }

        &.selected {
            background-color: #ad40ff;
            color: #ffffff;

            &:hover {
                background-color: #9930e6;
            }
        }
    }
</style>
